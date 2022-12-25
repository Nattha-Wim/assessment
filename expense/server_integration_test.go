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

func TestGetExpenseById(t *testing.T) {
	c := seedExpense(t)
	var latest Expense
	res := request(http.MethodGet, uri("expenses", strconv.Itoa(c.Id)), nil)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	assert.Equal(t, c.Id, latest.Id)
	assert.NotEmpty(t, latest.Title)
	assert.NotEmpty(t, latest.Amount)
	assert.NotEmpty(t, latest.Note)
	assert.NotEmpty(t, latest.Tags)

}

func TestCreateExpense(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title": "salmon don & water",
		"amount": 350.00,
		"note": "dinner with friend at friday night", 
		"tags": ["food", "beverage"]
		}`)

	var detailExp Expense
	detailTags := []string{"food", "beverage"}
	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&detailExp)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, detailExp.Id)
	assert.Equal(t, "salmon don & water", detailExp.Title)
	assert.Equal(t, 350.00, detailExp.Amount)
	assert.Equal(t, "dinner with friend at friday night", detailExp.Note)
	assert.Equal(t, detailTags, detailExp.Tags)

}

func seedExpense(t *testing.T) Expense {
	var exp Expense
	body := bytes.NewBufferString(`{
		"title": "salmon don & water",
		"amount": 350.00,
		"note": "dinner with friend at friday night", 
		"tags": ["food", "beverage"]
		}`)
	err := request(http.MethodPost, uri("expenses"), body).Decode(&exp)
	if err != nil {
		t.Fatal("can't create user", err)
	}
	return exp
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

type Response struct {
	*http.Response // มีของที่ response กลับมาหมดเลย
	err            error
}

// put things(user) that we want to interface{}
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
