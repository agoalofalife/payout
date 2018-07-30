package payout

import (
	"github.com/agoalofalife/payout/drivers/yandex"
	"log"
)

func Start() {
	balance := yandex.NewBalance(12)
	balance.Run()
	log.Println(balance.Balance())
	log.Println(balance.BalanceResponseXml.Balance)
}
