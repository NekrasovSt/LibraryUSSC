package models

type User struct {
	Id         int        `json:"id" gorm:"primaryKey"`
	FirstName  string     `json:"firstName"`
	SecondName string     `json:"secondName"`
	BookItems  []BookItem `gorm:"foreignKey:UserId"`
}
