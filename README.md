go-todo-webserver
=================

Small example WebServer created in Go. 
Connects to a remote MYSQL server. 
Please keep in mind that this might be a security hole. Only allow connection form specific ip's in your mysql settings.

CRUD Operation for a TODO List, which actually just has a Text, no "state"

curl -i -d '{"Test":"Do This"}' http://127.0.0.1:8080/todos

curl -i http://127.0.0.1:8080/todos/1

curl -i http://127.0.0.1:8080/todos

curl -i -X PUT -d '{"Text":"is updated"}' http://127.0.0.1:8080/todos/1

curl -i -X DELETE http://127.0.0.1:8080/todos/1
