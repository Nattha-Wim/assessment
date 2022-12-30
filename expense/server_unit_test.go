//go:build unit

package expense

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var expense = Expense{
	Id:     1,
	Title:  "mango juice",
	Amount: 80,
	Note:   "reading",
	Tags:   []string{"beverage"},
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestHomeExpense(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	h := handler{}
	c := e.NewContext(req, rec)

	// Act
	err := h.HomeExpenses(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestGetExpenseById(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses/1", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	db, mock := NewMock()
	repo := handler{db}

	query := "SELECT (.+) FROM expenses WHERE id = ?"
	rows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(expense.Id, expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags))

	mock.ExpectPrepare(query).ExpectQuery().WillReturnRows(rows)

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Act
	err := repo.GetExpenseById(c)

	expected := `{"Id":1,"Title":"mango juice","Amount":80,"Note":"reading","Tags":["beverage"]}`
	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestGetAllExpenses(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	db, mock := NewMock()
	repo := handler{db}

	tag := []string{"beverage"}
	rows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(1, "orange juice", 90, "no discount", pq.Array(tag)).
		AddRow(2, "apple juice", 90, "no discount", pq.Array(tag))

	mock.ExpectPrepare("^SELECT (.+) FROM expenses").ExpectQuery().WillReturnRows(rows)

	c := e.NewContext(req, rec)

	expected := `[{"Id":1,"Title":"orange juice","Amount":90,"Note":"no discount","Tags":["beverage"]},{"Id":2,"Title":"apple juice","Amount":90,"Note":"no discount","Tags":["beverage"]}]`

	// Act
	err := repo.GetAllExpenses(c)

	// Assertions
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

// 	db, mock := NewMock()
// 	repo := handler{db}

// 	query := "UPDATE expenses SET title=?, amount=?, note=?, tags=? WHERE id = ?"
// 	prep := mock.ExpectPrepare(query)
// 	prep.ExpectExec().WithArgs("orange juice", 20, "good", sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(0, 1))

// 	// mock.ExpectQuery(query).WillReturnRows(rows)

// 	c := e.NewContext(req, rec)
// 	c.SetParamNames("id")
// 	c.SetParamValues("1")

// 	// Act
// 	err := repo.UpdateExpenses(c)

// 	log.Println(strings.TrimSpace(rec.Body.String()))
// 	// Assertions
// 	if assert.NoError(t, err) {
// 		assert.Equal(t, http.StatusOK, rec.Code)
// 		//assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
// 	}
// }
