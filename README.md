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

## Архитектура


- **API Gateway** управляет всеми запросами клиентов:
  - `/files/upload` → отправка файлов на File Storing Service.
  - `/works/:typeWork/reports/` → получение отчетов через File Analysis Service.
  - Генерирует Word Cloud через QuickChart API на основе содержимого загруженных файлов.

- **File Storing Service**:
  - Сохраняет файлы в MinIO и метаданные в PostgreSQL.
  - Возвращает `work_id` для дальнейшей обработки.

- **File Analysis Service**:
  - Создает отчеты по файлам.
  - Сохраняет результаты анализа в MinIO Analytics.
  - Может возвращать последний отчет по типу работы.

- **WordCloudManager** (API Gateway):
  - Генерирует ссылку на облако слов через QuickChart API.
  - Используется при загрузке файла и создании отчета.

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

- **API Gateway** – [http://localhost:8082/swagger/index.html](http://localhost:8082/swagger/index.html)  
- **File Storing Service** – [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)  
- **File Analysis Service** – [http://localhost:8081/swagger/index.html](http://localhost:8081/swagger/index.html)  

Swagger позволяет тестировать эндпоинты и просматривать структуры запросов и ответов.

---

Пример запуска:
### ВАЖНО: ЕСЛИ ЗАПУСКАЛСЯ docker-compose.yml других микро-сервисов для их теста. То нужно удалить все контейнеры, образы, томы и т.д 
```bash
cd apiGateway
docker-compose up --build
```

###  У каждого микросервиса есть свйо readme
