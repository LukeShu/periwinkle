(function(){
	'use strict';

	angular
		.module('periwinkle')
		.controller('PeriwinkleController', ['$scope', PeriwinkleController]);
		
	function PeriwinkleController ($scope) {
		var resetHeader = function() {
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
		}
		$scope.resetHeader = resetHeader;
	}
})();
