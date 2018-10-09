package yandex

import (
	_ "github.com/joho/godotenv/autoload"
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"github.com/agoalofalife/payout/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// status response and list errors for human readable
const (
	statusSuccess    = 0
	statusInProgress = 1
	statusRejected   = 3

	errorSyntax                   = 10
	errorAgentId                  = 11
	errorSubAgentId               = 12
	errorCurrency                 = 14
	errorRequestDt                = 15
	errorDstAccount               = 16
	errorAmount                   = 17
	errorClientOrderId            = 18
	errorContract                 = 19
	errorForbiddenOperation       = 21
	errorNotUniqueClientOrderId   = 26
	errorBrokenPack               = 50
	errorInvalidSignature         = 51
	errorUnknowSignCert           = 53
	errorExpiredCert              = 55
	errorAccountClosed            = 40
	errorLockedYaWallet           = 41
	errorUnknowAccount            = 42
	errorOnceLimit                = 43
	errorPeriodLimit              = 44
	errorSmallBalance             = 45
	errorAmountTooSmall           = 46
	errorDepositionRequest        = 48
	errotLimitReceiverBalance     = 201
	errorYandexError              = 30
	errorReceiverRejectDeposition = 31
	errorPaymentExpiredTime       = 105
	errorReceiverPaymentRevert    = 110
)

var descriptionErrors = map[int]string{
	errorSyntax:                   "Ошибка синтаксического разбора XML-документа. Синтаксис документа нарушен или отсутствуют обязательные элементы XML.",
	errorAgentId:                  "Отсутствует или неверно задан идентификатор контрагента (agentId).",
	errorSubAgentId:               "Отсутствует или неверно задан идентификатор канала приема переводов (subagentId).",
	errorCurrency:                 "Отсутствует или неверно задана валюта (currency).",
	errorRequestDt:                "Отсутствует или неверно задано время формирования документа (requestDT).",
	errorDstAccount:               "Отсутствует или неверно задан идентификатор получателя средств (dstAccount).",
	errorAmount:                   "Отсутствует или неверно задана сумма (amount).",
	errorClientOrderId:            "Отсутствует или неверно задан номер транзакции (clientOrderId).",
	errorContract:                 "Отсутствует или неверно задано основание для зачисления перевода (contract).",
	errorForbiddenOperation:       "Запрашиваемая операция запрещена для данного типа подключения контрагента.",
	errorNotUniqueClientOrderId:   "Операция с таким номером транзакции (clientOrderId), но другими параметрами уже выполнялась.",
	errorBrokenPack:               "Невозможно открыть криптосообщение, ошибка целостности пакета.",
	errorInvalidSignature:         "АСП не подтверждена (данные подписи не совпадают с документом).",
	errorUnknowSignCert:           "Запрос подписан неизвестным Яндекс.Деньгам сертификатом.",
	errorExpiredCert:              "Истек срок действия сертификата в системе контрагента.",
	errorAccountClosed:            "Счет закрыт.",
	errorLockedYaWallet:           "Кошелек в Яндекс.Деньгах заблокирован. Данная операция для этого кошелька запрещена.",
	errorUnknowAccount:            "Счета с таким идентификатором не существует.",
	errorOnceLimit:                "Превышено ограничение на единовременно зачисляемую сумму.",
	errorPeriodLimit:              "Превышено ограничение на максимальную сумму зачислений за период времени.",
	errorSmallBalance:             "Недостаточно средств для проведения операции.",
	errorAmountTooSmall:           "Сумма операции слишком мала.",
	errorDepositionRequest:        "Ошибка запроса зачисления перевода на банковский счет, карту, мобильный телефон.",
	errotLimitReceiverBalance:     "Превышен лимит остатка на счете получателя.",
	errorYandexError:              "Технические проблемы на стороне Яндекс.Денег.",
	errorReceiverRejectDeposition: "Получатель перевода отклонил платеж (под получателем понимается сотовый оператор или процессинговый банк).",
	errorPaymentExpiredTime:       "Превышено допустимое время оплаты по данному коду платежа (при оплате наличными через терминалы, салоны связи и пр.).",
	errorReceiverPaymentRevert:    "Получатель перевода вернул платеж (под получателем понимается сотовый оператор или процессинговый банк).",
}

// load env variable in buffer
var (
	// short hostname server yandex
	hostName         = os.Getenv("YANDEX_MONEY_PAYOUT_HOST")
	yandexCertVerify = os.Getenv("YANDEX_CERT_VERIFY_RESPONSE")
	yandexSignCert   = os.Getenv("YANDEX_CERT_PATH")
	certPrivateKey   = os.Getenv("YANDEX_PRIVATE_KEY_PATH")
	certPassword     = os.Getenv("YANDEX_MONEY_PAYOUT_CERT_PASSWORD")
	agentId          = os.Getenv("YANDEX_MONEY_PAYOUT_AGENT_ID")
	contentType      = "application/pkcs7-mime"
)

type TypeRequest interface {
	// return string get type
	getType() string
	// get byte package for post in service
	getRequestPackage() io.Reader
	// post request
	Run()
	ErrorResponse
}

type ErrorResponse interface {
	IsError() bool
	GetMessageError() string
}

// deprecated
type Yandex struct {
	rawResponseData []byte
}

// deprecated
func (yandex Yandex) GetRawResponse() []byte {
	return yandex.rawResponseData
}

// function wrapper for create default client
func clientRequest() *http.Client {
	// Load client cert
	certificate, err := tls.LoadX509KeyPair(yandexSignCert, certPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{certificate},
		InsecureSkipVerify: true,
	}

	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return &http.Client{Transport: transport}
}

// TYPE REQUESTS YANDEX

// get balance

// helper constructor
func NewBalance(clientOrderId int) BalanceRequest {
	return BalanceRequest{clientOrderId, BalanceResponseXml{}, nil}
}

type BalanceRequest struct {
	ClientOrderId int // field clientOrderId
	BalanceResponseXml
	rawResponseData []byte
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

type TestDeposition struct {
	dstAccount    string
	clientOrderId int
	amount        float32
	contract      string
	currency      int
}

// Xml structures
type BaseXml struct {
	AgentId       int       `xml:"agentId,attr"`
	ClientOrderId int       `xml:"clientOrderId,attr"`
	RequestDT     time.Time `xml:"requestDT,attr"`
}
type balanceRequestXml struct {
	BaseXml
	XMLName xml.Name `xml:"balanceRequest"`
}
type MakeDepositionRequestXml struct {
	MakeDepositionRequest xml.Name `xml:"makeDepositionRequest"`
}
type BaseResponseXml struct {
	Status        int       `xml:"status,attr"`
	Error         int       `xml:"error,attr"`
	ProcessedDt   time.Time `xml:"processedDT,attr"`
	ClientOrderId int       `xml:"clientOrderId,attr"`
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
