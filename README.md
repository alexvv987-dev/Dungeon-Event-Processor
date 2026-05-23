# Dungeon challenge event processor

Обработка событий подземелья по конфигурации и файлу событий; вывод лога и финального отчёта. Go 1.22+.

## Требования

- Go 1.22 или новее
- `make` — опционально (Git Bash / WSL / Linux / macOS)

## Быстрый старт

```bash
make run
```

Или без make:

```bash
go run ./cmd/impulse config.json events
```

Пример входных файлов в корне репозитория: `config.json`, `events`.

## Команды

| Команда | Действие |
|---------|----------|
| `make test` | `go test ./...` |
| `make run` | запуск с `config.json` и `events` |
| `make build` | сборка в `bin/impulse` |
| `make clean` | удаление `bin/` |

```bash
go test ./...
go run ./cmd/impulse config.json events
go build -o bin/impulse ./cmd/impulse
```

## Структура проекта

```
cmd/impulse/          CLI
internal/config/      загрузка и валидация config.json
internal/timeclock/   время HH:MM:SS
internal/event/       парсинг строк событий
internal/domain/      игрок, расписание подземелья
internal/engine/      обработка событий
internal/output/      формат сообщений и отчёта
```

## Формат запуска

```text
impulse <config.json> <events>
```

Вывод — в stdout (лог событий и блок `Final report:`).
