# System Monitoring Demo (Go + Prometheus + Swagger)

Небольшой демонстрационный проект:
- HTTP API на `:8081`
- `GET /api/v1/ping` (документирован в Swagger UI)
- Prometheus метрики всех типов: `counter`, `gauge`, `histogram`, `summary`
- Prometheus поднимается рядом и доступен в браузере

## Быстрый старт

```bash
cd /Users/aleksandr.lekomtsev/code/courses
make tidy
make up
```

После запуска:
- API: http://localhost:8081/api/v1/ping?name=demo
- Swagger UI: http://localhost:8081/swagger/index.html
- Метрики приложения: http://localhost:8081/metrics
- Prometheus UI: http://localhost:9090
- Grafana UI: http://localhost:3000 (login `admin` / `admin`)

## Полезные команды

```bash
make ping
make ping-fail
make metrics
make logs
make down
make urls
```

## PromQL шпаргалка

Подробные примеры запросов с пояснениями: [PROMQL_EXAMPLES.md](/Users/aleksandr.lekomtsev/code/courses/PROMQL_EXAMPLES.md)

## Примеры запросов в Prometheus

- Общий RPS по ping:
```promql
demo:ping_requests:rate1m
```

- Ошибки в секунду:
```promql
demo:ping_errors:rate1m
```

- P95 latency (агрегированная):
```promql
demo:http_request_latency:p95
```

- Средняя latency:
```promql
demo:http_request_latency:avg
```

- Текущее число запросов в обработке:
```promql
demo_http_inflight_requests
```

## Эндпоинты API

- `GET /healthz`
- `GET /api/v1/ping?name=<value>&fail=<true|false>`
- `GET /metrics`
- `GET /swagger/index.html`

## Grafana

- Grafana поднимается через `docker-compose` вместе с сервисом и Prometheus.
- Datasource `Prometheus` и дашборд `System Monitoring Demo` импортируются автоматически через provisioning.
- Также автоматически импортируется runtime-дашборд `System Runtime CPU & Memory` (CPU/RSS/heap/goroutines/GC).
- После старта открой `http://localhost:3000` и зайди под `admin/admin`.

# Общий флоу

app -> Prometheus -> Grafana (Grafana к приложению напрямую не ходит)

1. Приложение отдаёт метрики на /metrics: internal/httpapi/server.go:29.
2. Prometheus скрейпит app:8081/metrics внутри docker-compose сети: deploy/prometheus/prometheus.yml:9.
3. Grafana поднимается в том же docker-compose и получает datasource через provisioning: docker-compose.yml:23, deploy/grafana/provisioning/datasources/datasource.yml:8.
4. Адрес datasource в Grafana: http://prometheus:9090 (имя сервиса Prometheus в сети compose), UI Grafana на http://localhost:3000 (admin/admin).
