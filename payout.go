package payout

import (
	"github.com/agoalofalife/payout/drivers/yandex"
	_ "github.com/joho/godotenv/autoload"
)

func Start() {
	balance := yandex.NewBalance(12)
}
