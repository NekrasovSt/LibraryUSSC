package models

import "time"

type BookItem struct {
	Id      int       `json:"id" gorm:"primaryKey"`
	ISBN    string    `json:"isbn"`
	Receipt time.Time `json:"receipt"`
	UserId  *int
}

func GetBookItems() []BookItem {
	return []BookItem{
		BookItem{Id: 0, ISBN: "978-5-4461-1155-8", Receipt: time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC)},
		BookItem{Id: 0, ISBN: "978-5-4461-1155-8", Receipt: time.Date(2021, 1, 1, 9, 2, 0, 0, time.UTC)},
		BookItem{Id: 0, ISBN: "978-5-4461-1155-8", Receipt: time.Date(2021, 1, 1, 9, 3, 0, 0, time.UTC)},
	}
}
