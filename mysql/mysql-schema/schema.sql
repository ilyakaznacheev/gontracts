CREATE TABLE `company` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `regcode` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `contract` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `sellerid` int(11) NOT NULL,
  `clientid` int(11) NOT NULL,
  `validfrom` date NOT NULL,
  `validto` date NOT NULL,
  `creditamount` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `contract_seller_company_FK` (`sellerid`),
  KEY `contract_client_company_FK` (`clientid`),
  CONSTRAINT `contract_client_company_FK` FOREIGN KEY (`clientid`) REFERENCES `company` (`id`),
  CONSTRAINT `contract_seller_company_FK` FOREIGN KEY (`sellerid`) REFERENCES `company` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `purchase` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `contractid` int(11) NOT NULL,
  `purchasedatetime` datetime NOT NULL,
  `creditspent` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `purchase_contract_FK` (`contractid`),
  CONSTRAINT `purchase_contract_FK` FOREIGN KEY (`contractid`) REFERENCES `contract` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;