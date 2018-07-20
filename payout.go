package payout

import (
	"github.com/agoalofalife/payout/drivers"
	_ "github.com/joho/godotenv/autoload"
)

func Start() {
	manager := new(drivers.Definer)
	driver := manager.Define(drivers.DriverYandex, drivers.NewBalance())
	driver.ExecutePayout()
}
