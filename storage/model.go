package storage

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type TranslateHistory struct {
	Id        int64          `gorm:"primaryKey;autoIncrement;comment:主键id"`
	Src       string         `gorm:"comment:原文"`
	Dst       string         `gorm:"comment:译文"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (t *TranslateHistory) InsertOne() (int64, error) {
	result := GetSqlite().Create(t)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, result.Error
}

/*
根据src判断是否已经翻译过
*/
func (t *TranslateHistory) FindBySrc() (bool, error) {
	result := GetSqlite().Where("src = ?", t.Src).First(t)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}
