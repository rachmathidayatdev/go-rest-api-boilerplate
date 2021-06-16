package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-rest-api-boilerplate/server/book/models"
	"github.com/go-rest-api-boilerplate/server/book/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/go-rest-api-boilerplate/docs"
)

//InitBookController struct
type InitBookController struct {
	Service *InitBookServiceInterface
}

//InitBookServiceInterface struct
type InitBookServiceInterface struct {
	Book service.BookService
}

//NewBookServer func
func NewBookServer(bookService service.BookService) {

	bookServer := &InitBookController{
		Service: &InitBookServiceInterface{
			Book: bookService,
		},
	}

	s := echo.New()

	s.Use(middleware.Logger())
	s.Use(middleware.Recover())

	http.Handle("/", s)

	// router
	s.GET("/book", bookServer.GetListBook)
	s.GET("/book/:id", bookServer.GetBook)
	s.POST("/book", bookServer.CreateBook)
	s.PUT("/book", bookServer.UpdateBook)
	s.DELETE("/book/:id", bookServer.DeleteBook)

	s.GET("/swagger/*", echoSwagger.WrapHandler)

	apiServicePort := os.Getenv("API_PORT")

	log.Printf("API Service listening on port %v", apiServicePort)

	apiServer := &http.Server{
		Addr:         ":" + apiServicePort,
		ReadTimeout:  20 * time.Minute,
		WriteTimeout: 20 * time.Minute,
	}

	err := apiServer.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		log.Println(err.Error())
	}
}

// GetListBook godoc
// @Summary Get list of book
// @Description Get list of book item
// @Tags book
// @Accept */*
// @Produce json
// @Success 200 {object} models.SuccessResponseList
// @Failure 500 {object} models.SuccessResponseObject
// @Router /book [get]
// GetListBook func
func (init *InitBookController) GetListBook(ctx echo.Context) error {
	books, err := init.Service.Book.ListBook()
	if err != nil {
		data := &models.SuccessResponseList{
			Status:  500,
			Message: "failed",
			Data:    make([]*models.Book, 0),
		}

		return ctx.JSON(http.StatusInternalServerError, data)
	}

	data := &models.SuccessResponseList{
		Status:  200,
		Message: "success",
		Data:    books,
	}

	return ctx.JSON(http.StatusOK, data)
}

// GetBook godoc
// @Summary Get a book
// @Description Get a book item
// @Tags book
// @Accept */*
// @Produce json
// @Param id path integer true "Book ID"
// @Success 200 {object} models.SuccessResponseObject
// @Failure 500 {object} models.SuccessResponseObject
// @Router /book/{id} [get]
// GetBook func
func (init *InitBookController) GetBook(ctx echo.Context) error {
	idParse := ctx.Param("id")
	id, _ := strconv.ParseInt(idParse, 10, 64)

	book, err := init.Service.Book.GetBook(id)
	if err != nil {
		data := &models.SuccessResponseObject{
			Status:  500,
			Message: "failed",
			Data:    new(models.Book),
		}

		return ctx.JSON(http.StatusInternalServerError, data)
	}

	data := &models.SuccessResponseObject{
		Status:  200,
		Message: "success",
		Data:    book,
	}

	return ctx.JSON(http.StatusOK, data)
}

// CreateBook godoc
// @Summary Create a book
// @Description Create a new book item
// @Tags book
// @Accept json
// @Produce json
// @Param book body models.Book true "Param Book"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.SuccessResponse
// @Router /book [post]
// CreateBook func
func (init *InitBookController) CreateBook(ctx echo.Context) error {
	var book *models.Book
	err := ctx.Bind(&book)
	if err != nil {
		data := &models.SuccessResponse{
			Status:  400,
			Message: err.Error(),
		}

		return ctx.JSON(http.StatusBadRequest, data)
	}

	err = book.Validate()
	if err != nil {
		data := &models.SuccessResponse{
			Status:  400,
			Message: err.Error(),
		}

		return ctx.JSON(http.StatusBadRequest, data)
	}

	book, err = init.Service.Book.CreateBook(book)
	if err != nil {
		data := &models.SuccessResponse{
			Status:  500,
			Message: err.Error(),
		}

		return ctx.JSON(http.StatusInternalServerError, data)
	}

	data := &models.SuccessResponse{
		Status:  201,
		Message: fmt.Sprintf("Create book success #%d", book.ID),
	}

	return ctx.JSON(http.StatusCreated, data)
}

// UpdateBook godoc
// @Summary Update a book
// @Description Update a book item
// @Tags book
// @Accept json
// @Produce json
// @Param book body models.Book true "Param Book"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.SuccessResponse
// @Failure 500 {object} models.SuccessResponse
// @Router /book [put]
// UpdateBook func
func (init *InitBookController) UpdateBook(ctx echo.Context) error {
	var book *models.Book
	err := ctx.Bind(&book)
	if err != nil {
		data := &models.SuccessResponse{
			Status:  400,
			Message: err.Error(),
		}

		return ctx.JSON(http.StatusBadRequest, data)
	}

	book, err = init.Service.Book.UpdateBook(book)
	if err != nil {
		data := &models.SuccessResponse{
			Status:  500,
			Message: err.Error(),
		}

		return ctx.JSON(http.StatusInternalServerError, data)
	}

	data := &models.SuccessResponse{
		Status:  200,
		Message: "success",
	}

	return ctx.JSON(http.StatusOK, data)
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Delete a book item
// @Tags book
// @Accept */*
// @Produce json
// @Param id path integer true "Book ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 500 {object} models.SuccessResponse
// @Router /book/{id} [delete]
// DeleteBook func
func (init *InitBookController) DeleteBook(ctx echo.Context) error {
	idParse := ctx.Param("id")
	id, _ := strconv.ParseInt(idParse, 10, 64)

	err := init.Service.Book.DeleteBook(id)
	if err != nil {
		data := &models.SuccessResponse{
			Status:  500,
			Message: "failed",
		}

		return ctx.JSON(http.StatusInternalServerError, data)
	}

	data := &models.SuccessResponse{
		Status:  200,
		Message: "success",
	}

	return ctx.JSON(http.StatusOK, data)
}
