package telemetry

import (
	"log"
	"os"
)

func InitLogs() {
	// Configurar o logger padrão para escrever no stdout
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Logs configurados com sucesso.")
}
