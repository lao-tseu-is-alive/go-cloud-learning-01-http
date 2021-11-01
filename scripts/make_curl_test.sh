#!/bin/bash
curl -H "Content-Type: application/json" 'http://localhost:8080/todos' |json_pp
curl -XPOST -H "Content-Type: application/json" -d '{"task":"learn Linux"}'  'http://localhost:8080/todos'
curl -H "Content-Type: application/json" 'http://localhost:8080/todos' |json_pp
curl -H "Content-Type: application/json" 'http://localhost:8080/todos/1' |json_pp
curl -H "Content-Type: application/json" 'http://localhost:8080/todos/3' |json_pp
curl -v -XPUT -H "Content-Type: application/json" -d '{"id": 3, "task":"learn Linux", "completed": true}' 'http://localhost:8080/todos/3'
curl -H "Content-Type: application/json" 'http://localhost:8080/todos/3' |json_pp
curl -v -XPUT -H "Content-Type: application/json" -d '{"id": 3, "task":"learn Linux", "completed": false}' 'http://localhost:8080/todos/3'
curl -H "Content-Type: application/json" 'http://localhost:8080/todos/3' |json_pp
curl -H "Content-Type: application/json" 'http://localhost:8080/todos/maxid'
curl -v -XDELETE -H "Content-Type: application/json" 'http://localhost:8080/todos/3'
curl -H "Content-Type: application/json" 'http://localhost:8080/todos/maxid'
curl -H "Content-Type: application/json" 'http://localhost:8080/todos' |json_pp
