package dbcommon

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/go-redis/redis"
	"github.com/xormplus/xorm"
)

/*************************
* author: Dev0026
* createTime: 19-4-2
* updateTime: 19-4-2
* description:
*************************/

var xormEngines []*xorm.Engine
var redisEngines []*redis.ClusterClient

func init() {
	xormEngines = make([]*xorm.Engine, 0)
	redisEngines = make([]*redis.ClusterClient, 0)
}

func Close() {
	for _, xormEngine := range xormEngines {
		xormEngine.Close()
	}
	log.Println("-------------- dataLayer db close success --------------")

	for _, cluster := range redisEngines {
		cluster.Close()
	}
	log.Println("-------------- dataLayer redis close success --------------")

	CloseKafka()
	log.Println("-------------- dataLayer kafka close success --------------")
}

func HandleXormConfigWithDefault(defaultCfg XormDefaultConfig, cfgs ...*XormGroupConfig) {
	for _, cfg := range cfgs {
		if cfg.GroupPolicy == "" {
			cfg.GroupPolicy = defaultCfg.GroupPolicy
		}

		if len(cfg.GroupPolicyWeight) == 0 {
			cfg.GroupPolicyWeight = defaultCfg.GroupPolicyWeight
		}

		for index, xormCfg := range cfg.XormConfigs {
			if xormCfg.Platform == "" {
				cfg.XormConfigs[index].Platform = defaultCfg.XormConfig.Platform
			}
			if xormCfg.DriverName == "" {
				cfg.XormConfigs[index].DriverName = defaultCfg.XormConfig.DriverName
			}

			if xormCfg.MaxIdleConns == 0 {
				cfg.XormConfigs[index].MaxIdleConns = defaultCfg.XormConfig.MaxIdleConns
			}
			if xormCfg.MaxOpenConns == 0 {
				cfg.XormConfigs[index].MaxOpenConns = defaultCfg.XormConfig.MaxOpenConns
			}
			if xormCfg.ShowSql == false {
				cfg.XormConfigs[index].ShowSql = defaultCfg.XormConfig.ShowSql
			}
			if xormCfg.LogLevel == 0 {
				cfg.XormConfigs[index].LogLevel = defaultCfg.XormConfig.LogLevel
			}
			if xormCfg.SqlMapOptionRootDir == "" {
				cfg.XormConfigs[index].SqlMapOptionRootDir = defaultCfg.XormConfig.SqlMapOptionRootDir
			}
			if xormCfg.SqlTemplateDir == "" {
				cfg.XormConfigs[index].SqlTemplateDir = defaultCfg.XormConfig.SqlTemplateDir
			}

			if xormCfg.ParamStr == "" {
				cfg.XormConfigs[index].ParamStr = defaultCfg.XormConfig.ParamStr
			}

			cfg.XormConfigs[index].DataSourceName = ""

			//log.Println("---------------> @ cfg.XormConfigs[index] ", cfg.XormConfigs[index])
		}
	}
}

func HandleRedisConfigWithDefault(defaultCfg RedisConfig, cfgs ...*RedisConfig) {
	for _, cfg := range cfgs {
		if len(cfg.Addrs) == 0 {
			cfg.Addrs = defaultCfg.Addrs
		}

		if cfg.MaxRedirects == 0 {
			cfg.MaxRedirects = defaultCfg.MaxRedirects
		}

		if cfg.ReadOnly == false {
			cfg.ReadOnly = defaultCfg.ReadOnly
		}
		if cfg.RouteByLatency == false {
			cfg.RouteByLatency = defaultCfg.RouteByLatency
		}
		if cfg.MaxRetries == 0 {
			cfg.MaxRetries = defaultCfg.MaxRetries
		}
		if cfg.Password == "" {
			cfg.Password = defaultCfg.Password
		}
		if cfg.IdleTimeout == 0 {
			cfg.IdleTimeout = defaultCfg.IdleTimeout
		}
		if cfg.PoolSize == 0 {
			cfg.PoolSize = defaultCfg.PoolSize
		}
		if cfg.LuaScriptDir == "" {
			cfg.LuaScriptDir = defaultCfg.LuaScriptDir
		}
		if cfg.IdleCheckFrequency == 0 {
			cfg.IdleCheckFrequency = defaultCfg.IdleCheckFrequency
		}
		if cfg.RouteRandomly == false {
			cfg.RouteRandomly = defaultCfg.RouteRandomly
		}
		cfg.Addrs = nil
		cfg.Password = ""
	}
}

