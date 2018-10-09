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

// get balance

// helper constructor
func NewBalance(clientOrderId int) BalanceRequest {
	return BalanceRequest{clientOrderId,nil, BalanceResponseXml{}}
}

type BalanceRequest struct {
	ClientOrderId int // field clientOrderId
	rawResponseData []byte
	BalanceResponseXml
}

func (request BalanceRequest) getType() string {
	return "balance"
}

// Get data request
func (request BalanceRequest) getRequestPackage() io.Reader {
	agentId, err := strconv.Atoi(agentId)
	if err != nil {
		log.Fatal(err)
	}

	baseXml := BaseXml{
		AgentId:       agentId,
		ClientOrderId: request.ClientOrderId,
		RequestDT:     time.Now(),
	}

	xmlStruct := balanceRequestXml{
		baseXml,
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

func (request *BalanceRequest) Run() {
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
	if request.BalanceResponseXml.isEmpty() {
		err := xml.Unmarshal(request.rawResponseData, &request.BalanceResponseXml)
		if err != nil {
			fmt.Printf("error: %v", err)
		}
	}
}

func (request BalanceRequest) Balance() float32 {
	return request.BalanceResponseXml.Balance
}


type balanceRequestXml struct {
	BaseXml
	XMLName xml.Name `xml:"balanceRequest"`
}
type BalanceResponseXml struct {
	BaseResponseXml
	Balance float32 `xml:"balance,attr"`
}
func (responseXml BalanceResponseXml) IsError() bool {
	return responseXml.Status == statusRejected
}
func (responseXml BalanceResponseXml) GetMessageError() string {
	if errorMessage, ok := descriptionErrors[responseXml.Error]; ok {
		return errorMessage
	}
	return "Missing description error"
}
func (responseXml BalanceResponseXml) isEmpty() bool {
	r := responseXml
	return r.Balance == 0.0 && r.Status == 0 && r.Error == 0 && r.ClientOrderId == 0
}
