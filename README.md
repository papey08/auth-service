# auth-service

## Описание

Данный проект представляет собой часть сервиса аутентификации, отвечающей за 
выдачу и обновление Access и Refresh токенов.

## Структура проекта

```text
├── cmd
│   └── server
│       └── main.go // точка входа в приложение
│
├── configs
│   └── config.yml // файл с конфигами
│
├── docs
│   ├── coverage
│   │   └── coverage.html // отчёт о покрытии тестами
│   └── swagger // swagger-документация
│
├── internal
│   ├── app // слой бизнес-логики (usecase)
│   │   ├── repo_mocks
│   │   ├── token_mocks
│   │   ├── app.go // реализация интерфейса приложения
│   │   ├── app_interface.go // интерфейс приложения
│   │   └── app_test.go
│   │
│   ├── model // слой сущностей (entities)
│   │   ├── errs.go
│   │   ├── tokens.go
│   │   └── user.go
│   │
│   ├── repo // слой БД
│   │   └── repo.go
│   │
│   └── server // сетевой слой (infrastructure)
│       ├── app_mocks
│       ├── handlers.go
│       ├── handlers_test.go
│       ├── presenters.go
│       ├── responses.go
│       ├── routes.go
│       ├── server.go
│       └── server_test.go
│
├── pkg
│   ├── crypto-tools // пакет для base64 и brypt-хеширования
│   └── tokenizer // пакет для генерации токенов
│
├── README.md
├── go.mod
└── go.sum
```

## Бизнес-логика

Пользователь может получить пару Access, Refresh токенов, используя свой 
идентификатор (GUID). 

Access токен имеет тип JWT, время действия 15 минут, в 
нём закодирован идентификатор пользователя. 

Refresh токен представляет собой случайную строку из цифр и букв латинского 
алфавита в обоих регистрах длиной 64 символа, время действия 30 дней, в базе 
данных хранится в виде bcrypt хеша, пользователь получает его в формате base64. 
Также в базе данных хранится время, когда токен становится недействительным, и 
идентификатор владельца токена.

Пока действителен Refresh токен, пользователь с его помощью может получить 
новую пару токенов. С помощью одного и того же Refresh токена новую пару можно 
получить только 1 раз. По истечении срока действия Refresh токена пользователю 
придётся проходить процесс аутентификации заново.

Также предусмотрено очищение базы данных от всех истёкших токенов раз в сутки.

## Запуск приложения

Самостоятельно сконфигурировать MongoDB, изменить файл [***config.yml***](https://github.com/papey08/auth-service/blob/master/configs/config.yml), 
после чего выполнить команды

```shell
$ go mod download
$ go run cmd/server/main.go
```

## Формат запросов

Swagger-документация доступна по адресу [http://localhost:8080/auth/v1/swagger/index.html](http://localhost:8080/auth/v1/swagger/index.html) 
либо в файле [***swagger.json***](https://github.com/papey08/auth-service/blob/master/docs/swagger/swagger.json)

### Sign In

* Метод: `POST`
* Эндпоинт: `http://localhost:8080/auth/v1/sign-in/6F9619FF-8B86-D011-B42D-00CF4FC964FF`
* Формат ответа:

```json
{
    "data": {
        "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTE5NzMwNjIsInN1YiI6IjZGOTYxOUZGLThCODYtRDAxMS1CNDJELTAwQ0Y0RkM5NjRGRiJ9.yWwWIb8SaVhCz063zo1CAqC1o1N6BSEEReBQjb4csxb9iK28IayXi13g_jIchfG9qWZvgZpAj8OMTedsYoWLAg",
        "refresh_token": "UkI2cjJVaGZKcnp1YXZka1p3a3NTVUdqZnZ0ZzNVTWVXRDJ5Rkh5Umhuek9Xb2ZBdWZFemtrODc5bXQ0NUtSQw=="
    },
    "error": null
}
```

### Refresh

* Метод: `POST`
* Эндпоинт: `http://localhost:8080/auth/v1/refresh`
* Формат запроса:

```json
{
    "refresh_token": "UkI2cjJVaGZKcnp1YXZka1p3a3NTVUdqZnZ0ZzNVTWVXRDJ5Rkh5Umhuek9Xb2ZBdWZFemtrODc5bXQ0NUtSQw=="
}
```

* Формат ответа:

```json
{
    "data": {
        "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTE5NzMxNjEsInN1YiI6IjZGOTYxOUZGLThCODYtRDAxMS1CNDJELTAwQ0Y0RkM5NjRGRiJ9.ISAHDzfM720BtZYHjDkZc1MFrsxTsuYpPrLFJT7JqrDd5h5EAwUuptfJVoV1GAmBGEWqjwzQsEYUsh5AJXSRgQ",
        "refresh_token": "OEpEUjJaTVp4VUQzajZERWFJcTRhTjc1cG9HNFNzVWk1UHZVUkpYamN4MEtFOUNzeW1xWkNmNVpTbDkwU3NDRw=="
    },
    "error": null
}
```

## Используемые технологии

* go 1.21
* MongoDB
* [Gin Web Framework](https://github.com/gin-gonic/gin)
