package FastTranslate

import (
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

// go test -v -run TestTrans
func TestTrans(t *testing.T) {
	// TransByServer("So, thanks for watching")
	TranslateSrt("/Users/zen/gitea/yt-whisper-bilingual/精修英文字幕/Busty Cougar Billie Jean Austin Naked Yoga_origin.srt", "")
}


