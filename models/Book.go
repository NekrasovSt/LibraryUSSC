package models

type Book struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ISBN        string `json:"isbn" gorm:"primaryKey"`
	Edition     int    `json:"edition"`
	Year        int    `json:"year"`
	PublisherId int
	Books       []BookItem `gorm:"foreignKey:ISBN"`
	Authors     []Author   `gorm:"many2many:book_authors;"`
}

var defaultBooks = []Book{
	{ISBN: "978-5-4461-1155-8",
		Title:       "Современные операционные системы",
		Description: "Эндрю Таненбаум представляет новое издание своего всемирного бестселлера, необходимое для понимания функционирования современных операционных систем. Оно существенно отличается от предыдущего и включает в себя сведения о последних достижениях в области информационных технологий. Например, глава о Windows Vista теперь заменена подробным рассмотрением Windows 8.1 как самой актуальной версии на момент написания книги. Появился объемный раздел, посвященный операционной системе Android. Был обновлен материал, касающийся Unix и Linux, а также RAID-систем. Гораздо больше внимания уделено мультиядерным и многоядерным системам, важность которых в последние несколько лет постоянно возрастает. Появилась совершенно новая глава о виртуализации и облачных вычислениях. Добавился большой объем нового материала об использовании ошибок кода, о вредоносных программах и соответствующих мерах защиты. В книге в ясной и увлекательной форме приводится множество важных подробностей, которых нет ни в одном другом издании.",
		Edition:     4,
		Year:        2019,
	},
}

func GetDefaultBooks() []Book {
	return defaultBooks
}
