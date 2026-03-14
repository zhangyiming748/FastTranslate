package util
import (
	"os/exec"

	"log"
	"strings"
	"time"
)
func TransWithTranslateShell(src, proxy string) (dst string) {
	if src == "" {
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
