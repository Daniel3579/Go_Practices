```curl -i localhost:8080/health```

```
curl -i -X POST localhost:8080/tasks \
-H "Content-Type: application/json" \
-d '{"title":"Купить молоко"}'
```

```curl -i localhost:8080/tasks```

```curl -i "localhost:8080/tasks?q=молоко"```

```curl -i localhost:8080/tasks/1```

```
curl -i -X POST localhost:8080/tasks \
-H "Content-Type: application/json"
```

```curl -i localhost:8080/tasks/abc```

```curl -i localhost:8080/tasks/9999```

```
curl -X PATCH localhost:8080/tasks/1 \
-H "Content-Type: application/json" \
-d '{"done":true}'
```

```curl -X DELETE localhost:8080/tasks/1```