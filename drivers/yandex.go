package drivers


type Yandex struct {}

// Get name driver
func (yandex Yandex) getName() string {
	return DRIVER_YANDEX
}

type BalanceRequest struct {

}