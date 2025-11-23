package FastTranslate

import (
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}
func TestTransOne(t *testing.T) {
	tc := TranslateConfig{}
	tc.SetKey("")
	tc.SourceSrtFile = "en_processed.srt"
	TransVideo(tc)
}
