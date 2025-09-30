package FastTranslate

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/zhangyiming748/FastTranslate/storage"
	"github.com/zhangyiming748/FastTranslate/util"
)

var (
	seed = rand.New(rand.NewSource(time.Now().Unix()))
)

func TransVideo(tc TranslateConfig) {
	r := seed.Intn(2000)
	tmpname := strings.Join([]string{strings.Replace(tc.SrtRoot, ".srt", "", 1), strconv.Itoa(r), ".srt"}, "")
	before := util.ReadInSlice(tc.SrtRoot)
	fmt.Println(before)
	after, _ := os.OpenFile(tmpname, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	defer func() {
		if err := recover(); err != nil {
			v := fmt.Sprintf("捕获到错误:%v\n", err)
			if strings.Contains(v, "index out of range") {
				fmt.Println("捕获到 index out of range 类型错误,忽略并继续执行重命名操作")
				{
					origin := strings.Join([]string{strings.Replace(tc.SrtRoot, ".srt", "", 1), "_origin", ".srt"}, "")
					err1 := os.Rename(tc.SrtRoot, origin)
					err2 := os.Rename(tmpname, tc.SrtRoot)
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
		log.Printf("翻译之前序号\"%s\"时间\"%s\"正文\"%s\"空行\"%s\"\n", before[i], before[i+1], before[i+2], before[i+3])
		log.SetPrefix(before[i])
		after.WriteString(before[i])
		after.WriteString(before[i+1])
		src := before[i+2]
		src = strings.Replace(src, "\n", "", 1)
		src = strings.Replace(src, "\r\n", "", 1)

		var dst string
		behind := new(storage.TranslateHistory)
		behind.Src = src
		if has, _ := behind.FindBySrc(); has {
			dst = behind.Dst
			fmt.Printf("在缓存中找到dst = %s\n", dst)
		} else {
			fmt.Println("未在缓存中找到")
			dst = Trans(src, tc.Proxy)
			dst = strings.Replace(dst, "\n", "", -1)
			randomNumber := util.GetSeed().Intn(401) + 100
			time.Sleep(time.Duration(randomNumber) * time.Millisecond) // 暂停 100 毫秒
			behind.Dst = dst
			if _, err := behind.InsertOne(); err != nil {
				log.Fatalf("字幕文件写入缓存出现错误:%v\n", err)
			}
		}

		fmt.Printf("翻译之后序号:\"%s\"时间:\"%s\"正文:\"%s\"空行:\"%s\"原文:\"%s\"\t译文\"%s\"\n", before[i], before[i+1], before[i+2], before[i+3], src, dst)
		after.WriteString(src)
		after.WriteString("\n")
		after.WriteString(dst)
		after.WriteString(before[i+3])
		after.WriteString(before[i+3])
		after.Sync()
	}
	after.Close()
	origin := strings.Join([]string{strings.Replace(tc.SrtRoot, ".srt", "", 1), "_origin", ".srt"}, "")
	err1 := os.Rename(tc.SrtRoot, origin)
	err2 := os.Rename(tmpname, tc.SrtRoot)
	if err1 != nil || err2 != nil {
		log.Fatalf("字幕文件重命名出现错误:%v%v\n", err1, err2)
	}
}
func Trans(src, proxy string) string {
	h := new(storage.TranslateHistory)
	h.Src = src
	if found, _ := h.FindBySrc(); found {
		return h.Dst
	}
	var cmd *exec.Cmd
	if proxy == "" {
		cmd = exec.Command("trans", "-brief", "-engine", "bing", ":zh-CN", src)
	} else {
		cmd = exec.Command("trans", "-brief", "-engine", "google", "-proxy", proxy, ":zh-CN", src)
	}
	log.Printf("命令 : %s\n", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return src
	}

	dst := string(output)
	dst = strings.ReplaceAll(dst, "\n", "") // 删除所有换行符
	dst = strings.ReplaceAll(dst, "\r", "") // 删除所有回车符
	if strings.Contains(dst, "error") {
		return src
	}
	h.Dst = dst
	h.InsertOne()
	return dst
}
