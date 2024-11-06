
## Tech Stack

**Go:** I use go version ^1.22

**Library:** I use OAuth2.0 Protocol and Auth0 for my Authentication Server

**Gin-Gonic:** Because Gin-Gonic is very popular, so there's no reason to not to use this powerful framework

I've implement Clean Architecture design and Factory design pattern. So, you could change my existing Library or framework seamlessly.

## Installation

Install this project using go

```bash
  cp .config.example.json config.json
  go mod download && go mod tidy
```

to run this project

```bash
go run main.go
```

To run, watch, and build this project

```bash
CompileDaemon -command="./go-sso"
```

To migrate the database

```bash
go run ./cmd/migration/main.go
```

Make sure you fill the required credentials
## Demo

You could test the SSO using two methods. The first one use JWT and make this repo as Authentication Server. The second one use Auth0 as Third-party Authentication Server.

## Using normal JWT

```bash
/api/user/login
```

## Using Auth0
Make sure to open it on web browser because it will redirect you to Auth0 login page. And before using this, make sure you add some users on Auth0 platform and add those users to your own database.
```bash
/oauth/login
```

