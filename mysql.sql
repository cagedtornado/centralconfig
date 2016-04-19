CREATE TABLE configitem (
  id int(11) NOT NULL AUTO_INCREMENT,
  application varchar(100) NOT NULL DEFAULT '*',
  name varchar(100) NOT NULL,
  value longtext NOT NULL,
  machine varchar(100) NOT NULL DEFAULT '',
  updated datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY id_UNIQUE (id),
  UNIQUE KEY app_name_machine (application,name,machine),
  KEY idx_application (application)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8;