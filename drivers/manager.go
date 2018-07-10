package drivers

// string names - driver
const DRIVER_YANDEX = "yandex";


/**
 / Contract is defined type driver and transfers control to the next
 */
type DriverDefined interface{
	define(driver string)
}


type Definer struct {
	driver Driver
}

// Define current driver
func (definer Definer) Define(driver string) Driver {
	switch driver {
	 case DRIVER_YANDEX:
		return &Yandex{}
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
	//executePayout()

	// get name Driver
	getName() string
}

// Builder data request
// example xml or json data
type ConstructorRequest interface {
	// build data
	 asBuild()
	 // get data
	 toBuild()
}

