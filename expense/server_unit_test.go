//go:build unit

package expense

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
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

func seedExpense(t *testing.T) *sqlmock.Rows {
	tag := []string{"beverage"}
	rows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(1, "orange juice", 90, "no discount", pq.Array(tag)).
		AddRow(2, "apple juice", 100, "no discount", pq.Array(tag))
	return rows
}

func request(method, url string, body io.Reader) *http.Request {
	req := httptest.NewRequest(http.MethodGet, uri("expenses", strconv.Itoa(1)), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return req
}
func TestHomeExpense(t *testing.T) {
	e := echo.New()
	req := request(http.MethodGet, uri(), nil)
	rec := httptest.NewRecorder()
	h := handler{}
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

	rows := seedExpense(t)
	mock.ExpectQuery("SELECT (.+) FROM expenses WHERE id = ?").WillReturnRows(rows)

	err := h.GetExpenseById(c)

	expected := `{"Id":1,"Title":"orange juice","Amount":90,"Note":"no discount","Tags":["beverage"]}`
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

	expected := `[{"Id":1,"Title":"orange juice","Amount":90,"Note":"no discount","Tags":["beverage"]},{"Id":2,"Title":"apple juice","Amount":100,"Note":"no discount","Tags":["beverage"]}]`
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

// func TestUpdateExpensestemp(t *testing.T) {
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPut, "/expenses/1", nil)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	rec := httptest.NewRecorder()
// 	//db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
// 	db, mock := NewMock()
// 	h := handler{db}

// 	query := "UPDATE expenses SET title=?, amount=?, note=?, tags=? WHERE id = ?"
// 	prep := mock.ExpectPrepare(query)

// 	prep.ExpectExec().WithArgs("orange juice", 20, "good", sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(1, 1))
// 	//log.Println("++++++++++++++", prep.ExpectExec().WithArgs("orange juice", 20, "good", sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(0, 5)))

// 	c := e.NewContext(req, rec)
// 	c.SetParamNames("id")
// 	c.SetParamValues("1")

// 	// Act
// 	err := h.UpdateExpenses(c)

// 	log.Println(strings.TrimSpace(rec.Body.String()))
// 	// Assertions
// 	if assert.NoError(t, err) {
// 		assert.Equal(t, http.StatusOK, rec.Code)
// 		//assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
// 	}
// }
