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

func setUpServer(e *echo.Echo) {
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

func TestITGetAll(t *testing.T) {
	// Setup server
	e := echo.New()
	go func() {
		setUpServer(e)
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

var expBodyUpdate = Expense{
	Title:  "apple smoothie",
	Amount: 89.00,
	Note:   "no discount",
	Tags:   []string{"beverage"},
}

var expBodyCreate = Expense{
	Title:  "strawberry smoothie",
	Amount: 79.00,
	Note:   "night market promotion discount 10 bath",
	Tags:   []string{"food", "beverage"},
}

func encodExpense(body Expense) *bytes.Buffer {
	payload, _ := json.Marshal(body)
	return bytes.NewBuffer(payload)
}
func TestITUpdate(t *testing.T) {
	// Setup server
	e := echo.New()
	go func() {
		setUpServer(e)
	}()
	setUpTimeOut()

	// Arrange
	exp := seedExpense(t)

	body := encodExpense(expBodyUpdate)
	var expense Expense
	expense.Id = exp.Id
	res := request(http.MethodPut, uri("expenses", strconv.Itoa(exp.Id)), body)
	err := res.Decode(&expense)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, expBodyUpdate.Title, expense.Title)
	assert.Equal(t, expBodyUpdate.Amount, expense.Amount)
	assert.Equal(t, expBodyUpdate.Note, expense.Note)
	assert.Equal(t, expBodyUpdate.Tags, expense.Tags)

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
		setUpServer(e)
	}()
	setUpTimeOut()

	// Arrange
	body := encodExpense(expBodyCreate)
	var expense Expense
	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&expense)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, expense.Id)
	assert.Equal(t, expBodyCreate.Title, expense.Title)
	assert.Equal(t, expBodyCreate.Amount, expense.Amount)
	assert.Equal(t, expBodyCreate.Note, expense.Note)
	assert.Equal(t, expBodyCreate.Tags, expense.Tags)

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
		setUpServer(e)
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
		setUpServer(e)
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
