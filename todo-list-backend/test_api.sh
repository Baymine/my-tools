#!/bin/bash

BASE_URL="http://localhost:8081"

# Test GET /todos
echo "Testing GET /todos"
curl -X GET $BASE_URL/todos

# Test POST /todos
echo -e "\n\nTesting POST /todos"
curl -X POST -H "Content-Type: application/json" -d '{"title":"Test Todo","completed":false,"priority":"medium"}' $BASE_URL/todos

# Test GET /todos again to see the new todo
echo -e "\n\nTesting GET /todos after POST"
curl -X GET $BASE_URL/todos

# Get the ID of the created todo
ID=$(curl -s -X GET $BASE_URL/todos | grep -o '"id":[0-9]*' | tail -1 | grep -o '[0-9]*')

# Test PUT /todos/{id}
echo -e "\n\nTesting PUT /todos/$ID"
curl -X PUT -H "Content-Type: application/json" -d '{"title":"Updated Test Todo","completed":true,"priority":"high"}' $BASE_URL/todos/$ID

# Test GET /todos again to see the updated todo
echo -e "\n\nTesting GET /todos after PUT"
curl -X GET $BASE_URL/todos

# Test DELETE /todos/{id}
echo -e "\n\nTesting DELETE /todos/$ID"
curl -X DELETE $BASE_URL/todos/$ID

# Test GET /todos one last time to confirm deletion
echo -e "\n\nTesting GET /todos after DELETE"
curl -X GET $BASE_URL/todos

echo -e "\n\nAPI tests completed"
