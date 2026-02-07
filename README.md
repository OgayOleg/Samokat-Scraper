Samokat Scraper
Парсер витрин Samokat.ru на Go.
Открывает категорию в браузере через Rod, дергает внутренний API Samokat, вытаскивает товары и сохраняет их в удобные .txt файлы (название, актуальная цена, URL).

Проект умеет:
запускаться без/с HTTP‑прокси (в том числе с user:pass@ip:port);
брать настройки из .env;
сохранять товары по подкатегориям в каталог data/.

Структура проекта:
```
samokat-scraper/
├── main.go
├── go.mod
├── go.sum
├── .env.example         # пример настроек
├── data/                # сюда падают .txt файлы (добавь в .gitignore)
└── internal/
    ├── scraper/
    │   └── scraper.go   # основная логика парсинга
    ├── models/
    │   └── dto.go       # структуры JSON-ответа Samokat API
    └── utils/
        ├── helpers.go   # Slugify, формат цен и т.п.
        └── txt.go       # сохранение данных в .txt
```

Настройка .env

-URL категории на сайте
CATEGORY_URL=https://samokat.ru/category/molochnoe-i-yaytsa

-Внутренний API URL, который возвращает JSON с витриной/категорией
-Это URL, который ты копируешь из DevTools → Network → запрос к api-web.samokat.ru
API_URL=https://api-web.samokat.ru/v2/showcases/44bd9fb0-421e-428e-8e45-4854fe9c9eec/categories/dc7a2b88-1957-4897-a115-f9af1b2369f4

-Токен авторизации Samokat (обязательно)
-Берётся из DevTools → Network → запросы к api-web.samokat.ru → заголовок Authorization
AUTH_TOKEN=Bearer eyJ...твой_реальный_токен...

-HTTP-прокси (опционально)
-Если не нужен прокси — оставь пустым
-Для простого прокси без авторизации:
-PROXY=http://ip:port
-Для прокси с логином и паролем:
-PROXY=http://user:pass@ip:port
PROXY=


Как получить AUTH_TOKEN и API_URL (пример)

Открой https://samokat.ru/category/molochnoe-i-yaytsa в браузере.
Включи DevTools (F12) → вкладка Network.
Обнови страницу, прокрути товары, найди запрос к https://api-web.samokat.ru/....
В этом запросе:
во вкладке Headers скопируй:
полный URL → это твой API_URL,
заголовок authorization: Bearer … → это твой AUTH_TOKEN.
Без актуального AUTH_TOKEN этот API отдаёт либо ошибку, либо HTML/заглушку.

