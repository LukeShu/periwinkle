// Copyright 2015 Richard Wisniewski
(function(){
	'use strict';

	angular
		.module('periwinkle.UserService', [])
		.service('UserService', ['$http', '$httpProvider', UserService]);
		
	function UserService ($http, $httpProvider) {
		var self = this;
		
		self.reset = function() {
			self.user_id = '';
			self.session_id = '';
			self.loginRedir = {
				has:		false,
				path:		null,
				message:	null
			};
		};
		self.reset();
		
		self.setCookie = function() = {
			$httpProvider.defaults.xsrfCookieName = 'session_id';
		}
		
		self.validate = function(success_cb, fail_cb, noSession_cb) {
			$http({
				method:	'GET',
				url:	'/v1/session',
				responseType: 'json'
			}).then(
				function success(response) {
					if(!response.data) {
						//the user isnt logged in
						self.reset();
						noSession_cb();
					}	else {
						//they are
						self.user_id = response.data.user_id;
						self.session_id = response.data.session_id;
						success_cb();
					}
				},
				function fail(response) {
					self.reset();
					fail_cb(response.status);
				}
			);
		};
	}
})();
