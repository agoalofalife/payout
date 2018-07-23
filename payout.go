package payout

import (
	"github.com/agoalofalife/payout/drivers"
	_ "github.com/joho/godotenv/autoload"
	"log"
)

func Start() {
	manager := new(drivers.Definer)
	driver := manager.Define(drivers.DriverYandex, drivers.NewBalance(12))
	driver.ExecutePayout()
	log.Println(driver.GetMessageError())
	log.Println(driver.IsError())
}
