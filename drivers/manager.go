package drivers

import (
	"io"
)

// FILE IS DEPRECATED

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
	RawResponse
	ErrorResponse
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

	getResponseXml([]byte) (XmlResponse, error)
}

type RawResponse interface {
	// get raw byte date from service
	GetRawResponse() []byte
}

type ErrorResponse interface {
	IsError() bool
	GetMessageError() string
}

type XmlResponse interface {
	XmlIsEmpty
	isError() bool
	getMessageError() string
}

// check is empty structure
type XmlIsEmpty interface {
	isEmpty() bool
}
