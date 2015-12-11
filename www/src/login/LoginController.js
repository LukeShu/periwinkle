// Copyright 2015 Richard Wisniewski
// Copyright 2015 Luke Shumaker
(function(){
	'use strict';

	angular
		.module('login')
		.controller('LoginController', ['$cookies', '$http', '$scope', '$interval', '$location', '$mdDialog', '$filter', 'UserService', '$timeout', 'focus', LoginController]);

	function LoginController($cookies, $http, $scope, $interval, $location, $mdDialog, $filter, userService, $timeout, focus) {
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
		$scope.title = 'LOGIN.LOGIN.LOGIN';
		self.warn = {
			exists: false,
			prefix:	'',
			message: ''
		};
		$scope.loading = false;
		var captcha_id = '';
		var captcha_key = '';

		//for login redir;
		if(userService.loginRedir.has == true) {
			self.warn.exists = true;
			self.warn.prefix = 'LOGIN.LOGIN.MESSAGE';
			self.warn.message = userService.loginRedir.message;
		} else {
			var cookie = userService.getSession();
			 ;
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
		};

		this.testCaptcha = function() {
			$http({
				method:	'POST',
				url:	'/v1/captcha',
				data: ["Something needs to be here so there!"]
			}).then(
				function success(response) {
					//store token
					captcha_id = response.data.value;
					 ;
					//show dialog
					$mdDialog.show({
						controller:				'CaptchaController',
						templateUrl:			'src/captcha/captcha.html',
						parent:					angular.element(document.body),
						clickOutsideToClose:	true,
						locals:	{
							'captcha_id': captcha_id
						}
					}).then(
						function hide (response) {
							//dialog return captcha key or nothing (or error)
						},
						function cancel() {

						}
					);
				},
				function fail(response) {
					//show error to user
					var status_code = response.status;
					var reason = response.data;
					$scope.loading = false;
					//show alert
					switch(status_code){
						case 500:
							$scope.showError('GENERAL.ERRORS.500.TITLE', 'GENERAL.ERRORS.500.CONTENT', reason, '#signup-button', '#signup-button');
							break;
						default:
							$scope.showError('GENERAL.ERRORS.DEFAULT.TITLE', 'GENERAL.ERRORS.DEFAULT.CONTENT', reason, '#signup-button', '#signup-button');
					}
				}
			);
		};

		this.signup = function(form, ev, foc) {
			//http signup api call
			focus(foc);
			if(!form.$valid)
				return;
			$scope.loading = true;
			if(false){ //captcha_key == '') {
			} else {
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
						//show error to user
						var status_code = response.status;
						var reason = response.data;
						$scope.loading = false;
						//show alert
						switch(status_code){
							case 409:
								$scope.showError('LOGIN.SIGNUP.ERRORS.409.TITLE', 'LOGIN.SIGNUP.ERRORS.409.CONTENT', '', '#signup-button', '#signup-button');
								break;
							case 500:
								$scope.showError('GENERAL.ERRORS.500.TITLE', 'GENERAL.ERRORS.500.CONTENT', reason, '#signup-button', '#signup-button');
								break;
							default:
								$scope.showError('GENERAL.ERRORS.DEFAULT.TITLE', 'GENERAL.ERRORS.DEFAULT.CONTENT', reason, '#signup-button', '#signup-button');
						}
					}
				);
			}
		};

		this.togleSignup = function () {
			self.isSignup = !self.isSignup;
			if(!(new RegExp("/^.+@.+\..+$/")).test(self.username))
				self.username = '';
			self.password = '';
			if(self.isSignup) {
				$scope.title = 'LOGIN.SIGNUP.SIGNUP';
			} else {
				$scope.title = 'LOGIN.LOGIN.LOGIN';
			}
		};
	}
})();
