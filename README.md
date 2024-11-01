# gRPC File Service

Этот проект реализует gRPC-сервис для загрузки, скачивания и просмотра списка файлов. Сервис написан на Go и использует gRPC для высокопроизводительного общения между клиентом и сервером.

## Структура проекта

```plaintext
grpc-file-service/
├── cmd/
│   └── server/
│       └── main.go              // Точка входа для сервера
├── internal/
│   ├── server/
│   │   └── server.go            // Реализация gRPC-сервера
│   └── storage/
│       └── storage.go           // Логика работы с файлами
├── proto/
│   └── file_service.proto       // gRPC-протофайл
├── storage/                     // Директория хранения файлов
├── go.mod                       // Go-модуль
└── README.md                    // Описание проекта
```

## Функциональные возможности

- **Загрузка файлов**: Клиенты могут загружать бинарные файлы (например, изображения) на сервер.
- **Просмотр списка файлов**: Сервис предоставляет список всех загруженных файлов с указанием даты создания и обновления.
- **Скачивание файлов**: Клиенты могут запрашивать файлы для скачивания.
- **Ограничение подключений**: Сервис ограничивает количество одновременных подключений:
    - До 10 запросов для загрузки и скачивания файлов.
    - До 100 запросов для просмотра списка файлов.

## Установка

1. Склонируйте репозиторий:

    ```bash
    git clone https://github.com/yourusername/grpc-file-service.git
    cd grpc-service
    ```

2. Установите необходимые зависимости:

    ```bash
    go mod tidy
    ```

3. Установите Protocol Buffers

4. Сгенерируйте Go-код из `proto` файла:

    ```bash
    protoc --go_out=. --go-grpc_out=. proto/file_service.proto
    ```

## Запуск сервера

Запустите сервер с помощью следующей команды:

```bash
go run cmd/server/main.go
```

