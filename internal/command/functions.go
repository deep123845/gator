package command

import (
	"context"
	"fmt"
	"strconv"
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

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Aggregate command expects 1 argument")
	}

	time_between_reqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Add Feed Command requires 2 arguments: Name, URL")
	}

	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
	}
	feed, err := s.DB.CreateFeed(context.Background(), newFeed)
	if err != nil {
		return err
	}

	HandlerFollow(s, Command{Name: "", Args: []string{feed.Url}}, user)

	fmt.Printf("%+v\n", feed)

	return nil
}

func HandlerFeeds(s *State, _ Command) error {
	feeds, err := s.DB.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		user, err := s.DB.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return err
		}

		fmt.Printf("Name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("User: %s\n\n", user.Name)
	}

	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Follow command expects only one argument")
	}

	feed, err := s.DB.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	new_feed_follow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	}
	feed_follow, err := s.DB.CreateFeedFollow(context.Background(), new_feed_follow)
	if err != nil {
		return err
	}

	fmt.Printf("User: %s followed feed: %s\n", feed_follow.UserName, feed_follow.FeedName)

	return nil
}

func HandlerFollowing(s *State, _ Command, user database.User) error {
	feed_follows, err := s.DB.GetFeedFollowsByUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("Feeds for user: %s\n", user.Name)
	for _, feed_follow := range feed_follows {
		fmt.Printf("- %s\n", feed_follow.FeedName)
	}

	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Unfollow command expects one argument")
	}

	feed, err := s.DB.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	feed_follow_to_delete := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.DB.DeleteFeedFollow(context.Background(), feed_follow_to_delete)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted feed: %s from user: %s's follow list", feed.Name, user.Name)

	return nil
}

func HandlerBrowse(s *State, cmd Command, user database.User) error {
	limit := 2
	if len(cmd.Args) > 0 {
		var err error
		limit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return err
		}
	}

	posts, err := s.DB.GetPostsForUser(context.Background(), database.GetPostsForUserParams{UserID: user.ID, Limit: int32(limit)})
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("%s\n\n", post.Description)
	}

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

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.DB.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}

const url_error = "pq: duplicate key value violates unique constraint \"posts_url_key\" (23505)"

func scrapeFeeds(s *State) error {
	next_feed, err := s.DB.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	err = s.DB.MarkFeedFetched(context.Background(), next_feed.ID)
	if err != nil {
		return err
	}

	url := next_feed.Url

	feed_content, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		return err
	}

	for _, item := range feed_content.Channel.Item {
		pub_time, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			return err
		}

		new_post := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: pub_time,
			FeedID:      next_feed.ID,
		}
		post, err := s.DB.CreatePost(context.Background(), new_post)
		if err != nil {
			if err.Error() != url_error {
				return err
			}
		} else {
			fmt.Printf("%+v\n", post)
		}
	}

	return nil
}
