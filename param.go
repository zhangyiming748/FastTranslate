package FastTranslate

type TranslateConfig struct {
	SrtRoot       string // 字幕文件路径
	Proxy         string // 查询时使用的代理
	MysqHost      string // mysql host
	MysqlPort     string // mysql port
	MysqlUser     string // mysql user
	MysqlPassword string // mysql password
}

func (tc *TranslateConfig) SetProxy(s string) {
	tc.Proxy = s
}

func (tc *TranslateConfig) SetMysqHost(s string) {
	tc.MysqHost = s
}

func (tc *TranslateConfig) SetMysqlPort(s string) {
	tc.MysqlPort = s
}

func (tc *TranslateConfig) SetMysqlUser(s string) {
	tc.MysqlUser = s
}

func (tc *TranslateConfig) SetMysqlPassword(s string) {
	tc.MysqlPassword = s
}
