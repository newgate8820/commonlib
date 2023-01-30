package dbcommon

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

/*************************
* author: Dev0026
* createTime: 19-4-2
* updateTime: 19-4-2
* description:
*************************/

func NewNatsClusterByCfg(cfg NatsConfig) (*nats.Conn, error) {
	var err error
	//集群模式 用了用户名和密码校验 接下来用
	//var servers = "nats://localhost:4241, nats://localhost:4242, nats://localhost:4243"
	conn, err := nats.Connect(cfg.Url, nats.UserInfo(cfg.User, cfg.Password),
		nats.MaxReconnects(cfg.MaxReconnects), nats.ReconnectWait(time.Duration(cfg.ReconnectWait)*time.Second), //设置重连时间
		nats.ReconnectBufSize(100*1024*1024),
		nats.PingInterval(5*time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			fmt.Printf("nats Got disconnected!error:%v\n", err)
		}), nats.ReconnectHandler(func(nc *nats.Conn) {
			fmt.Printf("nats Got reconnected to %v!\n", nc.ConnectedUrl())
		}), nats.ClosedHandler(func(nc *nats.Conn) {
			fmt.Printf("nats Connection closed. Reason: %q\n", nc.LastError())
		}),
	)
	if err != nil {
		log.Println("Can't connect: ", err, cfg.Url)
		return nil, fmt.Errorf("can't connect:%w", err)
	}
	return conn, err
}
