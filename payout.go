package payout


import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/agoalofalife/payout/drivers"
)

func Start()  {
	manager := new(drivers.Definer)
	driver := manager.Define("yandex", drivers.BalanceRequest{0})
	driver.ExecutePayout()
}