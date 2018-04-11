Загрузить Dep для vendor (https://golang.github.io/dep/docs/installation.html)

```go get -u github.com/golang/dep/cmd/dep```

Установить все зависимости приложения

```dep ensure```

Собрать приложение

```go build -o image_proxy```

Запустить приложение

```./image_proxy```

Проверить 

```curl http://localhost:8080/?url=aHR0cHM6Ly9qcGVnLm9yZy9pbWFnZXMvanBlZy1ob21lLmpwZw&width=200&height=200 > image.jpg```