func HandleRedisConfigEncryFields(cfgs ...*RedisConfig) error {
	for index, cfg := range cfgs {
		info := getRedisInfoFromLinkMap(cfg.Platform)
		if nil != info {
			cfgs[index].Password = info.Passwd
			cfgs[index].Addrs = strings.Split(info.HostInfo, ",")
		}
	}
	return nil
}

func HandleNatsConfigFields(cfgs ...*NatsConfig) error {
	for index, cfg := range cfgs {
		if info := getNatsInfoFromLinkMap(cfg.Platform); nil != info {
			cfgs[index].User = info.UserName
			cfgs[index].Password = info.Passwd
			cfgs[index].Url = info.HostInfo
		}
	}
	return nil
}

func InitLinkMap(cfg *XormGroupConfig) error {
	fmt.Println("[InitLinkMap] ############Init DB Connections################")
	// ############################################## init mysql
	var xormCfg XormConfig
	if 0 < len(cfg.XormConfigs) {
		xormCfg = cfg.XormConfigs[0]
	} else {
		return fmt.Errorf("xormConfig len is 0")
	}
	engine, err := xorm.NewEngine(xormCfg.DriverName, xormCfg.DataSourceName)
	if nil != err {
		return fmt.Errorf("initLinkMap new engine err: %v", err)
	}
	defer engine.Close()
	// init mysql
	//var mysqlInfos = make([]mysqlInfoConfig, 0, 10)
	//err = engine.SQL("select user_name, user_passwd, host_info, platform from mysql_info").Find(&mysqlInfos)
	//if nil != err {
	//	return fmt.Errorf("initLinkMap mysqlInfo err: %v", err)
	//}
	//for index := range mysqlInfos {
	//	setMysqlInfoToLinkMap(&mysqlInfos[index])
	//}
	err = InitializeMysqlLinkMap(xormCfg.DriverName, xormCfg.DataSourceName)
	if err != nil {
		return fmt.Errorf("initLinkMap mysqlInfo err: %v", err)
	}

	// ############################################## init redis
	var redisInfos = make([]redisInfoConfig, 0, 2)
	err = engine.SQL("select host_info, passwd, platform from redis_info").Find(&redisInfos)
	if nil != err {
		return fmt.Errorf("initLinkMap redisInfo err: %v", err)
	}
	for index := range redisInfos {
		setRedisInfoToLinkMap(&redisInfos[index])
	}

	// ############################################## init nats
	//var natsInfos = make([]natsInfoConfig, 0, 2)
	//err = engine.SQL("select user_name, passwd, host_info, platform from nats_info").Find(&natsInfos)
	//if nil != err {
	//	return fmt.Errorf("initLinkMap natsinfo err: %v", err)
	//}
	//for index := range natsInfos {
	//	setNatsInfoToLinkMap(&natsInfos[index])
	//}
	err = InitializeNatsLinkMap(xormCfg.DriverName, xormCfg.DataSourceName)
	if nil != err {
		return fmt.Errorf("initLinkMap natsinfo err: %v", err)
	}

	// ############################################## kafka info kafkaInfoConfig
	//fmt.Println("[InitLinkMap] ############Init to load kafka configs################")
	//var kafkaInfos = make([]kafkaInfoConfig, 0, 2)
	//err = engine.SQL("select kafka_info, zookeeper_info, platform from kafka_info").Find(&kafkaInfos)
	//if nil != err {
	//	fmt.Printf("initLinkMap kafkaInfos err: %v", err)
	//	return fmt.Errorf("initLinkMap kafkaInfos err: %v", err)
	//}
	//for index := range kafkaInfos {
	//	setKafkaInfoToLinkMap(&kafkaInfos[index])
	//}
	fmt.Println("linkinfo.[InitLinkMap]############################################## payserver_cfg info")
	var payserverCfgInfos = make([]payserverCfgConfig, 0, 2)
	err = engine.SQL("select pay_servers, platform from payserver_cfg").Find(&payserverCfgInfos)
	if nil != err {
		return fmt.Errorf("initLinkMap payserverCfgInfos err: %v", err)
	}
	for index := range payserverCfgInfos {
		setPayserverCfgToLinkMap(&payserverCfgInfos[index])
	}
	//fmt.Printf("[InitLinkMap] ############Init load kafka configs:%v\n", kafkaInfos)
	//var payserverCfgInfos = make([]payserverCfgConfig, 0, 2)
	//err = engine.SQL("select pay_servers, platform from payserver_cfg").Find(&payserverCfgInfos)
	//if nil != err {
	//	return fmt.Errorf("initLinkMap payserverCfgInfos err: %v", err)
	//}
	//for index := range payserverCfgInfos {
	//	setPayserverCfgToLinkMap(&payserverCfgInfos[index])
	//}
	//fmt.Printf("[InitLinkMap] ############Init load kafka configs:%v\n", kafkaInfos)
	return nil
}

