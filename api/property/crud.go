package property

/*
import (
	"etender/mysql"
	"fmt"

	"github.com/gin-gonic/gin"
)

func CreateOne(c *gin.Context) {

	mySql := mysql.MysqlDB()
	stmt, err := mySql.Prepare("INSERT INTO division(name) VALUES(?)")
	if err != nil {
		fmt.Printf("[Testsql] Error %v\n", err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec("division_updated")

	if err != nil {
		fmt.Printf("[Testsql] Insertion Error %v\n", err.Error())
	}

	defer mySql.Close()
}


//Index func to view all the records
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()

	var users []User

	result, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var post User
		err = result.Scan(&post.ID, &post.Name, &post.Country, &post.Number)
		if err != nil {
			panic(err.Error())
		}

		users = append(users, post)
	}
	json.NewEncoder(w).Encode(users)
	defer db.Close()
}

func insertUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	vars := mux.Vars(r)
	Name := vars["name"]
	Country := vars["country"]
	Number := vars["number"]

	// perform a db.Query insert
	stmt, err := db.Prepare("INSERT INTO users(name, country, number) VALUES(?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(Name, Country, Number)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New user was created")
	defer db.Close()
}
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	params := mux.Vars(r)

	// perform a db.Query insert
	stmt, err := db.Query("SELECT * FROM users WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	var post User
	for stmt.Next() {

		err = stmt.Scan(&post.ID, &post.Name, &post.Country, &post.Number)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(post)
	defer db.Close()
}
func delUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	params := mux.Vars(r)

	// perform a db.Query insert
	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	fmt.Fprintf(w, "User with ID = %s was deleted", params["id"])
	defer db.Close()
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	params := mux.Vars(r)
	Name := params["name"]
	Country := params["country"]
	Number := params["number"]

	// perform a db.Query insert
	stmt, err := db.Prepare("Update users SET name = ?, country = ?, number = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(Name, Country, Number, params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "User with ID = %s was updated", params["id"])
	defer db.Close()
}

*/
