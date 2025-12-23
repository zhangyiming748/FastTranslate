package FastTranslate

type TranslateConfig struct {
	SourceSrtFile string // 字幕文件路径
	Keyword           string // 查询时使用的暗号
}

func (tc *TranslateConfig) SetKey(s string) {
	tc.Keyword = s
}

func (tc *TranslateConfig) SetRoot(s string) {
	tc.SourceSrtFile = s
}

