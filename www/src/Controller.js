// Copyright 2015 Richard Wisniewski
(function(){
	'use strict';

	angular
		.module('periwinkle')
		.controller('PeriwinkleController', ['$scope', 'UserService', PeriwinkleController]);
		
	function PeriwinkleController ($scope, userService) {
		var reset = function() {
			$scope.sidenav = {
				exists: false,
				items: [],
				selected: NaN
			};
			$scope.toolbar = {
				exists: true,
				title: '',
				buttons: [],
				onclick: function(){}
			};
			$scope.expandMenu = {
				exists: false
			};
			$scope.loading = {
				is:	false
			};
			$scope.originalEvent = null;
		}
		$scope.reset = reset;
		
		$scope.openMenu = function($mdOpenMenu, ev) {
			$scope.originalEvent = ev;
			$mdOpenMenu(ev);
		};
		
		$scope.logout = function () {
			$scope.loading.is = true;
			if(userService.session_id && userService.session_id != "") {
				$http({
					method: 'DELETE',
					url: '/v1/session',
					headers: {
						'Content-Type': 'application/json'
					},
					data: {
						session_id: userService.session_id
					}
				}).then(
					function success(response) {
						//do work with response
						debugger;
						userService.reset();
						$location.path('/login').replace();
					},
					function fail(response) {
						//do work with response
						//show error to user
						debugger;
						$scope.loading.is = false;
						//show alert
					}
				);
			}
		};
	}
})();
