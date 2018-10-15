package mysql

import (
	"database/sql"
	"github.com/agoalofalife/payout/databases"
	"github.com/agoalofalife/payout/drivers/yandex"
)

type Mysql struct {
	databases.Migrator
	databases.NameDatabase
	databases.Commiter
}

var createTableStatements = []string{
	`CREATE TABLE IF NOT EXISTS request (
  id int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  type_transfer varchar(255) NOT NULL,
  dstAccount varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  clientOrderId bigint(11) NOT NULL,
  requestDT datetime NOT NULL,
  amount double(9,2) NOT NULL,
  currency int(11) DEFAULT NULL,
  agentId int(11) NOT NULL,
  contract varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
  paymentParams json DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;`,
`CREATE TABLE IF NOT EXISTS response (
  id int(11) NOT NULL AUTO_INCREMENT,
  status tinyint(3) UNSIGNED NOT NULL,
  error tinyint(3) UNSIGNED DEFAULT NULL,
  clientOrderId bigint(20) UNSIGNED NOT NULL,
  processedDT datetime NOT NULL,
  balance double(9,2) NOT NULL,
  techMessage varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
  identification varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
  request_id int(11) UNSIGNED NOT NULL,
  PRIMARY KEY (id),
  KEY reponse_request_id_foreign (request_id),
  CONSTRAINT reponse_request_id_foreign FOREIGN KEY (request_id) REFERENCES request (id)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;`,
}

func (m Mysql) Migrate() []string {
	return createTableStatements
};


func (m Mysql) GetType() string {
	return "mysql"
}

func (m Mysql) String() string {
	return "mysql"
}

func (m Mysql) RequestCommit(conn *sql.DB, req yandex.DepositionRequestXml, transferType databases.TypeTransfer) (sql.Result, error) {
	return conn.Exec("INSERT INTO request (type_transfer, dstAccount, clientOrderId, requestDT, amount, currency, agentId, contract) values (?,?,?,?,?,?,?,?)",
		transferType, req.DstAccount, req.ClientOrderId, req.RequestDT, req.Amount, req.Currency, req.AgentId, req.Contract)
}