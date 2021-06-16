package repository

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-rest-api-boilerplate/server/book/models"
	"github.com/go-rest-api-boilerplate/util/dbtrxn"
)

// BookRepository to get book data from databasesa
type BookRepository interface {
	List() ([]*models.Book, error)
	Find(id int64) (*models.Book, error)
	Insert(ctx context.Context, book *models.Book) (*models.Book, error)
	Update(ctx context.Context, book *models.Book) (*models.Book, error)
	Delete(ctx context.Context, id int64) error
}

//InitBookRepository struct
type InitBookRepository struct {
	connection *sql.DB
}

// NewBookRepository return new instance of BookRepository
func NewBookRepository(connection *sql.DB) BookRepository {
	return &InitBookRepository{
		connection: connection,
	}
}

//List func
func (init *InitBookRepository) List() (list []*models.Book, err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select(BookColumns...).From(bookTable).OrderBy("id ASC")

	rows, err := builder.RunWith(init.connection).Query()

	list = make([]*models.Book, 0)

	if err != nil {
		return list, err
	}

	for rows.Next() {
		var book *models.Book
		book, err = models.ScanBook(rows)
		if err != nil {
			return
		}
		list = append(list, book)
	}

	return list, err
}

//Find func
func (init *InitBookRepository) Find(id int64) (book *models.Book, err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select(BookColumns...).
		From(bookTable).
		Where(sq.Eq{idColumn: id})

	rows, err := builder.RunWith(init.connection).Query()
	if err != nil {
		return book, err
	}

	if rows.Next() {
		book, err = models.ScanBook(rows)
	}

	return book, err
}

//Insert func
func (init *InitBookRepository) Insert(ctx context.Context, book *models.Book) (*models.Book, error) {
	trxn, err := dbtrxn.Use(ctx, init.connection)

	query := sq.Insert(bookTable).
		Columns(bookTitleColumn, bookAuthorColumn).
		Values(book.Title, book.Author).
		Suffix("RETURNING \"id\"").
		RunWith(trxn.DB).
		PlaceholderFormat(sq.Dollar)

	err = query.QueryRow().Scan(&book.ID)
	if err != nil {
		trxn.SetError(err)
		return book, err
	}

	return book, err
}

//Update func
func (init *InitBookRepository) Update(ctx context.Context, book *models.Book) (*models.Book, error) {
	trxn, err := dbtrxn.Use(ctx, init.connection)

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Update(bookTable).
		Set(bookTitleColumn, book.Title).
		Set(bookAuthorColumn, book.Author).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: book.ID})

	_, err = builder.RunWith(trxn.DB).Exec()

	if err != nil {
		trxn.SetError(err)
		return book, err
	}

	return book, err
}

//Delete func
func (init *InitBookRepository) Delete(ctx context.Context, id int64) error {
	trxn, err := dbtrxn.Use(ctx, init.connection)

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Delete(bookTable).
		Where(sq.Eq{idColumn: id})

	_, err = builder.RunWith(trxn.DB).Exec()
	if err != nil {
		trxn.SetError(err)
		return err
	}

	return err
}
