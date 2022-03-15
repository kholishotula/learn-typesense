package rest

import (
	"encoding/json"
	"learn-typesense/internal/core/entity"
	"learn-typesense/internal/core/service"
	"net/http"

	"github.com/gorilla/mux"
)

type API struct {
	bookService service.Service
}

func (a *API) Handler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/book/create", a.createBook())
	r.HandleFunc("/book/get", a.getBookById())
	r.HandleFunc("/book/search", a.searchBooks())

	return r
}

func (a *API) createBook() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			WriteRespBody(w, NewErrorResp(NewErrMethodNotAllowed()))
		}
		// parse request body
		var rb createReqBody
		err := json.NewDecoder(r.Body).Decode(&rb)
		if err != nil {
			WriteRespBody(w, NewErrorResp(NewErrBadRequest(err.Error())))
			return
		}
		err = rb.Validate()
		if err != nil {
			WriteRespBody(w, NewErrorResp(err))
			return
		}

		book, err := a.bookService.CreateBook(r.Context(), entity.CreateInput{
			Title:   rb.Title,
			Author:  rb.Author,
			Year:    rb.Year,
			Summary: rb.Summary,
		})
		if err != nil {
			WriteRespBody(w, NewErrorResp(err))
			return
		}

		WriteRespBody(w, NewSuccessResp(book))
	}
}

func (a *API) getBookById() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			WriteRespBody(w, NewErrorResp(NewErrMethodNotAllowed()))
		}

		bookId := r.URL.Query().Get("id")
		if bookId == "" {
			WriteRespBody(w, NewErrorResp(NewErrBadRequest("invalid book id")))
			return
		}

		book, err := a.bookService.GetBookById(r.Context(), bookId)
		if err != nil {
			WriteRespBody(w, NewErrorResp(err))
			return
		}

		WriteRespBody(w, NewSuccessResp(book))
	}
}

func (a *API) searchBooks() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			WriteRespBody(w, NewErrorResp(NewErrMethodNotAllowed()))
		}

		books, err := a.bookService.SearchBooks(r.Context())
		if err != nil {
			WriteRespBody(w, NewErrorResp(err))
			return
		}

		WriteRespBody(w, NewSuccessResp(books))
	}
}

type APIConfig struct {
	BookService service.Service
}

func NewAPI(cfg APIConfig) *API {
	return &API{bookService: cfg.BookService}
}
