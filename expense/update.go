package expense

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) UpdateExpenses(c echo.Context) error {
	//id := c.Param("id")
	rowID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "id should be int " + err.Error()})
	}

	var exp Expense
	err = c.Bind(&exp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	stmt, err := h.db.Prepare("UPDATE expenses SET title=?, amount=?, note=?, tags=? WHERE id = ?")
	log.Println("====", stmt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if _, err := stmt.Exec(rowID, exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags)); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, exp)
}

// -------------
// func (h *handler) UpdateExpenses1(c echo.Context) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	rowID, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, Err{Message: "id should be int " + err.Error()})
// 	}
// 	var exp Expense
// 	err = c.Bind(&exp)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
// 	}

// 	query := "UPDATE expenses SET title=?, amount=?, note=?, tags=? WHERE id = ?"
// 	stmt, err := h.db.PrepareContext(ctx, query)
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()

// 	if _, err = stmt.ExecContext(ctx, exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags), rowID); err != nil {
// 		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
// 	}
// 	return c.JSON(http.StatusOK, exp)

// }
