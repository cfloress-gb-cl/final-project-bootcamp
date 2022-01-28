<div align="center">
<h1>final-project-bootcamp Globant 2022</h1>
</div>

## About

This final project of Gobant's Golang Bootcamp implements a REST API solution with [Go Kit](https://gokit.io/) - [gRPC](https://grpc.io/docs/languages/go/basics/).

This REST API has a USER entity to create, list, get, update and delete a user through CURLs calls or using Postman if you will.

The implemented Databases as repository are Mysql and MongoDB, to pick one just must set the env file in te project root.

## Docs

You can read more about final project definition at:

- https://github.com/josnelihurt/golang-bootcamp/wiki/Final-Project

Here a diagram to see the hole picture of the project.

<div align="center">
<img src="images/general_arch.png" />
</div>

There is a [JWT](https://jwt.io/) integration to validate the request token in the call to the endpoints. To do this the project uses a [COBRA](https://github.com/spf13/cobra) CLI app to validate and to get the token.

<div align="center">
<img src="images/Cobra-jwt_diagram.png" />
</div>

## Usage
To execute the calls here some examples to do it through CURL in the Linux command line:
#### Get Token
```sh
$ tkn-gb getTkn claudio.flores@globant.com
Token:  eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImNsYXVkaW8uZmxvcmVzQGdsb2JhbnQuY29tIiwiZXhwIjoxNjQzMzg1Mzg5LCJpYXQiOjE2NDMzODE3NTksIm5iZiI6MTY0MzM4MTc1OX0.RdG4g8Jy1YCDiCkYlTNfdyB9yuoQ0vgIWOt4VBT45aU
```
#### Create User
```sh
$ curl -H "Content-Type: application/json" -X POST 'http://localhost:8080/user/' -d'{"name":"Claudio Andres","lastname":"Flores Sapiain","email":"claudio.flores@globant.com","age":"39","status":"1"}'  -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImNsYXVkaW8uZmxvcmVzQGdsb2JhbnQuY29tIiwiZXhwIjoxNjQzMzg1Mzg5LCJpYXQiOjE2NDMzODE3NTksIm5iZiI6MTY0MzM4MTc1OX0.RdG4g8Jy1YCDiCkYlTNfdyB9yuoQ0vgIWOt4VBT45aU'
{"user":"claudio.flores@globant.com"}
```
#### List Users
```sh
$ curl -H "Content-Type: application/json" -X GET 'http://localhost:8080/user/' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImNsYXVkaW8uZmxvcmVzQGdsb2JhbnQuY29tIiwiZXhwIjoxNjQzMzg1Mzg5LCJpYXQiOjE2NDMzODE3NTksIm5iZiI6MTY0MzM4MTc1OX0.RdG4g8Jy1YCDiCkYlTNfdyB9yuoQ0vgIWOt4VBT45aU'
{"users":[{"id":1,"email":"claudio.flores@globant.cl","name":"Claudio Andres","lastname":"Flores Sapiain","age":"39","status":"1"},{"id":1,"email":"claudio.flores@globant.cll","name":"ca","lastname":"fs","age":"35","status":"0"},{"id":1,"email":"claudio.flores@globant.com","name":"Claudio Andres","lastname":"Flores Sapiain","age":"39","status":"1"}]}
```
#### Get User
```sh
$ curl -H "Content-Type: application/json" -X GET 'http://localhost:8080/user/claudio.flores@globant.com' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImNsYXVkaW8uZmxvcmVzQGdsb2JhbnQuY29tIiwiZXhwIjoxNjQzMzg1Mzg5LCJpYXQiOjE2NDMzODE3NTksIm5iZiI6MTY0MzM4MTc1OX0.RdG4g8Jy1YCDiCkYlTNfdyB9yuoQ0vgIWOt4VBT45aU'
{"user":{"id":1,"email":"claudio.flores@globant.com","name":"Claudio Andres","lastname":"Flores Sapiain","age":"39","status":"1"}}
```
#### Update User
```sh
$ curl -H "Content-Type: application/json" -X PUT 'http://localhost:8080/user/claudio.flores@globant.com' -d'{"name":"Claudio","lastname":"Flores","age":"39","status":"1"}'  -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImNsYXVkaW8uZmxvcmVzQGdsb2JhbnQuY29tIiwiZXhwIjoxNjQzMzg1Mzg5LCJpYXQiOjE2NDMzODE3NTksIm5iZiI6MTY0MzM4MTc1OX0.RdG4g8Jy1YCDiCkYlTNfdyB9yuoQ0vgIWOt4VBT45aU' 
{"user":"claudio.flores@globant.com"}
```
#### Delete User
```sh
$ curl -H "Content-Type: application/json" -X DELETE 'http://localhost:8080/user/1' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImNsYXVkaW8uZmxvcmVzQGdsb2JhbnQuY29tIiwiZXhwIjoxNjQzMzg1Mzg5LCJpYXQiOjE2NDMzODE3NTksIm5iZiI6MTY0MzM4MTc1OX0.RdG4g8Jy1YCDiCkYlTNfdyB9yuoQ0vgIWOt4VBT45aU'
{"user":"1"}
```

## Next Steps
- Generating Docker compose yaml file.
- Makefile to testing, build and deploy app.
- Implement gRPC server stream to getAll method.
- gRPC Gateway implementation.





