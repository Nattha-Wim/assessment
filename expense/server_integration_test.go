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

func seedExpense(t *testing.T) Expense {
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79.00,
		"note": "night market promotion discount 10 bath",
		"tags": ["food","beverage"]
		}`)

	var expense Expense
	res := request(http.MethodPost, uri("expenses"), body)
	res.Decode(&expense)
	return expense
}

func setUpDB(e *echo.Echo) {
	db, err := sql.Open("postgres", "postgres://vpovznnb:ayqqQAENpjSG6STGdF5CMxXGni5DAhj0@tiny.db.elephantsql.com/vpovznnb")
	if err != nil {
		log.Fatal(err)
	}

	h := NewApplication(db)
	e.GET("/", h.HomeExpenses)
	e.GET("/expenses", h.GetAllExpenses)
	e.POST("/expenses", h.CreateExpense)
	e.PUT("/expenses/:id", h.UpdateExpenses)
	e.GET("/expenses/:id", h.GetExpenseById)

	e.Start(fmt.Sprintf(":%d", serverPort))
}

func setUpTimeOut() {
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
}

// func shutDown() {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	err = eh.Shutdown(ctx)
// 	assert.NoError(t, err)
// }

func TestITGetAll(t *testing.T) {
	// Setup server
	e := echo.New()
	go func() {
		setUpDB(e)
	}()
	setUpTimeOut()

	// Arrange
	var expense []Expense
	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&expense)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Greater(t, len(expense), 0)
	}

	// shutdown
	if err := e.Shutdown(context.Background()); err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("Bye Bye")
}

func TestITUpdate(t *testing.T) {
	// Setup server
	e := echo.New()
	go func() {
		setUpDB(e)
	}()
	setUpTimeOut()

	// Arrange
	exp := seedExpense(t)

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

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, reqBody.Title, expense.Title)
	assert.Equal(t, reqBody.Amount, expense.Amount)
	assert.Equal(t, reqBody.Note, expense.Note)
	assert.Equal(t, reqBody.Tags, expense.Tags)

	// shutdown
	if err := e.Shutdown(context.Background()); err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("Bye Bye")
}
func TestITCreate(t *testing.T) {
	// Setup server
	e := echo.New()
	go func() {
		setUpDB(e)
	}()
	setUpTimeOut()

	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79.00,
		"note": "night market promotion discount 10 bath",
		"tags": ["food","beverage"]
		}`)

	var expense Expense
	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&expense)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, expense.Id)
	assert.Equal(t, "strawberry smoothie", expense.Title)
	assert.Equal(t, 79.00, expense.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", expense.Note)
	assert.Equal(t, []string{"food", "beverage"}, expense.Tags)

	// shutdown
	if err := e.Shutdown(context.Background()); err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("Bye Bye")
}

func TestITGetById(t *testing.T) {
	// Setup server
	e := echo.New()
	go func() {
		setUpDB(e)
	}()
	setUpTimeOut()

	// Arrange
	expense := seedExpense(t)

	var latest Expense
	res := request(http.MethodGet, uri("expenses", strconv.Itoa(expense.Id)), nil)
	err := res.Decode(&latest)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, expense.Id, latest.Id)
	assert.NotEmpty(t, latest.Title)
	assert.NotEmpty(t, latest.Amount)
	assert.NotEmpty(t, latest.Note)
	assert.NotEmpty(t, latest.Tags)

	// shutdown
	if err := e.Shutdown(context.Background()); err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("Bye Bye")
}

func TestHome(t *testing.T) {
	// Setup server
	e := echo.New()
	go func() {
		setUpDB(e)
	}()
	setUpTimeOut()

	// Arrange
	res := request(http.MethodGet, uri(), nil)

	// Assertions
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// shutdown
	if err := e.Shutdown(context.Background()); err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("Bye Bye")
}
