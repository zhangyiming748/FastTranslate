package util

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func ReadInSlice(fp string) []string {
	fileBytes, err := os.ReadFile(fp)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return []string{}
	}
	// 创建一个 bufio.Reader 对象
	reader := bufio.NewReader(bytes.NewReader(fileBytes))
	// 按行读取文件内容并存储到字符串切片中
	var lines []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		lines = append(lines, line)
	}
	return lines
}
