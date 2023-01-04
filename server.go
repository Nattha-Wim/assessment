package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"
	"github.com/nattha-wim/assessment/expense"
)

func setUpServer(e *echo.Echo) {
	serverPort := ":2565"

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	handler := expense.NewApplication(db)
	handler.InitDB()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "admin" && password == "45678" {
			return true, nil
		}
		return false, nil
	}))

	e.GET("/", handler.HomeExpenses)
	e.POST("/expenses", handler.CreateExpense)
	e.GET("/expenses/:id", handler.GetExpenseById)
	e.GET("/expenses", handler.GetAllExpenses)
	e.PUT("/expenses/:id", handler.UpdateExpenses)

	log.Println("server starting at " + serverPort)
	if err := e.Start(serverPort); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal("shutting down the server")
	}
}
func main() {
	e := echo.New()
	go func() {
		setUpServer(e)
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
