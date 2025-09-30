package main

import (
	"github.com/zhangyiming748/FastTranslate"
)

func main() {
	// 正确初始化TranslateConfig
	tc := FastTranslate.TranslateConfig{}
	tc.SetProxy("http://192.168.5.2:8889")
	tc.SetMysqHost("192.168.5.2")
	tc.SetMysqlPort("3306")
	tc.SetMysqlUser("root")
	tc.SetMysqlPassword("163453")
	tc.SrtRoot = "/data/en_processed.srt"

	FastTranslate.TransVideo(tc)

}
