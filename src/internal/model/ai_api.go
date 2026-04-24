package model

type AIAPI struct {
	Hash     string `gorm:"column:hash;primaryKey;type:text"`
	Requests int    `gorm:"column:requests;type:int;default:0"`
}

func (AIAPI) TableName() string {
	return "midray.ai_api"
}
