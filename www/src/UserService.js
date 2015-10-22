// Copyright 2015 Richard Wisniewski
(function(){
	'use strict';

	angular
		.module('periwinkle.UserService', [])
		.service('UserService', ['$scope', '$http', '$translate', UserService]);
		
	function UserService ($scope, $http, $translate) {
		var self = this;
		
		self.reset = function() {
			self.username = '';
			self.session_id = '';
			self.language = '';
		};
		self.reset();
		
		$scope.logout = self.logout = function () {
			if(session_id && session_id != "") {
				$http({
					method: 'DELETE',
					url: '/v1/session',
					headers: {
						'Content-Type': 'application/json'
					},
					data: {
						session_id: self.session_id
					}
				}).then(
					function success(data, status, headers, config) {
						//do work with response
						debugger;
						self.reset();
						$location.path('/login').replace();
					},
					function fail(data, status, headers, config) {
						//do work with response
						//show error to user
						debugger;
						$scope.loading.is = false;
						//show alert
					}
				);
			}
		};
	}
})();