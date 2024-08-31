
# KamiBooking REST API backend-приложение на Go

Простое REST API приложение для системы бронирования конференц-залов, использует go-chi/chi как основную библиотеку для роутинга и jackc/pgx для взаимодействия с базой данных (PostgreSQL)<br>

# Установка и запуск приложения
Для запуска приложения вам понадобится лишь *Docker* с *Docker Compsoe*, *Go* и *PostgreSQL* необходимы только в случае если вы собираетесь запускать приложение не с помощью *Docker Compose*

Копирование репозитория
```
git clone https://github.com/withoutsecondD/kamibooking
```

Переход в корневую папку проекта
```
cd kamibooking
```
Сборка приложения и последующий запуск с помощью Docker Compose
```
docker compose build
```
```
docker compose up
```

Данных шагов достаточно, чтобы собрать приложение и запустить его вместе с базой данных, используя Docker Compose

Экземпляр базы данных PostgreSQL слушает на порту <b>:5432</b>, а само приложение будет доступно на локальном хосте по порту <b>:3000</b>, также в приложении дополнительно предусмотрен Adminer, доступный по порту **:8080**, с помощью которого можно удобно проверить базу данных (название базы данных - **kami**)

> Экземпляр базы данных PostgreSQL запускается с именем пользователя **postgres** и паролем **postgres**, название сервера при входе в Adminer - **db**

## Запуск приложения без Docker Compose
Если вы не хотите использовать Docker Compose для запуска приложения, всё что вам нужно сделать это создать `.env` файл в корневой директории проекта, установить все необходимые переменные окружения и запустить скрипт создания таблицы
> Имейте в виду, что для запуска приложения таким образом на вашей машине должен быть установленный PostgreSQL и настроенный `PATH`, который позволяет вам запускать `psql` команды с любого места вашего компьютера.

`.env`
```
POSTGRESQL_DB_USER: <ИМЯ ПОЛЬЗОВАТЕЛЯ>
POSTGRESQL_DB_PASSWORD: <ПАРОЛЬ ПОЛЬЗОВАТЕЛЯ>
POSTGRESQL_DB_HOST: <ХОСТ POSTGRESQL>
POSTGRESQL_DB_PORT: <ПОРТ POSTGRESQL>
POSTGRESQL_DB_NAME: <ИМЯ БАЗЫ ДАННЫХ>
```
<br>

Запуск скрипта создания таблицы (команду необходимо выполнять в корневой директории проекта)
```
psql -h <ХОСТ> -U <ИМЯ ПОЛЬЗОВАТЕЛЯ> -d <ИМЯ БАЗЫ ДАННЫХ> -f ./db_scripts/create_reservation_table.sql
```
После выполнения данной команды psql попросит ввести пароль от пользователя, которого вы указали в команде, введите его
>Если по какой-то причине вы не смогли создать таблицу с помощью скрипта, вы можете создать её вручную, удостоверьтесь в том, что таблица называется `reservations` и вы создаёте её в той базе, которую указали в переменных окружения

После настройки переменных окружения и создания таблицы вы можете подтянуть зависимости, собрать приложение и запустить его следующими командами:
```
go mod download
```
```
go build
```
```
./kamibooking
```
После запуска приложения оно всё так же будет доступно по `localhost:3000`

# API

`GET /reservations/{roomId}`
Возвращает все брони по указанному залу

- Ответ:
  - `200 OK` Со списком бронирований в виде JSON
  - `400 Bad Request` При отправке запроса с неправильным аргументом roomId

`POST /reservations/`
Создаёт новую бронь

- Тело запроса:
  - `room_id` - Id зала, в котором будет создана бронь
  - `start_time` - Время начала брони в формате RFC3339
  - `end_time` - Время конца брони в формате RFC3339
  ```
  Пример тела запроса:
  {
      "room_id":  1,
      "start_time":  "2024-08-30T10:00:00Z",
      "end_time":  "2024-08-30T10:30:00Z"
  }
  ```

- Ответ:
  - `201 Created` При успешном создании брони и отсутствии конфликтов
  - `400 Bad Request` При попытке создания брони с некорректными данными (время начала брони позже времени конца брони)
  - `409 Conflict` При попытке создания брони, которая вызовет конфликт в данном зале

# Запуск тестов
В приложении предусмотрены тесты для основной логики приложения при создании броней, для запуска тестов достаточно запустить команду в корневой директории проекта:
```
go test -v ./...
```

# Использованные библиотеки

- `go-chi/chi/v5` - Для роутинга и написания API
- `jackc/pgx/v5` - Для работы с базой данных
- `stretchr/testify` - Для написания тестов и удобного создания моков
- `joho/godotenv` - Для чтения переменных окружения из файла .env
