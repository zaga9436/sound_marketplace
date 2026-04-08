Миграции лежат в этой директории.

- [001_init.sql](C:\Users\User\Desktop\project\migrations\001_init.sql) — начальная схема базы данных

В текущем dev/docker сценарии миграции применяет само приложение при старте, если `AUTO_APPLY_MIGRATIONS=true`.

Ручной запуск нужен только если:

- автоприменение отключено;
- нужно применить SQL отдельно от старта API;
- идет production deploy с отдельным migration step.
