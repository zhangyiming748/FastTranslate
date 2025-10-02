package storage

import (
	"time"
)

type TranslateHistory struct {
	Id        int64     `xorm:"pk autoincr notnull comment('主键id') INT(11)"`
	Src       string    `xorm:"Text comment(原文)"`
	Dst       string    `xorm:"Text comment(译文)"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

func (t *TranslateHistory) InsertOne() (int64, error) {
	return GetMysql().InsertOne(t)
}

func (t *TranslateHistory) FindBySrc() (bool, error) {
	return GetMysql().Where("src = ?", t.Src).Get(t)
}

func (t *TranslateHistory) InsertAll(histories []TranslateHistory) (int64, error) {
	return GetMysql().Insert(histories)
}
