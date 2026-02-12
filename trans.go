package FastTranslate

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"os/exec"
	"github.com/zhangyiming748/FastTranslate/util"
)

var (
	seed = rand.New(rand.NewSource(time.Now().Unix()))
	host = "http://127.0.0.1:6380/api/v1/translate"
)

/*
sourceSrtFile: 源文件
proxy: 代理地址
*/
func TranslateSrt(sourceSrtFile, proxy string) {
	r := seed.Intn(2000)
	tmpname := strings.Join([]string{strings.Replace(sourceSrtFile, ".srt", "", 1), strconv.Itoa(r), ".srt"}, "")
	before := util.ReadInSlice(sourceSrtFile)
	//fmt.Println(before)
	after, _ := os.OpenFile(tmpname, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	defer func() {
		if err := recover(); err != nil {
			v := fmt.Sprintf("捕获到错误:%v\n", err)
			if strings.Contains(v, "index out of range") {
				fmt.Println("捕获到 index out of range 类型错误,忽略并继续执行重命名操作")
				{
					origin := strings.Join([]string{strings.Replace(sourceSrtFile, ".srt", "", 1), "_origin", ".srt"}, "")
					err1 := os.Rename(sourceSrtFile, origin)
					err2 := os.Rename(tmpname, sourceSrtFile)
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
		dst = Trans(src,proxy)
		dst = strings.Replace(dst, "\n", "", -1)
		randomNumber := util.GetSeed().Intn(401) + 100
		time.Sleep(time.Duration(randomNumber) * time.Millisecond) // 暂停 100 毫秒
		fmt.Printf("trans.go的第61行输出src = %s\n", src)
		fmt.Printf("trans.go的第62行输出dst = %s\n", dst)
		after.WriteString(src)
		after.WriteString("\n")
		after.WriteString(dst)
		after.WriteString(before[i+3])
		after.WriteString(before[i+3])
		after.Sync()
	}
	after.Close()
	origin := strings.Join([]string{strings.Replace(sourceSrtFile, ".srt", "", 1), "_origin", ".srt"}, "")
	err1 := os.Rename(sourceSrtFile, origin)
	err2 := os.Rename(tmpname, sourceSrtFile)
	if err1 != nil || err2 != nil {
		log.Fatalf("字幕文件重命名出现错误:%v%v\n", err1, err2)
	}
}
func Trans(src,proxy string) (dst string) {
	dst = TransWithTranslateShell(src,proxy)
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
func TransWithTranslateShell(src ,proxy string) (dst string) {
	return TransByGoogle(src,proxy)
}
func TransByGoogle(src, proxy string) (dst string) {
	cmd := exec.Command("trans", "-brief", "-engine", "google", "-proxy", proxy, ":zh-CN", src)
	output, err := cmd.CombinedOutput()	
	result := string(output)
	result = strings.Replace(result, "\\r\\n", "", 1)
	result = strings.Replace(result, "\n", "", 1)
	result = strings.Replace(result, "\r\n", "", 1)
	if result == "" {
		result=TransByBing(src)
	}
	if err != nil || strings.Contains(string(output), "u001b") || strings.Contains(string(output), "Didyoumean") || strings.Contains(string(output), "Connectiontimedout") {
		log.Printf("google查询命令执行出错\t命令原文:%v\t错误原文:%v\n", cmd.String(), err.Error())
		time.Sleep(3 * time.Second)
		result=TransByBing(src)
	}
	return result
}
func TransByBing(src string) (dst string) {
	cmd := exec.Command("trans", "-brief", "-engine", "bing", ":zh-CN", src)
	log.Printf("查询命令:%s\n", cmd.String())
	output, err := cmd.CombinedOutput()
	result := string(output)
	result = strings.Replace(result, "\\r\\n", "", 1)
	result = strings.Replace(result, "\n", "", 1)
	result = strings.Replace(result, "\r\n", "", 1)
	if result == "" {
		return src
	}
	if err != nil || strings.Contains(string(output), "u001b") || strings.Contains(string(output), "Didyoumean") || strings.Contains(string(output), "Connectiontimedout") {
		log.Printf("bing查询命令执行出错\t命令原文:%v\t错误原文:%v\n", cmd.String(), err.Error())
		time.Sleep(3 * time.Second)
		TransByBing(src)
	}
	log.Printf("%s执行查询命令后的译文:%s\n", src, result)
	return result
}
