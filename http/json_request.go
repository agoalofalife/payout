package http

// Json Request from API

type BaseJsonRequest struct {
	ClientOrderId int
}

type DepositionRequest struct {
	BaseJsonRequest
	DstAccount int64
	Amount float64
	Contract string
}

type DepositionRequestPhone struct {
	DepositionRequest
	PaymentParamsPhone
}

func newDepositionJsonRequestPhone()  DepositionRequestPhone{
	return DepositionRequestPhone{}
}

type DepositionPaymentParams struct {
	PofOfferAccepted bool
}

type PaymentParamsPhone struct {
	DepositionPaymentParams
	Property1 uint16
	Property2 uint32
}

type PaymentParamsBankAccount struct {
	DepositionPaymentParams
	CustAccount uint64
	BankBIK uint32
	Payment_purpose string
	Pdr_lastName string
	Pdr_firstName string
	Pdr_middleName string
	Pdr_docNumber uint64
	Pdr_docIssueYear uint16
	Pdr_docIssueMonth uint8
	Pdr_docIssueDay uint8
	Pdr_address string
	PdrBirthDate string
	SmsPhoneNumber uint64

	BankName string
	BankCity string
	BankCorAccount string
}

type PaymentParamsBankCard struct {
	DepositionPaymentParams
	Skr_destinationCardSynonim string
	Pdr_lastName string
	Pdr_firstName string
	Pdr_middleName string
	Pdr_docNumber uint64
	Pdr_docIssueYear uint16
	Pdr_docIssueMonth uint8
	Pdr_docIssueDay uint8
	SmsPhoneNumber uint64
	Pdr_birthDate string
	Pdr_birthPlace string
	Pdr_docIssuedBy string
	Pdr_country string
	Pdr_city string
	Pdr_address string
	Pdr_postcode string
}
