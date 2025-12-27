### DevOps School BE+FE+DB project

DATABASE_URL - database DSN
LOG_LEVEL - info (default), warn, debug, error


#### API curl
POST /user 
`curl --header "Content-Type: application/json" --request POST --data '{"username":"xyz"}' http://localhost:8080/user`

GET /user
`curl --request GET http://localhost:8080/user`

DELETE /user/{name}
`curl --request DELETE http://localhost:8080/user/xyz`