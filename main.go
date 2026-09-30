package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type User struct {
	ID      int64  `json:"id"`   // Добавили теги, чтобы в JSON ключи
	Name    string `json:"name"` // были красивыми (маленькими буквами)
	Balance int64  `json:"balance"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	password := os.Getenv("POSTGRES_PASSWORD")

	conn, err := pgx.Connect(
		context.Background(),
		fmt.Sprintf("postgres://bonus:%s@localhost:5432/bonus", password),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
		if err != nil {
			http.Error(w, "Bad ID", http.StatusBadRequest)
			return
		}

		user, err := getUserByID(r.Context(), conn, userID)
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getUserByID(ctx context.Context, conn *pgx.Conn, userID int64) (User, error) {
	var user User
	err := conn.QueryRow(ctx, "SELECT id, name, balance FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Name, &user.Balance)
	if err != nil {
		return user, err
	}
	return user, nil
}
