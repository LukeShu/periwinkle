// Copyright 2015 Richard Wisniewski
(function(){
	'use strict';

	angular
		.module('periwinkle.UserService', [])
		.service('UserService', [UserService]);
		
	function UserService () {
		var self = this;
		
		self.reset = function() {
			self.username = '';
			self.session_id = '';
			self.language = '';
		};
		self.reset();
	}
})();