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
		$scope.toolbar.title = "USER.INFO.USER";
		//set up public fields

		$scope.toolbar.buttons = [{
			aria_label: "LogOut",
			label:	"GENERAL.SIGNOUT"
		}];
		$scope.toolbar.onclick = function(index) {
			if(index == 0) {
				$scope.logout();
			}
		};

		self.info = {
			status: {
				loading: 	true,
				error:		'',
				editing:	false
			},
			title:		'USER.INFO.TITLE',
			username:	'',
			addresses:	[],
			fullName:	{
				editing:	false,
				loading:	false,
				text:		'',
				new_text:	''
			},
			toggleEditing: function(event) {
				self.info.status.editing = !self.info.status.editing;
			},
			edit_fullName:	function() {
				self.info.fullName.editing = true;
				focus('edit_fullName');
			},
			set_fullName:	function() {
				self.info.fullName.editing = false;
				if(self.info.fullName.text !== self.info.fullName.new_text) {
					self.info.fullName.loading = true;
					$http({
						method: 'PATCH',
						url: '/v1/users/' + userService.user_id,
						headers: {
							'Content-Type': 'application/json-patch+json'
						},
						data: [
							{
								'op':		'replace',
								'path':		'/fullname',
								'value':	self.info.fullName.new_text
							}
						]
					}).then(
						function success (response) {
							debugger;
							self.info.load();
						},
						function fail (response) {
							debugger;
						}
					);
				}
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
							}
						]
					}).then(
						function success (response) {
							debugger;
							self.info.load();
						},
						function fail (response) {
							debugger;
						}
					);
				}
			},
			delete_address: function(index) {
				self.info.addresses[index].loading = true;
				$http({
					method: 'PATCH',
					url: '/v1/users/' + userService.user_id,
					headers: {
						'Content-Type': 'application/json-patch+json'
					},
					data: [
						{
							'op':		'remove',
							'path':		'/addresses/' + index + '/address',
							'value':	self.info.addresses[index].new_address
						}
					]
				}).then(
					function success (response) {
						debugger;
						self.info.load();
					},
					function fail (response) {
						debugger;
					}
				);
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
			newAddress:	function(ev) {
				$mdDialog.show({
					controller:				NewAddressController,
					templateUrl:			'src/user/new_address.html',
					parent:					angular.element(document.body),
					targetEvent:			ev,
					clickOutsideToClose:	true
				}).then(
					function (response) {
						//the dialog responded before closing
						self.info.load();
					}, function () {
						//the dialog was cancelled
					}
				);
			},
			delete: function(event) {
				$mdDialog.show(
					$mdDialog.confirm()
						.title('Are you sure?')
						.content('Are you sure you want to delete your account?')
						.targetEvent(event)
						.ok('Yes')
						.cancel('No')
						.openFrom('#info_menu')
						.closeTo('#info_menu')
				).then(
					function ok() {
						$scope.loading.is = true;
						$http({
							method: 'DELETE',
							url: '/v1/users/' + userService.user_id,
							headers: {
								'Content-Type': 'application/json'
							},
							data: {
								session_id: userService.session_id
							}
						}).then(
							// 410 is success - therefore success is a type
							// of fail so the real success is in the fail
							// block
							function success(response) {
								debugger;
							},
							function fail(response) {
								if(response.status == 410) {
									$location.path('/login');
								}
								debugger;
							}
						);
					},
					function cancel() {
						//do nothing?
					}
				)
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
						debugger;
						self.info.username = response.data.user_id;
						self.info.addresses = response.data.addresses;
						self.info.fullName.text = response.data.fullname;
						self.info.fullName.new_text = response.data.fullname;
						self.info.fullName.editing = false;
						self.info.fullName.loading = false;
						var i;
						for (i in self.info.addresses) {
							self.info.addresses[i].medium = self.info.addresses[i].medium.toUpperCase();
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
					userService.loginRedir.message = 'USER.REDIR';
					$location.path('/login');
				}
			);
		};

		self.load();
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
			$mdDialog.cancel();
		};
		self.change = function() {
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
					},
					{
						'op':		'replace',
						'path':		'/password',
						'value':	self.newPassword[0]
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

	function NewAddressController($scope, $mdDialog, $http, UserService) {
		var self = $scope.address = this;

		$scope.loading = false;
		$scope.title = 'New Address';
		$scope.errors = [];

		self.mediums = [
			{
				name:	'email',
				type:	'email'
			},
			{
				name:	'sms',
				type:	'tel'
			},
			{
				name:	'mms',
				type:	'tel'
			}
		];
		self.medium = 0;
		self.address = '';

		self.test = function() {
			debugger;
		};

		self.cancel = function() {
			$mdDialog.cancel();
		};
		self.create = function() {
			$scope.loading = true;
			$scope.title = 'Creating Address...';
			$http({
				method: 'PATCH',
				url: '/v1/users/' + UserService.user_id,
				headers: {
					'Content-Type': 'application/json-patch+json'
				},
				data: [
					{
						'op':		'add',
						'path':		'/addresses',
						'value':	{
							medium:	self.mediums[self.medium].name,
							address: self.address
						}
					}
				]
			}).then(
				function success (response) {
					debugger;
					self.info.load();
				},
				function fail (response) {
					debugger;
				}
			);
		};
	}

})();
