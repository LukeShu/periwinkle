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
			title:		'USER.INFO.USER',
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
							self.info.load();
						},
						function fail (response) {
							debugger;
							var status_code = response.status;
							var reason = response.data;
							$scope.loading.is = false;
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
					var i;
					var list = [];
					for (i in self.info.addresses) {
						var item = {
							medium:	self.info.addresses[i].medium.toLowerCase(),
							address: self.info.addresses[i].address
						};
						if(i == index) {
							item.address = self.info.addresses[i].new_address;
						}
						list.push(item);
					}
					$http({
						method: 'PATCH',
						url: '/v1/users/' + userService.user_id,
						headers: {
							'Content-Type': 'application/json-patch+json'
						},
						data: [
							{
								'op':		'replace',
								'path':		'/addresses',
								'value':	list
							}
						]
					}).then(
						function success (response) {
							self.info.load();
						},
						function fail (response) {
							debugger;
							var status_code = response.status;
							var reason = response.data;
							//show alert
							switch(status_code){
								case 500:
									$scope.showError('GENERAL.ERRORS.500.TITLE', 'GENERAL.ERRORS.500.CONTENT', reason, 'body', 'body');
									break;
								default:
									$scope.showError('GENERAL.ERRORS.DEFAULT.TITLE', 'GENERAL.ERRORS.DEFAULT.CONTENT', reason, 'body', 'body');
							}
						}
					);
				}
			},
			delete_address: function(index) {
				self.info.addresses[index].loading = true;
				var i;
				var list = [];
				for (i in self.info.addresses) {
					if(i != index) {
						var item = {
							medium:	self.info.addresses[i].medium.toLowerCase(),
							address: self.info.addresses[i].address
						};
						list.push(item);
					}
				}
				$http({
					method: 'PATCH',
					url: '/v1/users/' + userService.user_id,
					headers: {
						'Content-Type': 'application/json-patch+json'
					},
					data: [
						{
							'op':		'replace',
							'path':		'/addresses',
							/* TODO: value should send in an array [] of all of the
							 * values that are currently in it, minus the one below
							 * that the user clicked to delete
							 */
							'value':	list
						}
					]
				}).then(
					function success (response) {
						self.info.load();
					},
					function fail (response) {
						debugger;
						var status_code = response.status;
						var reason = response.data;
						//show alert
						switch(status_code){
							case 500:
								$scope.showError('GENERAL.ERRORS.500.TITLE', 'GENERAL.ERRORS.500.CONTENT', reason, 'body', 'body');
								break;
							default:
								$scope.showError('GENERAL.ERRORS.DEFAULT.TITLE', 'GENERAL.ERRORS.DEFAULT.CONTENT', reason, 'body', 'body');
						}
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
				}).then(
					function hide (response) {
						if(response !== "success") {
							//errors
							var status_code = response.status;
							var reason = response.data;
							//show alert
							switch(status_code){
								case 500:
									$scope.showError('GENERAL.ERRORS.500.TITLE', 'GENERAL.ERRORS.500.CONTENT', reason, '#user-info-menu', '#user-info-menu');
									break;
								default:
									$scope.showError('GENERAL.ERRORS.DEFAULT.TITLE', 'GENERAL.ERRORS.DEFAULT.CONTENT', reason, '#user-info-menu', '#user-info-menu');
							}
						} else {
							//succeeded
						}
					},
					function cancel() {
						//Do nothing?
					}
				);
			},
			newAddress:	function(ev) {
				$mdDialog.show({
					controller:				NewAddressController,
					templateUrl:			'src/user/new_address.html',
					parent:					angular.element(document.body),
					targetEvent:			ev,
					clickOutsideToClose:	true,
					locals:	{
						addresses: self.info.addresses
					}
				}).then(
					function (response) {
						//the dialog responded before closing
						if(response !== "success") {
							//errors
							var status_code = response.status;
							var reason = response.data;
							//show alert
							switch(status_code){
								case 500:
									$scope.showError('GENERAL.ERRORS.500.TITLE', 'GENERAL.ERRORS.500.CONTENT', reason, '#new-address-fab', '#new-address-fab');
									break;
								default:
									$scope.showError('GENERAL.ERRORS.DEFAULT.TITLE', 'GENERAL.ERRORS.DEFAULT.CONTENT', reason, '#new-address-fab', '#new-address-fab');
							}
						} else {
							//succeeded
							self.info.load();
						}
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
							function success(response) {
								$location.path('/login');
							},
							function fail(response) {
								var status_code = response.status;
								var reason = response.data;
								//show alert
								switch(status_code){
									case 500:
										$scope.showError('GENERAL.ERRORS.500.TITLE', 'GENERAL.ERRORS.500.CONTENT', reason, '#user-info-menu', '#user-info-menu');
										break;
									default:
										$scope.showError('GENERAL.ERRORS.DEFAULT.TITLE', 'GENERAL.ERRORS.DEFAULT.CONTENT', reason, '#user-info-menu', '#user-info-menu');
								}
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
						var status_code = response.status;
						var reason = response.data;
						//show alert
						switch(status_code){
							case 500:
								$scope.showError('GENERAL.ERRORS.500.TITLE', 'GENERAL.ERRORS.500.CONTENT', reason, 'body', 'body');
								break;
							default:
								$scope.showError('GENERAL.ERRORS.DEFAULT.TITLE', 'GENERAL.ERRORS.DEFAULT.CONTENT', reason, 'body', 'body');
						}
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
						if(response !== "success") {
							//errors
							var status_code = response.status;
							var reason = response.data;
							//show alert
							switch(status_code){
								case 500:
									$scope.showError('GENERAL.ERRORS.500.TITLE', 'GENERAL.ERRORS.500.CONTENT', reason, '#new-address-fab', '#new-address-fab');
									break;
								default:
									$scope.showError('GENERAL.ERRORS.DEFAULT.TITLE', 'GENERAL.ERRORS.DEFAULT.CONTENT', reason, '#new-address-fab', '#new-address-fab');
							}
						} else {
							//succeeded
							self.groups.load();
						}
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
						self.groups.status.loading = false;
					},
					function fail(response) {
						self.groups.status.loading = false;
						var status_code = response.status;
						var reason = response.data;
						//show alert
						switch(status_code){
							case 500:
								$scope.showError('GENERAL.ERRORS.500.TITLE', 'GENERAL.ERRORS.500.CONTENT', reason, 'body', 'body');
								break;
							default:
								$scope.showError('GENERAL.ERRORS.DEFAULT.TITLE', 'GENERAL.ERRORS.DEFAULT.CONTENT', reason, 'body', 'body');
						}
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
					var status_code = response.status;
					var reason = response.data;
					//show alert
					switch(status_code){
						case 500:
							$scope.showError('GENERAL.ERRORS.500.TITLE', 'GENERAL.ERRORS.500.CONTENT', reason, 'body', 'body');
							break;
						default:
							$scope.showError('GENERAL.ERRORS.DEFAULT.TITLE', 'GENERAL.ERRORS.DEFAULT.CONTENT', reason, 'body', 'body');
					}
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
		$scope.title = 'USER.NEW_GROUP.TITLE.MAIN';
		$scope.error = '';

		self.name = '';

		self.cancel = function() {
			$mdDialog.cancel();
		};
		self.create = function() {
			$scope.loading = true;
			$scope.title = 'USER.NEW_GROUP.TITLE.CREATING';
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
					$mdDialog.hide("success");
				},
				function fail(response) {
					if(response.status == 409) {
						$scope.loading = false;
						$scope.title = 'USER.NEW_GROUP.ERRORS.409.TITLE';
						$scope.error = 'USER.NEW_GROUP.ERRORS.409.CONTENT';
					} else {
						$mdDialog.hide(response);
					}
				}
			);
		};
	}

	function ChangePasswordController($scope, $mdDialog, $http, UserService) {
		var self = $scope.password = this;

		$scope.loading = false;
		$scope.title = 'USER.CHANGE_PASSWORD.TITLE.MAIN';
		$scope.error = '';

		self.oldPassword = '';
		self.newPassword = ['',''];
		var isCancel = false;

		self.cancel = function() {
			$mdDialog.cancel();
		};
		self.change = function() {
			$scope.loading = true;
			$scope.title = 'USER.CHANGE_PASSWORD.TITLE.CREATING';
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
					$mdDialog.hide("success");
				},
				function fail(response) {
					if(response.status == 409) {
						$scope.loading = false;
						$scope.title = 'USER.CHANGE_PASSWORD.ERRORS.409.TITLE';
						$scope.error = 'USER.CHANGE_PASSWORD.ERRORS.409.CONTENT';
						self.oldPassword = '';
						self.newPassword = ['', ''];
					} else {
						$mdDialog.hide(response);
					}
				}
			);
		};
	}

	function NewAddressController($scope, $mdDialog, $http, UserService, addresses) {
		var self = $scope.address = this;

		$scope.loading = false;
		$scope.title = 'USER.NEW_ADDRESS.TITLE.MAIN';
		$scope.errors = [];

		self.mediums = [
			'EMAIL',
			'SMS',
			'MMS'
		];
		self.medium = 0;
		self.medium_type = function() {
			switch(self.medium) {
				case 0:
					return 'email';
				case 1:
				case 2:
					return 'tel';
				default:
					//error
			}
		};
		self.email_address = '';
		self.tel_address = '';
		self.address = function() {
			switch(self.medium) {
				case 0:
					return self.email_address;
				case 1:
				case 2:
					return self.tel_address;
				default:
					//error
			}
		}

		self.cancel = function() {
			$mdDialog.cancel();
		};
		self.create = function() {
			$scope.loading = true;
			$scope.title = 'USER.NEW_ADDRESS.TITLE.CREATING';
			var i;
			var list = [];
			for (i in addresses) {
				var item = {
					medium:	addresses[i].medium.toLowerCase(),
					address: addresses[i].address
				};
				list.push(item);
			}
			var item = {
				medium:	self.mediums[self.medium].name.toLowerCase(),
				address: self.address()
			};
			list.push(item);
			$http({
				method: 'PATCH',
				url: '/v1/users/' + UserService.user_id,
				headers: {
					'Content-Type': 'application/json-patch+json'
				},
				data: [
					{
						'op':		'replace',
						'path':		'/addresses',
						'value':	list
					}
				]
			}).then(
				function success (response) {
					$mdDialog.hide("success");
				},
				function fail (response) {
					$mdDialog.hide(response);
				}
			);
		};
	}

})();
