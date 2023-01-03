//go:build unit

package expense

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host + "/"
	}
	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func ExpenseNewRows(t *testing.T) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(1, "orange juice", 90, "no discount", pq.Array([]string{"beverage"})).
		AddRow(2, "apple juice", 100, "no discount", pq.Array([]string{"beverage"}))
	return rows
}

func request(method, url string, body io.Reader) *http.Request {
	req := httptest.NewRequest(http.MethodGet, url, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return req
}

var expense = Expense{
	Title:  "orange juice",
	Amount: 120.00,
	Note:   "no discount",
	Tags:   []string{"beverage"},
}

func detailExpense() *bytes.Buffer {
	body, _ := json.Marshal(expense)
	return bytes.NewBuffer(body)
}

func TestHomeExpense(t *testing.T) {
	e := echo.New()
	req := request(http.MethodGet, uri(), nil)
	rec := httptest.NewRecorder()
	h := handler{nil}
	c := e.NewContext(req, rec)

	err := h.HomeExpenses(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestGetExpenseById(t *testing.T) {
	e := echo.New()
	req := request(http.MethodGet, uri("expenses", strconv.Itoa(1)), nil)
	rec := httptest.NewRecorder()
	db, mock := NewMock()
	h := handler{db}
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	rows := ExpenseNewRows(t)
	mock.ExpectQuery("SELECT (.+) FROM expenses WHERE id = ?").WillReturnRows(rows)

	err := h.GetExpenseById(c)

	expected := `{"id":1,"title":"orange juice","amount":90,"note":"no discount","tags":["beverage"]}`
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestGetAllExpenses(t *testing.T) {
	e := echo.New()
	req := request(http.MethodGet, uri("expenses"), nil)
	rec := httptest.NewRecorder()
	db, mock := NewMock()
	h := handler{db}
	c := e.NewContext(req, rec)

	rows := ExpenseNewRows(t)
	mock.ExpectPrepare("SELECT (.+) FROM expenses").ExpectQuery().WillReturnRows(rows)

	err := h.GetAllExpenses(c)

	expected := `[{"id":1,"title":"orange juice","amount":90,"note":"no discount","tags":["beverage"]},{"id":2,"title":"apple juice","amount":100,"note":"no discount","tags":["beverage"]}]`
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestCreateExpenses(t *testing.T) {
	e := echo.New()
	req := request(http.MethodPost, uri("expenses"), detailExpense())
	rec := httptest.NewRecorder()
	db, mock := NewMock()
	h := handler{db}
	c := e.NewContext(req, rec)

	query := "INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id"
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags)).WillReturnRows(rows)

	err := h.CreateExpense(c)

	expected := `{"id":1,"title":"orange juice","amount":120,"note":"no discount","tags":["beverage"]}`
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestUpdateExpenses(t *testing.T) {
	e := echo.New()
	req := request(http.MethodPut, uri("expenses", strconv.Itoa(1)), detailExpense())
	rec := httptest.NewRecorder()
	db, mock := NewMock()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := handler{db}

	query := "UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id = $1"
	prep := mock.ExpectPrepare(regexp.QuoteMeta(query))
	prep.ExpectExec().WithArgs(1, expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags)).WillReturnResult(sqlmock.NewResult(1, 1))

	err := h.UpdateExpenses(c)

	expected := `{"id":1,"title":"orange juice","amount":120,"note":"no discount","tags":["beverage"]}`
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}
