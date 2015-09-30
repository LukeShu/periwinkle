create table `group_address` (`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT, 
	`group_id` INT(11) UNSIGNED NOT NULL,
	`medium_id` INT(11) UNSIGNED NOT NULL,
	`address` VARCHAR(500) NOT NULL, 
	PRIMARY KEY (`id`),
        CONSTRAINT `group_id_constraint`
                FOREIGN KEY (`group_id`) REFERENCES groups (`id`),
        CONSTRAINT `medium_id_constraint`
                FOREIGN KEY (`medium_id`) REFERENCES medium (`id`)	
	) ENGINE=MyISAM;
