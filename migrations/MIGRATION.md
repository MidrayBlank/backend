## Миграции базы данных

### Управление миграциями через Goose

Для управления схемой БД используется `goose`. Все миграции лежат в папке `migrations/`.

### Установка goose

``` bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Доступные команды

```bash
# Посмотреть статус миграций
sh scripts/migrations.sh --status

# Создать новую миграцию
sh scripts/migrations.sh --new <название_миграции>
 
# Накатить все миграции
sh scripts/migrations.sh --up

# Откатить последнюю миграцию
sh scripts/migrations.sh --down

# Накатить до конкретной версии
sh scripts/migrations.sh --up <версия>

# Откатить до конкретной версии
sh scripts/migrations.sh --down <версия>

```

### Структура базы данных

Все таблицы находятся в схеме `midray`:

| Таблица | Описание |
|---------|----------|
| `midray.geo` | Муниципальные районы (иерархия) |
| `midray.rosstat` | Основная демографическая статистика |
| `midray.rosstat_age` | Возрастная структура населения |

### Структура миграция

``` text
migrations/
├── 001_create_geo_table.sql
├── 002_create_rosstat_table.sql
└── 003_create_rosstat_age_table.sql
```

### Запуск миграция

``` bash
# Запускаем docker compose
docker compose up -d test_db

# Накатываем миграции
sh scripts/migrations.sh --up

# Проверяем статус
sh scripts/migrations.sh --status
```