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
	//"github.com/google/uuid"
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

func seedExpense(t *testing.T) *sqlmock.Rows {
	tag := []string{"beverage"}
	rows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(1, "orange juice", 90, "no discount", pq.Array(tag)).
		AddRow(2, "apple juice", 100, "no discount", pq.Array(tag))
	return rows
}

func request(method, url string, body io.Reader) *http.Request {
	req := httptest.NewRequest(http.MethodGet, url, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return req
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	return json.NewDecoder(r.Body).Decode(v)
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
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	rows := seedExpense(t)
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

	rows := seedExpense(t)
	mock.ExpectPrepare("SELECT (.+) FROM expenses").ExpectQuery().WillReturnRows(rows)

	err := h.GetAllExpenses(c)

	expected := `[{"id":1,"title":"orange juice","amount":90,"note":"no discount","tags":["beverage"]},{"id":2,"title":"apple juice","amount":100,"note":"no discount","tags":["beverage"]}]`
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestCreateExpenses(t *testing.T) {
	body := strings.NewReader(`{
		"title": "orange juice",
		"amount": 90.00,
		"note": "no discount",
		"tags": ["beverage"]
		}`)

	e := echo.New()
	req := request(http.MethodPost, uri("expenses"), body)
	rec := httptest.NewRecorder()
	db, mock := NewMock()
	h := handler{db}
	c := e.NewContext(req, rec)

	query := "INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id"
	prep := mock.ExpectQuery(regexp.QuoteMeta(query))
	prep.WithArgs("orange juice", 90.00, "no discount", pq.Array([]string{"beverage"})).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err := h.CreateExpense(c)

	expected := `{"id":1,"title":"orange juice","amount":90,"note":"no discount","tags":["beverage"]}`
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestUpdateExpenses(t *testing.T) {
	expense := Expense{
		Title:  "orange juice",
		Amount: 120.00,
		Note:   "no discount",
		Tags:   []string{"beverage"},
	}
	payload, _ := json.Marshal(expense)

	e := echo.New()
	req := request(http.MethodPut, uri("expenses", strconv.Itoa(1)), bytes.NewBuffer(payload))
	rec := httptest.NewRecorder()
	db, mock := NewMock()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := handler{db}

	prep := mock.ExpectPrepare(regexp.QuoteMeta("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id = $1"))
	prep.ExpectExec().WithArgs(1, "orange juice", 120.00, "no discount", pq.Array([]string{"beverage"})).WillReturnResult(sqlmock.NewResult(1, 1))

	err := h.UpdateExpenses(c)

	expected := `{"id":1,"title":"orange juice","amount":120,"note":"no discount","tags":["beverage"]}`
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

// var u = Expense{
// 	Id:     1,
// 	Title:  "strawberry smoothie",
// 	Amount: 79.00,
// 	Note:   "night market promotion discount 10 bath",
// 	Tags:   []string{"food", "beverage"},
// }

// func TestFind(t *testing.T) {
// 	e := echo.New()
// 	req := request(http.MethodGet, uri("expenses"), nil)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	db, mock := NewMock()
// 	repo := handler{db}

// 	query := "SELECT id, name, title, note, tags FROM expenses"

// 	rows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
// 		AddRow(u.Id, u.Title, u.Amount, u.Note, pq.Array(u.Tags))

// 	mock.ExpectQuery(query).WillReturnRows(rows)

// 	err := repo.GetAll(c)

// 	expected := `[{"id":1,"title":"strawberry smoothie","amount":79,"note":"night market promotion discount 10 bath","tags":["food","beverage"]}]`
// 	if assert.NoError(t, err) {
// 		assert.Equal(t, http.StatusOK, rec.Code)
// 		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
// 	}
// }

// func TestUpdate(t *testing.T) {

// 	e := echo.New()
// 	req := request(http.MethodPut, uri("expenses/1"), nil)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	c.SetParamNames("id")
// 	c.SetParamValues("1")
// 	db, mock := NewMock()
// 	repo := handler{db}

// 	query := "UPDATE expenses SET title = \\?, amount = \\?, note = \\?, tags = \\? WHERE id = \\?"

// 	prep := mock.ExpectPrepare(query)
// 	prep.ExpectExec().WithArgs(u.Title, u.Amount, u.Note, u.Tags, u.Id).WillReturnResult(sqlmock.NewResult(0, 1))

// 	expected := `[{"id":1,"title":"strawberry smoothie","amount":79,"note":"night market promotion discount 10 bath","tags":["food","beverage"]}]`

// 	err := repo.Update(c)
// 	if assert.NoError(t, err) {
// 		assert.Equal(t, http.StatusOK, rec.Code)
// 		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
// 	}
// }

// func TestCreate(t *testing.T) {
// 	e := echo.New()
// 	req := request(http.MethodPost, uri("expenses1"), nil)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	db, mock := NewMock()
// 	repo := handler{db}

// 	query := "INSERT INTO expenses \\(id, title, amount, note, tags\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)"

// 	prep := mock.ExpectPrepare(query)
// 	prep.ExpectExec().WithArgs(u.Id, u.Title, u.Amount, u.Note, u.Tags).WillReturnResult(sqlmock.NewResult(0, 1))

// 	err := repo.Create(c)
// 	assert.NoError(t, err)
// }
