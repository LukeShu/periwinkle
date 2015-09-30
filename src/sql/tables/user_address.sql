create table `user_address` (`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT, 
	`user_id` INT(11) UNSIGNED NOT NULL,
	`medium_id` INT(11) UNSIGNED NOT NULL,
	`address` VARCHAR(50) NOT NULL,
	PRIMARY KEY (`id`),
	CONSTRAINT `user_id_constraint`
		FOREIGN KEY (`user_id`) REFERENCES users (`id`),
		#ON DELETE CASCADE
		#ON UPDATE RESTRICT 
	CONSTRAINT `medium_id_constraint`
		FOREIGN KEY (`medium_id`) REFERENCES medium (`id`)
	) ENGINE=MyISAM;
	

