create table `subscription` (`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT, 
	`address_id` INT(11) UNSIGNED NOT NULL,
	`group_id` INT(11) UNSIGNED NOT NULL,
	PRIMARY KEY (`id`)) ENGINE=MyISAM;
