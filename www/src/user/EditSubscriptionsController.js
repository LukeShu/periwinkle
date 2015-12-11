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
		if(group.subscriptions == null)
			group.subscriptions = [];
		var i,n;
		for(n in addresses) {
			for(i in addresses[n]){
				self.addresses.push({
					medium: n,
					address: addresses[n][i].address,
					is:	group.subscriptions.indexOf(addresses[n][i].address) !== -1
				});
			}
		}

		self.submit = function() {
			debugger;
			$scope.loading = true;
			var list = [];
			var i;
			for(i in self.addresses) {
				if(self.addresses[i].is){
					list.push({
						medium: self.addresses[i].medium,
						address: self.addresses[i].address
					});
				}
			}
			$http({
				method:	'PATCH',
				url:	'/v1/groups/' + group.groupname,
				headers: {
					'Content-Type': 'application/json-patch+json'
				},
				data: [
					{
						'op':		'replace',
						'path':		'/subscriptions',
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
			)
		};

		$scope.cancel = function() {
			$mdDialog.cancel();
		};
	}

})();
