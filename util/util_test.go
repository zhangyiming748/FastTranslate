package util

import (
    "testing"
)
// go test -v -timeout 2m -run TestTransWordsWithAI
func TestTransWordsWithAI(t *testing.T) {
	dst := TransWithLLM("You have the best cock to suck.")
	t.Logf("翻译后:%v\n", dst)
}