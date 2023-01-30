package dbcommon

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/go-redis/redis"
)

/*************************
* author: Dev0026
* createTime: 19-4-2
* updateTime: 19-4-2
* description:
*************************/

func NewRedisEngine(cfg RedisConfig) (*redis.ClusterClient, error) {
	if len(cfg.Addrs) == 0 && cfg.Platform != "" {
		HandleRedisConfigEncryFields(&cfg)
	}
	var opt = &redis.ClusterOptions{
		// A seed list of host:port addresses of cluster nodes.
		Addrs: cfg.Addrs,

		// The maximum number of retries before giving up. Command is retried
		// on network errors and MOVED/ASK redirects.
		// Default is 16.
		MaxRedirects: cfg.MaxRedirects,

		// Enables read-only commands on slave nodes.
		ReadOnly: cfg.ReadOnly,
		// Allows routing read-only commands to the closest master or slave node.
		RouteByLatency: cfg.RouteByLatency,

		RouteRandomly: cfg.RouteRandomly,

		// Following options are copied from Options struct.

		MaxRetries:  cfg.MaxRetries,
		Password:    cfg.Password,
		IdleTimeout: time.Duration(cfg.IdleTimeout) * time.Minute,
		// PoolSize applies per cluster node and not for the whole cluster.
		PoolSize:           cfg.PoolSize,
		IdleCheckFrequency: cfg.IdleCheckFrequency,
	}

	cluster := redis.NewClusterClient(opt)

	str, err := cluster.Ping().Result()

	if nil != err {
		er := fmt.Errorf("NewRedisEngine [ cfg %+v ] redis.NewClusterClient err %v ", cfg, err)
		return nil, er
	}
	fmt.Printf(" %+v redis cluster engine start success, ping %v \n", cfg.Addrs, str)

	redisEngines = append(redisEngines, cluster)

	return cluster, nil
}

func NewRedisScriptCli(cli *redis.ClusterClient) *RedisScriptCli {
	return &RedisScriptCli{
		Cli: cli,
	}
}

func (cc *RedisScriptCli) LoadScripts(dir string) error {
	cc.scriptDir = dir
	fInfo, err := os.Stat(dir)
	if nil != err {
		er := fmt.Errorf("*RedisScriptCli.LoadScript [ dir %v ] Stat err: %v", dir, err)
		return er
	} else if fInfo.IsDir() {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() && path == dir {
				return nil
			} else if info.IsDir() && path != dir {
				return cc.LoadScripts(path)
			} else {
				shaStr, err := cc.LoadScriptFile(path)
				if nil != err {
					return err
				} else {
					cc.sha1Map.Store(info.Name(), shaStr)
					return nil
				}
			}
		})
		return err
	} else {
		shaStr, err := cc.LoadScriptFile(dir)
		if nil != err {
			return err
		} else {
			cc.sha1Map.Store(dir, shaStr)
			return nil
		}
	}
}

func (cc *RedisScriptCli) LoadScriptFile(fileName string) (string, error) {
	bin, err := ioutil.ReadFile(fileName)
	if nil != err {
		er := fmt.Errorf("*RedisScriptCli.LoadScriptFile [ fileName %v ] ioutil.ReadFile err: %v", fileName, err)
		return "", er
	}

	scriptStr := string(bin)
	str, err := cc.Cli.ScriptLoad(scriptStr).Result()
	if nil != err {
		er := fmt.Errorf("*RedisScriptCli.LoadScript [ scriptStr %v ] ScriptLoad err: %v", scriptStr, err)
		return "", er
	} else {
		return str, nil
	}
}

func (cc *RedisScriptCli) EvalSha(luaFileName string, keys []string, args ...interface{}) (interface{}, error) {
	if val, ok := cc.sha1Map.Load(luaFileName); ok {
		sha1 := val.(string)
		//fmt.Println(sha1)
		rst, err := cc.Cli.EvalSha(sha1, keys, args...).Result()
		if nil != err && err.Error() == "NOSCRIPT No matching script. Please use EVAL." {
			return cc.Eval(cc.scriptDir, luaFileName, keys, args...)
		} else {
			if err != err {
				er := fmt.Errorf("EvalSha [ luaFileName %v keys [ %+v ] args [ %+v ] rst %v ] err: %v ", luaFileName, keys, args, rst, err)
				return rst, er
			}
			return rst, nil
		}
	} else {
		return cc.Eval(cc.scriptDir, luaFileName, keys, args...)
	}
}

func (cc *RedisScriptCli) Eval(dir, luaFileName string, keys []string, args ...interface{}) (interface{}, error) {
	var rst interface{}
	var fileExist bool
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path == cc.scriptDir {
			return nil
		} else if info.IsDir() {
			rst, err = cc.Eval(path, luaFileName, keys, args...)
			return err
		} else if luaFileName == info.Name() {
			fileExist = true
			bin, err := ioutil.ReadFile(path)
			if nil != err {
				er := fmt.Errorf("*RedisScriptCli.LoadScriptFile [ fileName %v ] ioutil.ReadFile err: %v", path, err)
				return er
			}

			scriptStr := string(bin)
			rst, err = cc.Cli.Eval(scriptStr, keys, args...).Result()
			return err
		} else {
			return nil
		}
	})
	if !fileExist {
		return rst, fmt.Errorf("*RedisScriptCli.Eval [ luaFileName %v ] file not found", luaFileName)
	}

	if err != err {
		er := fmt.Errorf("Eval [ luaFileName %v keys [ %+v ] args [ %+v ] rst %v ] err: %v ", luaFileName, keys, args, rst, err)
		return rst, er
	}

	return rst, err
}
