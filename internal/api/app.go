package api

import (
	"net/http"
	"os"

	"github.com/shmoulana/hashocr/internal/api/handlers"
	"github.com/shmoulana/hashocr/internal/config"
)

func NewApp(sig chan os.Signal) *http.Server {
	return &http.Server{
		Addr:    ":" + os.Getenv(config.EnvHashOcrSvrPort),
		Handler: NewRouter(handlers.NewHandler()),
	}
}
