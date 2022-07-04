# MODERATION SERVICE

* **POST /moderate** - проверяет комментарий на содержание запрещенных слов. Если такие слова присутствуют, возвращает 400 ошибку.

### Микросервис работает в связке с другими сервисами:
* gateway: https://github.com/MarySmirnova/api_gateway
* сервис комментариев: https://github.com/MarySmirnova/comments_service
* сервис - парсер новостей: https://github.com/MarySmirnova/news_reader

### .env example:

    API_MODER_LISTEN=:8083
	API_MODER_READ_TIMEOUT=30s
	API_MODER_WRITE_TIMEOUT=30s
