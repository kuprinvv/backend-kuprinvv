[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-22041afd0340ce965d47ae6ef1cefeee28c7c493a6346c4f15d667ab976d596c.svg)](https://classroom.github.com/a/xR-tWBKa)

## Запуск

```bash
make up    # БД + миграции + приложение
make seed  # тестовые данные
```

Приложение: `http://localhost:8080`
Swagger: `http://localhost:8080/swagger/index.html`


## Решения

**Слоты генерируются по запросу.** При первом обращении к `/rooms/{roomId}/slots/list?date=...` слоты создаются из расписания (30-минутные интервалы) и сохраняются в БД. Повторные запросы отдают из БД. `ON CONFLICT DO NOTHING` защищает от гонки.

**Конференц-ссылка не блокирует бронь.** Если внешний сервис недоступен, бронь создаётся без ссылки.
