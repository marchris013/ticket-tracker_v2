package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type Ticket struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) createTicket(w http.ResponseWriter, r *http.Request) {
	var ticket Ticket
	if err := json.NewDecoder(r.Body).Decode(&ticket); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := s.db.QueryRow(
		"INSERT INTO tickets (title, description) VALUES ($1, $2) RETURNING id, title, description, status, created_at",
		ticket.Title, ticket.Description,
	).Scan(&ticket.ID, &ticket.Title, &ticket.Description, &ticket.Status, &ticket.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to create ticket", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ticket)
}

func (s *Store) getTickets(w http.ResponseWriter, _ *http.Request) {
	rows, err := s.db.Query("SELECT id, title, description, status, created_at FROM tickets")
	if err != nil {
		http.Error(w, "Failed to get tickets", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	tickets := []Ticket{}
	for rows.Next() {
		var ticket Ticket
		rows.Scan(&ticket.ID, &ticket.Title, &ticket.Description, &ticket.Status, &ticket.CreatedAt)
		tickets = append(tickets, ticket)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tickets)
}

func (s *Store) getTicket(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tickets/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}
	var ticket Ticket
	err = s.db.QueryRow(
		"SELECT id, title, description, status, created_at FROM tickets WHERE id = $1", id,
	).Scan(&ticket.ID, &ticket.Title, &ticket.Description, &ticket.Status, &ticket.CreatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to get ticket", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticket)
}

func main() {
	connStr := "host=localhost port=5432 user=postgres password=admin dbname=ticket_tracker sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	store := NewStore(db)
	http.HandleFunc("/tickets", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			store.getTickets(w, r)
		case http.MethodPost:
			store.createTicket(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/tickets/", func(w http.ResponseWriter, r *http.Request) {
		store.getTicket(w, r)
	})
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
