package drivers

import "io"

// string names - driver
const DRIVER_YANDEX = "yandex"

/**
/ Contract is defined type driver and transfers control to the next
*/
type DriverDefined interface {
	define(driver string)
}

type Definer struct {
	driver Driver
}

// Define current driver
func (definer Definer) Define(driver string, payout TypePayout) Driver {
	switch driver {
	case DRIVER_YANDEX:
		return &Yandex{payout, nil}
	default:
		panic("Driver is not found!")
	}
}

/**
/ Base methods for implements drivers payouts
*/
type Driver interface {
	//// called prior to execution of the payment
	//pre() bool
	//// call after to execution of the payment
	//after() bool
	//
	ExecutePayout()

	// get name Driver
	GetName() string

	// response from service - driver
	ConstructorResponse
}

// Builder data request
// example xml or json data
type ConstructorRequest interface {
	// return build data for post request in service Yandex
	// in this situation  get xml
	getDataRequest() io.Reader
}

// types payout
// example
// credit bank, mobile phone, internet purse
type TypePayout interface {
	ConstructorRequest
	// get name type payout
	GetType() string
}

type ConstructorResponse interface {
	// get raw byte date from service
	GetRawResponse() []byte
}
