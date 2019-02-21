package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/namsral/flag"
	"github.com/sirupsen/logrus"

	"github.com/timcurless/eta/db"
	"github.com/timcurless/eta/db/postgres"
	"github.com/timcurless/eta/user"
	"github.com/timcurless/eta/vault"
)

var (
	vaultURL   string
	vaultToken string
	vc         *vault.VaultClient
	dbuser     string
	dbpassword string
	dbconn     bool
	dblease    string
)

const serviceName = "eta"

func init() {
	flag.StringVar(&vaultURL, "vault-url", "https://localhost:8200", "URL (TLS) of Vault Server, including port")
	flag.StringVar(&vaultToken, "vault-token", "", "Vault API Token")
	db.Register("postgres", &postgres.Postgres{})
}

func main() {
	flag.Parse()

	// Logging Domain
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Infof("Starting %v", serviceName)

	logrus.Infof("Connecting to Vault API at %v", vaultURL)
	var err error
	vc, err = vault.NewVaultClient(vaultURL, vaultToken)
	if err != nil {
		logrus.Panicf("Panic: %v", err)
	}

	// Service domain
	dbuser, dbpassword, dblease, err = vc.GetDatabaseCreds()
	if err != nil {
		logrus.Panicf("Panic: %v", err)
	}
	dbconn = false
	for !dbconn {
		err := db.Init(dbuser, dbpassword)
		if err != nil {
			if err == db.ErrNoDatabaseSelected {
				logrus.Fatal(err)
			}
			logrus.Print(err)
		} else {
			dbconn = true
		}
	}

	// Transport Domain
	r := gin.Default()

	r.Use(static.Serve("/", static.LocalFile("./views", true)))

	api := r.Group("/api")
	{
		api.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		api.GET("/health", HealthHandler)
		api.GET("/users", GetUsersHandler)
		api.POST("/users", PostUserHandler)
	}

	errc := make(chan os.Signal, 1)
	signal.Notify(errc, os.Interrupt)
	signal.Notify(errc, syscall.SIGTERM)
	go func() {
		<-errc
		os.Exit(1)
	}()
	r.Run(":3000")
}

func HealthHandler(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	hstatus, err := vc.GetVaultHealth()
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": err,
		})
		return
	}
	d := DatabaseStatus{
		Username:  dbuser,
		Connected: dbconn,
		Engine:    "Postgres",
		Lease:     dblease,
	}
	h := HealthResponse{Health: hstatus, Database: d}
	ctx.JSON(http.StatusOK, h)
}

func GetUsersHandler(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	users, err := db.DefaultDB.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func PostUserHandler(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	var newUser user.User
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	res, err := db.DefaultDB.PostUser(newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res})
}

type HealthResponse struct {
	Health   interface{}    `json:"health"`
	Database DatabaseStatus `json:"database"`
}

type DatabaseStatus struct {
	Username  string `json:"username"`
	Connected bool   `json:"connected"`
	Engine    string `json:"engine"`
	Lease     string `json:"lease"`
}
