package command

import "fmt"

func HandlerLogin(s *State, cmd Command) error {
	if num_args := len(cmd.Args); num_args != 1 {
		return fmt.Errorf("Login command expects one argument received %v", num_args)
	}

	err := s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Login failed, %v", err)
	}

	fmt.Printf("Login successful, user set to %v\n", cmd.Args[0])
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
