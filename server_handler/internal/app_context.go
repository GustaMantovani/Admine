package internal

import (
	"context"
	"fmt"
	"sync"

	"github.com/GustaMantovani/Admine/server_handler/internal/config"
	mcserver "github.com/GustaMantovani/Admine/server_handler/internal/mc_server"
)

// AppContext is the singleton application context
type AppContext struct {
	MainCtx         *context.Context
	Config          *config.Config
	MinecraftServer *mcserver.MinecraftServer
}

var (
	instance *AppContext
	once     sync.Once
)

// Init initializes the AppContext singleton with the YAML config path
func Init(configPath string, mainCtx *context.Context) (*AppContext, error) {
	var err error
	once.Do(func() {
		cfg, e := config.LoadConfig(configPath)
		if e != nil {
			err = fmt.Errorf("failed to load config: %w", e)
			return
		}

		mc, e := mcserver.CreateMinecraftServer(cfg.MinecraftServer)
		if e != nil {
			err = fmt.Errorf("failed to load config: %w", e)
			return
		}

		instance = &AppContext{
			MainCtx:         mainCtx,
			Config:          cfg,
			MinecraftServer: &mc,
		}
	})
	if instance == nil && err == nil {
		err = fmt.Errorf("failed to initialize AppContext")
	}
	return instance, err
}

// Get returns the initialized AppContext singleton
func Get() *AppContext {
	return instance
}
