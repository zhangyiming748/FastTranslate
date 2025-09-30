package FastTranslate

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
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

	storage.SetMysql(tc.MysqlUser, tc.MysqlPassword, tc.MysqlHost, tc.MysqlPort)

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
			dst = Trans(src, tc.Key)
			dst = strings.Replace(dst, "\n", "", -1)
			randomNumber := util.GetSeed().Intn(401) + 100
			time.Sleep(time.Duration(randomNumber) * time.Millisecond) // 暂停 100 毫秒
			behind.Dst = dst
			if _, err := behind.InsertOne(); err != nil {
				log.Fatalf("字幕文件写入缓存出现错误:%v\n", err)
			}
		}
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
	origin := strings.Join([]string{strings.Replace(tc.SrtRoot, ".srt", "", 1), "_origin", ".srt"}, "")
	err1 := os.Rename(tc.SrtRoot, origin)
	err2 := os.Rename(tmpname, tc.SrtRoot)
	if err1 != nil || err2 != nil {
		log.Fatalf("字幕文件重命名出现错误:%v%v\n", err1, err2)
	}
}
func Trans(src, key string) string {
	h := new(storage.TranslateHistory)
	h.Src = src
	if found, _ := h.FindBySrc(); found {
		return h.Dst
	}
RETRY:
	dst, err := TransByServer(src, key)
	if err != nil {
		time.Sleep(3 * time.Second)
		goto RETRY
	}

	dst = strings.ReplaceAll(dst, "\n", "") // 删除所有换行符
	dst = strings.ReplaceAll(dst, "\r", "") // 删除所有回车符
	if strings.Contains(dst, "error") {
		return src
	}
	h.Dst = dst
	h.InsertOne()
	return dst
}
func TransByServer(src, key string) (string, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	params := map[string]string{
		"text":        src,
		"source_lang": "auto",
		"target_lang": "zh",
	}
	host, _ := url.JoinPath(PREFIX, key, SUFFIX)
	b, err := util.HttpPostJson(headers, params, host)
	if err != nil {
		return "", err
	}
	fmt.Println(string(b))
	var d DeepLXTranslationResult
	if e := json.Unmarshal(b, &d); e != nil {
		return "", e
	}
	fmt.Printf("%+v\n", d)
	return d.Data, nil
}

const PREFIX = "https://api.deeplx.org"
const SUFFIX = "translate"

type DeepLXTranslationResult struct {
	Code         int      `json:"code"`
	ID           int64    `json:"id"`
	Message      string   `json:"message,omitempty"`
	Data         string   `json:"data"`         // The primary translated text
	Alternatives []string `json:"alternatives"` // Other possible translations
	SourceLang   string   `json:"source_lang"`
	TargetLang   string   `json:"target_lang"`
	Method       string   `json:"method"`
}
