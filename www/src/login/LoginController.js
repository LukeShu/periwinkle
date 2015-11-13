// Copyright 2015 Richard Wisniewski
// Copyright 2015 Luke Shumaker
(function(){
	'use strict';

	angular
		.module('login')
		.controller('LoginController', ['$cookies', '$http', '$scope', '$interval', '$location', '$mdDialog', '$filter', 'UserService', '$timeout', LoginController]);

	function LoginController($cookies, $http, $scope, $interval, $location, $mdDialog, $filter, userService, $timeout) {
		//gives us an anchor to the outer object from within sub objects or functions
		var self = this;
		//set up public fields
		self.username = '';
		self.password = '';
		self.email = "";
		self.comfirmEmail = '';
		self.confirmPassword = '';
		self.captcha_key = '';
		self.isSignup = false;
		//prep the toolbar
		self.title = 'LOGIN.LOGIN.LOGIN';
		self.warn = {
			exists: false,
			prefix:	'',
			message: ''
		};
		$scope.loading = false;

		//for login redir;
		if(userService.loginRedir.has == true) {
			self.warn.exists = true;
			self.warn.prefix = 'LOGIN.LOGIN.MESSAGE';
			self.warn.message = userService.loginRedir.message;
		} else {
			var cookie = userService.getSession();
			debugger;
			if(cookie != null && cookie != "") {
				//the user may have a session
				$scope.loading = true;
				userService.validate(
					function success() {
						//user is logged in
						$location.path('/user').replace();
					},
					function fail(status) {
						//TODO: uh oh!
					},
					function noSession() {
						//the user isn't logged in
						userService.reset();
						$scope.loading = false;
					}
				);
			}
		}

		//public functions
		this.login = function() {
			//http login api call
			$scope.loading = true;
			$http({
				method: 'POST',
				url: '/v1/session',
				headers: {
					'Content-Type': 'application/json'
				},
				data: {
					username: self.username,
					password: self.password
				}
			}).then(
				function success(response) {
					//do work with response
					userService.setSession(response.data.session_id);
					userService.user_id = response.data.user_id;
					if(userService.loginRedir.has) {
						var redir = userService.loginRedir.path;
						userService.loginRedir = null;
						$location.path(redir).replace();
					} else {
						$location.path('/user').replace();
					}
				},
				function fail(response) {
					//do work with response
					//show error to user
					var status_code = response.status;
					var reason = response.data;
					$scope.loading = false;
					//show alert
					switch(status_code){
						case 403:
							$scope.showError('LOGIN.LOGIN.ERRORS.403.TITLE', 'LOGIN.LOGIN.ERRORS.403.CONTENT', '', '#login-button', '#login-button');
							break;
						case 500:
							$scope.showError('GENERAL.ERRORS.500.TITLE', 'GENERAL.ERRORS.500.CONTENT', reason, '#login-button', '#login-button');
							break;
						default:
							$scope.showError('GENERAL.ERRORS.DEFAULT.TITLE', 'GENERAL.ERRORS.DEFAULT.CONTENT', reason, '#login-button', '#login-button');
					}
				}
			);
		}

		this.signup = function(ev) {
			//http signup api call
			$scope.loading = true;
			$http({
				method: 'POST',
				url: '/v1/users',
				headers: {
					'Content-Type': 'application/json'
				},
				data: {
					username: self.username,
					password: self.password,
					email: self.email
				}
			}).then(
				function success(response) {
					//do work with response
					self.login();
				},
				function fail(response) {
					//do work with response
					//show error to user
					var status_code = response.status;
					var reason = response.data;
					var $translate = $filter('translate');
					var $escape = $filter('escapeHTML');
					$scope.loading = false;
					//show alert
					switch(status_code){
						case 409:
							$scope.showError('LOGIN.SIGNUP.ERRORS.409.TITLE', 'LOGIN.SIGNUP.ERRORS.409.CONTENT', '', '#signup-button', '#signup-button');
							break;
						case 500:
							$scope.showError('GENERAL.ERRORS.500.TITLE', 'GENERAL.ERRORS.500.CONTENT', '', '#signup-button', '#signup-button');
							break;
						default:
							$scope.showError('GENERAL.ERRORS.DEFAULT.TITLE', 'GENERAL.ERRORS.DEFAULT.CONTENT', '', '#signup-button', '#signup-button');
					}
					$mdDialog.show(dialog);
				}
			);
		}

		this.togleSignup = function () {
			self.isSignup = !self.isSignup;
			if(!(new RegExp("/^.+@.+\..+$/")).test(self.username))
				self.username = '';
			self.password = '';
			if(self.isSignup) {
				self.title = 'LOGIN.SIGNUP.SIGNUP';
			} else {
				self.title = 'LOGIN.LOGIN.LOGIN';
			}
		}
	}
})();
