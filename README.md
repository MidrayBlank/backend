# backend

### Для удобного быстрого запуска всего приложения

```bash
docker compose up -d
```

- Порт бд для локального подключения 5432
- Порт бэка для локального подключения 8080

> Для завершения работы и удаления всего, что сгенерировал docker compose (удалит volumes, images и containers)

```bash
docker compose down -v --rmi all --remove-orphans
```

---

Пример `.env` файла:

```bash
DATABASE_USER=test_user
DATABASE_PASSWORD=123
DATABASE_DBNAME=test_db

SERVER_PORT=8080

DATABASE_HOST=db
DATABASE_PORT=5432
```

- PS если что DATABASE*HOST и DATABASE_PORT *не менять\_, это путь внутри docker по сети-docker, а порт всегда такой будет внутри контейнера
