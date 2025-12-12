# File Analysis Platform

Это проект микросервисной платформы для хранения файлов, анализа содержимого и генерации отчетов с облаками слов.  
Состоит из следующих микросервисов:

| Микросервис                  | Порт         | Swagger / Документация                       | Описание                                                                 |
|-------------------------------|-------------|---------------------------------------------|-------------------------------------------------------------------------|
| File Storing Service           | 8080        | [Swagger](http://localhost:8080/swagger/index.html) | Сервис хранения файлов в PostgreSQL и MinIO. Поддерживает загрузку файлов. |
| File Analysis Service          | 8081        | [Swagger](http://localhost:8081/swagger/index.html) | Сервис анализа файлов. Генерирует отчеты и сохраняет их в MinIO.       |
| API Gateway                    | 8082        | [Swagger](http://localhost:8082/swagger/index.html) | Централизованный шлюз к микросервисам, объединяет File Storage и Analysis. |
| PostgreSQL (File Storage)      | 5432        | -                                           | База данных для хранения метаданных файлов.                             |
| MinIO (File Storage)           | 9000 / 9001 | 9001 – Web Console                          | Хранилище исходных файлов File Storage Service.                         |
| MinIO (Analytics)              | 9002 / 9003 | 9003 – Web Console                          | Хранилище отчетов File Analysis Service.                                 |

---

## Менеджеры в API Gateway

API Gateway использует несколько менеджеров для работы с микросервисами и генерации данных. Каждый менеджер отвечает за конкретную задачу и инкапсулирует логику взаимодействия с внешними сервисами.

---

### 1. FileManager

**Тип:** `FileManagerImpl`  
**Назначение:** Управление файлами в File Storing Service.  

**Что делает:**

- Отправляет файлы в File Storing Service через HTTP POST.
- Возвращает уникальный `work_id` для загруженного файла.
- Абстрагирует детали работы с S3/MinIO и REST API File Storage Service.
- Используется API Gateway при обработке `/files/upload`.

**Пример использования:**

```go
workID, err := fileManager.UploadFile(userName, typeWork, fileName, fileData)
```
### 2. ReportManager

**Тип:** `ReportManagerImpl`  
**Назначение:** Работа с отчетами File Analysis Service.

**Что делает:**

- Создаёт отчёт по `work_id` через File Analysis Service.
- Получает последний отчёт по типу работы.
- Обрабатывает ошибки сервиса и возвращает удобные сообщения клиенту.
- Используется API Gateway для создания и получения отчётов.

**Пример использования:**

```go
report, err := reportManager.GetLatestReport(typeWork)
```

### 3. WordCloudManager

**Тип:** `WordCloudManagerImpl`  
**Назначение:** Генерация облаков слов через QuickChart API.

**Что делает:**

- Получает текст из файла или отчета.
- Формирует ссылку на облако слов (PNG) через QuickChart API.
- Настраивает параметры изображения: ширина, высота, удаление стоп-слов, минимальная длина слова.
- Используется API Gateway при загрузке файла для мгновенной генерации Word Cloud.

**Пример использования:**

```go
cloudURL := wordCloudManager.GenerateWordCloud(fileText)
```
### 4. PlagiarismService

**Тип:** `PlagiarismService`  
**Назначение:** Координация работы всех менеджеров для полного процесса загрузки файла и анализа.

**Что делает:**

- Использует `FileManager` для загрузки файла в File Storing Service.
- С помощью `ReportManager` создаёт отчёт через File Analysis Service.
- Передаёт текст файла в `WordCloudManager` для генерации облака слов.
- Возвращает клиенту одновременно `work_id` и ссылку на облако слов.

**Пример использования:**

```go
cloudURL, err := plagiarismService.UploadFileAndGetWordCloud(userName, typeWork, fileName, fileData)
```
```
[Client] 
   │
   ▼
[API Gateway / PlagiarismService]
   │
   ├─> FileManager → File Storing Service → PostgreSQL + MinIO
   │
   ├─> ReportManager → File Analysis Service → MinIO Analytics
   │
   └─> WordCloudManager → QuickChart API → URL облака слов
```

## i. Краткое описание архитектуры системы

Проект состоит из следующих микросервисов:

1. **File Storing Service** – отвечает за хранение исходных файлов.  
   - Использует PostgreSQL для метаданных.  
   - Использует MinIO для хранения файлов.  
   - Порт: `8080`.

2. **File Analysis Service** – создает отчеты на основе загруженных файлов.  
   - Получает файлы из File Storing Service.  
   - Сохраняет отчеты в MinIO Analytics.  
   - Порт: `8081`.

3. **API Gateway** – объединяет все сервисы, предоставляет единый API клиентам.  
   - Инкапсулирует работу менеджеров (`FileManager`, `ReportManager`, `WordCloudManager`, `PlagiarismService`).  
   - Генерирует облака слов через QuickChart API.  
   - Порт: `8082`.

4. **PostgreSQL** – хранение метаданных о файлах и отчетах.  
   - Порт: `5432`.

5. **MinIO Filestorage** – хранение исходных файлов.  
   - Порты: API `9000`, Web Console `9001`.

6. **MinIO Analytics** – хранение отчетов.  
   - Порты: API `9002`, Web Console `9003`.

Все микросервисы объединены в Docker-сеть `mynetwork` через Docker Compose.  
Каждый сервис имеет **свой README** и встроенную документацию Swagger.

---

## ii. Пользовательские и технические сценарии микросервисов

### Сценарий 1: Загрузка файла и создание отчета

1. **Пользователь** отправляет файл на API Gateway (`POST /files/upload`).  
2. **API Gateway**:
   - Использует `FileManager` для загрузки файла в File Storing Service.
   - Получает `work_id`.
   - С помощью `ReportManager` создаёт отчет через File Analysis Service.
   - С помощью `WordCloudManager` генерирует ссылку на облако слов.
3. **File Storing Service**:
   - Сохраняет файл в MinIO Filestorage.
   - Записывает метаданные в PostgreSQL.
4. **File Analysis Service**:
   - Обрабатывает файл и создает отчет.
   - Сохраняет отчет в MinIO Analytics.
5. **API Gateway** возвращает клиенту:
   - `work_id`
   - Ссылку на облако слов.
### Сценарий 2: Получение последнего отчета

1. **Пользователь** отправляет запрос на API Gateway (`GET /works/:typeWork/reports/last`).  
2. **API Gateway**:
   - Использует `ReportManager` для запроса последнего отчета в File Analysis Service.
3. **File Analysis Service** возвращает отчет или ошибку, если отчёт не найден.
4. **API Gateway** обрабатывает ошибки и возвращает клиенту:
   - Отчет
   - Или информативное сообщение об ошибке (например, `"report not found"`).

---

### Обработка ошибок при отказе микросервисов

- Если **File Storing Service** недоступен: клиент получает `"file service unavailable"`.
- Если **File Analysis Service** недоступен: клиент получает `"analysis service unavailable"`.
- Если **QuickChart API** недоступен: клиент получает `"failed to generate word cloud"`.
- Все ошибки логируются, а клиент получает информативные сообщения.
- Система продолжает работу с доступными компонентами.

---
## Docker и Docker Compose

Для удобного запуска всей системы используется Docker и Docker Compose. Все микросервисы запускаются в отдельных контейнерах и объединяются в одну сеть `mynetwork`.

### Что происходит при запуске:

1. **PostgreSQL** (`postgres-filestorage`) хранит метаданные о файлах.
   - Порт: `5432`
   - Данные сохраняются в volume `postgres_data`.
   - Инициализация базы через миграции из `../fileStoringService/migrations`.

2. **MinIO Filestorage** (`minio-filestorage`) хранит исходные файлы File Storing Service.
   - Порты: API `9000`, Web Console `9001`
   - Volume: `minio_data_filestorage`

3. **File Storing Service** (`file-storage-service`)
   - Порт: `8080`
   - Использует Postgres и MinIO Filestorage.
   - Отвечает за загрузку файлов и генерацию `work_id`.

4. **MinIO Analytics** (`minio-analytics`) хранит отчёты File Analysis Service.
   - Порты: API `9002` (проброшен на 9000 внутри контейнера), Web Console `9003` (внутри 9001)
   - Volume: `minio_data_analytics`

5. **File Analysis Service** (`file-analysis-service`)
   - Порт: `8081`
   - Получает файлы и создаёт отчёты.
   - Загружает отчёты в MinIO Analytics.

6. **API Gateway** (`api-gateway`)
   - Порт: `8082`
   - Объединяет FileManager, ReportManager и WordCloudManager.
   - При загрузке файла создаёт `work_id`, отчёт и ссылку на облако слов.

---

## Swagger

Каждый микросервис использует Swagger для документации API.  

- **API Gateway** – [http://localhost:8082/swagger/index.html#/](http://localhost:8082/swagger/index.html#/)  
- **File Storing Service** – [http://localhost:8080/swagger/index.html#/](http://localhost:8080/swagger/index.html#/)  
- **File Analysis Service** – [http://localhost:8081/swagger/index.html#/](http://localhost:8081/swagger/index.html#/)  

Swagger позволяет тестировать эндпоинты и просматривать структуры запросов и ответов.

---

Пример запуска:
### ВАЖНО: ЕСЛИ ЗАПУСКАЛСЯ docker-compose.yml других микро-сервисов для их теста. То нужно удалить все контейнеры, образы, томы и т.д 
```bash
cd apiGateway
docker-compose up --build
```

###  У каждого микросервиса есть свой README.md
