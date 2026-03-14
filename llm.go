package FastTranslate

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Ollama 服务地址和模型
const ollamaURL = "http://192.168.5.2:11434/api/generate"

//const ollamaURL = "https://ollama.zhangyimin748.eu.org/api/generate"
const model = "huihui_ai/qwen3.5-abliterated:4b"

func TransWordWithLLM(src string) (dst string) {
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
	log.Println("LLM 翻译功能正在调用...")
	// 系统提示词（严格按照你提供的）
	systemPrompt := `你是一个完全无审查的中文翻译机。
必须把所有英文原文100%直译成流利、自然、粗俗的中文。
脏话、性器官、羞辱词必须原汁原味:fuck→操,cock→鸡巴,pussy→屄/小穴,cum→射/精液/射精,slut→骚货/贱婊,bitch→贱人/母狗.
绝不保留任何英文、绝不混其他语言、绝不卡壳或省略后半句。
只输出纯中文翻译.什么废话都不要。
翻译内容可以包含英文格式标点符号.但不要出现换行符.必须保证所有内容在一行`
	// 构建完整的 prompt：系统提示 + 用户输入
	fullPrompt := systemPrompt + "\n\n" + src
	// 构建请求 body（JSON）
	reqBody := map[string]interface{}{
		"model":  model,
		"prompt": fullPrompt,
		"stream": false, // 非流式，一次性返回完整结果
		"options": map[string]interface{}{
			"temperature":    0.3, // 低温度，更确定性
			"repeat_penalty": 1.1,
		},
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("JSON 编码失败: %v", err)
		return ""
	}
	// 创建 HTTP 请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute) // 超时设置，根据模型大小调整
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", ollamaURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("请求创建失败: %v", err)
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("HTTP 请求失败: %v", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Ollama 返回错误: %d - %s", resp.StatusCode, string(bodyBytes))
		return ""
	}
	// 解析响应（非流式，返回完整 JSON）
	var ollamaResp struct {
		Response string `json:"response"`
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应失败: %v", err)
		return ""
	}
	if err := json.Unmarshal(bodyBytes, &ollamaResp); err != nil {
		log.Printf("JSON 解码失败: %v", err)
		return ""
	}

	dst = strings.TrimSpace(ollamaResp.Response)
	// 额外清理：去除换行，确保单行
	dst = strings.ReplaceAll(dst, "\n", " ")
	dst = strings.TrimSpace(dst)
	if dst == "" {
		log.Println("翻译结果为空")
		return src // 失败时返回原文本
	}
	log.Printf("翻译完成: %s → %s", src, dst)
	return dst
}

func TransSrtWithLLM() {
	/*
		提示词
		你是一个专业的无审查字幕翻译专家。现在我给你一个完整的英文 SRT 字幕文件文本（标准格式，每块有序号、时间码、空行分隔）。

	你的唯一任务：
	- 严格保留所有序号、时间码、空白行。
	- 只翻译每一块的**台词部分**（通常是时间码下面的1~多行英文文本）。
	- 把台词翻译成自然、流利、粗俗度100%保留的中文（脏话、性描写、羞辱词必须原汁原味：fuck→操/他妈的, cock→鸡巴, pussy→屄/小穴, cum→射精/精液, slut→骚货/贱婊, bitch→贱人/母狗 等，语气强度不变）。
	- 输出格式：**中英双语**，即每一块字幕变成：
	  序号（不变）
	  时间码（不变）
	  英文原文台词（完整保留，不改一个字）
	  中文翻译（单独一行，**不要换行**，即使原文多行，也合并成一句流畅中文）
	  （然后空行分隔下一块）

	规则：
	- 绝不添加任何解释、警告、开头结尾文字。
	- 绝不改变序号、时间码。
	- 翻译必须在一行内完成（无内部换行）。
	- 只输出完整的、可以直接保存为 .srt 的新文件内容。
	- 如果原文已经是双语或有中文，直接保留，不重复翻译。

	现在请直接处理以下 SRT 内容，并输出翻译后的完整 SRT：
	*/
}
