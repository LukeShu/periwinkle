// Copyright 2015 Richard Wisniewski
(function(){
	'use strict';

	angular
		.module('login')
		.controller('LoginController', ['$cookies', '$http', '$scope', '$interval', '$location', '$mdDialog', '$filter', LoginController]); 

	function LoginController($cookies, $http, $scope, $interval, $location, $mdDialog, $filter) {
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
		$scope.toolbar.buttons = [{
			label: "Signup",
			img_src: "assets/svg/phone.svg",
		}];
		$scope.toolbar.onclick = function(index) {
			if(index == 0) {
				self.togleSignup();
			}
		}; 
		
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
					debugger;
					$location.path('/user').replace();
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
				function success(data, status, headers, config) {
					//do work with response
					self.login();
				},
				function fail(data, status, headers, config) {
					//do work with response
					//show error to user
					debugger;
					var status_code = data.status;
					var reason = data.data;
					var $translate = $filter("translate");
					$scope.loading.is = false;
					//show alert
					switch(status_code){
						case 409:
							$mdDialog.show(
						      $mdDialog.alert()
						        .parent(angular.element(document.querySelector('#content')))
						        .clickOutsideToClose(true)
						        .title($translate("SIGNUP.ERRORS.409.TITLE"))
						        .content($translate("SIGNUP.ERRORS.409.CONTENT"))
						        .ariaLabel('User Creation Error')
						        .ok('Got it!')
						        .targetEvent(ev)
						    );
							break;
						case 500:
							$mdDialog.show(
							      $mdDialog.alert()
							        .parent(angular.element(document.querySelector('#content')))
							        .clickOutsideToClose(true)
							        .title($translate("SIGNUP.ERRORS.500.TITLE"))
							        .content($translate("SIGNUP.ERRORS.500.CONTENT"))
							        .ariaLabel('Server Error')
							        .ok('Got it!')
							        .targetEvent(ev)
							    );
							break;
						default:
							$mdDialog.show(
						      $mdDialog.alert()
						        .parent(angular.element(document.querySelector('#content')))
						        .clickOutsideToClose(true)
						        .title($translate("SIGNUP.ERRORS.DEFAULT.TITLE"))
						        .content($translate("SIGNUP.ERRORS.DEFAULT.CONTENT"))
						        .ariaLabel('Unexpected Error')
						        .ok('Got it!')
						        .targetEvent(ev)
						    );
					}
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
