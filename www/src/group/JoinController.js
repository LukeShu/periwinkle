// Copyright 2015 Richard Wisniewski
;(function(){
	'use strict';

	angular
	.module('group')
	.controller('JoinController', JoinController);

	function JoinController($scope, $mdDialog, $http, groupname) {
		var self = $scope.join = this;

		$scope.loading = false;
		$scope.title = 'Join ' + groupname;
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
				url: '/v1/users/' + self.name + '/subscriptions',
				data: {
					group_id:	groupname
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
