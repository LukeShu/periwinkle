// Copyright 2015 Richard Wisniewski
(function(){
	'use strict';

	angular
	.module('user')
	.controller('UserController', ['$cookies', '$http', '$scope', '$interval', 'UserService', '$location', '$mdDialog', '$timeout', 'focus', UserController]);

	function UserController($cookies, $http, $scope, $interval, userService, $location, $mdDialog, $timeout, focus) {
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
				self.info.addresses[index].new_address = self.info.addresses[index].address;
				self.info.addresses[index].editing = true;
				focus('edit_address');
			},
			save_edit_address:	function(index) {
				self.info.addresses[index].editing = false;
				if(self.info.addresses[index].address !== self.info.addresses[index].new_address) {
					self.info.addresses[index].loading = true;
					$http({
						method: 'PATCH',
						url: '/v1/users/' + userService.user_id,
						headers: {
							'Content-Type': 'application/json-patch+json'
						},
						data: [
							{
								'op':		'replace',
								'path':		'/addresses/' + index + '/address',
								'value':	self.info.addresses[index].new_address
							},
							{
								'op':		'test',
								'path':		'/addresses/' + index + '/address',
								'value':	self.info.addresses[index].new_address
							},
						]
					}).then(
						function success (response) {
							debugger;
						},
						function fail (response) {
							debugger;
						}
					);
				}
			},
			changePassword:	function(ev) {
				$mdDialog.show({
					controller:				ChangePasswordController,
					templateUrl:			'src/user/change_password.html',
					parent:					angular.element(document.body),
					targetEvent:			ev,
					clickOutsideToClose:	true
				});
			},
			load:	function() {
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
						var i;
						for (i in self.info.addresses) {
							self.info.addresses[i].new_address = '';
							self.info.addresses[i].editing = false;
							self.info.addresses[i].loading = false;
						}
						self.info.status.loading = false;
					},
					function fail(response) {
						//do work with response
						//show error to user
						self.info.status.loading = false;
					}
				);
			}
		};
		self.groups = {
			status:		{
				loading:	true,
				error:		''
			},
			list:		[],
			'new':	function(ev) {
				$mdDialog.show({
					controller:				NewGroupController,
					templateUrl:			'src/user/new_group.html',
					parent:					angular.element(document.body),
					targetEvent:			ev,
					clickOutsideToClose:	true
				}).then(
					function (response) {
						//the dialog responded before closing
						self.groups.load();
					}, function () {
						//the dialog was cancelled
					}
				);
			},
			'join':	function() {

			},
			edit:	function(index) {

			},
			load:	function() {
				self.groups.status.loading = true;
				$http({
					method:	'GET',
					url:	'/v1/groups'
				}).then(
					function success(response) {
						self.groups.list = response.data;
						debugger;
						self.groups.status.loading = false;
					},
					function fail(response) {
						debugger;
						self.groups.status.loading = false;
					}
				);
			}
		};

		//check and load
		self.load = function() {
			$scope.loading.is = true;
			userService.validate(
				function success() {
					$scope.loading.is = false;
					self.info.load();
					self.groups.load();
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

		$timeout(self.load);
	}

	function NewGroupController($scope, $mdDialog, $http) {
		var self = $scope.group = this;

		$scope.loading = false;
		$scope.title = 'New Group';
		$scope.errors = [];

		self.name = '';

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
					'groupname': self.name
				}
			}).then(
				function success(response) {
					debugger;
					$mdDialog.hide(self.name);
				},
				function fail(response) {
					debugger;
					$scope.loading = false;
					$scope.title = 'Fail';
				}
			);
		};
	}

	function ChangePasswordController($scope, $mdDialog, $http, UserService) {
		var self = $scope.password = this;

		$scope.loading = false;
		$scope.title = 'Change Password';
		$scope.errors = [];

		self.oldPassword = '';
		self.newPassword = ['',''];
		var isCancel = false;

		self.cancel = function() {
			isCancel = true;
			$mdDialog.cancel();
		};
		self.create = function() {
			if(isCancel)
				return;
			$scope.loading = true;
			$scope.title = 'Changing Password...';
			$http({
				method: 'PATCH',
				url: '/v1/users/' + UserService.user_id,
				headers: {
					'Content-Type': 'application/json-patch+json'
				},
				data: [
					{
						'op':		'test',
						'path':		'/password',
						'value':	self.oldPassword
					}
				]
			}).then(
				function success(response) {
					debugger;
					$mdDialog.hide("success");
				},
				function fail(response) {
					debugger;
				}
			);
		};
	}

})();
