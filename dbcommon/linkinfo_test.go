package dbcommon

import "testing"

/*************************
* author: Dev0026
* createTime: 19-4-2
* updateTime: 19-4-2
* description:
*************************/

func TestGetMysqlInfoFromLinkMap(t *testing.T) {
	err := InitLinkMapByDataSourceName("root:bsrt_123,./&^%@tcp(192.168.212.205:3306)/link_information?charset=utf8mb4")
	if nil != err {
		t.Fatal(err)
	}
	info := GetMysqlInfoFromLinkMap("pt_rmsg")

	t.Log(info)
}
