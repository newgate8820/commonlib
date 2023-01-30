package dbcommon

import "fmt"

/*************************
* author: Dev0026
* createTime: 19-4-2
* updateTime: 19-4-2
* description:
*************************/

func GetRedisInfoFromLinkMap(platform string) *redisInfoConfig {
	return getRedisInfoFromLinkMap(platform)
}

func getRedisInfoFromLinkMap(platform string) *redisInfoConfig {
	key := fmt.Sprintf("redisinfo_%v", platform)

	if val, ok := linkMap.Load(key); ok {
		return val.(*redisInfoConfig)
	} else {
		return nil
	}
}
