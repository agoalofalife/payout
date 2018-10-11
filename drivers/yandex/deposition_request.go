package yandex

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/agoalofalife/payout/utils"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

// make deposition

// helper constructor
func NewDeposition(clientOrderId int, dstAccount int64, amount float64, contract string) DepositionRequest {
	curreny, err := strconv.Atoi(currency)
	if err != nil {
		panic(err)
	}
	return DepositionRequest{clientOrderId,amount, dstAccount, contract, curreny,nil,DepositionResponseXml{}}
}

type DepositionRequest struct {
	ClientOrderId int // field clientOrderId
	Amount float64
	DstAccount int64
	Contract string // max 128 characters
	Currency int
	rawResponseData []byte
	DepositionResponseXml
}

func (request DepositionRequest) getType() string {
	return "testDeposition"
}

// Get data request
func (request DepositionRequest) getRequestPackage() io.Reader {
	agentId, err := strconv.Atoi(agentId)
	if err != nil {
		log.Fatal(err)
	}

	baseXml := BaseXml{
		AgentId:       agentId,
		ClientOrderId: request.ClientOrderId,
		RequestDT:     time.Now(),
	}

	xmlStruct := DepositionRequestXml{
		baseXml,
		fmt.Sprintf("%0.2f", request.Amount),
		request.Currency,
		request.Contract,
		request.DstAccount,
		xml.Name{},
	}

	buff := bytes.NewBuffer([]byte(xml.Header))

	enc := xml.NewEncoder(buff)
	enc.Indent("  ", "    ")

	if err := enc.Encode(xmlStruct); err != nil {
		fmt.Printf("error: %v\n", err)
	}

	dat, err := utils.EncryptPackagePKCS7(buff.Bytes(), yandexSignCert, certPrivateKey, certPassword)
	return bytes.NewBuffer(dat)
}

func (request *DepositionRequest) Run() {
	url := hostName + "/webservice/deposition/api/" + request.getType() // balance

	dataPKCS7 := request.getRequestPackage()

	resp, err := clientRequest().Post(url, contentType, dataPKCS7)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	request.rawResponseData, err = utils.DecryptPackagePKCS7(data, yandexCertVerify)
	if err != nil {
		log.Fatal(err)
	}

	// cache in memory structure
	if request.DepositionResponseXml.isEmpty() {
		err := xml.Unmarshal(request.rawResponseData, &request.DepositionResponseXml)
		if err != nil {
			fmt.Printf("error: %v", err)
		}
	}
}

type DepositionRequestXml struct {
	BaseXml
	Amount  string `xml:"amount,attr"`
	Currency int `xml:"currency,attr"`
	Contract string `xml:"contract,attr"`
	DstAccount int64 `xml:"dstAccount,attr"`
	XMLName xml.Name `xml:"testDepositionRequest"`
}

type DepositionResponseXml struct {
	BaseResponseXml
}
func (responseXml DepositionResponseXml) IsError() bool {
	return responseXml.Status == statusRejected
}
func (responseXml DepositionResponseXml) GetMessageError() string {
	if errorMessage, ok := descriptionErrors[responseXml.Error]; ok {
		return errorMessage
	}
	return "Missing description error"
}
func (responseXml DepositionResponseXml) isEmpty() bool {
	r := responseXml
	return r.Status == 0 && r.Error == 0 && r.ClientOrderId == 0
}
func (responseXml DepositionResponseXml) IsSuccess() bool {
	return responseXml.Status == statusSuccess
}
func (responseXml DepositionResponseXml) IsProgress() bool {
	return responseXml.Status == statusInProgress
}
