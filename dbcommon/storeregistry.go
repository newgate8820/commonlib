package dbcommon

import (
	"github.com/nats-io/nats.go"
	"sync"

	"github.com/go-redis/redis"
	"github.com/olivere/elastic"
	"github.com/xormplus/xorm"
)

//go:generate stringer -type=registryName -output=storeregistry_string.go
type registryName int32

const (
	_ registryName = iota + 4000

	MessagedbDatabase
	MessagedbShardDatabase
	MessagedbEncryDatabase
	MessagedbOldDatabase
	MessagedbPtsHandleDatabase

	ChanneldbDatabase
	ChanneldbMsgShardDatabase
	ChanneldbEncryDatabase
	ChanneldbOldDatabase
	ChanneldbPtsHandleDatabase

	PotatoDcConfigDatabase

	PotatoAccountDatabase
	PotatoUserDatabase
	PotatoChannelDatabase
	PotatoGroupDatabase
	PotatoMonitorDatabase
	PotatoCommonDatabase
	PotatoAppliedDatabase
	PotatoSecurityChatDatabase
	PotatoPrivateChatHistoryDatabase
	PotatoSecurityChatDiffMessageDatabase
	PotatoConfigDatabase
	PotatoProxyConfigDatabase
	PotatoUserGpsDatabase
	PotatoUserBot // schema potato_user_bot
	PotatoUserSessionDatabase
	PotatoEmailDatabase
	GeoDatabase
	OpenPlatformMgrDatabase
	GameApiDatabase
	ShortVideoDatabase
	PayServerMysql
	// ========================= redis =========================
	MessageRedis
	MessageShardRedis
	MessageRedisScri //= "messageRe"
	UserInfoRedis
	AuthorizationRedis // 用户session redis
	BasicGroupRedis    // 普通群、超级群信息专用redis
	MediaPushRedis
	MediaPullRedis
	MesageIncreaseRedis
	SecrityChatRedis                 = MessageRedis
	ChannelMessageRedis registryName = 4000 + iota
	ChannelMessageShardRedis
	GeoRedis
	OpenPlatformRedis
	GameApiRedis
	ShortVideoRedis
	PayServerRedis
	// ========================= nats =========================
	PtNats
	PtRouteNats
	PtMemNats
	PtMomentNats
	PtCircleFriend

	// ========================= es =========================
	PtRedisES
	Unknown
)

var (
	registry = sync.Map{}
)

// ************************ base function ******************
func GetDatabase(dbName registryName) *xorm.EngineGroup {
	if r, ok := registry.Load(dbName); ok {
		if result, ok := r.(*xorm.EngineGroup); ok {
			return &(*result)
		}
	}
	return nil
}

func GetRedisCluster(redisName registryName) *redis.ClusterClient {
	if r, ok := registry.Load(redisName); ok {
		if result, ok := r.(*redis.ClusterClient); ok {
			return result
		}
	}
	return nil
}

func GetRedisScriptCli(redisName registryName) *RedisScriptCli {
	if r, ok := registry.Load(redisName); ok {
		if result, ok := r.(*RedisScriptCli); ok {
			return result
		}
	}
	return nil
}

func GetNatsConn(natsName registryName) *nats.Conn {
	if n, ok := registry.Load(natsName); ok {
		if result, ok := n.(*nats.Conn); ok {
			return result
		}
	}
	return nil
}

func GetESConn(esName registryName) *elastic.Client {
	if n, ok := registry.Load(esName); ok {
		if result, ok := n.(*elastic.Client); ok {
			return result
		}
	}
	return nil
}

func RegisterDatabase(dbName registryName, db *xorm.EngineGroup) {
	registry.Store(dbName, db)
}

func RegisterRedis(redisName registryName, r *redis.Client) {
	registry.Store(redisName, r)
}

func RegisterRedisCluster(redisName registryName, r *redis.ClusterClient) {
	registry.Store(redisName, r)
}

func RegisterRedisScriptCli(redisName registryName, r *RedisScriptCli) {
	registry.Store(redisName, r)
}

func RegisterNats(natsName registryName, n *nats.Conn) {
	registry.Store(natsName, n)
}

func RegisterES(esName registryName, cli *elastic.Client) {
	registry.Store(esName, cli)
}
