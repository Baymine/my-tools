#!/bin/bash

BASE_URL="http://localhost:8081"

# Test GET /todos
echo "Testing GET /todos"
curl -X GET $BASE_URL/todos

# Test POST /todos
echo -e "\n\nTesting POST /todos"
curl -X POST -H "Content-Type: application/json" -d '{"title":"New Test Todo","completed":false,"priority":"high"}' $BASE_URL/todos

# Test GET /todos again to see the new todo
echo -e "\n\nTesting GET /todos after POST"
curl -X GET $BASE_URL/todos

echo -e "\n\nAPI tests completed"
