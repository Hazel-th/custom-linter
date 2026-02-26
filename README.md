# loglint

Линтер для проверки лог-записей в Go, совместимый с `golangci-lint`.
Проверяет сообщения в `log/slog` и `go.uber.org/zap` по правилам стиля и безопасности.

## Правила проверки

Линтер проверяет 4 правила:

1. Строчная буква в начале сообщения.
   - ❌ `slog.Info("Starting server")`
   - ✅ `slog.Info("starting server")`

2. Только английский язык в сообщении.
   - ❌ `slog.Error("ошибка подключения")`
   - ✅ `slog.Error("connection failed")`

3. Без спецсимволов и эмодзи.
   - ❌ `logger.Warn("connection failed!!! 🚀")`
   - ✅ `logger.Warn("connection failed")`

4. Без потенциально чувствительных данных.
   - ❌ `slog.Info("token: " + token)`
   - ✅ `slog.Info("token validated")`

Для нарушений доступны `SuggestedFixes` (автоисправление через `-fix`).

## Поддерживаемые логгеры

- `log/slog`
  - package-level: `Debug/Info/Warn/Error`
  - package-level: `*Context`
  - `*slog.Logger`: те же методы
- `go.uber.org/zap`
  - `*zap.Logger`: `Debug/Info/Warn/Error/DPanic/Panic/Fatal`
  - `*zap.SugaredLogger`: `Debug/Info/Warn/Error/...` и `*w`-методы

## Конфигурация

По умолчанию линтер читает файл `.loglint.json` из корня проекта (если файла нет, берутся значения по умолчанию).

Поддерживаемые поля:

- `sensitive_patterns`: дополнительные паттерны чувствительных данных.
- `custom_patterns`: карта regex-паттернов для детекта чувствительных данных.
- `auto_fix`: включить/отключить `SuggestedFixes`.
- `disabled_rules`: список отключенных правил (`lowercase`, `english`, `specialchars`, `sensitive`).

Пример:

```json
{
  "sensitive_patterns": ["password", "token", "secret"],
  "custom_patterns": {
    "email": "[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}",
    "order-id": "\\b\\d{4}\\b"
  },
  "auto_fix": true,
  "disabled_rules": []
}
```

Готовый шаблон: `.loglint.example.json`.

## Установка и запуск

### Запуск как standalone линтер

```bash
go test ./...
go run ./cmd/loglint ./...
```

С автоисправлением:

```bash
go run ./cmd/loglint -fix ./...
```

Показать diff без изменения файлов:

```bash
go run ./cmd/loglint -fix -diff ./...
```

### Интеграция с golangci-lint (module plugin)

1. Создать `.custom-gcl.yml`:

```yaml
version: v2.5.0
plugins:
  - module: github.com/victornechaev/loglint
    path: .
    import: github.com/victornechaev/loglint/pkg/loglint
```

2. Создать `.golangci.yml`:

```yaml
version: "2"
linters:
  default: none
  enable:
    - loglint

  settings:
    custom:
      loglint:
        type: module
        description: Log messages linter
        original-url: github.com/victornechaev/loglint
        settings:
          config-path: .loglint.json
          auto-fix: true
          disabled-rules: []
          sensitive-patterns:
            - refresh token
          custom-patterns:
            order-id: "\\b\\d{4}\\b"
```

3. Собрать кастомный бинарник и запустить:

```bash
golangci-lint custom
./custom-gcl run ./...
```

## Пример результата

```text
main.go:10:12: log message should start with a lowercase letter (loglint)
main.go:11:12: log message should contain only English language (loglint)
main.go:12:12: log message must not contain special symbols or emoji (loglint)
main.go:13:12: log message may contain sensitive data (loglint)
```

## Структура проекта

- `cmd/loglint` — запуск как standalone линтера.
- `internal/analyzer` — ядро анализатора и правила.
- `internal/config` — загрузка `.loglint.json`.
- `pkg/loglint` — публичный пакет и регистрация plugin для `golangci-lint`.
- `pkg/loglint/testdata` — тест-кейсы `analysistest`.

## Тестирование

```bash
go test ./...
```

## Как реализовывался проект

Проект делался поэтапно:

1. Собрал базовый анализатор на `golang.org/x/tools/go/analysis`.
2. Добавил поддержку `log/slog` и `go.uber.org/zap`.
3. Реализовал 4 правила проверки логов и `SuggestedFixes`.
4. Добавил конфигурацию через `.loglint.json` и настройки для `golangci-lint`.
5. Написал тесты на `analysistest` и отдельные unit-тесты для конфигурации.

Использованные материалы:

- документация `go/analysis` (`analysis`, `singlechecker`, `analysistest`);
- документация `golangci-lint` по Module Plugin System;
- пример module-плагина для `golangci-lint`;
- исходники `staticcheck` как ориентир по подходу к анализу кода.

## Использование нейросети

Нейросеть использовалась как инженерный помощник:

- для поиска и сверки технической информации в документации;
- для проектирования структуры проекта;
- для ускорения исправления багов;
- для подготовки и расширения набора тестов.

## CI

- GitHub Actions: `.github/workflows/ci.yml`
