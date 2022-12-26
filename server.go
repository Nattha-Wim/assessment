package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"
	"github.com/nattha-wim/assessment/expense"
)

func main() {

	expense.InitDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover()) // เผื่อ server เรา down
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "admin" && password == "45678" {
			return true, nil
		}
		return false, nil
	}))

	e.GET("/", expense.HomeExpenses)
	e.POST("/expenses", expense.CreateExpenses)
	e.GET("/expenses/:id", expense.GetExpensesById)
	e.GET("/expenses", expense.GetAllExpenses)
	e.PUT("/expenses/:id", expense.UpdateExpenses)

	go func() {
		log.Println("server starting at :2565")
		if err := e.Start(":2565"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("shuting down.....")
	if err := e.Shutdown(context.Background()); err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("Bye Bye")

}
