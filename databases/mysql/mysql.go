package mysql

import "github.com/agoalofalife/payout/databases"

type Mysql struct {
	databases.Migrator
	databases.NameDatabase
}

func (m Mysql) Migrate() string {
	return `CREATE TABLE IF NOT EXISTS request (
  id int(11) NOT NULL AUTO_INCREMENT,
  type_transfer varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  dstAccount varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  clientOrderId bigint(11) NOT NULL,
  requestDT datetime NOT NULL,
  amount double(9,2) NOT NULL,
  currency int(11) DEFAULT NULL,
  agentId int(11) NOT NULL,
  contract varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
  paymentParams json DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;`
};


func (m Mysql) GetType() string {
	return "mysql"
}

func (m Mysql) String() string {
	return "mysql"
}