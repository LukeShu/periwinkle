// Copyright 2015 Richard Wisniewski
(function(){
	'use strict';

	angular
		.module('user')
		.controller('UserController', ['$cookies', '$http', '$scope', '$interval', 'UserService', '$location', UserController]);

	function UserController($cookies, $http, $scope, $interval, userService, $location) {
		//gives us an anchor to the outer object from within sub objects or functions
		var self = this;
		//clears the toolbar and such so we can set it up for this view
		$scope.reset();
		$scope.toolbar.title = "USER.USER";
		//set up public fields

		$scope.toolbar.buttons = [{
			aria_label: "LogOut",
			label:	"SIGNOUT"
		}];
		$scope.toolbar.onclick = function(index) {
			if(index == 0) {
				$scope.logout();
			}
		};

		self.username = 'Richard Wisniewski';
		self.email = 'rwisniew@purdue.edu';
		self.sessionID = "0x1234567890";
		
		self.groups = [];
		
		var __load = function() {
			//http call point at /v1/users/"user_id"
			//on success set userData to reponse.data
			//fail : debugger ;
			//http user profile api call
			$scope.loading.is = true;
			$http({
				method: 'GET',
				url: '/v1/users/' + userService.user_id, 
				headers: {
					'Content-Type': 'application/json'
				},
				data: {
					session_id: userService.session_id
				},
				responseType: 'json'
			}).then(
				function success(response) {
					//do work with response
					self.username = response.data.user_id;
					self.addresses = response.data.addresses;
					$scope.loading.is = false;
				},
				function fail(response) {
					//do work with response
					//show error to user
					$scope.loading.is = false;
				}
			);
		};
		self.createGroup = function() {
			
		};
		self.joinGroup = function() {
			
		};
		
		//check and load
		self.load = function() {
			userService.validate(
				function success() {
					__load();
				},
				function fail(status) {
					debugger;
				},
				function noSession_cb() {
					debugger;
					userService.loginRedir.has = true;
					userService.loginRedir.path = $location.path();
					userService.loginRedir.message = "You will be redirected back to your user once you log in. ";
					$location.path('/login');
				}
			);
		};
	}
})();
