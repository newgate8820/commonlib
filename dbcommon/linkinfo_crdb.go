package dbcommon

import (
	"fmt"
	"github.com/xormplus/xorm"
)

/*************************
* author: Dev0026
* createTime: 19-5-9
* updateTime: 19-5-9
* description:
*************************/

type cockdbInfo struct {
	UserName   string
	UserPasswd string
	HostInfo   string
	Platform   string
}

func initCRDBInfo(engine *xorm.Engine) error {
	var crdbInfo = make([]*cockdbInfo, 0)
	err := engine.SQL("select user_name, user_passwd, host_info, platform from cockdb_info").Find(&crdbInfo)
	if nil != err {
		return fmt.Errorf("initCRDBInfo err: %v", err)
	}
	for _, one := range crdbInfo {
		err = decodeCRDBInfo(one)
		if nil != err {
			return err
		}
		key := getCRDBKey(one.Platform)
		linkMap.Store(key, one)
	}
	return nil
}

func GetCRDBInfo(platform string) (*cockdbInfo, bool) {
	key := getCRDBKey(platform)
	val, ok := linkMap.Load(key)
	if ok {
		return val.(*cockdbInfo), ok
	}
	return nil, ok
}

func getCRDBKey(platform string) string {
	return fmt.Sprintf("crdbinfo_%v", platform)
}

func decodeCRDBInfo(info *cockdbInfo) (err error) {
	info.UserName, err = RsaDecrypt(info.UserName)
	if nil != err {
		return err
	}
	if info.UserPasswd != "" {
		info.UserPasswd, err = RsaDecrypt(info.UserPasswd)
		if nil != err {
			return err
		}
	}
	info.HostInfo, err = RsaDecrypt(info.HostInfo)
	if nil != err {
		return err
	}
	return nil
}
