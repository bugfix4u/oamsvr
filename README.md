# **Online Account Manager API Server**

## Description
This is a simple example project to show how to build a basic API server backed by an SQL database, to manage online accounts. This includes information about the account, like name, URL, username, and password, as well as security validation questions that many online account require users to answer. 

## Resources
The project uses the following software and components:
- Golang - The Go programming language ( https://golang.org/ )
- GORM - The fantastic ORM library for Golang ( https://gorm.io/ )
- Gorilla Mux - Gorilla Web Toolkit ( https://www.gorillatoolkit.org/pkg/mux )
- JWT-GO - A Golang implementation of JSON Web Tokens ( https://github.com/dgrijalva/jwt-go )
- Golang Crypto (bcrypt) Package - BCrypt library for Go ( https://godoc.org/golang.org/x/crypto/bcrypt )
- Google UUID for Golang - Go package for UUIDs based on RFC 4122 and DCE 1.1 ( https://github.com/google/uuid )
- Postgres DB - Open source relational database ( https://www.postgresql.org/ )
- Docker - Docker container platform ( https://www.docker.com/ )
- Docker Compose - Docker tool for running multi-container Docker applications ( https://docs.docker.com/compose/ )
- Swagger - API documentation ( https://swagger.io/ )

## Building the Code
This code was orginally built using Go 1.13.4 but older versions will most likely work as well. If using an older version make use that you have enabled Go Modules (https://github.com/golang/go/wiki/Modules). You should be able to simply run the "go build" comand to build the "oamsvr" binary.

## Packaging
The Online Account Manager server is designed to be packaged as a set of microservices, each running as a seperate Docker container. To package and run this example, you will need to have a working Docker installation, as well as the Docker Compose tool installed.

## Running the server
To run the Online Account Managers server and the Postgres DB, simply run "sudo docker-compose up". You will need internet access to run this command, since it will be pulling down a Postgres docker image and a ubuntu image for packaging of the "oamsvr" binary. The Postgres DB is accessible from the standard 5432 port, while the Online Account Manager server is accessible through HTTPS port 8443. The Online Account Manager server uses self signed certificate so make sure any ReST calls, so make ignore self signed certificate errors. The default user for the Online Account Manager is "admin/H1r3M3N0W".

## Examples
Below are a few examples to get you going with the Online Account Manager servers. Please see the swagger documentation in the ./api directory for more info.

To get a authorization token, you must first authenticate with the server.
```bash
curl --request POST \
  --url https://localhost:8443/api/v1/authenticate \
  --header 'content-type: application/json' \
  --data '{  
	"username":"admin",
	"password":"H1r3M3N0W"
}'
```

Once authenticated, you can then make calls to the Account annd User APIs.
```bash
curl --request POST \
  --url https://localhost:8443/api/v1/accounts \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXNzd29yZCI6IiQyYSQxMiR5QkVDeUhyeC9KQTlxUks5bWVZNmouaGNod2dab0N2Y3puWXVRNjBRVDRpSGdKZDRqM0lmYSIsInVzZXJuYW1lIjoiYWRtaW4ifQ.Tce8qaOBVQpC2rLTgq0KxXB2srtiveehmhF9UEbkons' \
  --header 'content-type: application/json' \
  --data '{
	"name":"My Bank",
	"url":"www.mybank.com",
	"username":"bob@gmail.com",
	"password":"foobar",
	"questions":[{
		"name":"Question 1",
		"question":"What is you fathers middle name?",
		"answer":"James"
	}]
}'
```


## API Documentation
The API's for this example server have been documented using Swagger (https://swagger.io/). You can run the swagger-ui as a docker container using the following instructions:

This will start nginx with Swagger UI on port 80.

```bash
docker pull swaggerapi/swagger-ui
docker run -p 80:8080 swaggerapi/swagger-ui
```

To start swagger with the Online Account Manager swagger documentation, you can use a command similar to:

```bash
docker run -p 80:8080 -e SWAGGER_JSON=/swagger/swagger.yaml -v /path/to/oamsvr/api:/swagger swaggerapi/swagger-ui
```

## License
GPL v3 ( http://www.gnu.org/licenses/gpl-3.0.html )