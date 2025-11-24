package FastTranslate

type TranslateConfig struct {
	SourceSrtFile string // 字幕文件路径
	Key           string // 查询时使用的linux.do key
	Proxy         string // 如果使用代理，通过Google翻译 否则使用bing翻译
}

func (tc *TranslateConfig) SetKey(s string) {
	tc.Key = s
}

func (tc *TranslateConfig) SetRoot(s string) {
	tc.SourceSrtFile = s
}

func (tc *TranslateConfig) SetProxy(s string) {
	tc.Proxy = s
}
