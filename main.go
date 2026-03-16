package main

import (
	"log"
	"os"

	"github.com/deep123845/blogaggregator/internal/command"
	"github.com/deep123845/blogaggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	s := command.State{Config: &cfg}
	cmds := command.Commands{Command_mapping: make(map[string]func(*command.State, command.Command) error)}
	cmds.Register("login", command.HandlerLogin)

	args := os.Args

	if len(args) < 2 {
		log.Fatalf("expected at least 2 arguments, exiting")
	}

	cmd := command.Command{Name: args[1], Args: args[2:]}
	err = cmds.Run(&s, cmd)
	if err != nil {
		log.Fatalf("error running command: %v", err)
	}
}
