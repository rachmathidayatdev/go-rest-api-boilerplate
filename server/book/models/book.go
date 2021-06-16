package models

import (
	"database/sql"
	"time"

	validator "gopkg.in/go-playground/validator.v9"
)

// Book represented database model
type Book struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required"`
	Author    string    `json:"author" validate:"required"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}

// SuccessResponseList struct
type SuccessResponseList struct {
	Status  int64   `json:"status"`
	Message string  `json:"message"`
	Data    []*Book `json:"data"`
}

// SuccessResponseObject struct
type SuccessResponseObject struct {
	Status  int64  `json:"status"`
	Message string `json:"message"`
	Data    *Book  `json:"data"`
}

// SuccessResponse struct
type SuccessResponse struct {
	Status  int64  `json:"status"`
	Message string `json:"message"`
}

//ScanBook func
func ScanBook(rows *sql.Rows) (*Book, error) {
	var book Book
	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.UpdatedAt, &book.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// Validate book
func (b *Book) Validate() error {
	validate := validator.New()
	return validate.Struct(b)
}
