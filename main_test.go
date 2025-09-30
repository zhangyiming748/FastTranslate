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
	tc.SetMysqlHost("192.168.5.2")
	tc.SetMysqlPort("3306")
	tc.SetMysqlUser("root")
	tc.SetMysqlPassword("163453")
	tc.SrtRoot = "en_processed.srt"
	TransVideo(tc)
}
