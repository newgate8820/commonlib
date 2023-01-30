package dbcommon

import "fmt"

/*************************
* author: Dev0026
* createTime: 19-4-2
* updateTime: 19-4-2
* description:
*************************/

func HandleXormGroupConfigEncryFields(cfgs ...*XormGroupConfig) error {
	for i, cfg := range cfgs {
		for j, cfg := range cfg.XormConfigs {
			info := getMysqlInfoFromLinkMap(cfg.Platform)
			if nil != info {
				if cfg.ParamStr == "" {
					cfgs[i].XormConfigs[j].DataSourceName = fmt.Sprintf("%v:%v@tcp(%v)/%v",
						info.UserName, info.UserPasswd, info.HostInfo, cfg.SchemaName)
				} else {
					cfgs[i].XormConfigs[j].DataSourceName = fmt.Sprintf("%v:%v@tcp(%v)/%v?%v",
						info.UserName, info.UserPasswd, info.HostInfo, cfg.SchemaName, cfg.ParamStr)
				}
			}
		}
	}
	return nil
}

func HandleXormConfigEncryFields(cfgs ...*XormConfig) error {
	for index, cfg := range cfgs {
		info := getMysqlInfoFromLinkMap(cfg.Platform)
		if cfg.DriverName == "" {
			cfg.DriverName = "mysql"
		}
		if nil != info {
			if cfg.ParamStr == "" {
				cfgs[index].DataSourceName = fmt.Sprintf("%v:%v@tcp(%v)/%v?charset?utf8mb4",
					info.UserName, info.UserPasswd, info.HostInfo, cfg.SchemaName)
			} else {
				cfgs[index].DataSourceName = fmt.Sprintf("%v:%v@tcp(%v)/%v?%v",
					info.UserName, info.UserPasswd, info.HostInfo, cfg.SchemaName, cfg.ParamStr)
			}
		}
	}
	return nil
}
