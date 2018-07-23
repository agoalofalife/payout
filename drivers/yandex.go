package drivers

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
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
)

type Yandex struct {
	TypePayout
	rawResponseData []byte
}

func (yandex Yandex) verify(data []byte, pathCert string) ([]byte, error) {
	path := existCliCommand("openssl")
	cmd := exec.Command(path, "smime", "-verify", "-inform", "PEM", "-nointern", "-certfile", pathCert, "-CAfile", pathCert)
	stdin, err := cmd.StdinPipe()

	if err != nil {
		return nil, err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, string(data))
	}()

	return cmd.CombinedOutput()
}

// Get name driver
func (yandex Yandex) GetName() string {
	return DriverYandex
}

func (yandex *Yandex) ExecutePayout() {
	url := hostName + "/webservice/deposition/api/" + yandex.GetType()

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
	client := &http.Client{Transport: transport}

	xmlReader := yandex.getDataRequest()

	resp, err := client.Post(url, "application/pkcs7-mime", xmlReader)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	yandex.rawResponseData, err = yandex.verify(data, yandexCertVerify)
	if err != nil {
		log.Fatal(err)
	}

	//v := BalanceResponseXml{}
	//err = xml.Unmarshal(yandex.rawResponseData, &yandex)
	//if err != nil {
	//	fmt.Printf("error: %v", err)
	//	return
	//}
	//log.Println(yandex)
	////log.Println(string(out))
	//os.Exit(0)
}

func (yandex Yandex) GetRawResponse() []byte {
	return yandex.rawResponseData
}

// get bool flag is error from type payout
func (yandex Yandex) IsError() bool {
	xmlResponse, err := yandex.getResponseXml(yandex.rawResponseData)
	if err != nil {
		panic(err)
	}
	return xmlResponse.isError()
}

// Get string message from type payout
func (yandex Yandex) GetMessageError() string {
	xmlResponse, err := yandex.getResponseXml(yandex.rawResponseData)
	if err != nil {
		panic(err)
	}
	return xmlResponse.getMessageError()
}

// TYPE REQUESTS YANDEX

// get balance

// helper constructor
func NewBalance(clientOrderId int) TypePayout {
	return BalanceRequest{clientOrderId, BalanceResponseXml{}}
}

type BalanceRequest struct {
	ClientOrderId int // field clientOrderId
	BalanceResponseXml
}

// Get data request
func (request BalanceRequest) getDataRequest() io.Reader {
	agentId, err := strconv.Atoi(os.Getenv("YANDEX_MONEY_PAYOUT_AGENT_ID"))
	if err != nil {
		log.Fatal(err)
	}
	currentTime := time.Now()

	baseXml := BaseXml{
		AgentId:       agentId,
		ClientOrderId: request.ClientOrderId,
		RequestDT:     currentTime.Format(`2006-01-02T15:04:05.999Z`),
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
	dat, err := encryptPackage(buff.Bytes(), yandexSignCert, certPrivateKey, certPassword)
	return bytes.NewBuffer(dat)
}

func (request BalanceRequest) GetType() string {
	return "balance"
}

func (request BalanceRequest) getResponseXml(rawData []byte) (XmlResponse, error) {
	// cache in memory structure
	if request.BalanceResponseXml.isEmpty() {
		v := BalanceResponseXml{}
		err := xml.Unmarshal(rawData, &v)
		if err != nil {
			fmt.Printf("error: %v", err)
			return v, err
		}
		request.BalanceResponseXml = v
		return request.BalanceResponseXml, err
	}

	return request.BalanceResponseXml, nil
}

// Xml structures
type BaseXml struct {
	AgentId       int    `xml:"agentId,attr"`
	ClientOrderId int    `xml:"clientOrderId,attr"`
	RequestDT     string `xml:"requestDT,attr"`
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

func (responseXml BalanceResponseXml) isError() bool {
	return responseXml.Status == statusRejected
}
func (responseXml BalanceResponseXml) getMessageError() string {
	if errorMessage, ok := descriptionErrors[responseXml.Error]; ok {
		return errorMessage
	}
	return "Missing description error"
}
func (responseXml BalanceResponseXml) isEmpty() bool {
	r := responseXml
	return r.Balance == 0.0 && r.Status == 0 && r.Error == 0 && r.ClientOrderId == 0
}

// helper function
func existCliCommand(command string) string {
	path, err := exec.LookPath(command)
	if err != nil {
		log.Fatal("You need to install openssl!")
		os.Exit(-1)
	}
	return path
}

func encryptPackage(data []byte, cert string, privateKey string, certPassword string) ([]byte, error) {
	path := existCliCommand("openssl")
	cmd := exec.Command(path, "smime", "-sign", "-signer", cert, "-inkey", privateKey, "-nochain", "-nocerts", "-outform", "PEM", "-nodetach", "-passin", "pass:", certPassword)
	stdin, err := cmd.StdinPipe()

	if err != nil {
		return nil, err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, string(data))
	}()

	return cmd.CombinedOutput()
}
