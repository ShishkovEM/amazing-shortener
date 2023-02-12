package main

import (
	"github.com/ShishkovEM/amazing-shortener/internal/app/assembly"
	"github.com/ShishkovEM/amazing-shortener/internal/app/config"
)

var (
	lsc config.LinkServiceConfig
)

func main() {

	// Считываем конфигурацию приложения
	lsc.Parse()

	// Собираем и запускаем при
	if lsc.DatabaseDSN != "" {
		assembly.AssembleAndStartAppWithStandAloneDB(lsc)
	} else {
		assembly.AssembleAndStartAppWithFileStorage(lsc)
	}
}
