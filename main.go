package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func getEnv(key string) string {
	// Carrega as variáveis do .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}
	return os.Getenv(key)
}

func databaseConnection() (*sql.DB, error) {
	dbUser := getEnv("user")
	dbName := getEnv("database")
	dbPass := getEnv("password")
	dbPort := getEnv("port")
	sslMode := getEnv("sslmode")
	dbHost := getEnv("host")
	dbDriver := getEnv("driver")

	// Monta a string de conexão
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPass, dbHost, dbPort, dbName, sslMode,
	)

	// Abre a conexão com o banco
	db, err := sql.Open(dbDriver, connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão: %w", err)
	}

	// Testa a conexão
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("erro ao testar conexão: %w", err)
	}

	return db, nil
}

func main() {
	db, err := databaseConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Conexão com o banco de dados bem-sucedida!")
}
