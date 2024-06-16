package postgres

import (
	"database/sql"
	"fmt"
	"library"
	"strings"
)

type AuthorPg struct {
	db *sql.DB
}

func NewAuthorPg(db *sql.DB) *AuthorPg {
	return &AuthorPg{db: db}
}

func (r *AuthorPg) Create(author library.Author) (int, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int

	createAuthorQuery := fmt.Sprintf("INSERT INTO %s (first_name, last_name, bio, date_of_birth) VALUES ($1, $2, $3, $4) RETURNING id", authorsTable)

	row := tx.QueryRow(createAuthorQuery, author.FirstName, author.LastName, author.Bio, author.DateOfBirth)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()

}

func (r *AuthorPg) GetAll() ([]library.Author, error) {
	query := fmt.Sprintf("SELECT t.id, t.first_name, t.last_name, t.bio, t.date_of_birth FROM %s t", authorsTable)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []library.Author
	for rows.Next() {
		var author library.Author
		err := rows.Scan(&author.ID, &author.FirstName, &author.LastName, &author.Bio, &author.DateOfBirth)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

func (r *AuthorPg) GetById(authorId int) (library.Author, error) {
	var author library.Author

	query := fmt.Sprintf("SELECT t.id, t.first_name, t.last_name, t.bio, t.date_of_birth FROM %s t WHERE t.id = $1", authorsTable)

	row := r.db.QueryRow(query, authorId)
	err := row.Scan(&author.ID, &author.FirstName, &author.LastName, &author.Bio, &author.DateOfBirth)
	if err != nil {
		if err == sql.ErrNoRows {
			// Возвращаем пустого автора и ошибку, если запись не найдена
			return author, fmt.Errorf("author not found")
		}
		return author, err
	}

	return author, nil
}

func (r *AuthorPg) Delete(authorId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, authorsTable)
	_, err := r.db.Exec(query, authorId)
	return err
}

func (r *AuthorPg) Update(authorId int, input library.UpdateAuthor) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.FirstName != nil {
		setValues = append(setValues, fmt.Sprintf("first_name=$%d", argId))
		args = append(args, *input.FirstName)
		argId++
	}

	if input.LastName != nil {
		setValues = append(setValues, fmt.Sprintf("last_name=$%d", argId))
		args = append(args, *input.LastName)
		argId++
	}

	if input.Bio != nil {
		setValues = append(setValues, fmt.Sprintf("bio=$%d", argId))
		args = append(args, *input.Bio)
		argId++
	}
	if input.DateOfBirth != nil {
		setValues = append(setValues, fmt.Sprintf("date_of_birth=$%d", argId))
		args = append(args, *input.DateOfBirth)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", authorsTable, setQuery, argId)
	args = append(args, authorId)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *AuthorPg) Find(author library.Author) (*int, error) {

	query := fmt.Sprintf("SELECT t.id FROM %s t WHERE t.first_name = $1 and t.last_name = $2 and t.date_of_birth = $3", authorsTable)

	row := r.db.QueryRow(query, author.FirstName, author.LastName, author.DateOfBirth)
	var id int // для хранения id автора, если найдем
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
