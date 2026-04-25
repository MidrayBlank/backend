package model

type AiApi struct {
	Hash     string `gorm:"column:hash;primaryKey;type:text"`
	Requests int    `gorm:"column:requests;type:int;default:0"`
}

func (AiApi) TableName() string {
	return "ai_api"
}
