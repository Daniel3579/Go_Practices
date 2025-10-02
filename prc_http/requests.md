```
curl -i localhost:8080/health
```
![Screenshot](https://github.com/Daniel3579/Go_Practices/blob/main/prc_http/screenshots/Screenshot_1.png)

```
curl -i -X POST localhost:8080/tasks \
-H "Content-Type: application/json" \
-d '{"title":"Купить молоко"}'
```
![Screenshot](https://github.com/Daniel3579/Go_Practices/blob/main/prc_http/screenshots/Screenshot_2.png)

```
curl -i localhost:8080/tasks
```
![Screenshot](https://github.com/Daniel3579/Go_Practices/blob/main/prc_http/screenshots/Screenshot_3.png)

```
curl -i "localhost:8080/tasks?q=молоко"
```
![Screenshot](https://github.com/Daniel3579/Go_Practices/blob/main/prc_http/screenshots/Screenshot_4.png)

```
curl -i localhost:8080/tasks/1
```
![Screenshot](https://github.com/Daniel3579/Go_Practices/blob/main/prc_http/screenshots/Screenshot_5.png)

```
curl -i -X POST localhost:8080/tasks \
-H "Content-Type: application/json"
```
![Screenshot](https://github.com/Daniel3579/Go_Practices/blob/main/prc_http/screenshots/Screenshot_6.png)

```
curl -i localhost:8080/tasks/abc
```
![Screenshot](https://github.com/Daniel3579/Go_Practices/blob/main/prc_http/screenshots/Screenshot_7.png)

```
curl -i localhost:8080/tasks/9999
```
![Screenshot](https://github.com/Daniel3579/Go_Practices/blob/main/prc_http/screenshots/Screenshot_8.png)

```
curl -X PATCH localhost:8080/tasks/1 \
-H "Content-Type: application/json" \
-d '{"done":true}'
```
![Screenshot](https://github.com/Daniel3579/Go_Practices/blob/main/prc_http/screenshots/Screenshot_9.png)

```
curl -X DELETE localhost:8080/tasks/1
```
![Screenshot](https://github.com/Daniel3579/Go_Practices/blob/main/prc_http/screenshots/Screenshot_10.png)
