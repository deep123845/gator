package command

import (
	"github.com/deep123845/blog-aggregator/internal/config"
	"github.com/deep123845/blog-aggregator/internal/database"
)

type State struct {
	Config *config.Config
	DB     *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Command_mapping map[string]func(*State, Command) error
}
