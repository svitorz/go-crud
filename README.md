# CRUD simples em GO.

Utilizei o pacote (pq)[https://pkg.go.dev/github.com/lib/pq] para trabalhar com o banco Postgres.

Para utilizar, você precisa de um arquivo .env.

Recomendação de .env:

```
database=go_crud # nome do banco
port= 5432 #porta (padrão 5342)
user=postgres #nome do usuário no banco
password=root #senha do banco
host=localhost #host
driver=postgres #driver do pacote pq
sslmode=disable #recomendação apenas para ambiente de desenvolvimento
```

Script de criação da tabela do banco de dados -- Bando de dados POSTGRES (em outros bancos pode haver erros):

````
```sql
CREATE TABLE USERS(
  ID SERIAL PRIMARY KEY,
  USERNAME VARCHAR(152),
  PASSWORD VARCHAR(255)
);
```
````
