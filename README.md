# Gator
Gator is a cli tool which aggregates and displays blog posts from rss feeds.

# Installation
## Requirements
Go 1.26.0
Local Postgres database
## Instructions
Run the command:

```bash
go install https://github.com/deep123845/gator@latest
```

Create a `.gatorconfig.json` file in your home directory with the following format

```json
{
  "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
}
```

# Usage
```bash
gator {command} {args}
```

## Commands
```bash
gator reset
```
Clears the DB of all users, feeds, and posts
___
```bash
gator register {name}
```
Creates a new user with the given name and logs in to it
___
```bash
gator login {name}
```
Logs in to the given user
___
```bash
gator users
```
Displays a list of all users in the DB
___
```bash
gator addfeed {url}
```
Adds the feed found at the given url
___
```bash
gator agg {seconds}s
```
Starts the aggregator waiting the amount of given seconds betweens blogs
___
```bash
gator browse {limit}
```
View the posts that have been aggregated
___


# References
Built while going through a bootdev course
