package main

import (
	"log"

	"bpf.com/api"
	"bpf.com/pkg/core"
)

func main() {
	if err := core.InitConfig("./config.yaml"); err != nil {
		log.Fatalf("Init config.yaml fail: %v", err)
	}

	if err := core.InitLogger(); err != nil {
		log.Fatalf("Init logger fail: %v", err)
	}

	if err := core.InitDatabase(); err != nil {
		log.Fatalf("Init database fail: %v", err)
	}

	if err := core.InitCache(); err != nil {
		log.Fatalf("Init cache fail: %v", err)
	}

	router := core.InitGin()
	api.SetupRoutes(router)

	app := core.NewApplication(router)
	app.Run()
}
