package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

var serviceName = "http-service"
var servicePort = "8080"
var hostname,_ = os.Hostname()

type Response struct{
	ServiceName string `json:"service-name"`
	Message string `json:"message"`
	Hostname string `json:"hostname"`
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("received a DELETE request on /items %v\n", r)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(Response{
		ServiceName: serviceName,
		Message:     "received delete items request",
		Hostname: hostname,
	})
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("received a POST request on /items %v\n", r)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(Response{
		ServiceName: serviceName,
		Message:     "update items request received",
		Hostname: hostname,
	})
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("received a PUT request on /items %v\n", r)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{
		ServiceName: serviceName,
		Message:     "put items request received",
		Hostname: hostname,
	})
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("received a GET request on /items %v\n", r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		ServiceName: serviceName,
		Message:     "list items request received",
		Hostname: hostname,
	})
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("received a GET request on / %v\n", r)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		ServiceName: serviceName,
		Message:     fmt.Sprintf("This is a http service[%s]! If you see this then the " +
			"service is deployed as working as expected :)", serviceName),
		Hostname: hostname,
	})
}

func main() {
	if name, ok := os.LookupEnv("SERVICE_NAME"); ok {
		serviceName = name
	}
	if port, ok := os.LookupEnv("SERVICE_PORT"); ok {
		servicePort = port
	}
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second * 15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler).Methods("GET")
	r.HandleFunc("/items", listHandler).Methods("GET")
	r.HandleFunc("/items", putHandler).Methods("PUT")
	r.HandleFunc("/items", updateHandler).Methods("POST")
	r.HandleFunc("/items", deleteHandler).Methods("DELETE")

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", servicePort),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler: r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}