// Copyright 2015 Richard Wisniewski
;(function(){
	'use strict';

	angular
	.module('user')
	.controller('ChangePasswordController', ChangePasswordController);

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

})();
