// Copyright 2015 Richard Wisniewski
(function(){
	'use strict';

	angular
	.module('user')
	.controller('NewAddressController', NewAddressController);

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
			var i; var n;
			var so_index = 0;
			var list = [];
			 ;
			for (n in addresses) {
				for(i in addresses[n]) {
					var item = {
						medium:	n,
						address: addresses[n][i].address,
						sort_order: so_index
					};
					so_index++;
					list.push(item);
				}
				if(n == self.mediums[self.medium].toLowerCase()) {
					var item = {
						medium:	self.mediums[self.medium].toLowerCase(),
						address: self.address(),
						sort_order: so_index
					};
					so_index++;
					list.push(item);
				}
			}
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
