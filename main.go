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
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id       int
	username string
	password string
}

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
			users, err := listar(db)
			if err != nil {
				log.Fatalf("Erro ao listar usuários: %v", err)
			}

			// Melhorar a visualização dos dados
			fmt.Println("ID\tUsuário")
			fmt.Println("--------------------")
			for _, user := range users {
				fmt.Printf("%d\t%s\n", user.id, user.username)
			}
		case 3:
			fmt.Println("Insira o ID do usuário que deseja editar:")
			input, err := reader.ReadString('\n')
			input = strings.TrimSuffix(input, "\n")
			id, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Erro no sistema: ", err)
			}
			editar(db, id)
		case 4:
			fmt.Println("Insira o ID do usuário que deseja excluir:")
			input, err := reader.ReadString('\n')
			input = strings.TrimSuffix(input, "\n")
			id, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Erro no sistema: ", err)
			}
			excluir(db, id)
		default:
			fmt.Println("Insira um valor válido.")
		}
	}
}

/**
* @param Uma string que vai ser criptografada
* @return a criptografia da string
* */
func HashPassword(password string) string {
	// gera um hash baseado no algoritmo bcrypt
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

/**
* @param senha criptografada e senha sem criptografia
* @return a verificação se as senhas conferem
* */
func VerifyPassword(db *sql.DB, id int) bool {
	fmt.Println("Insira sua senha para poder realizar a edição:")
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Erro ao ler senha: ", err)
	}
	password = strings.TrimSpace(password)

	// Verificar se a senha corresponde ao hash
	var userPassword string
	query := `SELECT password FROM USERS WHERE id = $1;`

	// Executa a consulta e armazena o resultado na variável userPassword
	err = db.QueryRow(query, id).Scan(&userPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("nenhum usuário encontrado com ID %d", id)
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))
	return err == nil
}

/**
* @return O retorno deve ser os valores inseridos pelo usuário após tratamento, em caso de erros a função nao deve retornar nada para nao comprometer o sql.
* */
func getValues() ([2]string, error) {
	// cria um novo leitor para o console
	reader := bufio.NewReader(os.Stdin)

	// solicita o username e o trata retirando os espaços. Valores diferentes de String retornam erros.
	fmt.Println("Insira seu nome de usuário:")
	username, err := reader.ReadString('\n')
	if err != nil {
		return [2]string{}, fmt.Errorf("Erro ao ler nome de usuário: %w", err)
	}
	username = strings.TrimSpace(username)

	fmt.Println("Insira sua senha:")

	// Tratamento da senha, valores diferentes de String retornam erros.
	password, err := reader.ReadString('\n')
	if err != nil {
		return [2]string{}, fmt.Errorf("Erro ao ler senha: %w", err)
	}
	password = strings.TrimSpace(password)

	// criptografia da senha
	hashedPassword := HashPassword(password)

	// retorna o nome do usuário e a senha já criptografada. para nao causar erros no banvo de dados, o retorno pode ser nulo.
	return [2]string{username, hashedPassword}, nil
}

/**
* @param
* */
func inserir(db *sql.DB) bool {
	// pega os valores inseridos pelo usuário
	values, err := getValues()
	if err != nil {
		fmt.Println("Erro ao obter valores:", err)
		return false
	}
	// inicia o banco de dados.
	txn, err := db.Begin()
	if err != nil {
		fmt.Println("Erro ao iniciar a transação:", err)
		return false
	}

	// Preparar a declaração
	stmt, err := txn.Prepare("INSERT INTO USERS(USERNAME, PASSWORD) VALUES ($1, $2)")
	if err != nil {
		fmt.Println("Erro ao preparar declaração:", err)
		// em caso de erro, o rollback desfaz a inserção no banco.
		txn.Rollback()
		return false
	}
	defer stmt.Close()

	// Executar a inserção
	_, err = stmt.Exec(values[0], values[1])
	if err != nil {
		fmt.Println("Erro ao inserir:", err)
		txn.Rollback()
		return false
	}

	// Confirmar a transação
	err = txn.Commit()
	if err != nil {
		fmt.Println("Erro ao confirmar a transação:", err)
		return false
	}

	fmt.Println("Inserido com sucesso.")
	return true
}

func listar(db *sql.DB) ([]User, error) {
	// query para selecionar todos os usuários
	rows, err := db.Query("SELECT id,username FROM USERS")
	if err != nil {
		return nil, err
	}
	var users []User
	defer rows.Close()

	for rows.Next() {
		var row User
		if err := rows.Scan(&row.id, &row.username); err != nil {
			return users, err
		}
		users = append(users, row)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}
	return users, nil
}

func editar(db *sql.DB, id int) bool {
	if !VerifyPassword(db, id) {
		fmt.Println("Você não está autorizado a realizar esta operação")
		return false
	}
	values, err := getValues()
	if err != nil {
		fmt.Println("Erro ao pegar valores")
	}

	txn, err := db.Begin()
	if err != nil {
		fmt.Println("Erro ao inicializar banco de dados")
	}

	stmt, err := txn.Prepare("UPDATE USERS SET USERNAME = $1, PASSWORD = $2 WHERE ID = $3")
	if err != nil {
		fmt.Println("Erro ao preparar declaração:", err)
		txn.Rollback()
	}
	defer stmt.Close()

	// Executar a edição
	_, err = stmt.Exec(values[0], values[1], id)
	if err != nil {
		fmt.Println("Erro ao inserir:", err)
		txn.Rollback()
	}

	err = txn.Commit()
	if err != nil {
		fmt.Println("Erro ao confirmar a transação:", err)
	}

	fmt.Println("Alterado com sucesso com sucesso.")
	return true
}

func excluir(db *sql.DB, id int) bool {
	if !VerifyPassword(db, id) {
		fmt.Println("Você não está autorizado a realizar esta operação")
		return false
	}
	// inicia o banco de dados.
	txn, err := db.Begin()
	if err != nil {
		fmt.Println("Erro ao iniciar a transação:", err)
		return false
	}

	// Preparar a declaração
	stmt, err := txn.Prepare("DELETE FROM USERS WHERE ID = $1")
	if err != nil {
		fmt.Println("Erro ao preparar declaração:", err)
		// em caso de erro, o rollback desfaz a inserção no banco.
		txn.Rollback()
		return false
	}
	defer stmt.Close()

	// Executar a inserção
	_, err = stmt.Exec(id)
	if err != nil {
		fmt.Println("Erro ao inserir:", err)
		txn.Rollback()
		return false
	}

	// Confirmar a transação
	err = txn.Commit()
	if err != nil {
		fmt.Println("Erro ao confirmar a transação:", err)
		return false
	}

	fmt.Println("Usuário excluído com sucesso.")

	return true
}
