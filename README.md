# File Analysis Platform

Проект микросервисной платформы для хранения файлов, анализа содержимого и генерации отчетов с облаками слов.

---

## Архитектура системы

Система состоит из следующих компонентов:

| Микросервис                  | Порт         | Swagger / Документация                       | Описание                                                                 |
|-------------------------------|-------------|---------------------------------------------|-------------------------------------------------------------------------|
| File Storing Service           | 8080        | [Swagger](http://localhost:8080/swagger/index.html) | Сервис хранения файлов в PostgreSQL и MinIO. Поддерживает загрузку файлов. |
| File Analysis Service          | 8081        | [Swagger](http://localhost:8081/swagger/index.html) | Сервис анализа файлов. Генерирует отчеты и сохраняет их в MinIO.       |
| API Gateway                    | 8082        | [Swagger](http://localhost:8082/swagger/index.html) | Централизованный шлюз к микросервисам, объединяет File Storage и Analysis. |
| PostgreSQL (File Storage)      | 5432        | -                                           | База данных для хранения метаданных файлов.                             |
| MinIO (File Storage)           | 9000 / 9001 | 9001 – Web Console                          | Хранилище исходных файлов File Storing Service.                         |
| MinIO (Analytics)              | 9002 / 9003 | 9003 – Web Console                          | Хранилище отчетов File Analysis Service.                                 |

Каждый микросервис имеет **свой README** и встроенную документацию Swagger.

Все сервисы объединены в Docker-сеть `mynetwork` через Docker Compose.

---

## Менеджеры в API Gateway

API Gateway использует несколько менеджеров для работы с микросервисами и генерации данных. Каждый менеджер отвечает за конкретную задачу и инкапсулирует логику взаимодействия с внешними сервисами.

### 1. FileManager

**Тип:** `FileManagerImpl`  
**Назначение:** Управление файлами в File Storing Service.

**Функции:**

- Загружает файлы через HTTP POST.
- Возвращает уникальный `work_id`.
- Абстрагирует детали работы с S3/MinIO.
- Используется при обработке `/files/upload`.

---

### 2. ReportManager

**Тип:** `ReportManagerImpl`  
**Назначение:** Работа с отчетами File Analysis Service.

**Функции:**

- Создает отчет по `work_id`.
- Получает последний отчет по типу работы.
- Обрабатывает ошибки и возвращает удобные сообщения клиенту.

---

### 3. WordCloudManager

**Тип:** `WordCloudManagerImpl`  
**Назначение:** Генерация облаков слов через QuickChart API.

**Функции:**

- Принимает текст из файла или отчета.
- Формирует ссылку на облако слов (PNG) через QuickChart API.
- Настраивает параметры: ширина, высота, удаление стоп-слов, минимальная длина слова.

---

### 4. PlagiarismService

**Тип:** `PlagiarismService`  
**Назначение:** Координация работы всех менеджеров.

**Функции:**

- Загружает файл через `FileManager`.
- Создает отчет через `ReportManager`.
- Генерирует облако слов через `WordCloudManager`.
- Возвращает клиенту `work_id` и ссылку на облако слов.

**Схема взаимодействия менеджеров:**

```text
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
## Пользовательские и технические сценарии

### 1. Загрузка файла и создание отчета

1. Пользователь отправляет файл на API Gateway (`POST /files/upload`).

2. **API Gateway**:
   - Загружает файл через `FileManager`.
   - Создает отчет через `ReportManager`.
   - Генерирует Word Cloud через `WordCloudManager`.

3. **File Storing Service**:
   - Сохраняет файл в MinIO Filestorage.
   - Записывает метаданные в PostgreSQL.

4. **File Analysis Service**:
   - Создает отчет и сохраняет его в MinIO Analytics.

5. **API Gateway** возвращает клиенту:
   - Ссылку на облако слов

---

### 2. Получение последнего отчета

1. Пользователь делает запрос на API Gateway (`GET /works/{typeWork}/reports/last`).

2. **API Gateway** использует `ReportManager` для получения последнего отчета.

3. **File Analysis Service** возвращает отчет или ошибку (`report not found`).

4. **API Gateway** обрабатывает ошибки и возвращает информативное сообщение клиенту.

---

### Обработка ошибок

- **File Storing Service недоступен:** `"file service unavailable"`.
- **File Analysis Service недоступен:** `"analysis service unavailable"`.
- **QuickChart API недоступен:** `"failed to generate word cloud"`.
- Все ошибки логируются, клиент получает информативное сообщение.
- Система продолжает работу с доступными сервисами.

---

## Docker и Docker Compose

Для запуска всей системы используется Docker и Docker Compose:

```bash
cd apiGateway
docker-compose up --build
```
## Сценарий запуска сервисов

1. **PostgreSQL** – хранит метаданные файлов  
   - Порт: `5432`  
   - Volume: `postgres_data`

2. **MinIO Filestorage** – хранит исходные файлы  
   - API: `9000`  
   - Web Console: `9001`  
   - Volume: `minio_data_filestorage`

3. **File Storing Service** – загружает файлы и создает `work_id`  
   - Порт: `8080`

4. **MinIO Analytics** – хранит отчеты  
   - API: `9002` (внутри контейнера 9000)  
   - Web Console: `9003` (внутри контейнера 9001)  
   - Volume: `minio_data_analytics`

5. **File Analysis Service** – создает отчеты и загружает их в MinIO Analytics  
   - Порт: `8081`

6. **API Gateway** – объединяет `FileManager`, `ReportManager` и `WordCloudManager`  
   - Порт: `8082`

---

## Swagger

Каждый микросервис имеет собственную документацию API через Swagger:

- **API Gateway** – [http://localhost:8082/swagger/index.html#/](http://localhost:8082/swagger/index.html#/)  
- **File Storing Service** – [http://localhost:8080/swagger/index.html#/](http://localhost:8080/swagger/index.html#/)  
- **File Analysis Service** – [http://localhost:8081/swagger/index.html#/](http://localhost:8081/swagger/index.html#/)  


