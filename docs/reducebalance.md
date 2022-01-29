## Уменьшить баланс пользователя

**Метод**: `POST`

**URL**: `/balance/reduce`

### Формат запроса

```json
{
  "userId": "(строка) UUID пользователя",
  "amount": "(число с плавающей точкой, больше нуля) сумма для списания",
  "currency": "(строка, опционально) 3х буквенный код валюты, в которой указана сумма для списания (по умолчанию RUB)",
  "description": "(строка, опционально) описание (по умолчанию пустая строка)"
}
```

### Формат ответа (успех)

Код ответа: `200 OK`

Тело _отсутствует_

### Формат ответа (ошибка)

Ошибка может возникнуть в следующих случаях:

* Некорректно составлен запрос
* Внутренняя ошибка сервера
* У пользователя недостаточно средств

Примеры:
* Указана неподдерживаемая валюта (поле `code` равно `100`)

  Код ответа: `400 Bad Request`
  ```json
  {
    "error": {
      "message": "specified currency not supported",
      "code": 100
    },
    "payload": null
  }
  ```

* У пользователя недостаточно средств (поле `code` равно `101`)

  Код ответа: `400 Bad Request`
   ```json
  {
    "error": {
        "message": "user has not enough money",
        "code": 101
    },
    "payload": null
  }
  ```
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