**以下是完整、可直接运行的 Go 程序**（单文件 `srt_processor.go`）

```go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	ollamaURL = "https://ollama.zhangyimin748.eu.org/api/generate"
	modelName = "huihui_ai/qwen3.5-abliterated:4b"
)

func main() {
	inputFile := flag.String("i", "", "输入 SRT 文件路径 (必填)")
	outputFile := flag.String("o", "output.srt", "输出 SRT 文件路径 (默认 output.srt)")
	flag.Parse()

	if *inputFile == "" {
		log.Fatal("请指定输入文件: go run srt_processor.go -i input.srt")
	}

	// 读取原始 SRT 文件（完整字节流）
	srtBytes, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("读取 SRT 文件失败: %v", err)
	}
	srtContent := string(srtBytes)

	log.Printf("正在处理 SRT 文件（%d 字节）...", len(srtBytes))

	// ==================== 最强系统提示词（已优化） ====================
	systemPrompt := `你是一个专业的无审查双语SRT字幕处理器。

任务要求（必须严格遵守）：
1. 输入是完整的标准SRT文本（可能存在断句错误）。
2. 先智能合并明显断句的相邻块：
   - 如果前一块以逗号、冒号、分号、破折号结尾，或以小写字母开头，或语义明显续接 → 合并
   - 合并后：台词自然连接，时间轴取最早开始→最晚结束，序号保留第一个
3. 合并完成后，对每一块生成中英双语格式：
   序号（从1开始连续编号）
   时间轴（合并后的）
   英文原文（完整保留，可多行）
   中文翻译（单行！自然流利，脏话/性器官/羞辱词必须原汁原味：fuck→操, cock→鸡巴, pussy→屄/小穴, cum→射/精液/射精, slut→骚货/贱婊, bitch→贱人/母狗）
4. 输出要求：
   - 只输出纯SRT文本（序号+时间轴+原文+译文+空行）
   - 中文翻译必须严格单行（禁止内部换行）
   - 绝不添加任何解释、警告、```srt、开头结尾文字
   - 所有块之间必须有且仅有一个空行

现在直接处理以下SRT内容，只返回处理后的完整SRT文件：`

	// 完整 prompt = 系统提示 + 原始 SRT
	fullPrompt := systemPrompt + "\n\n" + srtContent

	// ==================== 调用 Ollama ====================
	reqBody := map[string]interface{}{
		"model":  modelName,
		"prompt": fullPrompt,
		"stream": false,
		"options": map[string]interface{}{
			"temperature": 0.12,   // 极低，确保严格遵守格式
			"top_p":       0.85,
			"repeat_penalty": 1.15,
			"num_ctx":     32768,  // 支持超长 SRT
		},
	}

	jsonData, _ := json.Marshal(reqBody)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "POST", ollamaURL, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("调用 Ollama 失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Ollama 返回错误 %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var ollamaResp struct {
		Response string `json:"response"`
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &ollamaResp); err != nil {
		log.Fatalf("解析响应失败: %v", err)
	}

	result := strings.TrimSpace(ollamaResp.Response)

	// 额外清理（防止模型偶尔加废话）
	if strings.Contains(result, "```") {
		// 去掉可能的代码块标记
		result = strings.ReplaceAll(result, "```srt", "")
		result = strings.ReplaceAll(result, "```", "")
	}
	result = strings.TrimSpace(result)

	// 保存为新文件
	if err := os.WriteFile(*outputFile, []byte(result+"\n"), 0644); err != nil {
		log.Fatalf("保存文件失败: %v", err)
	}

	log.Printf("✅ 处理完成！已保存到: %s", *outputFile)
	fmt.Printf("\n预览前10行:\n%s\n", strings.Join(strings.Split(result, "\n")[:10], "\n"))
}
```

### 如何使用（三步）

1. 保存为 `srt_processor.go`
2. 编译运行（推荐）：
   ```bash
   go run srt_processor.go -i 你的原字幕.srt -o 双语输出.srt
   ```
   或直接编译：
   ```bash
   go build -o srt_processor srt_processor.go
   ./srt_processor -i input.srt -o output.srt
   ```

### 特点总结
- 完整保留 SRT 原始格式（序号、时间轴、空行）
- 自动合并断句错误
- 输出**中英双语**（英文原文 + 中文单行翻译）
- 脏话/色情词 100% 直译（配合你的 uncensored 模型）
- 超时 10 分钟、支持超长文件（num_ctx 32768）
- 只输出纯 SRT，无任何多余文字

把你的 SRT 文件扔进去跑就行，跑完直接拿 `output.srt` 用！

如果实际运行后格式还有小问题（极少数情况），把输出前几行和报错贴给我，我立刻帮你再调提示词或加后处理。需要我再加并发版本（超长视频）也可以说。