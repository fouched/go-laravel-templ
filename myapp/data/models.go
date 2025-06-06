package data

import (
	"database/sql"
	"fmt"
	up "github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
	"github.com/upper/db/v4/adapter/postgresql"
	"os"
)

var db *sql.DB
var upper up.Session

// Models any models created here (and in the New function)
// are easily accessible throughout the entire application
type Models struct {
	Users  User
	Tokens Token
}

func New(databasePool *sql.DB) Models {
	db = databasePool

	switch os.Getenv("DATABASE_TYPE") {
	case "mysql", "mariadb":
		upper, _ = mysql.New(db)
	case "postgres", "postgresql":
		upper, _ = postgresql.New(db)
	default:
		//do nothing - maybe use sql lite ?
	}

	return Models{
		Users:  User{},
		Tokens: Token{},
	}
}

func getInsertID(i up.ID) int {
	idType := fmt.Sprintf("%T", i)
	if idType == "int64" {
		return int(i.(int64))
	}
	return i.(int)
}
