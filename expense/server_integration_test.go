//go:build integration
// +build integration

package expense

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	//"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const serverPort = 2565

// func TestITGetById(t *testing.T) {
// 	exp := seedExpense(t)

// 	// Setup server
// 	eh := echo.New()
// 	go func(e *echo.Echo) {
// 		db, err := sql.Open("postgres", "postgres://vpovznnb:ayqqQAENpjSG6STGdF5CMxXGni5DAhj0@tiny.db.elephantsql.com/vpovznnb")
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		h := NewApplication(db)

// 		e.GET("/expenses/:id", h.GetExpenseById)
// 		e.Start(fmt.Sprintf(":%d", serverPort))
// 	}(eh)
// 	for {
// 		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 70*time.Second)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		if conn != nil {
// 			conn.Close()
// 			break
// 		}
// 	}

// 	reqBody := ``
// 	var latest Expense
// 	res := request(http.MethodGet, uri("expenses", strconv.Itoa(exp.Id)), strings.NewReader(reqBody))
// 	err := res.Decode(&latest)

// 	assert.Nil(t, err)
// 	assert.Equal(t, http.StatusOK, res.StatusCode)
// 	assert.Equal(t, exp.Id, latest.Id)
// 	assert.NotEmpty(t, latest.Title)
// 	assert.NotEmpty(t, latest.Amount)
// 	assert.NotEmpty(t, latest.Note)
// 	assert.NotEmpty(t, latest.Tags)

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	err = eh.Shutdown(ctx)
// 	assert.NoError(t, err)

// }

func TestITGetAll(t *testing.T) {

	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgres://vpovznnb:ayqqQAENpjSG6STGdF5CMxXGni5DAhj0@tiny.db.elephantsql.com/vpovznnb")
		if err != nil {
			log.Fatal(err)
		}

		h := NewApplication(db)

		e.GET("/expenses", h.GetAllExpenses)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 70*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	// Arrange
	reqBody := ``
	var expense []Expense
	res := request(http.MethodGet, uri("expenses"), strings.NewReader(reqBody))
	err := res.Decode(&expense)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Greater(t, len(expense), 0)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestITUpdate(t *testing.T) {
	exp := seedExpense(t)
	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgres://vpovznnb:ayqqQAENpjSG6STGdF5CMxXGni5DAhj0@tiny.db.elephantsql.com/vpovznnb")
		if err != nil {
			log.Fatal(err)
		}

		h := NewApplication(db)

		e.PUT("/expenses/:id", h.UpdateExpenses)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	// Arrange

	reqBody := Expense{
		Title:  "apple smoothie",
		Amount: 89.00,
		Note:   "no discount",
		Tags:   []string{"beverage"},
	}
	payload, _ := json.Marshal(reqBody)
	var expense Expense
	res := request(http.MethodPut, uri("expenses", strconv.Itoa(exp.Id)), bytes.NewBuffer(payload))
	err := res.Decode(&expense)
	expense.Id = exp.Id
	log.Println("expense seed update : ", expense)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, reqBody.Title, expense.Title)
	assert.Equal(t, reqBody.Amount, expense.Amount)
	assert.Equal(t, reqBody.Note, expense.Note)
	assert.Equal(t, reqBody.Tags, expense.Tags)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}
func TestITCreate(t *testing.T) {
	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgres://vpovznnb:ayqqQAENpjSG6STGdF5CMxXGni5DAhj0@tiny.db.elephantsql.com/vpovznnb")
		if err != nil {
			log.Fatal(err)
		}

		h := NewApplication(db)

		e.POST("/expenses", h.CreateExpense)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 100*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79.00,
		"note": "night market promotion discount 10 bath",
		"tags": ["food","beverage"]
		}`)

	var expense Expense
	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&expense)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, expense.Id)
	assert.Equal(t, "strawberry smoothie", expense.Title)
	assert.Equal(t, 79.00, expense.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", expense.Note)
	assert.Equal(t, []string{"food", "beverage"}, expense.Tags)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestHome(t *testing.T) {
	eh := echo.New()
	go func(e *echo.Echo) {

		h := NewApplication(nil)

		e.GET("/", h.HomeExpenses)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	res := request(http.MethodGet, uri(), nil)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	eh.Shutdown(ctx)

}

func seedExpense(t *testing.T) Expense {
	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgres://vpovznnb:ayqqQAENpjSG6STGdF5CMxXGni5DAhj0@tiny.db.elephantsql.com/vpovznnb")
		if err != nil {
			log.Fatal(err)
		}

		h := NewApplication(db)

		e.POST("/expenses", h.CreateExpense)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79.00,
		"note": "night market promotion discount 10 bath",
		"tags": ["food","beverage"]
		}`)

	var expense Expense
	res := request(http.MethodPost, uri("expenses"), body)
	res.Decode(&expense)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	eh.Shutdown(ctx)
	return expense
}

