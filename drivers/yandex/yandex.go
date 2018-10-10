package yandex

import (
	"crypto/tls"
	"encoding/xml"
	_ "github.com/joho/godotenv/autoload"
	"io"
	"log"
	"net/http"
	"os"
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
	currency         = os.Getenv("YANDEX_MONEY_PAYOUT_CURRENCY")
	contentType      = "application/pkcs7-mime"
)
// hm... so deprecated, i think..
type TypeRequest interface {
	// return string get type
	getType() string
	// get byte package for post in service
	getRequestPackage() io.Reader
	// post request
	Run()
	ErrorResponse
}
// deprecated
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

type MakeDepositionRequestXml struct {
	MakeDepositionRequest xml.Name `xml:"makeDepositionRequest"`
}
type BaseResponseXml struct {
	Status        int       `xml:"status,attr"`
	Error         int       `xml:"error,attr"`
	ProcessedDt   time.Time `xml:"processedDT,attr"`
	ClientOrderId int       `xml:"clientOrderId,attr"`
}


