create table `user_address` (`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT, 
	`user_id` INT(11) UNSIGNED NOT NULL,
	`medium_id` INT(11) UNSIGNED NOT NULL,
	`address` VARCHAR(50) NOT NULL, PRIMARY KEY (`id`)) ENGINE=MyISAM;
