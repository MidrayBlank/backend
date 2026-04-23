package model

// AIAPI - счётчик запросов по api key, хранящемуся в виде хеша (50 в день, 20 в минуту)
type AIAPI struct {
	Hash     string `gorm:"column:hash;primaryKey;type:text"`
	Requests int    `gorm:"column:requests;type:int;default:0"`
}
