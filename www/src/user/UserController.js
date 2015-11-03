// Copyright 2015 Richard Wisniewski
(function(){
	'use strict';

	angular
		.module('user')
		.controller('UserController', ['$cookies', '$http', '$scope', '$interval', 'UserService', '$location', '$mdDialog', UserController]);

	function UserController($cookies, $http, $scope, $interval, userService, $location, $mdDialog) {
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

		self.info = {
			status: {
				loading: true,
				error:	''
			},
			title:		'USER.INFO.TITLE',
			username:	'',
			addresses:	[],
			fullName:	'',
			set_fullName:	function() {
				//open dialogue
			},
			edit_address:	function(index) {
				
			},
			save_edit_address:	function(index) {
				
			}
		};
		self.groups = {
			status:		{
				loading:	true,
				error:		''
			},
			list:		[],
			'new':		function() {
				
			},
			new_data:	{
				//new group data
			}
		};
		
		self.groups.new = function(ev) {
			$mdDialog.show({
				controller:				NewGroupController,
				templateUrl:			'src/user/new_group.html',
				parent:					angular.element(document.body),
				targetEvent:			ev,
				clickOutsideToClose:	true
			});
		};
		
		self.groups.join = function() {
		
		};
		
		var __load = function() {
			//http call point at /v1/users/"user_id"
			//on success set userData to reponse.data
			//fail : debugger ;
			//http user profile api call
			self.info.status.loading = true;
			$http({
				method: 'GET',
				url: '/v1/users/' + userService.user_id, 
				headers: {
					'Content-Type': 'application/json'
				},
				data: {
					session_id: userService.session_id
				}
			}).then(
				function success(response) {
					//do work with response
					self.info.username = response.data.user_id;
					self.info.addresses = response.data.addresses;
					self.info.status.loading = false;
				},
				function fail(response) {
					//do work with response
					//show error to user
					self.info.status.loading = false;
				}
			);
			self.groups.status.loading = true;
			$http({
				method:	'GET',
				url:	'/v1/groups'
			}).then(
				function success(response) {
					self.groups.list = reponse.data;
					debugger;
					self.groups.status.loading = false;
				},
				function fail(response) {
					debugger;
					self.groups.status.loading = false;
				}
			);
		};
		self.createGroup = function() {
			
		};
		self.joinGroup = function() {
			
		};
		
		//check and load
		self.load = function() {
			$scope.loading.is = true;
			userService.validate(
				function success() {
					$scope.loading.is = false;
					__load();
				},
				function fail(status) {
					debugger;
				},
				function noSession_cb() {
					userService.loginRedir.has = true;
					userService.loginRedir.path = $location.path();
					userService.loginRedir.message = "You will be redirected back to your user once you log in. ";
					$location.path('/login');
				}
			);
		};
	}
	
	function NewGroupController($scope, $mdDialog, $http) {
		var self = $scope.group = this;
		
		$scope.loading = false;
		$scope.title = 'New Group';
		$scope.errors = [];
		
		self.groupname = '';
		
		self.cancel = function() {
			$mdDialog.cancel();
		};
		self.create = function() {
			$scope.loading = true;
			$scope.title = 'Creating Group...';
			$http({
				method: 'POST',
				url: '/v1/groups',
				headers: {
					'Content-Type': 'application/json'
				},
				data: {
					'grouponame': self.groupname
				}
			}).then(
				function success(response) {
					$mdDialog.hide(self.groupname);
				},
				function fail(response) {
					debugger;
					$scope.loading = false;
					$scope.title = 'Fail';
				}
			);
		};
	}
	
	function ChangePasswordController($scope, $mdDialog, $http) {
		var self = this;
		
		$scope.loading = false;
		$scope.title = 'New Group';
		$scope.errors = [];
		
		self.oldPassword = '';
		self.newPassword = ['',''];
	}
	
})();
