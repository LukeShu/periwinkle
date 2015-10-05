create table `subscription` (`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT, 
	`address_id` INT(11) UNSIGNED NOT NULL,
	`group_id` INT(11) UNSIGNED NOT NULL,
	PRIMARY KEY (`id`),
        CONSTRAINT `user_address_id_constraint`
                FOREIGN KEY (`address_id`) REFERENCES user_address (`id`),
        CONSTRAINT `group_id_constraint`
	                FOREIGN KEY (`group_id`) REFERENCES groups (`id`)
	) ENGINE=MyISAM;
