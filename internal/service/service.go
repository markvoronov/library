package service

import (
	"library"
	"library/internal/storage"
)

type Author interface {
	Create(author library.Author) (int, error) // Возвращаем айди созданного автора
	GetAll() ([]library.Author, error)
	GetByID(authorId int) (library.Author, error)
	Delete(authorId int) error
	Update(authorId int, input library.UpdateAuthor) error
}

type Book interface {
	Create(book library.Book) (int, error) // Возвращаем айди созданного автора
	GetAll() ([]library.Book, error)
	GetByID(bookId int) (library.Book, error)
	Delete(bookId int) error
	Update(bookId int, input library.UpdateBook) error
	UpdateBookAndAuthor(bookId int, authorId int, input library.UpdateAuthorBook) error
}

type Service struct {
	Author
	Book
}

func NewService(repos *storage.Repository) *Service {
	return &Service{
		Author: NewAuthorService(repos.Author),
		Book:   NewBookService(repos.Book),
	}
}
