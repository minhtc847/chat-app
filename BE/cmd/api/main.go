package main

import (
	"BE-chat-app/internal/data"
	"BE-chat-app/internal/data/mailer"
	"BE-chat-app/internal/jsonlog"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	cors struct {
		trustedOrigins []string
	}
	jwt struct {
		secret string
	}
	redis struct {
		addr     string
		password string
		db       int
	}
}
type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	redis  *RedisLocal
}

func main() {

	if os.Getenv("RENDER") == "" { // Check if not running on Render
		if err := godotenv.Load("./.env"); err != nil {
			fmt.Println("No .env file found, using environment variables")
		}
	}
	dbDsn := os.Getenv("DB_DSN")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		fmt.Printf("Invalid PORT value: %s\n", smtpPort)
		smtpPort = 4000 // Use a default value if conversion fails
	}
	smtpUsername := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASS")
	smtpSender := os.Getenv("SMTP_SENDER")

	fmt.Printf("SMTP_PORT: %d\n", smtpPort)
	fmt.Printf("SMTP_USERNAME: %s\n", smtpUsername)

	corsTrustOrigin := os.Getenv("CORS_TRUST_ORIGIN")
	jwtSecret := os.Getenv("JWT_SECRET")
	var cfg config
	//go run ./cmd/api -help for more information
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", dbDsn, "PostgreSQL DSN")
	fmt.Printf("DB DSN: %s\n", cfg.db.dsn)
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.StringVar(&cfg.smtp.host, "smtp-host", smtpHost, "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", smtpPort, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", smtpUsername, "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", smtpPassword, "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", smtpSender, "SMTP sender")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		if len(val) == 0 {
			cfg.cors.trustedOrigins = strings.Fields(corsTrustOrigin)
		}
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	flag.StringVar(&cfg.jwt.secret, "jwt-secret", jwtSecret, "JWT secret")

	flag.StringVar(&cfg.redis.addr, "redis-addr", os.Getenv("REDIS_ADDR"), "Redis address")
	flag.StringVar(&cfg.redis.password, "redis-password", os.Getenv("REDIS_PASSWORD"), "Redis password")
	flag.IntVar(&cfg.redis.db, "redis-db", -1, "Redis database")
	if cfg.redis.db == -1 {
		cfg.redis.db, err = strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			fmt.Printf("Invalid DB value: %s\n", cfg.redis.db)
			cfg.redis.db = 0
		}
	}
	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	//redisClient, err := openRedis(cfg)
	//if err != nil {
	//	logger.PrintFatal(err, nil)
	//}
	//defer redisClient.Close()
	//logger.PrintInfo("Redis connection established", nil)

	redisLocal := NewRedisLocal()

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(
			cfg.smtp.host,
			cfg.smtp.port,
			cfg.smtp.username,
			cfg.smtp.password,
			cfg.smtp.sender),
		redis: redisLocal,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  cfg.env,
	})
	err = srv.ListenAndServe()
	logger.PrintFatal(err, nil)

}
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func openRedis(cfg config) (*redis.Client, error) {
	//opt := &redis.Options{
	//	Addr:        cfg.redis.addr,
	//	Password:    cfg.redis.password, // Leave empty for no password
	//	DB:          cfg.redis.db,       // Default DB
	//	DialTimeout: 10 * time.Second,
	//	Username:    "default",
	//}

	opt := &redis.Options{
		Addr:     "redis-15133.c252.ap-southeast-1-1.ec2.redns.redis-cloud.com:15133",
		Username: "default",
		Password: "1JMxgARJe3EDrZlRo7NDQBFsYN1C4JOT",
		DB:       0,
	}
	client := redis.NewClient(opt)

	fmt.Printf("Redis connection established at %s", opt.Addr)
	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
