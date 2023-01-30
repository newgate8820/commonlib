package dbcommon

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
	"log"
	"strings"
	"sync"
)

/*************************
* author: Dev0026
* createTime: 19-4-2
* updateTime: 19-4-2
* description:
*************************/

var (
	mysqlLinks sync.Map
	natsLinks  sync.Map
	etcdLinks  sync.Map
)

func InitializeMysqlLinkMap(driver, src string) error {
	engine, err := xorm.NewEngine(driver, src)
	if err != nil {
		return err
	}
	defer engine.Close()

	infos := make([]*LinkInfo, 0)
	err = engine.SQL("select user_name, user_passwd, host_info, platform from mysql_info").Find(&infos)
	if err != nil {
		return err
	}

	for _, info := range infos {
		err = add(info)
		if err != nil {
			log.Println("DDDDDDDDDDDDdd  info ", info.Platform)
			return err
		}
	}

	return nil
}

func InitializeNatsLinkMap(driver, src string) error {
	engine, err := xorm.NewEngine(driver, src)
	if err != nil {
		return err
	}
	defer engine.Close()

	infos := make([]*NatsLinkInfo, 0)
	err = engine.SQL("select user_name, passwd, host_info, platform from nats_info").Find(&infos)
	if err != nil {
		return err
	}

	for _, info := range infos {
		err = addNats(info)
		if err != nil {
			return err
		}
	}

	return nil
}

func InitializeEtcdLinkMap(driver, source string) error {
	engine, err := xorm.NewEngine(driver, source)
	if err != nil {
		return err
	}

	infos := make([]*EtcdLinkInfo, 0)
	engine.SQL("SELECT host_info, platform FROM etcd_info;").Find(&infos)

	for _, info := range infos {
		info.HostInfo, err = RsaDecrypt(info.HostInfo)
		if err != nil {
			return err
		}
		etcdLinks.Store(info.Platform, info)
	}
	return nil
}

func GetEtcdPlatformInfo(platform string) ([]string, error) {
	info, ok := etcdLinks.Load(platform)
	if ok {
		return strings.Split(info.(*EtcdLinkInfo).HostInfo, ","), nil
	}

	return nil, errors.New("get etcd host info error.")
}

func GetMysqlPlatformInfo(platform string) (*LinkInfo, bool) {
	key := fmt.Sprintf("mysqlinfo_%s", platform)
	if val, ok := mysqlLinks.Load(key); ok {
		return val.(*LinkInfo), true
	} else {
		return nil, false
	}
}

func GetNatsPlatformInfo(platform string) (*NatsLinkInfo, bool) {
	key := fmt.Sprintf("natsinfo_%s", platform)
	if val, ok := natsLinks.Load(key); ok {
		return val.(*NatsLinkInfo), true
	} else {
		return nil, false
	}
}

func add(info *LinkInfo) (err error) {
	key := fmt.Sprintf("mysqlinfo_%s", info.Platform)
	info.HostInfo, err = RsaDecrypt(info.HostInfo)
	if err != nil {
		return err
	}

	info.UserPasswd, err = RsaDecrypt(info.UserPasswd)
	if err != nil {
		return err
	}

	info.UserName, err = RsaDecrypt(info.UserName)
	if err != nil {
		return err
	}

	mysqlLinks.Store(key, info)
	return nil
}

func addNats(info *NatsLinkInfo) (err error) {
	key := fmt.Sprintf("natsinfo_%s", info.Platform)
	info.HostInfo, err = RsaDecrypt(info.HostInfo)
	if err != nil {
		return err
	}

	info.Passwd, err = RsaDecrypt(info.Passwd)
	if err != nil {
		return err
	}

	info.UserName, err = RsaDecrypt(info.UserName)
	if err != nil {
		return err
	}

	natsLinks.Store(key, info)
	return nil
}

// 序列化
func MakeMarshal(model interface{}) (interface{}, error) {
	bin, err := json.Marshal(model)
	return bin, err
}

// 反序列化
func MakeUnmarshal(bin []byte, model interface{}) error {
	err := json.Unmarshal(bin, model)
	return err
}

func GetZrangeBy(min, max string, offset, limit int64) redis.ZRangeBy {
	return redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  limit,
	}
}

func RsaDecrypt(str string) (string, error) {
	ciphertext := []byte(str)
	block, _ := pem.Decode(privateKey) //将密钥解析成私钥实例
	if block == nil {
		return "", fmt.Errorf("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes) //解析pem.Decode（）返回的Block指针实例
	if err != nil {
		return "", err
	}
	rst, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext) //RSA算法解密
	if nil != err {
		return "", err
	}
	return string(rst), nil
}

var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
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
B2/X8LVUbWjJ
-----END RSA PRIVATE KEY-----
`)
