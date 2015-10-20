// Copyright 2015 Richard Wisniewski
(function(){
	'use strict';

	angular
		.module('dashboard')
		.controller('DashboardController', ['$cookies', '$http', '$scope', '$interval', DashboardController]); 

	function DashboardController($cookies, $http, $scope, $interval) {
		$scope.resetHeader();
	}
})();
