#!/bin/bash

BASE_URL="http://localhost:8081"

# 测试用户注册
echo "Testing user registration"
REGISTER_RESPONSE=$(curl -s -X POST -H "Content-Type: application/json" -d '{"username":"testuser","email":"test@example.com","password":"password123"}' $BASE_URL/register)
echo $REGISTER_RESPONSE

# 从注册响应中提取令牌
TOKEN=$(echo $REGISTER_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo -e "\nToken: $TOKEN"

# 如果注册失败，尝试登录
if [ -z "$TOKEN" ]; then
  echo -e "\nRegistration failed, trying login"
  LOGIN_RESPONSE=$(curl -s -X POST -H "Content-Type: application/json" -d '{"email":"test@example.com","password":"password123"}' $BASE_URL/login)
  echo $LOGIN_RESPONSE
  TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
  echo -e "\nToken from login: $TOKEN"
fi

# 使用令牌测试获取待办事项
echo -e "\n\nTesting GET /todos with token"
curl -X GET -H "Authorization: Bearer $TOKEN" $BASE_URL/todos

# 使用令牌测试创建待办事项
echo -e "\n\nTesting POST /todos with token"
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{"title":"New Test Todo","completed":false,"priority":"high"}' $BASE_URL/todos

# 获取新创建的待办事项ID
echo -e "\n\nGetting todo ID for further tests"
TODO_ID=$(curl -s -X GET -H "Authorization: Bearer $TOKEN" $BASE_URL/todos | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
echo "Using todo ID: $TODO_ID"

# 使用令牌测试更新待办事项
echo -e "\n\nTesting PUT /todos/$TODO_ID with token"
curl -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{"title":"Updated Test Todo","completed":true,"priority":"medium"}' $BASE_URL/todos/$TODO_ID

# 使用令牌再次测试获取待办事项
echo -e "\n\nTesting GET /todos after PUT with token"
curl -X GET -H "Authorization: Bearer $TOKEN" $BASE_URL/todos

# 使用令牌测试删除待办事项
echo -e "\n\nTesting DELETE /todos/$TODO_ID with token"
curl -X DELETE -H "Authorization: Bearer $TOKEN" $BASE_URL/todos/$TODO_ID

# 使用令牌再次测试获取待办事项以确认删除
echo -e "\n\nTesting GET /todos after DELETE with token"
curl -X GET -H "Authorization: Bearer $TOKEN" $BASE_URL/todos

echo -e "\n\nAPI tests completed"


