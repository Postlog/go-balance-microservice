# Тестовое задание на позицию стажера Avito
Репозиторий содержит реализацию [тестового задания](https://github.com/avito-tech/autumn-2021-intern-assignment) на позицию стажера (бекенд)

#### Выполнено
- Основное задание
- Дополнительное задание 1 (работа с балансом в разных валютах)
- Дополнительное задание 2 (история транзакций)
- Докеризация
- Человеко-читабельные описания ошибок API
#### Что использовалось
- Язык программирования: Go
- База данных: PostgreSQL
- API для получения курса валют: [https://exchangeratesapi.io](https://exchangeratesapi.io)
- Фреймворк для работы с HTTP-подключениями: [https://github.com/gofiber/fiber](https://github.com/gofiber/fiber)
#### Запуск
Для запуска требуется выполнить следующие шаги:
1. Скачать исходный код и перейти в директорию с проектом
    ```
    git clone https://github.com/Postlog/go-balance-microservice.git
    cd go-balance-microservice
    ```
2. Настроить переменные окружения. В корне проекта нужно создать `.env` файл следующего содержания:
    ```dotenv
    POSTGRES_USER=<имя_пользователя_которое_будет_использовать_бд>
    POSTGRES_PASSWORD=<пароль_пользователя_который_будет_использовать_бд>
    POSTGRES_DB=<имя_бд>
    POSTGRES_PORT=<порт_для_внешних_покдлючений>
    
    MICROSERVICE_DSN="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@127.0.0.1:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"
    MICROSERVICE_PORT=<порт_на_котором_будут_прослушиваться_HTTP_подключения>
    MICROSERVICE_ER_API_KEY=<API_ключ_ExchangeRatesAPI>
    ```
3. Создать папку для логов:
    ```
    mkdir logs
    ```
4. Для развертывания dev-среды:
   * Запустить контейнер с базой данных
       ```
       docker-compose up
       ```
   * Запустить сервер
       ```
       go run ./cmd/server/
       ```
5. Для развертывания prod-среды:
   * Запустить контейнеры с базой данных и сервером
       ```
       docker-compose -f docker-compose.prod.yml up --build 
       ```
6. Создать таблицы в базе данных, SQL файл с созданием необходимых таблиц находится [тут](sql/tables.sql)

#### Описание API

* [Получить баланс пользователя](docs/getbalance.md) (`POST balance/get`)
* [Пополнить баланс пользователя](docs/addtobalance.md) (`POST balance/add`)
* [Списать с баланса пользователя](docs/reducebalance.md) (`POST balance/reduce`)
* [Перевести средства от одного пользователя к другому](docs/transferfounds.md) (`POST balance/transfer`)
* [Получить историю транзакций](docs/gettransactions.md) (`POST transactions/get`)

#### Детали реализации

Для хранения данных о балансе в базе данных используется тип `numeric(32,2)` (см. [tables.sql](sql/tables.sql))

Для обеспечения устойчивости данных о балансе пользователя использовались следующие подходы:

- На уровне базы данных установлен запрет на хранение отрицательных балансов (см. [tables.sql](sql/tables.sql))
- Для избежания race condition при изменении баланса в базе данных используются `SELECT FOR UPDATE` запросы и транзакции
(см. [balance repository](dataservice/balance/databaserepository/databaserepository.go))

Чтобы избежать ситуации, когда баланс пользователя изменен, а информацию об этом изменении не удалось сохранить в историю
транзакций, на уровне API-хендлеров используется обертка, откатывающая изменение
баланса (см [balance handlers](adapter/http/handlers/balance/balance.go))

При попытке получения баланса пользователя, записи о котором еще нет в базе данных, возвращается 0

При попытке списания средств с баланса пользователя (а так же при переводе средств с баланса), записи о котором еще нет в базе данных, возвращается ошибка "недостаточно средств"

При попытке получить историю транзакций пользователя, записи о котором еще нет в базе данных, возвращается пустой массив 

#### Структура проекта

```
go-balance-microservice/
├─ dataservice/             слой, отвечающий за работу с данными 
│  ├─ balance/              
│  ├─ currency/             
│  ├─ transaction/          
│  └─ models/
│               
├─ service/                 слой бизнес-логики
│  ├─ balance/
│  ├─ transaction/
│  ├─ currency/
│  └─ constants/
│
├─ adapter/                 адаптеры для взаимодействия с бизнес-логикой
│  └─ http/
│     ├─ handlers/
│     ├─ middleware/
│     └─ protocol/          структуры для взаимодействия с микросервисом
│
├─ config/                  конфигурация проекта
├─ cmd/                     главные программы для взаимодействия с приложением
│  └─ server/
│
├─ docs/                    документация
├─ sql/                     скрипты для создания необходимых таблиц
└─ pkg/                     общие пакеты
   ├─ database/
   ├─ logger/
   ├─ errors/
   ├─ types/
   └─  utils/

```
