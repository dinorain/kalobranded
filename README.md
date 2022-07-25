### Checkoutaja: Go REST example for product checkout service

#### What have been used:
* [sqlx](https://github.com/jmoiron/sqlx) - Extensions to database/sql.
* [pgx](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit for Go
* [JWT](https://github.com/golang-jwt/jwt) - A Go implementation of JSON Web Tokens.
* [viper](https://github.com/spf13/viper) - A Go configuration with fangs
* [go-redis](https://github.com/go-redis/redis) - Redis client for Golang
* [zap](https://github.com/uber-go/zap) - Logger
* [validator](https://github.com/go-playground/validator) - Go Struct and Field validation
* [migrate](https://github.com/golang-migrate/migrate) - Database migrations. CLI and Golang library.
* [testify](https://github.com/stretchr/testify) - Testing toolkit
* [gomock](https://github.com/golang/mock) - Mocking framework
* [CompileDaemon](https://github.com/githubnemo/CompileDaemon) - Compile daemon for Go
* [Docker](https://www.docker.com/) - Docker

#### Docker compose files:
    docker-compose.local.yml
    docker-compose.dev.yml

### Docker development usage:
    make develop

### Local development usage:
    make local
    make run

### Swagger:

http://localhost:5001/swagger/

### Test Accounts:

#### Admin
```sh
curl -X POST                                                   \
    -d '{
        	"email": "admin@gmail.com",
        	"password": "admin"
        }' \
    http://139.162.55.156:5001/swagger/index.html#/Users/post_user_login
```

#### Seller
```sh
curl -X POST                                                   \
    -d '{
        	"email": "seller@gmail.com",
        	"password": "seller"
        }' \
    http://139.162.55.156:5001/swagger/index.html#/Sellers/post_user_login
```

#### Buyer
```sh
curl -X POST                                                   \
    -d '{
        	"email": "djourdan555@gmail.com",
        	"password": "hello"
        }' \
    http://139.162.55.156:5001/swagger/index.html#/Users/post_user_login
```

