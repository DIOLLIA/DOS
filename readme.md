### DevOps School BE+FE+DB project

#### Stack

| Backend |         Frontend          |  Database   |
|:-------:|:-------------------------:|:-----------:|
| go 1.25 | React + TypeScript + Vite | postgres 16 |

#### Environment variables

See the list of env vars that is used by BE and DB

| Environment Variable |        description        |     default value     |
|:--------------------:|:-------------------------:|:---------------------:|
|      LOG_LEVEL       | info, warn, debug, error  |         info          |
|       APP_PORT       |                           |         8080          |
|      PG_DB_URL       | backend will connect to * |       localhost       |
|      PG_DB_PORT      |                           |         5432          |
|      PG_DB_NAME      |                           |         appdb         |
|    PG_DB_USERNAME    |                           |         admin         |
|    PG_DB_PASSWORD    |                           |         test          |
|     FE_CORS_URL      |                           | http://localhost:5173 |

* Backend will connect to database using provided dsn which build as:

  postgres://${PG_DB_USERNAME}:${PG_DB_PASSWORD}@${PG_DB_URL}:${PG_DB_PORT}/${PG_DB_NAME}?sslmode=disable"

#### API curl

POST /user
`curl --header "Content-Type: application/json" --request POST --data '{"username":"xyz"}' http://localhost:8080/user`

GET /user
`curl --request GET http://localhost:8080/user`

DELETE /user/{name}
`curl --request DELETE http://localhost:8080/user/xyz`

POST /entries
`curl --header "Content-Type: application/json" --request POST --data '{"value":"white-kakao"}' http://localhost:8080/entries`

GET /entries
`curl --request GET http://localhost:8080/entries`

DELETE /entries/{id}
`curl --request DELETE http://localhost:8080/entry/11`

[//]: # (TODO)
delete .env file . App should run without it
check local launch with non existent db (without docker)
check launch with FE
checkk dockerized launch
