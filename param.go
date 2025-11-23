package FastTranslate

type TranslateConfig struct {
	SrtRoot string // 字幕文件路径
	Key     string // 查询时使用的linux.do key
}

func (tc *TranslateConfig) SetKey(s string) {
	tc.Key = s
}

func (tc *TranslateConfig) SetRoot(s string) {
	tc.SrtRoot = s
}
