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
		self.exists_prms = [
			{text: 'Public', server: 'public'},
			{text: 'Confirmed', server: 'confirmed'},
			{text: 'Member', server: 'member'}
		];
		self.exists = 0;
		self.read_prms = [
			{text: 'Public', server: 'public'},
			{text: 'Confirmed', server: 'confirmed'},
			{text: 'Member', server: 'member'}
		];
		self.read = 2;
		self.post_prms = [
			{text: 'Public', server: 'public'},
			{text: 'Confirmed Members', server: 'confirmed'},
			{text: 'Moderator', server: 'moderator'}
		];
		self.post = 1;
		self.join_prms = [
			{text: 'Auto Allow', server: 'auto'},
			{text: 'Require Confirmation', server: 'confirm'}
		];
		self.join = 0;

		self.permissions = {
			post : {
				public: 'bounce',
				confirmed: 'moderate',
				member: 'allow'
			},
			join : {
				public: 'bounce',
				confirmed: 'moderate',
				member: 'allow'
			},
			read: {
				public: 'no',
				confirmed: 'no'
			},
			exists: {
				public: 'yes',
				confirmed: 'yes'
			}
		};

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
					'groupname': self.name,
					'existence': self.exists_prms[self.exists].server,
					'read': self.read_prms[self.read].server,
					'post': self.post_prms[self.post].server,
					'join': self.join_prms[self.join].server
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
