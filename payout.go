package payout

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/agoalofalife/payout/drivers"
	"log"
)

func Start() {
	manager := new(drivers.Definer)
	driver := manager.Define(drivers.DriverYandex, drivers.NewBalance(12))
	driver.ExecutePayout()
	log.Println(string(driver.GetRawResponse()))
}
