package service

import (
	"context"

	"github.com/go-rest-api-boilerplate/server/book/models"
	"github.com/go-rest-api-boilerplate/server/book/repository"
	"github.com/go-rest-api-boilerplate/util/dbtrxn"
)

//BookService interface
type BookService interface {
	ListBook() ([]*models.Book, error)
	GetBook(id int64) (*models.Book, error)
	CreateBook(book *models.Book) (*models.Book, error)
	UpdateBook(book *models.Book) (*models.Book, error)
	DeleteBook(id int64) error
}

//InitBookService struct
type InitBookService struct {
	Repository *InitBookRepositoryInterface
}

//InitBookRepositoryInterface struct
type InitBookRepositoryInterface struct {
	Book repository.BookRepository
}

// NewBookService return new instance of BookRepository
func NewBookService(bookRepository repository.BookRepository) BookService {
	return &InitBookService{
		Repository: &InitBookRepositoryInterface{
			Book: bookRepository,
		},
	}
}

//ListBook func
func (init *InitBookService) ListBook() ([]*models.Book, error) {
	books, err := init.Repository.Book.List()

	if err != nil {
		return books, err
	}

	return books, err
}

//GetBook func
func (init *InitBookService) GetBook(id int64) (*models.Book, error) {
	book, err := init.Repository.Book.Find(id)

	if err != nil {
		return book, err
	}

	return book, err
}

//CreateBook func
func (init *InitBookService) CreateBook(book *models.Book) (*models.Book, error) {
	//start transaction
	ctx := context.Background()
	defer dbtrxn.Begin(&ctx)()

	book, err := init.Repository.Book.Insert(ctx, book)

	//transaction commit or rollback if error
	dbtrxn.Error(ctx)

	return book, err
}

//UpdateBook func
func (init *InitBookService) UpdateBook(book *models.Book) (*models.Book, error) {
	//start transaction
	ctx := context.Background()
	defer dbtrxn.Begin(&ctx)()

	book, err := init.Repository.Book.Update(ctx, book)

	//transaction commit or rollback if error
	dbtrxn.Error(ctx)

	return book, err
}

//DeleteBook func
func (init *InitBookService) DeleteBook(id int64) error {
	//start transaction
	ctx := context.Background()
	defer dbtrxn.Begin(&ctx)()

	err := init.Repository.Book.Delete(ctx, id)

	//transaction commit or rollback if error
	dbtrxn.Error(ctx)

	return err
}
