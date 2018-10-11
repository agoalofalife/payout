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

func newDepositionJsonRequest()  DepositionRequest{
	return DepositionRequest{nil, nil, nil, ""}
}