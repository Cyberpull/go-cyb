package env

import (
	"os"

	"cyberpull.com/go-cyb/log"

	"github.com/joho/godotenv"
)

func init() {
	_, err := os.Stat(".env")

	if err != nil {
		return
	}

	err = godotenv.Load()

	if err != nil {
		log.Println(err)
	}
}
