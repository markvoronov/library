package library

type Author struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Bio         string `json:"bio"`
	DateOfBirth string `json:"date_of_birth" db:"date_of_birth" binding:"required"` // Будем работать с датой как со строкой,
	// чтобы не заморачиваться со временем. Был вариант еще сделать доп поле уже типа time.Time, где хранить приведенную
	// к этому типу дату. Но для базы данных этого не нужно, а идей где еще она понадобится кроме входной проверки я не нашел
}

type UpdateAuthor struct {
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	Bio         *string `json:"bio"`
	DateOfBirth *string `json:"date_of_birth"`
}

type Book struct {
	ID       string `json:"id"`
	Title    string `json:"title" db:"title" binding:"required"`
	AuthorID int    `json:"author_id" db:"author" binding:"required"`
	Year     int    `json:"year_of_publication" db:"year_of_publication" binding:"required"`
	ISBN     string `json:"isbn" binding:"required"`
}

type UpdateBook struct {
	Title    *string `json:"title"`
	AuthorID *int    `json:"author_id"`
	Year     *int    `json:"year_of_publication"`
	ISBN     *string `json:"isbn"`
}

type UpdateAuthorBook struct {
	UpdateAuthor `json:"author"`
	UpdateBook   `json:"book"`
}
