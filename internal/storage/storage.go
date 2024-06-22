package storage

import (
	"database/sql"
	"library"
	"library/internal/storage/postgres"
)

type Author interface {
	Create(author library.Author) (int, error)
	GetAll() ([]library.Author, error)
	GetById(authorId int) (library.Author, error)
	Delete(authorId int) error
	Update(authorId int, input library.UpdateAuthor) error
	Find(author library.Author) (*int, error)
}

type Book interface {
	Create(book library.Book) (int, error)
	GetAll() ([]library.Book, error)
	GetById(bookId int) (library.Book, error)
	Delete(bookId int) error
	Update(bookId int, input library.UpdateBook) error
	UpdateBookAndAuthor(bookId int, authorId int, input library.UpdateAuthorBook) error
	Find(book library.Book) (*int, error)
	Db() *sql.DB
}

type Repository struct {
	Author Author
	Book   Book
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Author: postgres.NewAuthorPg(db),
		Book:   postgres.NewBookPg(db),
	}
}
