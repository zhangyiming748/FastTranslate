# FastTranslate
```go
import(
"github.com/zhangyiming748/FastTranslate"
)
type TranslateConfig struct {
	SrtRoot       string // 字幕文件路径
	Proxy         string // 查询时使用的代理
	MysqHost      string // mysql host
	MysqlPort     string // mysql port
	MysqlUser     string // mysql user
	MysqlPassword string // mysql password
}
FastTranslate.TransVideo(tc)
```
