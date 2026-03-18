package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/deep123845/blogaggregator/internal/command"
	"github.com/deep123845/blogaggregator/internal/config"
	"github.com/deep123845/blogaggregator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}

	dbQueries := database.New(db)

	s := command.State{Config: &cfg, DB: dbQueries}
	cmds := command.Commands{Command_mapping: make(map[string]func(*command.State, command.Command) error)}
	cmds.Register("login", command.HandlerLogin)
	cmds.Register("register", command.HandlerRegister)
	cmds.Register("reset", command.HandlerReset)
	cmds.Register("users", command.HandlerUsers)
	cmds.Register("agg", command.HandlerAgg)
	cmds.Register("addfeed", command.MiddlewareLoggedIn(command.HandlerAddFeed))
	cmds.Register("feeds", command.HandlerFeeds)
	cmds.Register("follow", command.MiddlewareLoggedIn(command.HandlerFollow))
	cmds.Register("following", command.MiddlewareLoggedIn(command.HandlerFollowing))

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
