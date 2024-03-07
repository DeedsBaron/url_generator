# UrlGenerator
Сервис отвечает за генерацию ссылки для доступа к пользовательскому объекту
# Запуск
```
make install-go-deps
make vendor-proto
make generate
make run
```
## POST /v1/url/create

Генерирует новую ссылку до тех пор пока она не будет уникальной и отдает пользователю

Request
```
{
    "input_string": "test"
}
```

Response
```
{
    "url": "http://localhost:8081/v1/url/get/WcdxfknkZ1m9U5u8AgIGBnhsk"
}
```

## GET /v1/url/get/{id}

Возвращает пользователю обьект по ссылке, ссылку нельзя использовать дважды 

Request
```
```

Success Response
```
{
  "resultString":"test"
}
```
Err Respone 
```
{
    "code": 3,
    "message": "urlGenerator: can't get string by url: url is already used",
    "details": []
}
```

Сервис доступен ко по gRPC так и по http