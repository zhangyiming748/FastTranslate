package FastTranslate

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/zhangyiming748/FastTranslate/util"
)

var (
	seed = rand.New(rand.NewSource(time.Now().Unix()))
)

func TranslateSrt(tc TranslateConfig) {
	r := seed.Intn(2000)
	tmpname := strings.Join([]string{strings.Replace(tc.SourceSrtFile, ".srt", "", 1), strconv.Itoa(r), ".srt"}, "")
	before := util.ReadInSlice(tc.SourceSrtFile)
	fmt.Println(before)
	after, _ := os.OpenFile(tmpname, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	defer func() {
		if err := recover(); err != nil {
			v := fmt.Sprintf("捕获到错误:%v\n", err)
			if strings.Contains(v, "index out of range") {
				fmt.Println("捕获到 index out of range 类型错误,忽略并继续执行重命名操作")
				{
					origin := strings.Join([]string{strings.Replace(tc.SourceSrtFile, ".srt", "", 1), "_origin", ".srt"}, "")
					err1 := os.Rename(tc.SourceSrtFile, origin)
					err2 := os.Rename(tmpname, tc.SourceSrtFile)
					if err1 != nil || err2 != nil {
						log.Fatalf("字幕文件重命名出现错误:%v%v\n", err1, err2)
					}
				}
				return
			} else {
				log.Fatalf("捕获到其他错误:%v\n", v)
			}
		}
	}()
	for i := 0; i < len(before); i += 4 {
		if i+3 > len(before) {
			continue
		}
		after.WriteString(before[i])
		after.WriteString(before[i+1])
		src := before[i+2]
		src = strings.Replace(src, "\n", "", 1)
		src = strings.Replace(src, "\r\n", "", 1)
		var dst string
		dst = Trans(src, tc)
		dst = strings.Replace(dst, "\n", "", -1)
		randomNumber := util.GetSeed().Intn(401) + 100
		time.Sleep(time.Duration(randomNumber) * time.Millisecond) // 暂停 100 毫秒
		fmt.Printf("src = %s\n", src)
		fmt.Printf("dst = %s\n", dst)
		after.WriteString(src)
		after.WriteString("\n")
		after.WriteString(dst)
		after.WriteString(before[i+3])
		after.WriteString(before[i+3])
		after.Sync()
	}
	after.Close()
	origin := strings.Join([]string{strings.Replace(tc.SourceSrtFile, ".srt", "", 1), "_origin", ".srt"}, "")
	err1 := os.Rename(tc.SourceSrtFile, origin)
	err2 := os.Rename(tmpname, tc.SourceSrtFile)
	if err1 != nil || err2 != nil {
		log.Fatalf("字幕文件重命名出现错误:%v%v\n", err1, err2)
	}
}
func Trans(src string, tc TranslateConfig) (dst string) {
	dst = TransByServer(src, tc)
	dst = strings.ReplaceAll(dst, "\n", "") // 删除所有换行符
	dst = strings.ReplaceAll(dst, "\r", "") // 删除所有回车符
	if strings.Contains(dst, "error") {
		return src
	}
	return dst
}

/*
curl --location --request POST 'http://trans.zhangyiming748.eu.org/api/v1/translate' \
--header 'Content-Type: application/json' \
--data-raw '{
"src":"hello",
"proxy":"http://127.0.0.1:8889"
}'
*/
func TransByServer(src string, tc TranslateConfig) (dst string) {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	params := map[string]string{
		"src": src,
	}
	if tc.Proxy != "" {
		params["proxy"] = tc.Proxy
	}
	b, err := util.HttpPostJson(headers, params, HOST)
	if err != nil {
		log.Printf("获取翻译服务响应失败,等待3秒后重试:%v\n", err)
		time.Sleep(3 * time.Second)
		TransByServer(src, tc)
	}
	fmt.Println(string(b))
	var r Req
	if e := json.Unmarshal(b, &r); e != nil {
		log.Printf("解析翻译内容失败,等待3秒后重试:%v\n", e)
		time.Sleep(3 * time.Second)
		TransByServer(src, tc)
	}
	fmt.Printf("请求服务返回的结构体是%+v\n", r)
	return r.Dst
}

const HOST = "http://trans.zhangyiming748.eu.org/api/v1/translate"

type Req struct {
	Dst string `json:"dst"`
}
