# File Analysis Service
Микросервис автоматической проверки работ на плагиат

## Архитектура проекта

```
fileAnalysisService/
│
├── cmd/
│   └── main.go                     # Точка входа, запуск HTTP-сервера
│
├── docs/                           # Swagger-документация (генерируемая)
│   └── docs.go
│
├── internal/
│   ├── antiplagiat/
│   │   └── antiplagiat.go          # Интерфейс движка сравнения файлов
│   │
│   ├── api/
│   │   └── analysis_handler.go      # HTTP API (Gin)
│   │
│   ├── apllication/
│   │   └── manager/
│   │       └── manager.go           # Бизнес-логика анализа
│   │
│   ├── config/
│   │   └── config.go                # Конфигурация и переменные окружения
│   │
│   ├── domain/
│   │   ├── antiplagiat/             # Типы/интерфейсы движка анализа
│   │   ├── entities/                # DTO и модели отчётов
│   │   │   ├── analysis_report.go
│   │   │   └── compare_result.go
│   │   └── repository/              # Интерфейсы хранилищ
│   │
│   └── repository/                  # Реализации хранилищ
│       ├── httpFileStorageClient.go # HTTP клиент FileStorageService
│       └── storage.go               # Работа с MinIO (анализ отчётов)
│
├── docker-compose.yml               # Сборка окружения
├── dockerfile                       # Dockerfile сервиса
├── go.mod
├── go.sum
└── README.md
```
## Логика работы микросервиса

Микросервис выполняет автоматическую проверку студенческих работ по определённому типу (например: лабораторная, контрольная, реферат) и формирует отчёты о плагиате.

### Процесс анализа

1. **Получение списка всех работ по `typeWork`**  
   Сервис обращается в FileStorageService и получает ID работы, имя файла, имя пользователя и тип работы.

2. **Загрузка файлов**  
   Каждый файл скачивается через FileStorage API.

3. **Попарное сравнение файлов**  
   Используется интерфейс:
   ```go
   Compare(a, b []byte) (float64, error)
   ```
   Он возвращает число от `0.0` (разные файлы) до `1.0` (идентичные).

4. **Определение плагиата**  
   ```
   isPlagiarism = similarity > 0.8
   ```

5. **Формирование отчёта**  
   Включает UUID, дату, тип работы, результаты сравнений.

6. **Сохранение отчёта в MinIO**  

## Docker Compose

### Запуск
Сначала нужно удалить контейнеры потом запустить, это если запускали какой-то другой микросервис
```
docker stop $(docker ps -aq)
docker rm $(docker ps -aq)
docker compose up --build
```

### Сервисы

| Сервис                | Описание                                |
|-----------------------|-------------------------------------------|
| file-storage-service  | хранилище файлов                          |
| postgres              | база данных для file-storage-service      |
| minio                 | S3-хранилище                              |
| file-analysis-service | сервис анализа                            |


## API (Swagger)

```
http://localhost:8081/swagger/index.html
```

## API Эндпоинты

### Создать новый отчёт  
`POST /works/{typeWork}/reports`

```json
{
  "reportID": "547e1ad8-4b94-4b0b-9e4d-f28fae4fc551",
  "typeWork": "lab",
  "createdAt": "2025-01-12T18:32:11Z",
  "results": [
    {
      "workA": "bd777",
      "workB": "a8333",
      "nameUserA": "Иванов",
      "nameUserB": "Петров",
      "nameFileA": "lab1.docx",
      "nameFileB": "lab1_new.docx",
      "similarity": 0.92,
      "isPlagiarism": true
    }
  ]
}
```

### Получить последний отчёт  
`GET /works/{typeWork}/reports/latest`

| Код | Описание             |
|-----|----------------------|
| 200 | отчёт найден         |
| 404 | отчёта нет           |
| 500 | ошибка хранилища     |

## Пример отчёта

```json
{
  "reportID": "547e1ad8-4b94-4b0b-9e4d-f28fae4fc551",
  "typeWork": "lab",
  "createdAt": "2025-01-12T18:32:11Z",
  "results": [
    {
      "workA": "bd777",
      "workB": "a8333",
      "nameUserA": "Иванов",
      "nameUserB": "Петров",
      "nameFileA": "lab1.docx",
      "nameFileB": "lab1_new.docx",
      "similarity": 0.92,
      "isPlagiarism": true
    }
  ]
}
```


