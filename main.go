package main

import (
	"fmt"
	"log"

	"github.com/deep123845/blogaggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Config: %+v\n", cfg)

	err = cfg.SetUser("deep")
	if err != nil {
		log.Fatalf("error writing config: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Config: %+v\n", cfg)
}
