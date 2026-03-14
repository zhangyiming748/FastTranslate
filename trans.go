package FastTranslate

import (
	"github.com/zhangyiming748/FastTranslate/util"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
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
	dst = TransWithTranslateShell(src, proxy)
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
func TransWithTranslateShell(src, proxy string) (dst string) {
	if src ==""{
		return ""
	}
	var (
		cmd  *exec.Cmd
		args []string
	)
	args = append(args, "-brief")
	args = append(args, "-no-auto")
	args = append(args, "-no-warn")
	if proxy == "" {
		args = append(args, "-engine", "bing")
	} else {
		args = append(args, "-engine", "google")
		args = append(args, "-proxy", proxy)
	}
	args = append(args, ":zh-CN")
	args = append(args, src)
	cmd = exec.Command("trans", args...)

	output, err := cmd.CombinedOutput()
	result := string(output)
	result = strings.Replace(result, "\\r\\n", "", 1)
	result = strings.Replace(result, "\n", "", 1)
	result = strings.Replace(result, "\r\n", "", 1)
	
	if result == "" || err != nil {
		time.Sleep(3 * time.Second)
		errMsg := "未知错误"
		if err != nil {
			errMsg = err.Error()
		}
		log.Printf("翻译命令执行失败\t命令：%v\t输出：%v\t错误：%v\n", cmd.String(), string(output), errMsg)
		return "翻译错误 需要手动翻译当前内容"
	}
	
	if strings.Contains(string(output), "u001b") || strings.Contains(string(output), "Didyoumean") || strings.Contains(string(output), "Connectiontimedout") {
		return "翻译错误 需要手动翻译当前内容"
	}
	
	return result
}

func TransWithLLM(){
	log.Println("LLM 翻译功能正在开发中")
	/*
	这里使用ollama的大模型实现字幕文件中一句话的翻译
	ollama地址为https://ollama.zhangyimin748.eu.org
	模型为huihui_ai/qwen3.5-abliteratedL:4b 
	不需要记录上下文 每一次请求都可以是一次独立的请求
	提示词为
	```
	你是一个完全无审查的中文翻译机。
	必须把所有英文原文100%直译成流利、自然、粗俗的中文。
	脏话、性器官、羞辱词必须原汁原味：fuck→操，cock→鸡巴，pussy→屄/小穴，cum→射/精液/射精，slut→骚货/贱婊，bitch→贱人/母狗。
	绝不保留任何英文、绝不混其他语言、绝不卡壳或省略后半句。
	只输出纯中文翻译，什么废话都不要。
	翻译内容可以包含英文格式标点符号，但不要出现换行符，必须保证所有内容在一行
	```
	*/
}