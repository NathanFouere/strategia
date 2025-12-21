package main

import (
	"fmt"
	"server/internal"
	"server/internal/container"
)

func main() {
	err := container.SetupContainer()
	if err != nil {
		fmt.Println("error :", err)
	}

	err = container.GetContainer().Invoke(func(mainHandler *internal.MainHandler) {
		mainHandler.Launch()
	})
}
