# Ozon Tech Task

Этот проект представляет собой простое веб-приложение для управления URL-ссылками. Оно поддерживает два режима хранения данных:
1. **In-Memory** (в памяти).
2. **PostgreSQL** (в базе данных).

Проект реализован на Go и использует Docker для удобного развертывания.

---

## **Содержание**
1. [Запуск проекта](#запуск-проекта)
   - [In-Memory режим](#in-memory-режим)
   - [PostgreSQL режим](#postgresql-режим)
2. [Примеры запросов](#примеры-запросов)
3. [Структура проекта](#структура-проекта)
4. [Технологии](#технологии)

---

## **Запуск проекта**

### **Требования**
- Установленный Docker и Docker Compose.
- Go (если вы хотите запускать проект локально без Docker).

---

### **In-Memory режим**

1. Перейдите в директорию проекта:
   ```bash
   cd ozon_tz
   ```

2. Запустите проект в режиме In-Memory:
   ```bash
   docker-compose --profile inmemory up
   ```

3. Приложение будет доступно по адресу:
  ```bash
  http://localhost:8080/query
  ```

---

### **PostgreSQL режим**

1. Перейдите в директорию проекта:
   ```bash
   cd ozon_tz
   ```

2. Запустите проект в режиме In-Memory:
   ```bash
   docker-compose --profile postgres up
   ```

3. Приложение будет доступно по адресу:
  ```bash
  http://localhost:8080/query
  ```
4. База данных PostgreSQL будет доступна на порту 5432.

---

### **Примеры запросов**

1. Создание поста
```json
{
  "query": "mutation CreatePost($title: String!, $content: String!, $authorId: String!, $allowComments: Boolean!) { createPost(title: $title, content: $content, authorId: $authorId, allowComments: $allowComments) { id title content authorId allowComments createdAt } }",
  "variables": {
    "title": "Test title",
    "content": "Test content",
    "authorId": "user-1",
    "allowComments": true
  }
}
```

Ответ
```json
{
    "data": {
        "createPost": {
            "allowComments": true,
            "authorId": "user-1",
            "content": "Test content",
            "createdAt": "2025-02-10 20:24:25.835852719 +0000 UTC",
            "id": "post-2",
            "title": "Test title"
        }
    }
}
```

2. Добавление коментария
```json
{
  "query": "mutation AddComment($postId: String!, $parentId: String, $authorId: String!, $text: String!) { addComment(postId: $postId, parentId: $parentId, authorId: $authorId, text: $text) { id postId parentId authorId text createdAt } }",
  "variables": {
    "postId": "post-2",
    "parentId": null,
    "authorId": "user-2",
    "text": "Comment"
  }
}
```

Ответ
```json
{
    "data": {
        "addComment": {
            "authorId": "user-2",
            "createdAt": "2025-02-10 20:24:55.277382189 +0000 UTC",
            "id": "com-1",
            "parentId": null,
            "postId": "post-2",
            "text": "Comment"
        }
    }
}
```

3. Получение поста по ID
```json
{
  "query": "query GetPost($id: String!) { post(id: $id) { id title content authorId allowComments createdAt lastComment { id text authorId createdAt } } }",
  "variables": {
  "id": "post-1739216059778261130"
  }
}
```

Ответ
```json
{
    "data": {
        "post": {
            "allowComments": true,
            "authorId": "user-1",
            "content": "Test content",
            "createdAt": "2025-02-10 20:24:25.835852719 +0000 UTC",
            "id": "post-2",
            "lastComment": {
                "authorId": "user-2",
                "createdAt": "2025-02-10 20:24:55.277382189 +0000 UTC",
                "id": "com-1",
                "text": "Comment"
            },
            "title": "Test title"
        }
    }
}
```

4. Получение всех постов
```json
{
  "query": "query { posts { id title content authorId allowComments createdAt } }"
}
```

Ответ
```json
{
    "data": {
        "posts": [
            {
                "allowComments": true,
                "authorId": "user-1",
                "content": "Test content",
                "createdAt": "2025-02-10 20:24:25.835852719 +0000 UTC",
                "id": "post-2",
                "title": "Test title"
            },
            {
                "allowComments": true,
                "authorId": "user-1",
                "content": "Test content",
                "createdAt": "2025-02-10 20:02:15.264586991 +0000 UTC",
                "id": "post-1",
                "title": "Test title"
            }
        ]
    }
}
```

5. Получение коментариев с пагинацией
```json
{
  "query": "query GetComments($postId: String!, $after: String) { comments(postId: $postId, after: $after) { id text authorId parentId createdAt } }",
  "variables": {
    "postId": "post-2",
    "after": null
  }
}
```

Ответ
```json
{
    "data": {
        "comments": [
            {
                "authorId": "user-2",
                "createdAt": "2025-02-10 20:24:55.277382189 +0000 UTC",
                "id": "com-1",
                "parentId": null,
                "text": "Comment"
            }
        ]
    }
}
```

---

### **Структура проекта**

```
ozon_tz/
├── cmd/
│   └── main.go          # Точка входа в приложение
├── internal/
│   ├── handler/         # Обработчики HTTP-запросов
│   ├── service/         # Бизнес-логика приложения
│   └── storage/         # Реализация хранилищ (In-Memory и PostgreSQL)
├── migrations/          # SQL-миграции для PostgreSQL
├── go.mod               # Файл зависимостей Go
├── go.sum               # Файл зависимостей Go
├── Dockerfile           # Dockerfile для сборки образа
└── docker-compose.yml   # Docker Compose для запуска контейнеров
```

---

### **Технологии**

- # Go — язык программирования.
- # PostgreSQL — реляционная база данных.
- # Docker — контейнеризация приложения.
- # SQL Migrate — управление миграциями базы данных.
