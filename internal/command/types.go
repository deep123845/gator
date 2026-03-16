package command

import (
	"github.com/deep123845/blogaggregator/internal/config"
)

type State struct {
	Config *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Command_mapping map[string]func(*State, Command) error
}