func InitLinkMapByDataSourceName(dataSourceName string) error {
	fmt.Println("[InitLinkMapByDataSourceName] ############Init LinkInfomation DB Connections################")
	//linkinfo.[InitLinkMap] ############################################## init mysql
	driverName := "mysql"

	engine, err := xorm.NewEngine(driverName, dataSourceName)
	if nil != err {
		return fmt.Errorf("initLinkMap new engine err: %v", err)
	}
	defer engine.Close()

	err = InitializeEtcdLinkMap(driverName, dataSourceName)
	if err != nil {
		return fmt.Errorf("initLinkMap  InitializeEtcdLinkMap mysqlInfo err: %v", err)
	}

	// init mysql
	//var mysqlInfos = make([]mysqlInfoConfig, 0, 10)
	//err = engine.SQL("select user_name, user_passwd, host_info, platform from mysql_info").Find(&mysqlInfos)
	//if nil != err {
	//	return fmt.Errorf("initLinkMap mysqlInfo err: %v", err)
	//}
	//for index := range mysqlInfos {
	//	setMysqlInfoToLinkMap(&mysqlInfos[index])
	//}
	err = InitializeMysqlLinkMap(driverName, dataSourceName)
	if err != nil {
		return fmt.Errorf("initLinkMap InitializeMysqlLinkMap  mysqlInfo err: %v", err)
	}

	// ############################################## init redis
	var redisInfos = make([]redisInfoConfig, 0, 2)
	err = engine.SQL("select host_info, passwd, platform from redis_info").Find(&redisInfos)
	if nil != err {
		return fmt.Errorf("initLinkMap redisInfo err: %v", err)
	}
	for index := range redisInfos {
		setRedisInfoToLinkMap(&redisInfos[index])
	}

	// ############################################## init nats
	//var natsInfos = make([]natsInfoConfig, 0, 2)
	//err = engine.SQL("select user_name, passwd, host_info, platform from nats_info").Find(&natsInfos)
	//if nil != err {
	//	return fmt.Errorf("initLinkMap natsinfo err: %v", err)
	//}
	//for index := range natsInfos {
	//	setNatsInfoToLinkMap(&natsInfos[index])
	//}
	err = InitializeNatsLinkMap(driverName, dataSourceName)
	if nil != err {
		return fmt.Errorf("initLinkMap natsinfo err: %v", err)
	}

	// ############################################## kafka info kafkaInfoConfig
	fmt.Println("[InitLinkMap] ############Init to load kafka configs################")
	var kafkaInfos = make([]kafkaInfoConfig, 0, 2)
	err = engine.SQL("select kafka_info, zookeeper_info, platform from kafka_info").Find(&kafkaInfos)
	if nil != err {
		fmt.Printf("initLinkMap kafkaInfos err: %v", err)
		return fmt.Errorf("initLinkMap kafkaInfos err: %v", err)
	}
	for index := range kafkaInfos {
		setKafkaInfoToLinkMap(&kafkaInfos[index])
	}
	// ############################################## payserver_cfg info
	var payserverCfgInfos = make([]payserverCfgConfig, 0, 2)
	err = engine.SQL("select pay_servers, platform from payserver_cfg").Find(&payserverCfgInfos)
	if nil != err {
		return fmt.Errorf("initLinkMap payserverCfgInfos err: %v", err)
	}
	for index := range payserverCfgInfos {
		setPayserverCfgToLinkMap(&payserverCfgInfos[index])
	}
	// fmt.Printf("[InitLinkMap] ############Init load kafka configs:%v\n", kafkaInfos)
	//var payserverCfgInfos = make([]payserverCfgConfig, 0, 2)
	//err = engine.SQL("select pay_servers, platform from payserver_cfg").Find(&payserverCfgInfos)
	//if nil != err {
	//	return fmt.Errorf("initLinkMap payserverCfgInfos err: %v", err)
	//}
	//for index := range payserverCfgInfos {
	//	setPayserverCfgToLinkMap(&payserverCfgInfos[index])
	//}
	//fmt.Printf("[InitLinkMap] ############Init load kafka configs:%v\n", kafkaInfos)

	// ############################################## cockdb info
	err = initCRDBInfo(engine)
	if nil != err {
		return err
	}

	return nil
}

