# PromQL Examples for System Monitoring Demo

Ниже примеры запросов для `http://localhost:9090` (вкладка **Graph**).

## 1) Трафик и ошибки

### 1.1 RPS по ping (успех + ошибка)
```promql
sum(rate(demo_ping_requests_total[1m]))
```
Показывает текущую общую интенсивность входящих ping-запросов.

### 1.2 RPS по статусу результата
```promql
sum by (result) (rate(demo_ping_requests_total[1m]))
```
Разделяет поток на `success` и `error`, удобно смотреть долю ошибок.

### 1.3 Error rate (ошибки в секунду)
```promql
rate(demo_ping_errors_total[1m])
```
Скорость роста счетчика ошибок.

### 1.4 Error ratio (доля ошибок)
```promql
sum(rate(demo_ping_requests_total{result="error"}[5m]))
/
sum(rate(demo_ping_requests_total[5m]))
```
Доля ошибочных запросов от общего числа за окно 5 минут.

## 2) Latency (histogram)

### 2.1 Средняя latency
```promql
sum(rate(demo_http_request_duration_seconds_sum[1m]))
/
sum(rate(demo_http_request_duration_seconds_count[1m]))
```
Среднее время ответа по всем запросам.

### 2.2 P95 latency
```promql
histogram_quantile(
  0.95,
  sum by (le) (rate(demo_http_request_duration_seconds_bucket[5m]))
)
```
Оценка 95-го перцентиля времени ответа.

### 2.3 P95 по endpoint
```promql
histogram_quantile(
  0.95,
  sum by (le, path) (rate(demo_http_request_duration_seconds_bucket[5m]))
)
```
Позволяет сравнить `/healthz` и `/api/v1/ping`.

## 3) Нагрузка в реальном времени (gauge)

### 3.1 In-flight requests
```promql
demo_http_inflight_requests
```
Сколько запросов обрабатывается прямо сейчас.

### 3.2 Время последнего успешного ping (unix)
```promql
demo_last_ping_unix_timestamp
```
Когда последний раз был успешный `ping`.

### 3.3 Сколько секунд прошло с последнего ping
```promql
time() - demo_last_ping_unix_timestamp
```
Быстрый индикатор «живости» трафика.

## 4) Summary метрики

### 4.1 Среднее бизнес-время ping из summary
```promql
rate(demo_ping_processing_seconds_sum[5m])
/
rate(demo_ping_processing_seconds_count[5m])
```
Среднее время обработки ping на уровне бизнес-логики.

### 4.2 Средний размер payload
```promql
rate(demo_ping_payload_size_bytes_sum[5m])
/
rate(demo_ping_payload_size_bytes_count[5m])
```
Средний размер ответа в байтах.

## 5) Записанные агрегаты (recording rules)

### 5.1 Готовый RPS
```promql
demo:ping_requests:rate1m
```
Запрос из recording rules, быстрее для дашбордов.

### 5.2 Готовый error rate
```promql
demo:ping_errors:rate1m
```
Предагрегированная скорость ошибок.

### 5.3 Готовый P95
```promql
demo:http_request_latency:p95
```
Предрасчитанный p95 по histogram.

### 5.4 Готовая средняя latency
```promql
demo:http_request_latency:avg
```
Предрасчитанное среднее время ответа.

## 6) Полезные фильтры

### 6.1 Только ping endpoint
```promql
sum(rate(demo_http_request_duration_seconds_count{path="/api/v1/ping"}[1m]))
```
Фокус на основной бизнес-эндпоинт.

### 6.2 Только 500 ответы
```promql
sum(rate(demo_http_request_duration_seconds_count{status="500"}[5m]))
```
Скорость HTTP 500 по сервису.

---

Подсказка: для «живых» графиков ставь диапазон `Last 15 minutes`, шаг `1s` или `5s`.
