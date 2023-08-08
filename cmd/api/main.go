package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jplindgren/rpg-vault/internal/clients"
	"github.com/jplindgren/rpg-vault/internal/jsonlog"
	"github.com/jplindgren/rpg-vault/internal/services"
)

type config struct {
	port    int
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	cors struct {
		trustedOrigins []string
	}
	aws struct {
		key    string
		secret string
		region string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	//models   adapters.Models
	services services.Services
}

const port = 4000

func main() {
	fmt.Println(os.Getenv("AWS_ACCESS_KEY_ID"))
	fmt.Println(os.Getenv("AWS_SECRET_ACCESS_KEY"))
	fmt.Println(os.Getenv("AWS_DEFAULT_REGION"))

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")

	// Create command line flags to read the setting values into the config struct.
	// Notice that we use true as the default for the 'enabled' setting?
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.aws.key, "aws-key", os.Getenv("AWS_ACCESS_KEY_ID"), "Aws key")
	flag.StringVar(&cfg.aws.secret, "aws-secret", os.Getenv("AWS_SECRET_ACCESS_KEY"), "Aws secret")
	flag.StringVar(&cfg.aws.region, "aws-region", os.Getenv("AWS_DEFAULT_REGION"), "Aws region")

	// Use the flag.Func() function to process the -cors-trusted-origins command line
	// flag. In this we use the strings.Fields() function to split the flag value into a
	// slice based on whitespace characters and assign it to our config struct.
	// Importantly, if the -cors-trusted-origins flag is not present, contains the empty
	// string, or contains only whitespace, then strings.Fields() will return an empty
	// []string slice.
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	flag.Parse()

	// Initialize a new logger which writes messages to the standard out stream,
	// prefixed with the current date and time.
	//logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	app := &application{
		logger: logger,
		config: cfg,
		services: services.NewServices(
			clients.GetDynamodbClient(cfg.aws.key, cfg.aws.secret, cfg.aws.region),
			clients.GetS3Client(cfg.aws.key, cfg.aws.secret, cfg.aws.region),
		),
	}

	router := app.routes()
	composerHandler := app.recoverPanic(app.enabledCORS(app.rateLimit(app.authenticate(router))))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      composerHandler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.PrintInfo("starting server on %s", map[string]string{
		"port": strconv.Itoa(port),
	})
	err := srv.ListenAndServe()
	logger.PrintFatal(err, nil)
}
