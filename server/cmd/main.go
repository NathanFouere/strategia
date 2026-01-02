package main

import (
	"server/internal"
	"server/internal/container"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Couldn't load .env")
	}

	err = container.SetupContainer()
	if err != nil {
		panic("Error during container setup")
	}

	err = container.GetContainer().Invoke(func(mainHandler *internal.MainHandler) {
		mainHandler.Launch()
	})
	if err != nil {
		panic("Error while launching game")
	}
}
