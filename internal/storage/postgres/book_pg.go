package postgres

import (
	"database/sql"
	"fmt"
	"library"
	"strings"
)

type BookPg struct {
	db *sql.DB
}

func NewBookPg(db *sql.DB) *BookPg {
	return &BookPg{db: db}
}

func (r *BookPg) Create(book library.Book) (int, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createBookQuery := fmt.Sprintf("INSERT INTO %s (title, author, year_of_publication, isbn) VALUES ($1, $2, $3, $4) RETURNING id", booksTable)
	row := tx.QueryRow(createBookQuery, book.Title, book.AuthorID, book.Year, book.ISBN)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()

}

func (r *BookPg) GetAll() ([]library.Book, error) {
	var books []library.Book

	query := fmt.Sprintf("SELECT id, title, author, year_of_publication, isbn FROM %s", booksTable)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var book library.Book
		err := rows.Scan(&book.ID, &book.Title, &book.AuthorID, &book.Year, &book.ISBN)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *BookPg) GetById(bookId int) (library.Book, error) {
	var book library.Book

	query := fmt.Sprintf("SELECT id, title, author, isbn, year_of_publication FROM %s WHERE id = $1", booksTable)
	err := r.db.QueryRow(query, bookId).Scan(&book.ID, &book.Title, &book.AuthorID, &book.ISBN, &book.Year)
	if err != nil {
		if err == sql.ErrNoRows {
			return book, fmt.Errorf("book not found")
		}
		return book, err
	}

	return book, nil
}

func (r *BookPg) Delete(bookId int) error {
	query := fmt.Sprintf(`DELETE FROM %s t WHERE t.id = $1`, booksTable)
	_, err := r.db.Exec(query, bookId)
	return err
}

func (r *BookPg) Update(bookId int, input library.UpdateBook) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Year != nil {
		setValues = append(setValues, fmt.Sprintf("year_of_publication=$%d", argId))
		args = append(args, *input.Year)
		argId++
	}

	if input.ISBN != nil {
		setValues = append(setValues, fmt.Sprintf("isbn=$%d", argId))
		args = append(args, *input.ISBN)
		argId++
	}

	if input.AuthorID != nil {
		setValues = append(setValues, fmt.Sprintf("author=$%d", argId))
		args = append(args, *input.AuthorID)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", booksTable, setQuery, argId)
	args = append(args, bookId)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *BookPg) UpdateBookAndAuthor(bookId int, authorId int, input library.UpdateAuthorBook) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Обновление таблицы books
	bookSetValues := make([]string, 0)
	bookArgs := make([]interface{}, 0)
	bookArgId := 1

	updateBook := input.UpdateBook
	if updateBook.Title != nil {
		bookSetValues = append(bookSetValues, fmt.Sprintf("title=$%d", bookArgId))
		bookArgs = append(bookArgs, *updateBook.Title)
		bookArgId++
	}

	if updateBook.Year != nil {
		bookSetValues = append(bookSetValues, fmt.Sprintf("year_of_publication=$%d", bookArgId))
		bookArgs = append(bookArgs, *updateBook.Year)
		bookArgId++
	}

	if updateBook.ISBN != nil {
		bookSetValues = append(bookSetValues, fmt.Sprintf("isnb=$%d", bookArgId))
		bookArgs = append(bookArgs, *updateBook.ISBN)
		bookArgId++
	}

	if updateBook.AuthorID != nil {
		bookSetValues = append(bookSetValues, fmt.Sprintf("author=$%d", bookArgId))
		bookArgs = append(bookArgs, *input.AuthorID)
		bookArgId++
	}

	bookSetQuery := strings.Join(bookSetValues, ", ")

	bookQuery := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", booksTable, bookSetQuery, bookArgId)
	bookArgs = append(bookArgs, bookId)

	if _, err := tx.Exec(bookQuery, bookArgs...); err != nil {
		return err
	}

	// Обновление таблицы authors
	authorSetValues := make([]string, 0)
	authorArgs := make([]interface{}, 0)
	authorArgId := 1

	updateAuthor := input.UpdateAuthor
	if updateAuthor.FirstName != nil {
		authorSetValues = append(authorSetValues, fmt.Sprintf("first_name=$%d", authorArgId))
		authorArgs = append(authorArgs, *updateAuthor.FirstName)
		authorArgId++
	}

	if updateAuthor.LastName != nil {
		authorSetValues = append(authorSetValues, fmt.Sprintf("last_name=$%d", authorArgId))
		authorArgs = append(authorArgs, *updateAuthor.LastName)
		authorArgId++
	}

	if updateAuthor.Bio != nil {
		authorSetValues = append(authorSetValues, fmt.Sprintf("bio=$%d", authorArgId))
		authorArgs = append(authorArgs, *updateAuthor.Bio)
		authorArgId++
	}

	if updateAuthor.DateOfBirth != nil {
		authorSetValues = append(authorSetValues, fmt.Sprintf("date_of_birth=$%d", authorArgId))
		authorArgs = append(authorArgs, *updateAuthor.DateOfBirth)
		authorArgId++
	}

	if len(authorSetValues) > 0 {
		authorSetQuery := strings.Join(authorSetValues, ", ")
		authorQuery := fmt.Sprintf("UPDATE authors SET %s WHERE id=$%d", authorSetQuery, authorArgId)
		authorArgs = append(authorArgs, authorId)

		if _, err := tx.Exec(authorQuery, authorArgs...); err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}

func (r *BookPg) Find(book library.Book) (*int, error) {

	query := fmt.Sprintf("SELECT t.id FROM %s t WHERE t.isbn = $1", booksTable)

	row := r.db.QueryRow(query, book.ISBN)
	var id int // для хранения id книги, если найдем
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// Возвращаем пустой идентификатор, если запись не найдена
			return nil, nil
		}
		return &id, err // Если вдруг ошибка не связана с отсустствием результата
	}

	return &id, nil

}
