create table `message` (`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT, 
	`group_id` INT(11) UNSIGNED NOT NULL,
	`filename` VARCHAR(500) NOT NULL, 
	PRIMARY KEY (`id`),
        CONSTRAINT `group_id_constraint`
	                FOREIGN KEY (`group_id`) REFERENCES groups (`id`)
	) ENGINE=MyISAM;