func uri(paths ...string) string {

	host := fmt.Sprintf("http://localhost:%d", serverPort)
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
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

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

// func TestHomeExpenses(t *testing.T) {
// 	res := request(http.MethodGet, uri(), nil)
// 	assert.Equal(t, http.StatusOK, res.StatusCode)
// }

// func TestGetAllExpenses(t *testing.T) {

// 	seedExpense(t)

// 	var expense []Expense
// 	res := request(http.MethodGet, uri("expenses"), nil)
// 	err := res.Decode(&expense)

// 	assert.Nil(t, err)
// 	assert.Equal(t, http.StatusOK, res.StatusCode)
// 	assert.Greater(t, len(expense), 0)
// }

// func TestGetExpenseById(t *testing.T) {

// 	exp := seedExpense(t)

// 	var latest Expense
// 	res := request(http.MethodGet, uri("expenses", strconv.Itoa(exp.Id)), nil)
// 	err := res.Decode(&latest)

// 	assert.Nil(t, err)
// 	assert.Equal(t, http.StatusOK, res.StatusCode)
// 	assert.Equal(t, exp.Id, latest.Id)
// 	assert.NotEmpty(t, latest.Title)
// 	assert.NotEmpty(t, latest.Amount)
// 	assert.NotEmpty(t, latest.Note)
// 	assert.NotEmpty(t, latest.Tags)

// }

// func TestUpdateExpense(t *testing.T) {
// 	exp := seedExpense(t)
// 	expense := Expense{
// 		Title:  "apple smoothie",
// 		Amount: 89.00,
// 		Note:   "no discount",
// 		Tags:   []string{"beverage"},
// 	}
// 	payload, _ := json.Marshal(expense)

// 	var latest Expense
// 	res := request(http.MethodPut, uri("expenses", strconv.Itoa(exp.Id)), bytes.NewBuffer(payload))
// 	//defer res.Body.Close()
// 	err := res.Decode(&latest)
// 	latest.Id = exp.Id

// 	assert.Nil(t, err)
// 	assert.Equal(t, http.StatusOK, res.StatusCode)
// 	assert.Equal(t, expense.Title, latest.Title)
// 	assert.Equal(t, expense.Amount, latest.Amount)
// 	assert.Equal(t, expense.Note, latest.Note)
// 	assert.Equal(t, expense.Tags, latest.Tags)

// }

// func TestCreateExpense(t *testing.T) {
// 	body := bytes.NewBufferString(`{
// 		"title": "strawberry smoothie",
// 		"amount": 79.00,
// 		"note": "night market promotion discount 10 bath",
// 		"tags": ["food","beverage"]
// 		}`)

// 	var expense Expense
// 	res := request(http.MethodPost, uri("expenses"), body)
// 	err := res.Decode(&expense)

// 	assert.Nil(t, err)
// 	assert.Equal(t, http.StatusCreated, res.StatusCode)
// 	assert.NotEqual(t, 0, expense.Id)
// 	assert.Equal(t, "strawberry smoothie", expense.Title)
// 	assert.Equal(t, 79.00, expense.Amount)
// 	assert.Equal(t, "night market promotion discount 10 bath", expense.Note)
// 	assert.Equal(t, []string{"food", "beverage"}, expense.Tags)
// }
