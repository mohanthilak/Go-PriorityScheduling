package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {

	PORT := 8001
	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(4 * time.Second)
		w.WriteHeader(http.StatusAccepted)
		// response := map[string]string{"success": "true"}
		// json.NewEncoder(w).Encode(response)
		w.Write([]byte("response"))
	})

	server := http.Server{
		Addr:         "127.0.0.1:" + strconv.Itoa(PORT),
		Handler:      handlers.CORS(handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}), handlers.AllowCredentials(), handlers.AllowedHeaders([]string{"Access-Control-Allow-Origin", "X-Requested-With", "Authorization", "Content-Type"}), handlers.AllowedOrigins([]string{"http://localhost:8000"}))(router),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Println("server starting at port: ", PORT)
		err := server.ListenAndServe()
		if err != nil {
			log.Panicf("error while serving the server. %s", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c

	log.Printf("Gracefully Ending the Server after receiving signal %s", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}
