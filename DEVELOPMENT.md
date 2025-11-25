# Development Guide

Этот документ описывает, как развернуть окружение для разработки, сгенерировать моки, прогнать тесты и запустить 
проект локально.

**Требуемые версии:**
- golang >= v1.23.4
- go.uber.org/mock >= v0.5.2

## Установка окружения.
1. Установите Go 1.23.4+ — https://go.dev/dl/
2. Установите генератор моков. 
```bash
go install go.uber.org/mock/mockgen@latest
```

## Генерация моков

Мы используем [go.uber.org/mock](https://pkg.go.dev/go.uber.org/mock) для генерации моков.

Комментарий в любом файле пакета, рядом с интерфейсами. Сохраняем в пакете mocks, внутри пакета для моков.
```go
//go:generate mockgen -source=configservice.go -destination=./mocks/configservice_mock.go -package=mocks
```

Моки генерируются командой:
```bash
go generate ./...
```

> Не забывайте пересоздавать моки после изменения интерфейсов! Все созданные файлы должны быть закоммитчены в репозиторий.

## Запуск тестов
```bash
go test ./... -race -coverprofile=coverage.out
go tool cover -html=coverage.out
```