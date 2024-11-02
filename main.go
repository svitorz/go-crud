package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Insira o número da operação que deseja realizar: \n1-Inserir\n2-Listar\n3-Editar\n4-Excluir\n0-Parar\n")

		input, err := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")
		operacao, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Erro no sistema%s", err)
		}

		switch operacao {
		case 0:
			os.Exit(2)
		case 1:
			inserir(db)
		case 2:
			listar(db)
		case 3:
			fmt.Println("Insira o ID do usuário que deseja editar:")
			input, err := reader.ReadString('\n')
			input = strings.TrimSuffix(input, "\n")
			id, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Erro no sistema%s", err)
			}
			editar(db,id)
		case 4:
			fmt.Println("Insira o ID do usuário que deseja excluir:")
			input, err := reader.ReadString('\n')
			input = strings.TrimSuffix(input, "\n")
			id, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Erro no sistema%s", err)
			}
			excluir(db,id)
		default:
			fmt.Println("Insira um valor válido.")
		}
	}
}

func HashPassword(password string) (string, error) {
	// Gerar um hash da senha usando bcrypt
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func getValues(){
	var username, password string
	var err error 
	fmt.Println("Insira seu nome de usuário:")
	username,err = reader.ReadString('\n')
	username = strings.TrimSuffix(username, '\n')
	fmt.Println("Insira sua senha:")
	password,err = reader.ReadString('\n')
	password = strings.TrimSuffix(password,'\n')
	if err != nil{
		fmt.Println("Erro no sistema: %s", err)
	}
	return [2]string{username,HashPassword(password)}
}
func inserir(db){

	insertion := db.Query("INSERT INTO USERS(USERNAME, PASSWORD) VALUES ($1,$2)")
}
func listar(){}func editar(){}func excluir(){}
