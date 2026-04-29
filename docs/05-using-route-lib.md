# `5` Использование Route библиотеки

Для упрощения работы с роутингом вы можете использовать стороннюю библиотеку [`max-router`](https://github.com/LZTD1/max-router).

## Установка

```bash
go get github.com/LZTD1/max-router@v1.1.1
```
Роутер позволяет гибко настраивать маршруты, избавляя от волокиты swtich-case, так же роутер позволяет добавлять middlewares.
Все взаимодействия внутри хендлеров сводится к использованию специального обьекта - `maxrouter.Context`, который позволяет удобно и быстро общатся с пользователем.
## Простой пример использования

```go
func adminAccess(next maxrouter.RouteHandler) maxrouter.RouteHandler {
	return maxrouter.HandlerFunc(func(c maxrouter.Context) error {
		if c.UserID() == 777 {
			return next.ServeContext(c)
		}
		log.Printf("[adminAccess] access denied for user %d", c.UserID())
		return nil
	})
}

func main() {
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
    defer stop()
    
    api, _ := maxbot.New(
        os.Getenv("MAX_TOKEN"),
    )
    
    // Инициализация роутера
    r := maxrouter.NewRouter(api)
    r.Use(middleware.Recover()) // Защита от паники
    r.Use(middleware.DefaultLogger())
    
    // --- Регистрация хендлера на команду /start ---
    r.HandleCommand("/start", func(ctx maxrouter.Context) error {
        return ctx.Send("Привет!")
    })
    
    // --- Обработка текста ---
    r.HandleText("admin", func(ctx maxrouter.Context) error {
        return ctx.Reply("Я бот!")
    })
	
    r.With(adminAccess).HandleCommand("/ap", func(c maxrouter.Context) error {
		log.Printf("user %d accessed /ap endpoint", c.UserID())
		return c.Send("Admin!")
	})
    
    // --- Обработка неизвестных команд (NotFound) ---
    r.NotFound(func(c maxrouter.Context) error {
        return c.Send("Извините, такую команду я еще не умею обрабатывать")
    })
    
    // --- Start polling ---
    for update := range api.GetUpdates(ctx) {
        r.Handle(update, ctx)
    }
}
```
Больше примеров в репозитории