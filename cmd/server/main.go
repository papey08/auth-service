package main

import (
	_ "auth-service/docs/swagger"
	"auth-service/internal/app"
	"auth-service/internal/repo"
	"auth-service/internal/server"
	"auth-service/pkg/tokenizer"
	"context"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

// InitConfig initializes configuration file
func InitConfig() error {
	viper.SetConfigFile("configs/config.yml")
	return viper.ReadInConfig()
}

//	@title			auth-service
//	@version		1.0
//	@description	Часть сервиса аутентификации, состоит из двух маршрутов

//	@host		localhost:8080
//	@BasePath	/auth/v1

func main() {
	ctx := context.Background()

	// initializing config info
	if err := InitConfig(); err != nil {
		log.Fatal("config file read error: ", err.Error())
	}

	// connecting to the database
	repoURL := fmt.Sprintf("mongodb://%s:%d",
		viper.GetString("repo.host"),
		viper.GetInt("repo.port"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(repoURL))
	if err != nil {
		log.Fatal("database connection error: ", err.Error())
	}
	defer func(ctx context.Context, client *mongo.Client) {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal("user_repo disconnect error:", err)
		}
	}(ctx, client)
	collection := client.Database(viper.GetString("repo.dbname")).Collection(viper.GetString("repo.collection"))

	// configuring app
	r := repo.New(ctx, collection)
	t := tokenizer.New(randomdata.Paragraph())
	a := app.New(r, t)

	// configuring server
	s := server.NewHTTPServer(a, viper.GetString("server.host"), viper.GetInt("server.port"))

	// preparing graceful shutdown
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Println("Starting http server")
		if err = s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal("can't listen and serve server:", err.Error())
		}
	}()

	// waiting for Ctrl+C
	<-osSignals

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // 30s timeout to finish all active connections
	defer cancel()

	if err = s.Shutdown(ctx); err != nil {
		log.Fatal("Server graceful shutdown failed:", err.Error())
	}
	log.Println("Server was gracefully stopped")
}
