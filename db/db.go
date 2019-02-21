package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/namsral/flag"

	"github.com/timcurless/eta/user"
)

type Database interface {
	Init(dbuser, dbpassword string) error
	GetUsers() ([]user.User, error)
	PostUser(newUser user.User) (user.User, error)
}

var (
	database              string
	DefaultDB             Database
	DBTypes               = map[string]Database{}
	ErrNoDatabaseFound    = "No database with name %v found"
	ErrNoDatabaseSelected = errors.New("No User DB Selected")
)

func init() {
	flag.StringVar(&database, "user_database", os.Getenv("user_database"), "Database to use for User Service")
}

func Init(dbuser, dbpassword string) error {
	if database == "" {
		return ErrNoDatabaseSelected
	}
	err := Set()
	if err != nil {
		return err
	}
	return DefaultDB.Init(dbuser, dbpassword)
}

func Register(name string, db Database) {
	DBTypes[name] = db
}

func Set() error {
	if v, ok := DBTypes[database]; ok {
		DefaultDB = v
		return nil
	}
	return fmt.Errorf(ErrNoDatabaseFound, database)
}
