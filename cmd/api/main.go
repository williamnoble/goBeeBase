package main

import (
	"flag"
	"fmt"
	"goBeeBase/cmd/data"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "0.1"

type Config struct {
	env string
	port int
	db struct {
		user string
		pass string
		host string
		port string
		database string
		ssldisable string
	}
}

type application struct {
	config Config
	infoLogger *log.Logger
	errorLogger *log.Logger
	models data.Models
}


func main() {
	var cfg Config
	flag.StringVar(&cfg.env, "app-env", "development", "App development environment")
	flag.IntVar(&cfg.port, "app-port", 4000, "API Server Port")
	flag.StringVar(&cfg.db.user, "db-user", os.Getenv("GANT_DB_USER"), "Gant Database DSN")
	flag.StringVar(&cfg.db.pass, "db-pass", os.Getenv("GANT_DB_PASS"), "Gant Database DSN")
	flag.StringVar(&cfg.db.host, "db-host", os.Getenv("GANT_DB_HOST"), "Gant Database DSN")
	flag.StringVar(&cfg.db.port, "db-port", os.Getenv("GANT_DB_PORT"), "Gant Database DSN")
	flag.StringVar(&cfg.db.database, "db-dsn", os.Getenv("GANT_DB_NAME"), "Gant Database DSN")
	flag.StringVar(&cfg.db.ssldisable, "db-ssl-enabled", os.Getenv("GANT_DB_SSL_ENABLED "), "Gant Database SSL Mode Enabled" )

	flag.Parse()



	infoLoggger := log.New(os.Stdout, "INFO:", log.Ldate | log.Ltime)
	errorLogger := log.New(os.Stdout, "ERRROR:", log.Ldate | log.Ltime | log.Llongfile)

	db, err := openDB(cfg)

	app := application{
		config: cfg,
		infoLogger: infoLoggger,
		errorLogger: errorLogger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}


	err = srv.ListenAndServe()
	errorLogger.Fatal(err)
}

func openDB(cfg Config) (*gorm.DB, error) {
	var sslmode string
	if cfg.db.ssldisable == "true" {
		sslmode = "enable"
	} else {
	sslmode = "disable"
	}

	dsnEnv := fmt.Sprintf("host=%s user=%s dbname=%s port=%d sslmode=%s", cfg.db.host, cfg.db.user, cfg.db.port, sslmode)
	_ = dsnEnv
	dsn := "host=localhost user=postgres dbname=experiments port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&data.Keeper{}, &data.Bee{})
	if err != nil {
		fmt.Println("an error occured: ", err)
	}
	return db, nil
}
