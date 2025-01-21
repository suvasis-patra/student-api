package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/suvasis-patra/student-api/internal/config"
	"github.com/suvasis-patra/student-api/internal/http/handlers/student"
	"github.com/suvasis-patra/student-api/internal/storage/sqlite"
)

func main() {
	// load config file
	cfg := config.MustLoad()
	// sutup and connect with database
	db,err:=sqlite.New(cfg)
	if err!=nil{
		log.Fatal(err)
	}
	slog.Info("server started")
	// setup routes
	route := http.NewServeMux()
	route.HandleFunc("POST /api/students", student.New(db))
	route.HandleFunc("GET /api/student/{id}",student.GetStudentById(db))
	route.HandleFunc("GET /api/students",student.GetAllStudents(db))
	route.HandleFunc("PUT /api/student/{id}",student.UpdateStudentDetailsById(db))
	route.HandleFunc("DELETE /api/student/{id}",student.DeleteStudentById(db))
	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: route,
	}
	fmt.Println("Server started!")
	// shutdown the server gracefully

	// 1. create a GO channel to catch the interrupt signals by the system of size 1
	done := make(chan os.Signal, 1)
	// 2. Notify the channel about listed signals
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// 3. Start the server in a different GO routine
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to start the server!: %v",err)
		}
	}()
	// 4. unblock the execution after getting interrupt signal
	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown the server", slog.String("error", err.Error()))
	}
	slog.Info("server shutdown")
}
