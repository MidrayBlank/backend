package model

type RosstatAge struct {
	RosstatID    int  `gorm:"column:rosstat_id;type:int;primaryKey"`
	Age          int  `gorm:"column:age;type:int;primaryKey"`
	MaleAmount   *int `gorm:"column:male_amount;type:int;nullable"`
	FemaleAmount *int `gorm:"column:female_amount;type:int;nullable"`
}

func (RosstatAge) TableName() string {
	return "rosstat_age"
}
