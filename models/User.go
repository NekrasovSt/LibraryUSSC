package models

type User struct {
	Id         int        `json:"id" gorm:"primaryKey"`
	FirstName  string     `json:"firstName"`
	SecondName string     `json:"secondName"`
	BookItems  []BookItem `gorm:"foreignKey:UserId"`
}

func GetDefaultUsers() []User {
	return []User{
		{FirstName: "Иван", SecondName: "Иванов"},
		{FirstName: "Петр", SecondName: "Петров"},
		{FirstName: "Василий", SecondName: "Васильев"},
	}
}
