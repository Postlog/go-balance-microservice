## Получить баланс пользователя

**Метод**: `POST`

**URL**: `/balance/get`

### Формат запроса

```json
{
  "userId": "(строка) UUID пользователя"
}
```

### Формат ответа (успех)

Код ответа: `200 OK`

Тело

```json
{
  "error": null,
  "payload": {
    "balance": {
      "userId": "(строка) UUID пользователя",
      "value": "(число с плавающей точкой) баланс пользователя в рублях"
    }
  }
}
```

### Формат ответа (ошибка)

Ошибка может возникнуть в следующих случаях:

* Некорректно составлен запрос
* Внутренняя ошибка сервера

Примеры:

* Какие-либо поля переданы в неверном формате

  Код ответа: `400 Bad Request`
   ```json
  {
    "error": {
      "message": "unexpected JSON schema",
      "code": null
    },
    "payload": null
  }
  ```

* Какие-либо обязательные поля отсутствуют

  Код ответа: `400 Bad Request`
  ```json
  {
    "error": {
      "message": "amount: cannot be blank.",
      "code": null
    },
    "payload": null
  }
  ```
* Непредвиденная ошибка на сервере

  Код ответа: `500 Internal Server Error`
  ```json
  {
    "error": {
      "message": "unexpected internal server error",
      "code": null
    },
    "payload": null
  }
  ```