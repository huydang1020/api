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
	JwtSecretKey   string
	RedisAddr      string
	RedisPassword  string
	RedisDb        string
}

var config *Configs

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env:", err)
	}
	config = &Configs{
		Port:           os.Getenv("PORT"),
		PermGrpcServer: os.Getenv("PERM_GRPC_SERVER"),
		UserGrpcServer: os.Getenv("USER_GRPC_SERVER"),
		JwtSecretKey:   os.Getenv("JWT_SECRET_KEY"),
		RedisAddr:      os.Getenv("REDIS_ADDR"),
		RedisPassword:  os.Getenv("REDIS_PASSWORD"),
		RedisDb:        os.Getenv("REDIS_DB"),
	}
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
