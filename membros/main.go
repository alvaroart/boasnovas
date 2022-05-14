package main //Baseado no artigo https://medium.com/baixada-nerd/criando-um-crud-simples-em-go-3640d3618a67

import (
	"database/sql" //
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Membros struct {
	Id    int
	Name  string
	Email string
}

//Função
func dbConn() (db *sql.db) {
	dbDriver := "mysql"
	dbUser := "jarvis"
	dbPass := ""
	dbName := "crudgo"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

//Função usada para renderizar o arquivo Index
func Index(w http.ResponseWriter, r *http.Request) {
	//Abrir conexão com o banco
	db := dbConn()
	//Realiza consulta com o banco e trata erros
	selDB, err := db.Query("SELECT * FROM names ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}

	//monta struct para utilizar no template
	n := Names{}

	//mont array para guardar os valores da struct
	res := []Names{}

	//Realiza a estrutura de repetição pegando todos os valores do banco
	for selDB.Next() {
		//armazena os valores nas variaveis
		var id int
		var name, email string

		//Faz o scan do select
		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		//envia o resultado pra o struct
		n.id = id
		n.name = name
		n.Email = email

		//junta struct com array
		res = append(res, n)
	}

	//open pagina index e exibe resultados na tela
	tmpl.ExecuteTemplate(w, "Index", res)

	//fecha conexão
	defer db.Close()

}

func main() {

	var tmpl = template.Must(template.ParseGlob("tmpl/*"))

	//Exibe mensagem que o servidor foi iniciado
	log.Println("Servidor iniciado em: http://localhost:9000")

	//Gerencia URLs
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)

	//Ações
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/upgrade", Update)
	http.HandleFunc("/delete", Delete)

	//inicia o servidor na orta 9000
	http.HandleFunc(":9000", nil)
}
