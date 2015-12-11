// Copyright 2015 Richard Wisniewski
;(function(){
	'use strict';

	angular
	.module('user')
	.controller('JoinGroupController', JoinGroupController);

	function JoinGroupController($scope, $mdDialog, $http, username) {
		var self = $scope.join = this;

		$scope.loading = false;
		$scope.title = 'Join';
		$scope.error = '';

		self.name = '';

		self.cancel = function() {
			$mdDialog.cancel();
		};
		self.change = function() {
			$scope.loading = true;
			$scope.title = 'USER.CHANGE_PASSWORD.TITLE.CREATING';
			$http({
				method: 'POST',
				url: '/v1/users/' + username + '/subscriptions',
				data: {
					group_id:	self.name
				}
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
