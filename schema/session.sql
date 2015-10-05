create table `session` (`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT, 
	`fullname` VARCHAR(20) NOT NULL,
	`passwordHash` VARCHAR(40) NOT NULL,
	PRIMARY KEY (`id`)) ENGINE=MyISAM;
