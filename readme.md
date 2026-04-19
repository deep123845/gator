# Gator
Gator is a cli tool which aggregates and displays blog posts from rss feeds.

# Installation
## Requirements
Go 1.26.0
## Instructions
Run the command:
> go install https://github.com/deep123845/gator@latest:w

# Usage
> gator {command} {args}

## Commands
> gator reset
Clears the DB of all users, feeds, and posts

> gator register {name}
Creates a new user with the given name and logs in to it

> gator login {name}
Logs in to the given user

> gator users
Displays a list of all users in the DB

# References
Built while going through a bootdev course