//func setMysqlInfoToLinkMap(info *mysqlInfoConfig) {
//	var (
//		key = fmt.Sprintf("mysqlinfo_%v", info.Platform)
//		err error
//	)
//	info.HostInfo, err = RsaDecrypt(info.HostInfo)
//	if nil != err {
//		log.Fatalf("RsaDecrypt mysqlInfo [ info %+v ]hostInfo err: %v", info, err)
//	}
//
//	info.UserPasswd, err = RsaDecrypt(info.UserPasswd)
//	if nil != err {
//		log.Fatalf("RsaDecrypt mysqlInfo [ info %+v ] UserPasswd err: %v", info, err)
//	}
//
//	info.UserName, err = RsaDecrypt(info.UserName)
//	if nil != err {
//		log.Fatalf("RsaDecrypt mysqlInfo [ info %+v ] UserName err: %v", info, err)
//	}
//
//	linkMap.Store(key, info)
//}

func setRedisInfoToLinkMap(info *redisInfoConfig) {
	var (
		key = fmt.Sprintf("redisinfo_%v", info.Platform)
		err error
	)
	info.HostInfo, err = RsaDecrypt(info.HostInfo)
	if nil != err {
		log.Fatalf("RsaDecrypt redisInfo [ info %+v ] HostInfo err: %v", info, err)
	}

	if info.Passwd != "" {
		info.Passwd, err = RsaDecrypt(info.Passwd)
		if nil != err {
			log.Fatalf("RsaDecrypt redisInfo [ info %+v ] Passwd err: %v", info, err)
		}
	}

	linkMap.Store(key, info)
}

