# Event Service

## Зона ответственности
- Управление созданием, обновлением и удалением постов.
- Управление комментариями, включая возможность оставлять комментарии на комментарии.
- Предоставление данных о постах и комментариях через REST API.

## Границы сервиса
- Использует собственную базу данных PostgreSQL (или Cassandra) для хранения постов, комментариев и дополнительных данных.
- Сервис изолирован от других, взаимодействует с API Gateway по HTTP.
