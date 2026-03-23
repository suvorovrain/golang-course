# Домашнее задание №3
Как домашнее задание №2, только лучше

---

## Задача

Продолжаем развивать наше приложение из второго задания. В рамках этой задачи предстоит реализовать "скелет" будущего полноценного микросервисного приложения, а именно, разделить сервис Collector на Collector и Processor. Ну и не забыть про API Gateway

Все сервисы должны быть реализованы с соблюдением принципов Чистой архитектуры (Clean Architecture).

В итоге должно получиться 4 сервиса:

1. Сервис API Gateway
- Является REST-сервером и одновременно gRPC-клиентом.
- Принимает внешние HTTP-запросы.
- Пробрасывает запрос в сервис Processor через gRPC.
- Предоставляет спецификацию Swagger (OpenAPI) для тестирования запросов и веб интерфейса.

2. Сервис Processor
- Является gRPC-сервером и gRPC-клиентом.
- Некоторый посредник между API Gateway и Collector. В будущем он будет накапливать всю полученную информацию о репозиториях и отдавать ее в API Gateway по запросу.
- На данный момент, просто передает запрос от API Gateway в Collector и возвращает результат обратно, без какой-либо дополнительной логики.

3. Сервис Collector
- Является gRPC-сервером и REST-клиентом.
- Инкапсулирует логику работы с GitHub API (используя наработки из ДЗ №1).
- Принимает запрос от Processor с данными репозитория (owner/repo) и возвращает информацию о нем.

4. Сервис Subscriber
- Он вам дан для примера, с ним **ничего делать не нужно**.

---

## Необходимая функциональность

Необходимо реализовать 2 endpoint'a.
- `GET /api/ping` --- отправить ping запрос из API Gateway в сервисы Processor и Subscriber и получить от них ответ. Выдать пользователю `200 OK` и JSON с информацией о статусе сервисов в формате:
```
{
  "status": "ok",
  "services": [
    {
      "name": "processor",
      "status": "up"
    },
    {
      "name": "subscriber",
      "status": "up"
    }
  ]
}
```
Или `503 Service Unavailable`, если хотя бы один из сервисов недоступен
```
{
  "status": "degraded",
  "services": [
    {
      "name": "processor",
      "status": "down"
    },
    {
      "name": "subscriber",
      "status": "up"
    }
  ]
}
```
- `GET /api/repositories/info?url=<github_repo_url>` --- получить базовую информацию о репозитории.

Пример запроса
```
GET /api/repositories/info?url=https://github.com/golang/go

Ответ

{
  "full_name": "golang/go",`
  "description": "The Go programming language",
  "stars": 123456,
  "forks": 12345,
  "created_at": "2009-11-10T23:00:00Z",
}
```

---

## Требования

- Взаимодействие: Строго gRPC между сервисами.
- Архитектура: Разделение на слои (Use Cases, Domain, Controller, Adapter) в каждом сервисе.
- Swagger: Автоматическая или ручная генерация документации, доступная по эндпоинту (например, /swagger/index.html).
- Обработка ошибок:
    - Корректные статус-коды gRPC.
    - Маппинг ошибок gRPC в соответствующие HTTP-коды на уровне Gateway (например, 404 если репозиторий не найден).
- Наличие dockerfile для всех сервисов и соответствующие изменения в compose.yaml.

---

## Формат сдачи задания

Необходимо добавить ревьюеров (если ранее не были добавлены) в collaborators репозитория: Settings -> Collaborators -> Add people
Ревьюеры:
- https://github.com/suvorovrain
- https://github.com/Dabzelos
- https://github.com/vacmannnn

Работу над заданием необходимо вести в отдельной ветке.

Вы должны использовать предоставленный шаблон. Вы можете спокойно менять написанную логику, если считаете нужным. Однако не должно быть изменений в тестах.

В конце работы необходимо открыть PR из вашей ветки в main **вашего** форка и отметить ревьюеров в разделе Reviewers.

Задание засчитывается, если в CI проходят тесты.

---

### Полезные материалы

- [Чистая архитектура в Go (статья)](https://pavel-v-p.medium.com/clean-architecture-in-go-2708304217f2)
- [Документация gRPC для Go](https://grpc.io/docs/languages/go/basics/)
- [Библиотека Swag для генерации Swagger](https://github.com/swaggo/swag)
- [Примеры чистой архитектуры](https://github.com/golang-school/evolution/tree/main/6-layers-ddd)
- [Про Dockerfile](https://docs.docker.com/build/concepts/dockerfile)
- [Про Docker compose](https://docs.docker.com/compose/gettingstarted)