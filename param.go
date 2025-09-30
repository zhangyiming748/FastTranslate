package FastTranslate

type TranslateConfig struct {
	SrtRoot       string // 字幕文件路径
	Key           string // 查询时使用的linux.do key
	MysqlHost     string // mysql host
	MysqlPort     string // mysql port
	MysqlUser     string // mysql user
	MysqlPassword string // mysql password
}

func (tc *TranslateConfig) SetKey(s string) {
	tc.Key = s
}

func (tc *TranslateConfig) SetMysqlHost(s string) {
	tc.MysqlHost = s
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
