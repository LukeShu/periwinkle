// Copyright 2015 Richard Wisniewski
(function(){
	'use strict';

	angular
		.module('periwinkle.UserService', [])
		.service('UserService', [UserService]);
		
	function UserService () {
		var self = this;
		
		self.reset = function() {
			self.user_id = '';
			self.session_id = '';
			debugger;
		};
		self.reset();
		
		self.validate = function(sucess_cb, fail_cb) {
			$http({
				method:	'GET',
				url:	'/v1/session'
			}).then(
				function success(response) {
				
					success_cb();
				},
				function fail(response) {
					self.reset();
					fail_cb(reponse.status);
				}
			);
		};
	}
})();
