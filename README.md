### Kalobranded: Go REST example for branded product checkout service

#### Assumptions:
* Each brand has it own `pickup_address`
* Validations for related resource are done in delivery layer.  e.g. `brand_id` in product create.
* Two roles are available for table `users`, which are "admin" and "user". Anyway, records of either roles can be created from guest http API. 
* Token-based authentication, and save auth session too

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
```sh
docker-compose.local.yml
```
or
```sh
docker-compose.dev.yml
```

### Docker development usage:
```sh
make develop
```

### Local development usages:
```sh
make local
```
or
```sh
make run
```
    
#### How to
After deployment, run 
```sh
make migrate_up
```

Register admin and buyer user from http://localhost:5001/user/create

### Swagger:

http://localhost:5001/swagger/ or http://139.162.7.112:5001/swagger/ (test)

### Test Accounts:

#### Admin
```sh
curl -X POST                                                   \
    -d '{
        	"email": "admin@gmail.com",
        	"password": "admin"
        }' \
    http://139.162.7.112:5001/swagger/index.html#/Users/post_user_login
```

#### Buyer
```sh
curl -X POST                                                   \
    -d '{
        	"email": "djourdan555@gmail.com",
        	"password": "hello"
        }' \
    http://139.162.7.112:5001/swagger/index.html#/Users/post_user_login
```

