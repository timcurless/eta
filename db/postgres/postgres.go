package postgres

import (
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/namsral/flag"
	"github.com/sirupsen/logrus"
	"github.com/twinj/uuid"
	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/timcurless/eta/user"
)

var (
	addr string
	port string
)

func init() {
	flag.StringVar(&addr, "db-host", "127.0.0.1", "Address of Database Server")
	flag.StringVar(&port, "db-port", "5432", "Port of Database Server")

}

type Postgres struct {
	pDB *gorm.DB
}

func (db *Postgres) Init(dbuser, dbpassword string) error {
	var err error
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/users?sslmode=disable", dbuser, dbpassword, addr, port)
	logrus.Infof("Connecting to DB with connection string: %s", connStr)
	db.pDB, err = gorm.Open("postgres", connStr)
	if err != nil {
		return err
	}
	db.pDB.AutoMigrate(&user.User{})
	return nil
}

func (db *Postgres) GetUsers() ([]user.User, error) {
	user := []user.User{}
	if db.pDB == nil {
		return user, fmt.Errorf("Connection to DB not established")
	}
	db.pDB.Find(&user)
	return user, nil
}

func (db *Postgres) PostUser(newUser user.User) (user.User, error) {
	if db.pDB == nil {
		return user.User{}, fmt.Errorf("Connection to DB not established")
	}
	newUser.ID = uuid.NewV4().String()
	newUser.RewardsID = strconv.Itoa(rand.IntnRange(10000000, 99999999))

	tx := db.pDB.Begin()
	tx = tx.Create(&newUser)
	if tx.Error != nil {
		logrus.Errorf("Error creating User: %v", tx.Error.Error())
		return user.User{}, tx.Error
	}
	tx = tx.Commit()
	if tx.Error != nil {
		logrus.Errorf("Error committing User: %v", tx.Error.Error())
		return user.User{}, tx.Error
	}
	logrus.Infof("Successfully created new user")
	return newUser, nil
}
