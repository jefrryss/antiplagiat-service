# File Storing Service

## Описание микросервиса

`File Storing Service` — микросервис для хранения и выдачи файлов, присланных студентами на проверку.  
Сервис отвечает за:

- Загрузку файлов с информацией о пользователе и типе работы.  
- Сохранение файлов в MinIO (S3-совместимое хранилище).  
- Сохранение метаданных о работе в PostgreSQL (кто загрузил, когда, тип работы, оригинальное имя файла, MIME-тип).  
- Получение списка файлов по типу работы(работы по контрольной, домашки и т.д.).  
- Скачивание файлов с оригинальными именами и корректными MIME-типами.

Сервис является частью более широкой архитектуры, включающей:

## Архитектура
fileStoringService/
├─ docs/ - нужен для swagger
├─ cmd/
│ └─ main.go # Точка входа, настройка сервера и роутов
├─ internal/
│ ├─ api/ # HTTP-обработчики (handlers) для взаимодействия с клиентом
│ ├─ application/ # Бизнес-логика, менеджеры, сервисы
│ │ └─ manager/ # Менеджер для работы с файлами (сохраняет, получает)
│ ├─ config/ # Конфигурация сервиса (env, параметры подключения)
│ ├─ domain/ # Сущности (entities) и интерфейсы репозиториев
│ └─ repository/ # Реализация репозиториев. Работа с базой данных и хранилищем файлов (Postgres, MinIO)
├─ migrations/ # SQL миграции для создания таблиц
├─ docker-compose.yml # Сборка и запуск контейнеров
├─ Dockerfile # Сборка Docker-образа сервиса
├─ go.mod / go.sum # Модули Go
└─ README.md # Документация проекта
### Слои архитектуры

1. **Domain**
   - Сущности: `Work`, `File`.
   - Интерфейсы репозиториев: `Repo`, `Storage`.
   - Отвечает только за модель данных и контракты.

2. **Repository**
   - `PostgresDB` — работа с PostgreSQL.
   - `MinioStorage` — работа с MinIO (S3-совместимое хранилище).
   - Реализует интерфейсы, определенные в `domain`.

3. **Application**
   - Менеджеры (`ManagerFileStorage`) объединяют логику работы с репозиториями.
   - Содержат все бизнес-правила (сохранение файла, генерация ID, контроль ошибок).

4. **API**
   - HTTP-обработчики (`FileHandler`) предоставляют REST API.
   - Обрабатывают запросы, вызывают методы менеджеров и формируют ответы.
   - Swagger-документация встроена через комментарии.

5. **Cmd**
   - Точка входа приложения (`main.go`).
   - Настройка сервера Gin, роутеров и подключений к репозиториям.

6. **Migrations**
   - SQL-скрипты для создания таблиц `works` и `files` в PostgreSQL.

---

### Поток данных

1. Клиент отправляет файл через `POST /upload`.
2. `FileHandler` получает файл и параметры.
3. `ManagerFileStorage`:
   - Генерирует идентификаторы работы и файла.
   - Сохраняет файл в MinIO.
   - Сохраняет метаданные в PostgreSQL.
4. Для получения списка файлов или скачивания используется `GET /files/list/{typeWork}` и `GET /files/download/{work_id}`.
5. Метаданные и файл возвращаются клиенту в корректном формате.

---

### Docker и контейнеризация

- Каждый компонент (PostgreSQL, MinIO, File Storing Service) запускается в отдельном контейнере.
- Взаимодействие между контейнерами настроено через Docker Compose.
- Все сервисы можно запустить командой:

```bash
docker compose up --build
## Архитектура проекта

Проект построен по принципам **чистой архитектуры** с разделением ответственности по слоям:

```mermaid
flowchart LR
    Client["Клиент (студент / преподаватель)"] -->|HTTP REST| API["FileHandler (API)"]
    API --> Manager["ManagerFileStorage (Application / бизнес-логика)"]
    Manager --> RepoDB["PostgresDB (Repository)"]
    Manager --> S3Storage["MinIO (S3-совместимое хранилище)"]
    
    RepoDB -->|CRUD| DB["PostgreSQL"]
    S3Storage -->|Put/Get/Delete| FileStorage["Файлы в MinIO"]

    style Client fill:#f9f,stroke:#333,stroke-width:2px
    style API fill:#bbf,stroke:#333,stroke-width:2px
    style Manager fill:#bfb,stroke:#333,stroke-width:2px
    style RepoDB fill:#ffb,stroke:#333,stroke-width:2px
    style S3Storage fill:#fbb,stroke:#333,stroke-width:2px
