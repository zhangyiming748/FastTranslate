# FastTranslate

FastTranslate 是一个轻量级的字幕翻译工具，通过调用 DeepLX 提供的免费翻译 API 实现对 `.srt` 字幕文件的高效自动翻译。该项目提供了一个简单的 Go 库，可以集成到其他项目中使用。

## 功能特性

- 支持 `.srt` 格式字幕文件的读取与解析
- 调用 DeepLX API 进行多语言翻译
- 使用 SQLite 数据库存储翻译缓存，避免重复请求
- 自动备份原始字幕文件为 `*_origin.srt`
- 将翻译结果写回原文件路径

## 安装

```bash
go get github.com/zhangyiming748/FastTranslate
```

## 使用方法

### 基本用法

```go
package main

import "github.com/zhangyiming748/FastTranslate"

func main() {
    tc := FastTranslate.TranslateConfig{
        Key:           "your-deeplx-api-key",      // API密钥
        SourceSrtFile: "example.srt",              // 源字幕文件路径
        Proxy:         "http://127.0.0.1:8889",   // 可选：代理地址
    }
    
    FastTranslate.TranslateSrt(tc)
}
```

### 参数说明

[TranslateConfig](file:///Users/zen/Github/FastTranslate/param.go#L2-L6) 结构体包含以下字段：

- `SourceSrtFile`: 源字幕文件路径（必填）
- `Key`: DeepLX API 密钥（可选）
- `Proxy`: 代理地址（可选）

## 主要函数

### TranslateSrt(tc TranslateConfig)

这是主要的翻译函数，用于翻译整个 SRT 字幕文件。

- 读取源字幕文件内容
- 遍历字幕片段并逐条翻译
- 将翻译结果写入临时文件
- 备份原文件并替换为翻译后的文件
- 包含错误处理机制，防止程序崩溃

### Trans(src string, tc TranslateConfig) string

中间处理函数，用于预处理翻译结果。

- 调用 [TransByServer](file:///Users/zen/Github/FastTranslate/trans.go#L77-L98) 函数获取翻译结果
- 清理返回的翻译文本（删除换行符和回车符）
- 如果翻译结果包含 "error"，则返回原文本
- 返回处理后的翻译文本

### TransByServer(src string, tc TranslateConfig) string

向翻译服务器发送请求的函数。

- 构建 HTTP POST 请求到翻译服务器
- 发送源文本进行翻译
- 如果设置了代理，会在请求中包含代理信息
- 处理服务器响应并解析翻译结果
- 包含错误重试机制（如果请求失败或解析失败，会等待3秒后重试）

## 注意事项

- 该库会自动在工作目录下生成 `translation_cache.db` 作为 SQLite 缓存数据库
- 翻译过程中会创建临时文件，翻译完成后会自动清理
- 原始字幕文件会被重命名为 `original_filename_origin.srt`
- 翻译过程中会在每个翻译请求之间添加随机延迟，避免请求过于频繁

## 技术依赖

- Go 1.19+
- gorm.io/gorm
- github.com/glebarez/sqlite

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](file:///Users/zen/Github/FastTranslate/LICENSE) 文件。