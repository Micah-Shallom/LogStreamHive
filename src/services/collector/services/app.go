package services

type App struct {
	Config Config
}

func NewApp(cfg Config) *App {
	return &App{Config: cfg}
}
