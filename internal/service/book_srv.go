package service

import (
	"fmt"
	"library"
	"library/internal/storage"
)

type BookService struct {
	repo storage.Book
}

func NewBookService(repo storage.Book) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) Create(book library.Book) (int, error) {
	// Сначала нужно проверить, нет ли уже такой книги. Идентифицировать будем по isbn
	id, err := s.repo.Find(book)
	if err != nil {
		err = fmt.Errorf("Не удалось проверить существование книги в базе", err)
		return -1, err
	}
	if id != nil {
		err = fmt.Errorf("Книга уже существует под id %d, создание прервано", *id)
		return -1, err
	}

	return s.repo.Create(book)
}

func (s *BookService) GetAll() ([]library.Book, error) {
	return s.repo.GetAll()
}

func (s *BookService) GetByID(bookId int) (library.Book, error) {
	return s.repo.GetById(bookId)
}

func (s *BookService) Delete(bookId int) error {

	_, err := s.repo.GetById(bookId)
	if err != nil {
		err = fmt.Errorf("Нет возможности удалить книгу ", err)
		return err
	}
	return s.repo.Delete(bookId)
}

func (s *BookService) Update(bookId int, input library.UpdateBook) error {
	_, err := s.repo.GetById(bookId)
	if err != nil {
		err = fmt.Errorf("Нет возможности обновить данные о книге", err)
		return err
	}
	return s.repo.Update(bookId, input)
}

func (s *BookService) UpdateBookAndAuthor(bookId int, authorId int, input library.UpdateAuthorBook) error {

	_, err := s.repo.GetById(bookId)
	if err != nil {
		err = fmt.Errorf("Нет возможности обновить данные о книге, операция прервана ", err)
		return err
	}

	return s.repo.UpdateBookAndAuthor(bookId, authorId, input)
}
