package dbcommon

import (
	"fmt"
)

/*************************
* author: Dev0026
* createTime: 19-4-2
* updateTime: 19-4-2
* description:
*************************/

func GetMysqlInfoFromLinkMap(platform string) *mysqlInfoConfig {
	return getMysqlInfoFromLinkMap(platform)
}

func getMysqlInfoFromLinkMap(platform string) *mysqlInfoConfig {
	info, ok := GetMysqlPlatformInfo(platform)
	if ok {
		return (*mysqlInfoConfig)(info)
	} else {
		return nil
	}
	key := fmt.Sprintf("mysqlinfo_%v", platform)
	if val, ok := linkMap.Load(key); ok {
		return val.(*mysqlInfoConfig)
	} else {
		return nil
	}
}
