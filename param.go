package FastTranslate

type TranslateConfig struct {
	SourceSrtFile string // 字幕文件路径
}

func (tc *TranslateConfig) SetRoot(s string) {
	tc.SourceSrtFile = s
}
