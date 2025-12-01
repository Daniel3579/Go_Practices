curl -s -X POST http://localhost:8080/api/v1/login \
 -H "Content-Type: application/json" \
 -d '{"Email":"admin@example.com","Password":"secret123"}'
# => {"token":"<JWT>"}

TOKEN=<скопируйте токен>
curl -s http://localhost:8080/api/v1/me -H "Authorization: Bearer $TOKEN"
curl -s http://localhost:8080/api/v1/admin/stats -H "Authorization: Bearer $TOKEN"


TOKEN_USER=$(curl -s -X POST http://localhost:8080/api/v1/login \
 -H "Content-Type: application/json" -d '{"Email":"user@example.com","Password":"secret123"}' | jq -r .token)
curl -i http://localhost:8080/api/v1/admin/stats -H "Authorization: Bearer $TOKEN_USER"  # ожидаем 403
