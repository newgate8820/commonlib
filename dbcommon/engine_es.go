package dbcommon

import (
	"context"
	"fmt"

	"github.com/olivere/elastic"
)

/*************************
* author: Dev0026
* createTime: 19-4-2
* updateTime: 19-4-2
* description:
*************************/

func NewESClientByCfg(cfg ESConfig) (cli *elastic.Client, err error) {
	cli, err = elastic.NewClient(elastic.SetURL(cfg.Url), elastic.SetBasicAuth(cfg.User, cfg.Password))
	if nil != err {
		return nil, fmt.Errorf("registerDatabase [ cfg %+v ] RedisESConfig err: %v",
			cfg, err)
	}
	_, _, err = cli.Ping(cfg.Url).Do(context.Background())
	if nil != err {
		return nil, fmt.Errorf("registerDatabase [ cfg %+v ] ping err: %v",
			cfg, err)
	}
	return
}
