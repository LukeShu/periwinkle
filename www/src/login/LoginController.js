// Copyright 2015 Richard Wisniewski
(function(){
	'use strict';

	angular
		.module('login')
		.controller('LoginController', ['$cookies', '$http', '$scope', '$interval', '$location', LoginController]); 

	function LoginController($cookies, $http, $scope, $interval, $location) {
		//gives us an anchor to the outer object from within sub objects or functions
		var self = this;
		//clears the toolbar and such so we can set it up for this view
		$scope.resetHeader();
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
				$scope.login.togleSignup();
			}
		} */
		
		//check if already loged in
		var sessionID = $cookies.get("sessionID");
		if(sessionID !== undefined) {
			//show loading spinner
			//validate sessionID
			//if valid take down spinner
			//if invalid (or error)
				//show login screen
				//clear cookie
				//show warning bar
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
				function success(data, status, headers, config) {
					//do work with response
					alert(data);
					$scope.loading.is = false;
					$location.path('/user').replace();
				},
				function fail(data, status, headers, config) {
					//do work with response
					//show error to user
					alert(data);
					$scope.loading.is = false;
					//show alert
				}
			);
		}
		
		this.signup = function() {
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
				function success(data, status, headers, config) {
					//do work with response
					alert(data);
					$scope.loading.is = false;
					$location.path('/user').replace();
				},
				function fail(data, status, headers, config) {
					//do work with response
					//show error to user
					alert(data);
					$scope.loading.is = false;
					//show alert
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
