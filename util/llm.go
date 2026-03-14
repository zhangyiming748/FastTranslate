package util

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
const ollamaURL = "http://192.168.5.4:11434/api/generate"

//const ollamaURL = "https://ollama.zhangyimin748.eu.org/api/generate"
const model = "huihui_ai/qwen3.5-abliterated:4b"

func TransWithLLM(src string) (dst string) {
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

