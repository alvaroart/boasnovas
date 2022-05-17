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

//Função show exibe apenas um resultado
func Show(w http, ResponseWriter, r *http.Request) {
	db := dbConn()

	//pegar o ID do parametro da URL
	nId := r.URL.Query().Get("id")

	//usa o ID para fazer a consulta e tratar erros
	selDB, err := db.Query("SELECT * FROM names WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}

	//monta a struct para ser utulizada no template
	n := Names{}

	// executa a estrutura de repetição pegando todos os valores do banco
	for selDB.Next() {
		//armazena os valore nas variaveis
		var id int
		var name, email string

		//faz o scan do select
		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		// Envia os resultados para a struct
		n.Id = id
		n.Name = name
		n.Email = email
	}

	// Mostra o template
	tmpl.ExecuteTemplate(w, "Show", n)

	// Fecha a conexão
	defer db.Close()

}

// Função New apenas exibe o formulário para inserir novos dados
func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

// Função Edit, edita os dados
func Edit(w http.ResponseWriter, r *http.Request) {
	// Abre a conexão com banco de dados
	db := dbConn()

	// Pega o ID do parametro da URL
	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT * FROM names WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}

	// Monta a struct para ser utilizada no template
	n := Names{}

	// Realiza a estrutura de repetição pegando todos os valores do banco
	for selDB.Next() {
		//Armazena os valores em variaveis
		var id int
		var name, email string

		// Faz o Scan do SELECT
		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		// Envia os resultados para a struct
		n.Id = id
		n.Name = name
		n.Email = email
	}

	// Mostra o template com formulário preenchido para edição
	tmpl.ExecuteTemplate(w, "Edit", n)

	// Fecha a conexão com o banco de dados
	defer db.Close()
}

// Função Insert, insere valores no banco de dados
func Insert(w http.ResponseWriter, r *http.Request) {

	//Abre a conexão com banco de dados usando a função: dbConn()
	db := dbConn()

	// Verifica o METHOD do fomrulário passado
	if r.Method == "POST" {

		// Pega os campos do formulário
		name := r.FormValue("name")
		email := r.FormValue("email")

		// Prepara a SQL e verifica errors
		insForm, err := db.Prepare("INSERT INTO names(name, email) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}

		// Insere valores do formulario com a SQL tratada e verifica errors
		insForm.Exec(name, email)

		// Exibe um log com os valores digitados no formulário
		log.Println("INSERT: Name: " + name + " | E-mail: " + email)
	}

	// Encerra a conexão do dbConn()
	defer db.Close()

	//Retorna a HOME
	http.Redirect(w, r, "/", 301)
}

// Função Update, atualiza valores no banco de dados
func Update(w http.ResponseWriter, r *http.Request) {

	// Abre a conexão com o banco de dados usando a função: dbConn()
	db := dbConn()

	// Verifica o METHOD do formulário passado
	if r.Method == "POST" {

		// Pega os campos do formulário
		name := r.FormValue("name")
		email := r.FormValue("email")
		id := r.FormValue("uid")

		// Prepara a SQL e verifica errors
		insForm, err := db.Prepare("UPDATE names SET name=?, email=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}

		// Insere valores do formulário com a SQL tratada e verifica erros
		insForm.Exec(name, email, id)

		// Exibe um log com os valores digitados no formulario
		log.Println("UPDATE: Name: " + name + " |E-mail: " + email)
	}

	// Encerra a conexão do dbConn()
	defer db.Close()

	// Retorna a HOME
	http.Redirect(w, r, "/", 301)
}

// Função Delete, deleta valores no banco de dados
func Delete(w http.ResponseWriter, r *http.Request) {

	// Abre conexão com banco de dados usando a função: dbConn()
	db := dbConn()

	nId := r.URL.Query().Get("id")

	// Prepara a SQL e verifica errors
	delForm, err := db.Prepare("DELETE FROM names WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	// Insere valores do form com a SQL tratada e verifica errors
	delForm.Exec(nId)

	// Exibe um log com os valores digitados no form
	log.Println("DELETE")

	// Encerra a conexão do dbConn()
	defer db.Close()

	// Retorna a HOME
	http.Redirect(w, r, "/", 301)
}

//TESTE
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
