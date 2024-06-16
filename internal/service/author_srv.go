package service

import (
	"fmt"
	"library"
	"library/internal/storage"
	"time"
)

type AuthorService struct {
	repo storage.Author
}

func NewAuthorService(repo storage.Author) *AuthorService {
	return &AuthorService{repo: repo}
}

func (s *AuthorService) Create(Author library.Author) (int, error) {

	//Проверм корректность переданной даты
	if err := checkDateOfBirthday(Author.DateOfBirth); err != nil {
		return -1, err
	}

	// Сначала нужно проверить, нет ли уже такого автора. Идентифицировать будем по 3 полям: имя, фамилия, дата рождения
	id, err := s.repo.Find(Author)
	if err != nil {
		err = fmt.Errorf("Не удалось проверить существование автора в базе", err)
		return -1, err
	}
	if id != nil {
		err = fmt.Errorf("Автор уже существует под id %d, создание прервано", *id)
		return -1, err
	}
	return s.repo.Create(Author)
}

func (s *AuthorService) GetAll() ([]library.Author, error) {
	return s.repo.GetAll()
}

func (s *AuthorService) GetByID(authorId int) (library.Author, error) {
	return s.repo.GetById(authorId)
}

func (s *AuthorService) Delete(authorId int) error {
	_, err := s.repo.GetById(authorId)
	if err != nil {
		err = fmt.Errorf("Нет возможности удалить данные о авторе ", err)
		return err
	}
	return s.repo.Delete(authorId)
}

func (s *AuthorService) Update(authorId int, input library.UpdateAuthor) error {
	_, err := s.repo.GetById(authorId)
	if err != nil {
		err = fmt.Errorf("Нет возможности обновить данные о авторе ", err)
		return err
	}
	return s.repo.Update(authorId, input)
}

func checkDateOfBirthday(checkedDate string) error {

	dateOfBirth, err := time.Parse("2006-01-02", checkedDate)

	if err != nil {
		err = fmt.Errorf("Не удалось распарсить дату рождения из переданного JSON", err)
		return err
	}

	// Получим текущую дату и вычислим 5 лет назад (я допускаю, что 6 летний ребенок может написать какое-то письмо)
	// А в дальнейшем можно расширить учет книг на такие письма и прочие фельетоны ))
	now := time.Now()
	fiveYearsAgo := now.AddDate(-5, 0, 0)

	// Проверка, что дата рождения не превышает 5 лет назад от текущей даты
	if dateOfBirth.After(fiveYearsAgo) {
		err := fmt.Errorf("Дата рождения автора не может быть меньше, чем 5 лет от текущей даты")
		return err
	}

	return nil

}
