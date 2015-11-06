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
		//clears the toolbar and such so we can set it up for this view
		$scope.reset();
		//set up public fields
		self.username = '';
		self.password = '';
		self.email = "";
		self.comfirmEmail = '';
		self.confirmPassword = '';
		self.captcha_key = '';
		self.isSignup = false;
		//prep the toolbar
		$scope.toolbar.title = 'LOGIN';
		/*$scope.toolbar.buttons = [{
			label: "Signup",
			img_src: "assets/svg/phone.svg",
		}];
		$scope.toolbar.onclick = function(index) {
			if(index == 0) {
				self.togleSignup();
			}
		}; */

		//for login redir;
		if(userService.loginRedir.has == true) {
			$scope.toolbar.warn.exists = true;
			$scope.toolbar.warn.prefix = "YOU ARE NOT LOGGED IN";
			$scope.toolbar.warn.message = userService.loginRedir.message;
		} else {
			var cookie = userService.getSession();
			debugger;
			if(cookie != null && cookie != "") {
				//the user may have a session
				$timeout(function() {
					$scope.loading.is = true;
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
							$scope.loading.is = false;
						}
					);
				});
			}
		}

		//public functions
		this.login = function() {
			//http login api call
			$scope.loading.is = true;
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
					var $translate = $filter('translate');
					var $escape = $filter('escapeHTML');
					$scope.loading.is = false;
					//show alert
					switch(status_code){
						case 403:
							$scope.showError('LOGIN.ERRORS.403.TITLE', 'LOGIN.ERRORS.403.CONTENT', '', '#login-button', '#login-button');
							break;
						case 500:
							$scope.showError('ERRORS.500.TITLE', 'ERRORS.500.CONTENT', reason, '#login-button', '#login-button');
							break;
						default:
							$scope.showError('ERRORS.DEFAULT.TITLE', 'ERRORS.DEFAULT.CONTENT', reason, '#login-button', '#login-button');
					}
				}
			);
		}

		this.signup = function(ev) {
			//http signup api call
			$scope.loading.is = true;
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
					$scope.loading.is = false;
					//show alert
					var dialog = null;
					debugger;
					switch(status_code){
						case 409:
							dialog = $mdDialog.alert()
								.parent(angular.element(document.querySelector('#content')))
								.clickOutsideToClose(true)
								.title($translate('SIGNUP.ERRORS.409.TITLE'))
								.content($translate('SIGNUP.ERRORS.409.CONTENT')+"<pre>"+$escape(reason)+"</pre>")
								.ariaLabel('User Creation Error')
								.ok('Got it!')
						        .targetEvent(ev);
							break;
						case 500:
							dialog = $mdDialog.alert()
								.parent(angular.element(document.querySelector('#content')))
								.clickOutsideToClose(true)
								.title($translate('ERRORS.500.TITLE'))
								.content($translate('ERRORS.500.CONTENT')+"<pre>"+$escape(reason)+"</pre>")
								.ariaLabel('Server Error')
								.ok('Got it!')
						        .targetEvent(ev);
							break;
						default:
							dialog = $mdDialog.alert()
								.parent(angular.element(document.querySelector('#content')))
								.clickOutsideToClose(true)
								.title($translate('ERRORS.DEFAULT.TITLE'))
								.content($translate('ERRORS.DEFAULT.CONTENT')+status_code+"<pre>"+$escape(reason)+"</pre>")
								.ariaLabel('Unexpected Response from Server')
								.ok('Got it!')
						        .targetEvent(ev);
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
				$scope.toolbar.title = 'SIGNUP.SIGNUP';
			} else {
				$scope.toolbar.title = 'LOGIN';
			}
		}
	}
})();
