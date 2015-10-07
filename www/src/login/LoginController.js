(function(){
	'use strict';

	angular
		.module('login')
		.controller('LoginController', ['$cookies', '$http', '$scope', '$interval', LoginController]); 

	function LoginController($cookies, $http, $scope, $interval) {
		//gives us an anchor to the outer object from within sub objects or functions
		var self = this;
		//clears the toolbar and such so we can set it up for this view
		$scope.resetHeader();
		//adds the public fields and methods of this object to the model ($scope)
		$scope.login = self;
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
				url: '/session',
				headers: {
					'Content-Type': 'application/json'
				},
				data: {
					username: self.username,
					password: self.password
				}
			}).then(
				function success() {
					//do work with response
					$scope.loading.is = false;
				},
				function fail() {
					//do work with response
					//show error to user
					$scope.loading.is = false;
				}
			);
		}
		
		this.signup = function() {
			//http signup api call
			alert(self.username);
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
