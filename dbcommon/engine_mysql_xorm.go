package dbcommon

import (
	"fmt"
	"github.com/xormplus/xorm/log"

	"github.com/xormplus/xorm"
	"time"
)

/*************************
* author: Dev0026
* createTime: 19-4-2
* updateTime: 19-4-2
* description:
*************************/

//go:generate stringer -type=XormEngineGroupPolicy
type XormEngineGroupPolicy int

const (
	RandomPolicy           XormEngineGroupPolicy = iota // 随机访问负载策略
	WeightRandomPolicy                                  // 权重随机访问负载策略
	RoundRobinPolicy                                    // 轮询访问负载策略
	WeightRoundRobinPolicy                              // 权重轮询访问负载策略
)

func NewXormEngineGroup(cfg XormGroupConfig) (*xorm.EngineGroup, error) {
	if len(cfg.XormConfigs) == 0 {
		return nil, fmt.Errorf("NewXormEngineGroup [ %+v ] len cfg.XormConfigs is 0", cfg)
	}

	var engineGroup *xorm.EngineGroup
	var err error
	var master *xorm.Engine
	var slaves []*xorm.Engine

	for index, xormCfg := range cfg.XormConfigs {
		if 0 == index {
			master, err = NewXormEngine(xormCfg)
			if nil != err {
				return nil, err
			}
		} else {
			slave, err := NewXormEngine(xormCfg)
			if nil != err {
				return nil, err
			} else {
				slaves = append(slaves, slave)
			}
		}
	}

	// 负载均衡策略
	switch cfg.GroupPolicy {
	case "":
		engineGroup, err = xorm.NewEngineGroup(master, slaves)
	case RandomPolicy.String():
		engineGroup, err = xorm.NewEngineGroup(master, slaves, xorm.RandomPolicy())
	case WeightRandomPolicy.String():
		engineGroup, err = xorm.NewEngineGroup(master, slaves, xorm.WeightRandomPolicy(cfg.GroupPolicyWeight))
	case RoundRobinPolicy.String():
		engineGroup, err = xorm.NewEngineGroup(master, slaves, xorm.RoundRobinPolicy())
	case WeightRoundRobinPolicy.String():
		engineGroup, err = xorm.NewEngineGroup(master, slaves, xorm.WeightRoundRobinPolicy(cfg.GroupPolicyWeight))
	default:
		err = fmt.Errorf("NewXormEngineGroup not support GroupPolicy %v", cfg.GroupPolicy)
	}
	if nil != err {
		er := fmt.Errorf("NewXormEngineGroup [ cfgs %+v ] xorm.NewEngineGroup err: %v", cfg, err)
		return nil, er
	}
	return engineGroup, nil
}

func NewXormEngine(cfg XormConfig) (*xorm.Engine, error) {
	if cfg.DataSourceName == "" && cfg.Platform != "" {
		err := HandleXormConfigEncryFields(&cfg)
		if nil != err {
			return nil, err
		}
	}
	engine, err := xorm.NewEngine(cfg.DriverName, cfg.DataSourceName)
	if nil != err {
		err := fmt.Errorf("NewXormEngine [ cfg %+v ] err: %v", cfg, err)
		return nil, err
	}

	if cfg.MaxIdleConns > 0 {
		engine.SetMaxIdleConns(cfg.MaxIdleConns)
	}

	if cfg.MaxOpenConns > 0 {
		engine.SetMaxOpenConns(cfg.MaxOpenConns)
	}

	if cfg.ShowSql {
		engine.ShowSQL(true)
	}

	if "" != cfg.SqlMapOptionRootDir {
		err = engine.RegisterSqlMap(xorm.Xml(cfg.SqlMapOptionRootDir, ".xml"))
		if nil != err {
			err := fmt.Errorf("NewXormEngine [ cfg %+v ] RegisterSqlMap err: %v", cfg, err)
			return nil, err
		}
		engine.SqlTemplate = &xorm.JetTemplate{}
	}

	if "" != cfg.SqlTemplateDir {
		err = engine.RegisterSqlTemplate(xorm.Jet(cfg.SqlTemplateDir, ".jet"))
		if nil != err {
			err := fmt.Errorf("NewXormEngine [ cfg %+v ] RegisterSqlTemplate err: %v", cfg, err)
			return nil, err
		}
	}

	//err = engine.StartFSWatcher()
	//if nil != err {
	//	err := fmt.Errorf("NewXormEngine [ cfg %+v ] StartFSWatcher err: %v", cfg, err)
	//	return nil, err
	//}

	engine.SetLogLevel(log.LogLevel(cfg.LogLevel))
	engine.SetConnMaxLifetime(time.Minute * 5)

	err = engine.Ping()
	if nil != err {
		err := fmt.Errorf("NewXormEngine [ cfg %+v ] Ping err: %v", cfg, err)
		return nil, err
	} else {
		//fmt.Printf("%s database %s start success \n", cfg.DriverName, strings.SplitN(cfg.DataSourceName, "@", 2)[1])
	}
	xormEngines = append(xormEngines, engine)
	return engine, nil
}
