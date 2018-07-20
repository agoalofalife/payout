## Переменные окружения
Следуя принцапам [12-факторных приложений](https://12factor.net/ru/), используются
переменные окружения. 
В корне вашего проекта, вам неоходимо создать файл `.env`(пример есть в иходниках`.env.example`),
и заполнить его необходимыми данными.

#### Переменные окружения Yandex Driver

##### YANDEX_MONEY_PAYOUT_HOST
Хост сервиса яндекс(массовых платежей)
Например для тестового окружения используется(`https://bo-demo02.yamoney.ru:9094`)

##### YANDEX_MONEY_PAYOUT_AGENT_ID
Ваш `agentId` выдается сотрудниками Yandex

##### YANDEX_MONEY_PAYOUT_CURRENCY
Код валюты(подробнее в [документации](https://tech.yandex.ru/money/doc/payment-solution/reference/datatypes-docpage/))

##### YANDEX_CERT_PATH
Сертификат который вам вернул Yandex, в формате `pem`.

##### YANDEX_PRIVATE_KEY_PATH
Ваш приватный сертификат в формате `pem`.

##### YANDEX_CERT_VERIFY_RESPONSE
Сертификат для верификации пакетных данных от Yandex.

##### YANDEX_MONEY_PAYOUT_CERT_PASSWORD
Пароль(если есть), от вашего сертификта.