package models

type Author struct {
	Id         int    `json:"id" gorm:"primaryKey"`
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
	Books      []Book `gorm:"many2many:book_authors;"`
}

func GetDefaultAuthors() []Author {
	return []Author{
		{FirstName: "Эндрю", SecondName: "Таненбаум"},
		{FirstName: "Херберт", SecondName: "Бос"},
		{FirstName: "Тодд", SecondName: "Остин"},
		{FirstName: "Адитья", SecondName: "Бхаргава"},
	}
}
