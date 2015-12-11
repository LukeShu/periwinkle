// Copyright 2015 Richard Wisniewski
(function(){
	'use strict';

	angular
		.module('periwinkle')
		.controller('PeriwinkleController', ['$scope', '$http', 'UserService', '$location', '$mdDialog', '$filter', PeriwinkleController]);

	function PeriwinkleController ($scope, $http, userService, $location, $mdDialog, $filter) {
		var self = this;

		$scope.title = '';
		$scope.loading = false;

		$scope.window_title = function () {
			var title = 'Periwinkle';
			var $translate = $filter('translate');
			if($scope.title != '') {
				debugger;
				title += " &mdash; ";
				title += $translate($scope.title);
			}
			return title;
		};

		$scope.openMenu = function($mdOpenMenu, ev) {
			$scope.originalEvent = ev;
			$mdOpenMenu(ev);
		};

		$scope.logout = function () {
			$scope.loading = true;
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
						userService.reset();
						$location.path('/login').replace();
					},
					function fail(response) {
						//do work with response
						//show error to user
						$scope.loading = false;
						//show alert
					}
				);
			}
		};

		$scope.showError = function(title, body, more, from, to) {
			debugger;
			$mdDialog.show({
				controller:				ErrorDialogController,
				templateUrl:			'src/error_dialog.html',
				parent:					angular.element(document.body),
				clickOutsideToClose:	true,
				openFrom:				from,
				closeTo:				to,
				locals:	{
					'title': title,
					'body': body,
					'more': more
				}
			});
		}
	}

	function ErrorDialogController ($scope, $mdDialog, title, body, more) {
		$scope.title = title;
		$scope.body = body;
		$scope.more = more;
		$scope.showMore = false;
		$scope.details = function() {
			$scope.showMore = ! $scope.showMore;
		};
		$scope.finish = function() {
			$mdDialog.hide();
		}
	}
})();
