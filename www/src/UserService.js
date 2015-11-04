// Copyright 2015 Richard Wisniewski
(function(){
	'use strict';

	angular
		.module('periwinkle.UserService', [])
		.service('UserService', ['$http', '$cookies', UserService]);

	function UserService ($http, $cookies) {
		var self = this;

		self.setSession = function(session) {
			self.session_id = session;
			var expireDate = new Date();
			expireDate.setDate(expireDate.getDate() + 1);
			$cookies.put('app_set_session_id', session, {'expires': expireDate});
		};
		self.getSession = function() {
			return $cookies.get('app_set_session_id');
		};

		self.reset = function() {
			self.user_id = '';
			self.loginRedir = {
				has:		false,
				path:		null,
				message:	null
			};
		};
		self.reset();

		self.validate = function(success_cb, fail_cb, noSession_cb) {
			$http({
				method:	'GET',
				url:	'/v1/session'
			}).then(
				function success(response) {
					if(!response.data) {
						//the user isnt logged in
						self.reset();
						noSession_cb();
					}	else {
						//they are
						self.user_id = response.data.user_id;
						self.setSession(response.data.session_id);
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
