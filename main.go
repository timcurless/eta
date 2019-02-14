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

	"github.com/timcurless/eta/vault"
)

var (
	vaultURL   string
	vaultToken string
	vc         *vault.VaultClient
)

const serviceName = "eta"

func init() {
	flag.StringVar(&vaultURL, "vault-url", "https://localhost:8200", "URL (TLS) of Vault Server, including port")
	flag.StringVar(&vaultToken, "vault-token", "", "Vault API Token")
	flag.Parse()
}

func main() {
	// Logging Domain
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Infof("Starting %v", serviceName)

	logrus.Infof("Connecting to Vault API at %v", vaultURL)
	var err error
	vc, err = vault.NewVaultClient(vaultURL, vaultToken)
	if err != nil {
		logrus.Panicf("Panic: %v", err)
	}

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
	health, err := vc.GetVaultHealth()
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"message": err,
		})
	}
	ctx.JSON(http.StatusOK, health)
}
