# Документация

## Общая архитектура

Архитектура проекта на контекстном уровне сгенерированная при помощи likec4:

![Диаграмма](architecture.png)

Исходный код для диаграммы находится в файле [architecture.c4](architecture.c4). Посмотреть архитектуру проекта на уровне контейнеров можно использовав код из файла.

## Таблицы для сервисов

### Сервис пользователей

```mermaid
erDiagram
    USERS {
        uuid id PK "Идентификатор пользователя"
        uuid role_id FK "Идентификатор роли"
        string username "Уникальное имя пользователя"
        string password_hash "Хэш пароля"
        string email "Электронная почта"
        datetime created_at "Дата регистрации"
        datetime updated_at "Дата обновления профиля"
    }

    ROLES {
        uuid id PK "Идентификатор роли"
        string name "Название роли"
        text description "Описание роли"
        datetime created_at "Дата создания"
        datetime updated_at "Дата обновления"
    }

    USER_PROFILES {
        uuid id PK "Идентификатор профиля"
        uuid user_id FK "Ссылка на пользователя"
        string first_name "Имя"
        string last_name "Фамилия"
        text bio "Информация о пользователе"
        datetime created_at "Дата создания профиля"
    }

    USERS ||--|| USER_PROFILES : ""

    USERS ||--|| ROLES : ""
```

### Сервис событий

```mermaid
erDiagram
    POSTS {
        uuid id PK "Идентификатор поста"
        uuid user_id "Идентификатор автора"
        string title "Заголовок поста"
        text content "Содержимое поста"
        datetime created_at "Дата создания"
        datetime updated_at "Дата обновления"
    }

    COMMENTS {
        uuid id PK "Идентификатор комментария"
        uuid post_id FK "Идентификатор поста"
        uuid user_id "Идентификатор автора комментария"
        text content "Содержимое комментария"
        datetime created_at "Дата создания"
        datetime updated_at "Дата обновления"
    }
    
    COMMENT_REPLIES {
        uuid id PK "Идентификатор ответа"
        uuid parent_comment_id FK "Идентификатор основного комментария"
        uuid user_id "Идентификатор автора ответа"
        text content "Содержимое ответа"
        datetime created_at "Дата создания"
        datetime updated_at "Дата обновления"
    }

    POST_TAGS {
        uuid id PK "Идентификатор тега"
        uuid post_id FK "Идентификатор поста"
        string tag "Название тега"
        datetime created_at "Дата добавления тега"
    }

    POSTS ||--|| POST_TAGS : ""

    POSTS ||--|| COMMENTS : ""

    COMMENTS ||--|| COMMENT_REPLIES : ""
```

### Сервис статистики

```mermaid
erDiagram
    COMMENT_STATISTICS {
        uuid id PK "Идентификатор записи"
        uuid comment_id "Идентификатор комментария"
        int likes_count "Количество лайков"
        int comments_count "Количество ответов"
        datetime updated_at "Дата обновления статистики"
    }

    POST_STATISTICS {
        uuid id PK "Идентификатор записи"
        uuid post_id "Идентификатор поста"
        int likes_count "Количество лайков"
        int views_count "Количество просмотров"
        int comments_count "Количество комментариев"
        datetime updated_at "Дата обновления статистики"
    }

    USER_STATISTICS {
        uuid id PK "Идентификатор записи"
        uuid user_id "Идентификатор пользователя"
        int posts_count "Количество постов"
        int likes_given "Количество поставленных лайков"
        int comments_count "Количество комментариев"
        datetime updated_at "Дата обновления"
    }

    EVENT_LOGS {
        uuid id PK "Идентификатор события"
        string event_type "Тип события (LIKE, VIEW, COMMENT)"
        string entity "Тип сущности (POST, COMMENT)"
        uuid entity_id "Идентификатор поста/комментария"
        uuid user_id "Идентификатор пользователя, сгенерировавшего событие"
        datetime created_at "Дата события"
    }

    COMMENT_STATISTICS ||--|| EVENT_LOGS : ""

    USER_STATISTICS ||--|| EVENT_LOGS : ""

    POST_STATISTICS ||--|| EVENT_LOGS : ""
```

### Сервис маршрутизации запросов

```mermaid
erDiagram
    API_ENDPOINTS {
        uuid id PK "Идентификатор эндпоинта"
        string path "URL-путь эндпоинта (например, /users)"
        string method "HTTP метод (GET, POST, и т.д.)"
        string service "Имя сервиса-назначения (куда направляется запрос)"
        string description "Описание назначения эндпоинта"
        datetime created_at "Дата создания записи"
        datetime updated_at "Дата обновления записи"
    }
    
    API_LOGS {
        uuid id PK "Идентификатор лога"
        uuid endpoint_id FK "Ссылка на API_ENDPOINTS (какому эндпоинту принадлежит запрос)"
        uuid session_id FK "Ссылка на API_SESSIONS (при наличии)"
        string method "HTTP метод запроса"
        int response_status "HTTP статус ответа"
        datetime timestamp "Время запроса"
        string message "Сообщение лога (детали запроса/ошибки)"
    }
    
    API_SESSIONS {
        uuid id PK "Идентификатор сессии"
        uuid user_id "Идентификатор пользователя, создавшего сессию"
        string session_token "Токен сессии для аутентификации"
        datetime created_at "Дата начала сессии"
        datetime expires_at "Дата истечения сессии"
        string user_agent "User Agent клиента (информация о браузере/устройстве)"
    }
    
    API_ENDPOINTS ||--|| API_LOGS : ""

    API_SESSIONS ||--|| API_LOGS : ""


```
