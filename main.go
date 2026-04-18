package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

type Configs struct {
	Port           string
	PermGrpcServer string
	UserGrpcServer string
	// PartnerGrpcServer string
	VoucherGrpcServer string
	ProductGrpcServer string
	JwtSecretKey      string
	RedisAddr         string
	RedisPassword     string
	RedisDb           string
	CloudinaryName    string
	CloudinaryApiKey  string
	CloudinarySecret  string
	AdminRole         string
}

var config *Configs

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: Error loading .env file:", err)
		} else {
			log.Println("Loaded .env file for local development")
		}
	} else {
		log.Println("No .env file found, using system environment variables")
	}

	config = &Configs{
		Port:           getEnv("PORT", "8080"),
		PermGrpcServer: getEnv("PERM_GRPC_SERVER", "localhost:7000"),
		UserGrpcServer: getEnv("USER_GRPC_SERVER", "localhost:6000"),
		// PartnerGrpcServer: getEnv("PARTNER_GRPC_SERVER", "localhost:5000"),
		VoucherGrpcServer: getEnv("VOUCHER_GRPC_SERVER", "localhost:4000"),
		ProductGrpcServer: getEnv("PRODUCT_GRPC_SERVER", "localhost:8000"),
		JwtSecretKey:      getEnv("JWT_SECRET_KEY", ""),
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		RedisDb:           getEnv("REDIS_DB", "0"),
		CloudinaryName:    getEnv("CLOUDINARY_NAME", ""),
		CloudinaryApiKey:  getEnv("CLOUDINARY_API_KEY", ""),
		CloudinarySecret:  getEnv("CLOUDINARY_SECRET", ""),
		AdminRole:         getEnv("ADMIN_ROLE", "admin"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func appRoot() error {
	app := cli.NewApp()

	app.Action = func(c *cli.Context) error {
		return errors.New("Wow, ^.^ dumb")
	}
	app.Commands = []*cli.Command{
		{Name: "start", Action: func(ctx *cli.Context) error {
			NewRouter(config)
			return nil
		}},
	}
	return app.Run(os.Args)
}

func main() {
	runtime.GOMAXPROCS(2)
	go freeMemory()
	if err := appRoot(); err != nil {
		panic(err)
	}
}

func freeMemory() {
	for {
		fmt.Println("run gc")
		start := time.Now()
		runtime.GC()
		debug.FreeOSMemory()
		elapsed := time.Since(start)
		fmt.Printf("gc took %s\n", elapsed)
		time.Sleep(2 * time.Minute)
	}
}
