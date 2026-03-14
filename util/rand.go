package util

import "math/rand"

// GetSeed 返回全局随机数生成器
// 注意：此函数每次调用返回同一个实例，避免重复创建
var globalRand = rand.New(rand.NewSource(42)) // 使用固定种子确保可重现性

func GetSeed() *rand.Rand {
	return globalRand
}
