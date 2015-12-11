// Copyright 2015 Richard Wisniewski
;(function(){
	'use strict';

	angular
	.module('user')
	.controller('EditSubscriptionsController', EditSubscriptionsController);

	function EditSubscriptionsController($scope, $mdDialog, $http, group, addresses) {
		var self = $scope.subs = this;

		$scope.loading = false;
		$scope.title = 'USER.SUBSCRIPTIONS.TITLE.MAIN';
		$scope.error = '';

		self.groupname = group.groupname;
		self.addresses = [];
		var i,n;
		for(n in addresses) {
			for(i in addresses[n]){
				self.addresses.push({
					media: n,
					address: addresses[n][i].address,
					is:	group.subscriptions.indexOf(addresses[n][i].address) !== -1
				});
			}
		}
	}

})();
