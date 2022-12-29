package expense

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) GetExpenseById(c echo.Context) error {
	id := c.Param("id")
	rowID, err := strconv.Atoi(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "id should be int " + err.Error()})
	}
	// stmt, err := h.db.Prepare("SELECT * FROM expenses WHERE id = $1")
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query user statment" + err.Error()})
	// }

	res := h.db.QueryRow("SELECT * FROM expenses WHERE id = $1", rowID)
	detailExp := Expense{}
	err = res.Scan(&detailExp.Id, &detailExp.Title, &detailExp.Amount, &detailExp.Note, pq.Array(&detailExp.Tags))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't sacan user: !!! " + err.Error()})
	}
	return c.JSON(http.StatusOK, detailExp)

}
