package command

import (
	"context"
	"fmt"
	"time"

	"github.com/deep123845/blogaggregator/internal/database"
	"github.com/deep123845/blogaggregator/internal/rss"
	"github.com/google/uuid"
)

func HandlerLogin(s *State, cmd Command) error {
	if num_args := len(cmd.Args); num_args != 1 {
		return fmt.Errorf("Login command expects one argument received %v", num_args)
	}

	_, err := s.DB.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Login Failed, %v", err)
	}

	err = s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Login failed, %v", err)
	}

	fmt.Printf("Login successful, user set to %v\n", cmd.Args[0])
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if num_args := len(cmd.Args); num_args != 1 {
		return fmt.Errorf("Register command expects one argument received %v", num_args)
	}

	new_user := database.CreateUserParams{ID: uuid.New(), Name: cmd.Args[0], CreatedAt: time.Now(), UpdatedAt: time.Now()}
	user, err := s.DB.CreateUser(context.Background(), new_user)
	if err != nil {
		return fmt.Errorf("Failed to register user, %v", err)
	}

	err = HandlerLogin(s, Command{Name: "", Args: []string{user.Name}})
	if err != nil {
		return err
	}

	fmt.Printf("User Created with information %v\n", user)
	return nil
}

func HandlerReset(s *State, _ Command) error {
	return s.DB.Reset(context.Background())
}

func HandlerUsers(s *State, _ Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Could not retreive users: %v", err)
	}

	for _, user := range users {
		if user.Name == s.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func HandleAgg(_ *State, _ Command) error {
	url := "https://www.wagslane.dev/index.xml"

	feed, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", feed)

	return nil
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.Command_mapping[cmd.Name]
	if !ok {
		return fmt.Errorf("Command: %v not found", cmd.Name)
	}

	return handler(s, cmd)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Command_mapping[name] = f
}
