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
	TranslateSrt("/Users/zen/Github/FastTranslate/source/Abella Danger.srt", "http://192.168.5.2:6380")
}
