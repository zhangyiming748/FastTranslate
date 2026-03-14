package FastTranslate

import (
	"log"
	"math/rand"
	"runtime"

	"os"
	"strconv"
	"strings"
	"time"

	"github.com/zhangyiming748/FastTranslate/util"
)

var (
	seed = rand.New(rand.NewSource(time.Now().Unix()))
)

/*
sourceSrtFile: 源文件
proxy: 代理地址
*/
func TranslateSrt(sourceSrtFile, proxy string) {
	r := seed.Intn(2000)
	tmpname := strings.Join([]string{strings.Replace(sourceSrtFile, ".srt", "", 1), strconv.Itoa(r), ".srt"}, "")
	log.Printf("开始处理字幕文件：%s, 临时文件名：%s\n", sourceSrtFile, tmpname)
	before := util.ReadInSlice(sourceSrtFile)
	log.Printf("读取源文件完成，共 %d 行\n", len(before))
	//fmt.Println(before)
	after, err := os.OpenFile(tmpname, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		log.Fatalf("无法创建临时字幕文件 [路径:%s]:%v\n", tmpname, err)
	}
	log.Printf("临时文件创建成功：%s\n", tmpname)
	for i := 0; i < len(before); i += 4 {
		if i+3 > len(before) {
			continue
		}
		log.Printf("正在处理第 %d 组字幕 (总共 %d 行)\n", i/4+1, len(before))
		after.WriteString(before[i])
		after.WriteString(before[i+1])
		src := before[i+2]
		src = strings.TrimSpace(src)
		var dst string
		dst = Trans(src, proxy)
		dst = strings.TrimSpace(dst)
		randomNumber := util.GetSeed().Intn(401) + 100
		time.Sleep(time.Duration(randomNumber) * time.Millisecond) // 暂停 100 毫秒
		after.WriteString(src)
		after.WriteString("\n")
		after.WriteString(dst)
		after.WriteString(before[i+3])
		after.Sync()
		log.Printf("第 %d 组字幕写入完成\n", i/4+1)
	}
	if err := after.Close(); err != nil {
		log.Fatalf("关闭临时字幕文件失败 [路径:%s]:%v\n", tmpname, err)
	}
	log.Printf("所有字幕处理完成，关闭文件并开始重命名\n")
	origin := strings.Join([]string{strings.Replace(sourceSrtFile, ".srt", "", 1), "_origin", ".srt"}, "")
	err1 := os.Rename(sourceSrtFile, origin)
	err2 := os.Rename(tmpname, sourceSrtFile)
	if err1 != nil || err2 != nil {
		log.Fatalf("字幕文件重命名出现错误:%v%v\n", err1, err2)
	}
}
func Trans(src, proxy string) (dst string) {
	
	switch runtime.GOOS {
	case "linux":
		dst = util.TransWithTranslateShell(src, proxy)
	default:
		dst =util.TransWithLLM(src)
	}
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
