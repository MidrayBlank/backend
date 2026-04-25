package domain

type Geo struct {
	Code       int
	ParentCode *int
	Name       string
	Level      int
}
