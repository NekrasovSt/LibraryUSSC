package models

type Publisher struct {
	Id    int    `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Books []Book `gorm:"foreignKey:PublisherId"`
}

func GetPublishers() []Publisher {
	return []Publisher{
		{
			Name:  "Питер",
			Id:    1,
			Email: "piter@piter.com",
		},
	}
}
