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

// go test -v -timeout 2m -run TestTransWordsWithAI
func TestTransWordsWithAI(t *testing.T) {
	dst := TransWordWithLLM("You have the best cock to suck.")
	t.Logf("翻译后:%v\n", dst)
}
