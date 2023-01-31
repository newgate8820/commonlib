package dbcommon

import (
	"sync"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap/zapcore"
)

/*************************
* author: Dev0026
* createTime: 19-4-2
* updateTime: 19-4-2
* description:
*************************/
type RedisScriptCli struct {
	Cli       *redis.ClusterClient
	sha1Map   sync.Map
	scriptDir string
}

type LinkInfo struct {
	UserName   string // 用户名
	UserPasswd string // 密码
	HostInfo   string // 链接信息，格式ip:port
	Platform   string // 所属业务
}

type NatsLinkInfo struct {
	UserName string
	Passwd   string
	HostInfo string
	Platform string
}

type ZapLoggerConfig struct {
	Level zapcore.Level `yaml:"Level"`
	Dir   string        `yaml:"Dir"`
}

type EtcdLinkInfo struct {
	Platform string `yaml:"Platform"`
	HostInfo string `yaml:"HostInfo"`
}

type XormConfig struct {
	SchemaName          string `yaml:"SchemaName"`
	Platform            string `yaml:"Platform"`
	DriverName          string `yaml:"DriverName"`          // 数据库类型 mysql tidb
	DataSourceName      string `yaml:"DataSourceName"`      // 数据库连接
	MaxIdleConns        int    `yaml:"MaxIdleConns"`        // 最大空闲连接数
	MaxOpenConns        int    `yaml:"MaxOpenConns"`        // 最大连接数
	ShowSql             bool   `yaml:"ShowSql"`             // 是否打印sql
	LogLevel            int    `yaml:"LogLevel"`            // xorm 日志级别  LOG_DEBUG: 0 LOG_INFO: 1 LOG_WARNING: 2  LOG_ERR: 3 LOG_OFF: 4 LOG_UNKNOWN: 5
	SqlMapOptionRootDir string `yaml:"SqlMapOptionRootDir"` // sql 文件根目录
	SqlTemplateDir      string `yaml:"SqlTemplateDir"`      // sql 文件目录
	ParamStr            string `yaml:"ParamStr"`            // 连接参数 exp: chatset=utf8mb4
}

type XormDefaultConfig struct {
	GroupPolicy       string     `yaml:"GroupPolicy"`       // 负载策略
	GroupPolicyWeight []int      `yaml:"GroupPolicyWeight"` // 权重负载策略 对应值
	XormConfig        XormConfig `yaml:"XormConfigs"`       // xorm具体配置
}

type EtcdConfig struct {
	// platform
	Platform string `yaml:"Platform"`
	// Endpoints is a list of URLs.
	Endpoints []string `yaml:"Endpoints"`

	// AutoSyncInterval is the interval to update endpoints with its latest members.
	// 0 disables auto-sync. By default auto-sync is disabled.
	AutoSyncInterval time.Duration `yaml:"AutoSyncInterval"`

	// DialTimeout is the timeout for failing to establish a connection.
	DialTimeout time.Duration `yaml:"DialTimeout"`

	// Username is a username for authentication.
	Username string `yaml:"Username"`

	// Password is a password for authentication.
	Password string `yaml:"Password"`
}

type MicroserviceConfig struct {
	ServerName string `yaml:"ServerName"`
	MaxMsgSize int    `yaml:"MaxMsgSize"`
}

type XormGroupConfig struct {
	GroupPolicy       string       `yaml:"GroupPolicy"`       // 负载策略
	GroupPolicyWeight []int        `yaml:"GroupPolicyWeight"` // 权重负载策略 对应值
	XormConfigs       []XormConfig `yaml:"XormConfigs"`       // xorm具体配置
}

type NatsConfig struct {
	Platform      string `yaml:"Platform"`
	Url           string `yaml:"Url"`
	User          string `yaml:"User"`
	Password      string `yaml:"Password"`
	MaxReconnects int    `yaml:"MaxReconnects"`
	ReconnectWait int    `yaml:"ReconnectWait"`
}

type RedisConfig struct {
	redis.ClusterOptions
	Platform           string        `yaml:"Platform"`
	Addrs              []string      `yaml:"Addrs"`
	MaxRedirects       int           `yaml:"MaxRedirects"`
	ReadOnly           bool          `yaml:"ReadOnly"`
	RouteByLatency     bool          `ymal:"RouteByLatency"`
	RouteRandomly      bool          `yaml:"RouteRandomly"`
	MaxRetries         int           `yaml:"MaxRetries"`
	IdleTimeout        int64         `yaml:"IdleTimeout"` // 过期时间 单位 分钟 默认5分钟
	Password           string        `yaml:"Password"`
	PoolSize           int           `yaml:"PoolSize"`
	LuaScriptDir       string        `yaml:"LuaScriptDir"`       // lua 脚本路径
	IdleCheckFrequency time.Duration `yaml:"IdleCheckFrequency"` // 1h1m1.003002s

}

type ESConfig struct {
	Url      string `yaml:"Url"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	Open     bool   `yaml:"Open"`
}
