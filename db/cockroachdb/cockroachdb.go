package cockroachdb

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/namsral/flag"
	"github.com/sirupsen/logrus"

	"github.com/timcurless/eta/user"
)

var (
	addr string
)

func init() {
	flag.StringVar(&addr, "db-host", "127.0.0.1:5679", "Address of Database Server")
}

type Cockroach struct {
	crDB *gorm.DB
}

func (db *Cockroach) Init() error {
	logrus.Infof("Connecting to DB with connection string: %v", addr)
	var err error
	db.crDB, err = gorm.Open("postgres", addr)
	if err != nil {
		return err
	}
	db.crDB.AutoMigrate(&user.User{})
	return nil
}

func (db *Cockroach) GetUsers() ([]user.User, error) {
	user := []user.User{}
	if db.crDB == nil {
		return user, fmt.Errorf("Connection to DB not established")
	}
	db.crDB.Find(&user)
	return user, nil
}
