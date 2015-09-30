create table `group_address` (`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT, 
						`group_id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
						`medium_id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
						`address` VARCHAR(500) NOT NULL, 
						PRIMARY KEY (`id`)) ENGINE=MyISAM;
