package main

import (
	"github.com/EslRain/leaf-segment/wire"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load(".conf.yaml")
}

func main() {
	h := wire.InitHandler()

	h.Run()
}
