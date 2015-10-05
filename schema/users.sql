create table `users` (`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT, 
	`login` VARCHAR(20) NOT NULL,
	`fullname` VARCHAR(40) NOT NULL,
	`password` VARCHAR(40) NOT NULL,
	`email` VARCHAR(40) NOT NULL,
	PRIMARY KEY (`id`)) ENGINE=MyISAM;
