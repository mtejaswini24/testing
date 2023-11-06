package main

import (
	"context"
	"fmt"
	"job-portal-api/internal/auth"
	"job-portal-api/internal/database"
	"job-portal-api/internal/handlers"
	"job-portal-api/internal/repository"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func main() {
	err := startApp() //calling the start app
	if err != nil {
		log.Panic().Err(err).Send()
	}
	log.Info().Msg("hello this is job portal app")
}

func startApp() error {
	log.Info().Msg("main : Started : Initializing authentication support")
	//message to the developer
	privatePEM, err := os.ReadFile("private.pem")
	//reading the file and returning the byte format
	if err != nil {
		return fmt.Errorf("reading auth private key %w", err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM) //converting the byte format and saving the address of privatekey
	if err != nil {
		return fmt.Errorf("parsing auth private key %w", err)
	}

	publicPEM, err := os.ReadFile("pubkey.pem") //reading the file and returning the byte format
	if err != nil {
		return fmt.Errorf("reading auth public key %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		return fmt.Errorf("parsing auth public key %w", err)
	}

	a, err := auth.NewAuth(privateKey, publicKey)
	if err != nil {
		return fmt.Errorf("constructing auth %w", err)
	}
	//connection of DB
	log.Info().Msg("main : Started : Initializing db support")
	db, err := database.Open()
	if err != nil {
		return fmt.Errorf("connecting to db %w", err)
	}
	pg, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = pg.PingContext(ctx) //verifies the database connection is there are not
	if err != nil {
		return fmt.Errorf("db is not connected %w", err)
	}
	//initialize conn layer support
	ms, err := repository.NewRepo(db)
	if err != nil {
		return err
	}
	// err = AutoMigrate()
	// if err != nil {
	// 	return err
	// }
	api := http.Server{
		Addr:         ":8082",
		ReadTimeout:  8000 * time.Second,
		WriteTimeout: 800 * time.Second,
		IdleTimeout:  800 * time.Second,
		Handler:      handlers.API(a, ms),
	}
	// channel to store any errors while setting up the service
	serverErrors := make(chan error, 1)
	go func() {
		log.Info().Str("port", api.Addr).Msg("main: API listening")
		serverErrors <- api.ListenAndServe()
	}()
	shutdown := make(chan os.Signal, 1)
	//shutdown is just an empty channel
	signal.Notify(shutdown, os.Interrupt)
	//notify gives the value to the shutdown channel if interrupt occur(interrupt occurs when we click ctrl+C)
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error %w", err)
	case sig := <-shutdown:
		log.Info().Msgf("main: Start shutdown %s", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		//ctx is the input of context
		defer cancel()
		//after timeout cancel function wil work
		err := api.Shutdown(ctx)
		//api.Shutdown is graceful shutdown
		//shutdown is taking context
		if err != nil {
			err = api.Close() // forcing shutdown
			return fmt.Errorf("could not stop server gracefully %w", err)
		}
	}
	return nil
}
