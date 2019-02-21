package cockroachdb

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

func (db *Cockroach) PostUser(newUser user.User) (user.User, error) {
	if db.crDB == nil {
		return user.User{}, fmt.Errorf("Connection to DB not established")
	}
	newUser.ID = uuid.NewV4().String()
	newUser.RewardsID = strconv.Itoa(rand.IntnRange(10000000, 99999999))

	tx := db.crDB.Begin()
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
