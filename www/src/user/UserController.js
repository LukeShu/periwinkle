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
		$scope.title = "USER.INFO.USER";
		$scope.loading = false;

		self.info = {
			status: {
				loading: 	true,
				error:		'',
				editing:	false
			},
			title:		'USER.INFO.USER',
			username:	'',
			addresses:	{},
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
							$scope.loading = false;
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
			edit_address:	function(name, index) {
				self.info.addresses[name][index].new_address = self.info.addresses[name][index].address;
				self.info.addresses[name][index].editing = true;
				focus('edit_address');
			},
			save_edit_address:	function(name, index) {
				self.info.addresses[name][index].editing = false;
				if(self.info.addresses[name][index].address !== self.info.addresses[name][index].new_address) {
					self.info.addresses[name][index].loading = true;
					var i; var n;
					var list = [];
					for (n in self.info.addresses) {
						for(i in self.info.addresses[n]) {
							var item = {
								medium:	n,
								address: self.info.addresses[n][i].address
							};
							if(i == index && n == name) {
								item.address = self.info.addresses[n][i].new_address;
							}
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
			address_orderChanged: function() {
				self.info.status.loading = true;
				var i; var n;
				var list = [];
				for (n in self.info.addresses) {
					for(i in self.info.addresses[n]) {
						var item = {
							medium:	n,
							address: self.info.addresses[n][i].address
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
							'value':	list
						}
					]
				}).then(
					function success (response) {
						self.info.load();
					},
					function fail (response) {
						self.info.status.loading = false;
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
			delete_address: function(name, index) {
				self.info.addresses[name][index].loading = true;
				var i; var n;
				var list = [];
				for (n in self.info.addresses) {
					for(i in self.info.addresses[n]) {
						if(i != index && n != name) {
							var item = {
								medium:	n,
								address: self.info.addresses[n][i].address
							};
							list.push(item);
						}
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
					controller:				'ChangePasswordController',
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
					controller:				'NewAddressController',
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
						$scope.loading = true;
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
						//self.info.addresses = response.data.addresses;
						self.info.fullName.text = response.data.fullname;
						self.info.fullName.new_text = response.data.fullname;
						self.info.fullName.editing = false;
						self.info.fullName.loading = false;
						var i;
						self.info.addresses = {
							email: [],
							sms: [],
							mms: []
						};
						for (i in response.data.addresses) {
							self.info.addresses[response.data.addresses[i].medium].push({
								address:		response.data.addresses[i].address,
								new_address:	'',
								editing:		false,
								loading:		false
							});
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
					controller:				'NewGroupController',
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
			$scope.loading = true;
			userService.validate(
				function success() {
					$scope.loading = false;
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

})();