//func setNatsInfoToLinkMap(info *natsInfoConfig) {
//	var (
//		key = fmt.Sprintf("natsinfo_%v", info.Platform)
//		err error
//	)
//	info.UserName, err = RsaDecrypt(info.UserName)
//	if nil != err {
//		log.Fatalf("RsaDecrypt natsInfo [ info %+v ] UserName err: %v", info, err)
//	}
//
//	info.HostInfo, err = RsaDecrypt(info.HostInfo)
//	if nil != err {
//		log.Fatalf("RsaDecrypt natsInfo [ info %+v ] HostInfo err: %v", info, err)
//	}
//
//	if info.Passwd != "" {
//		info.Passwd, err = RsaDecrypt(info.Passwd)
//		if nil != err {
//			log.Fatalf("RsaDecrypt natsInfo [ info %+v ] Passwd err: %v", info, err)
//		}
//	}
//
//	linkMap.Store(key, info)
//}

func getNatsInfoFromLinkMap(platform string) *natsInfoConfig {
	info, ok := GetNatsPlatformInfo(platform)
	if ok {
		return (*natsInfoConfig)(info)
	} else {
		return nil
	}
	key := fmt.Sprintf("natsinfo_%v", platform)
	if val, ok := linkMap.Load(key); ok {
		return val.(*natsInfoConfig)
	} else {
		return nil
	}
}

func setKafkaInfoToLinkMap(info *kafkaInfoConfig) {
	var (
		key = fmt.Sprintf("kafka_%v", info.Platform)
		err error
	)
	info.KafkaInfo, err = RsaDecrypt(info.KafkaInfo)
	if nil != err {
		log.Fatalf("RsaDecrypt KafkaInfo [ info %+v ] KafkaInfo err: %v", info, err)
	}

	info.ZookeeperInfo, err = RsaDecrypt(info.ZookeeperInfo)
	if nil != err {
		log.Fatalf("RsaDecrypt KafkaInfo [ info %+v ] ZookeeperInfo err: %v", info, err)
	}

	linkMap.Store(key, info)
}

func GetKafkaInfoFromLinkMap(platform string) *kafkaInfoConfig {
	key := fmt.Sprintf("kafka_%v", platform)
	if val, ok := linkMap.Load(key); ok {
		return val.(*kafkaInfoConfig)
	} else {
		return nil
	}
}

func setPayserverCfgToLinkMap(info *payserverCfgConfig) {
	var (
		key = fmt.Sprintf("payservercfg_%v", info.Platform)
		err error
	)
	info.PayServers, err = RsaDecrypt(info.PayServers)
	if nil != err {
		log.Fatalf("RsaDecrypt payserverCfgConfig [ info %+v ] PayServers err: %v", info, err)
	}

	linkMap.Store(key, info)
}

func GetPayserverCfgFromLinkMap(platform string) *payserverCfgConfig {
	key := fmt.Sprintf("payservercfg_%v", platform)
	if val, ok := linkMap.Load(key); ok {
		return val.(*payserverCfgConfig)
	} else {
		return nil
	}
}

//type mysqlInfoConfig struct {
//	UserName   string // 用户名
//	UserPasswd string // 密码
//	HostInfo   string // 链接信息，格式ip:port
//	Platform   string // 所属业务
//}

type mysqlInfoConfig LinkInfo

type redisInfoConfig struct {
	HostInfo string // redis集群链接信息,格式ip:port,ip:port...
	Passwd   string // redis密码
	Platform string // 所属业务
}

//	type natsInfoConfig struct {
//		UserName string
//		Passwd   string
//		HostInfo string
//		Platform string
//	}
type natsInfoConfig NatsLinkInfo

type kafkaInfoConfig struct {
	KafkaInfo     string
	ZookeeperInfo string
	Platform      string
}

type payserverCfgConfig struct {
	PayServers string
	Platform   string
}

var (
	linkMap sync.Map
)

//var privateKey = []byte(`
//-----BEGIN RSA PRIVATE KEY-----
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
//B2/X8LVUbWjJ
//-----END RSA PRIVATE KEY-----
//`)
