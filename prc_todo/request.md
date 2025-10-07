```
curl -i localhost:8080/health
```
![Screenshot](./screenshots/Screenshot_1.png)

```
curl -i -X POST localhost:8080/api/v1/tasks \
-H "Content-Type: application/json" \
-d '{"title":"Выучить chi"}'
```
![Screenshot](./screenshots/Screenshot_2.png)

```
curl -i localhost:8080/api/v1/tasks
```
![Screenshot](./screenshots/Screenshot_3.png)

```
curl -i localhost:8080/api/v1/tasks/3
```
![Screenshot](./screenshots/Screenshot_4.png)

```
curl -i -X PUT localhost:8080/api/v1/tasks/7 \
-H "Content-Type: application/json" \
-d '{"title":"Выучить chi глубже", "done":true}'
```
![Screenshot](./screenshots/Screenshot_5.png)

```
curl -i -X DELETE localhost:8080/api/v1/tasks/4
```
![Screenshot](./screenshots/Screenshot_6.png)

```
curl -i localhost:8080/api/v1/tasks?done=true&page=2&limit=1
```
![Screenshot](./screenshots/Screenshot_7.png)