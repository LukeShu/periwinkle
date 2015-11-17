// Copyright 2015 Richard Wisniewski
(function(){
	'use strict';

	angular
	.module('user')
	.controller('NewGroupController', NewGroupController);

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

})();
