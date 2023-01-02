//go:build integration
// +build integration

package expense

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const serverPort = "2565"

func TestHomeExpenses(t *testing.T) {
	res := request(http.MethodGet, uri(), nil)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetAllExpenses(t *testing.T) {
	seedExpense(t)

	var expense []Expense
	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&expense)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(expense), 0)
}

func TestGetExpenseById(t *testing.T) {
	exp := seedExpense(t)

	var latest Expense
	res := request(http.MethodGet, uri("expenses", strconv.Itoa(exp.Id)), nil)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, exp.Id, latest.Id)
	assert.NotEmpty(t, latest.Title)
	assert.NotEmpty(t, latest.Amount)
	assert.NotEmpty(t, latest.Note)
	assert.NotEmpty(t, latest.Tags)

}
func TestUpdateExpense(t *testing.T) {
	exp := seedExpense(t)
	expense := Expense{
		Title:  "apple smoothie",
		Amount: 89.00,
		Note:   "no discount",
		Tags:   []string{"beverage"},
	}
	payload, _ := json.Marshal(expense)

	var latest Expense
	res := request(http.MethodPut, uri("expenses", strconv.Itoa(exp.Id)), bytes.NewBuffer(payload))
	defer res.Body.Close()
	err := res.Decode(&latest)
	latest.Id = exp.Id

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, expense.Title, latest.Title)
	assert.Equal(t, expense.Amount, latest.Amount)
	assert.Equal(t, expense.Note, latest.Note)
	assert.Equal(t, expense.Tags, latest.Tags)

}

func TestCreateExpense(t *testing.T) {
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
}

func seedExpense(t *testing.T) Expense {
	var exp Expense
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
    	"amount": 79.00,
    	"note": "night market promotion discount 10 bath",
   		"tags": ["food","beverage"]
		}`)

	err := request(http.MethodPost, uri("expenses"), body).Decode(&exp)
	if err != nil {
		t.Fatal("can't create expense", err)
	}
	return exp
}

func uri(paths ...string) string {
	host := "http://localhost:" + serverPort
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
