package model

type Geo struct {
	Code       int    `gorm:"column:code;primaryKey;type:int"`
	ParentCode int    `gorm:"column:parent_code;type:int;nullable"`
	Name       string `gorm:"column:name;type:text;not null"`
	AdminType  int    `gorm:"column:admin_type;type:int;not null"`
}
