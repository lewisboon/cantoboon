package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lewisboon/cantoboon/internal/db"
	"github.com/lewisboon/cantoboon/internal/queries"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

type server struct {
	pool *pgxpool.Pool
}

func (s *server) homeHandler(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "home.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *server) lookupHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	entries, err := queries.LookupExact(r.Context(), s.pool, q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Query   string
		Entries []queries.Entry
	}{Query: q, Entries: entries}

	if err := templates.ExecuteTemplate(w, "lookup.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func main() {
	ctx := context.Background()

	pool, err := db.NewPool(ctx)
	if err != nil {
		log.Fatalf("connecting to database: %v", err)
	}
	defer pool.Close()

	s := &server{pool: pool}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.homeHandler)
	mux.HandleFunc("GET /lookup", s.lookupHandler)
	mux.HandleFunc("GET /healthz", healthHandler)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
