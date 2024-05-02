package main

import (
	"github.com/Arian-p1/spb/src"
	"github.com/Arian-p1/spb/src/database"
	"github.com/Arian-p1/spb/src/objectstorage"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

  	err = objectstorage.Init()
	if err != nil {
		panic(err)
	}

	engine := src.Init()
	err = database.DatabaseConnection()
	if err != nil {
		panic(err)
	}
	err = database.Migration()
	if err != nil {
		panic(err)
	}

	engine.Run(":1234")
}
