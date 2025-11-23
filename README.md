# FastTranslate

FastTranslate 是一个基于 DeepLX API 的字幕文件自动翻译工具。它支持 `.srt` 格式的字幕文件，并具有翻译缓存功能，避免重复翻译相同的文本。

## 特性

- 使用 DeepLX API 进行翻译
- 自动缓存翻译结果到 SQLite 本地数据库
- 支持 `.srt` 格式字幕文件
- 避免重复翻译，提高效率

## 安装

```bash
go get github.com/zhangyiming748/FastTranslate
```

## 使用方法

```go
import (
    "github.com/zhangyiming748/FastTranslate"
)

func main() {
    tc := FastTranslate.TranslateConfig{}
    tc.Key = "your-deeplx-api-key"           // 设置你的 DeepLX API 密钥
    tc.SourceSrtFile = "path/to/subtitle.srt" // 设置源字幕文件路径
    
    FastTranslate.TransVideo(tc)
}
```

## 配置说明

- `Key`: DeepLX API 密钥
- `SourceSrtFile`: 需要翻译的 `.srt` 字幕文件路径

## 工作原理

1. 程序读取指定的 `.srt` 字幕文件
2. 对于每个需要翻译的文本片段：
   - 首先检查 SQLite 本地数据库中是否有缓存的翻译结果
   - 如果有缓存则直接使用缓存结果
   - 如果没有缓存则调用 DeepLX API 进行翻译，并将结果存入数据库
3. 将翻译后的字幕保存到原文件，原文件会被重命名为 `*_origin.srt`

## 依赖

- [DeepLX](https://github.com/OwO-Network/DeepLX) - 提供翻译API服务
- SQLite - 用于本地缓存翻译结果