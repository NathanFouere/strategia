package main

import (
	"server/internal"
	"server/internal/container"
)

func main() {
	err := container.SetupContainer()
	if err != nil {
		panic("Error during container setup")
	}

	err = container.GetContainer().Invoke(func(mainHandler *internal.MainHandler) {
		mainHandler.Launch()
	})
}
