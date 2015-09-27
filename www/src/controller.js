(function(){
	
	angular
		.module('periwinkle', [])
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
				title: ''
			};
			$scope.expandMenu = {
				exists: false
			};
		}
		$scope.resetHeader = resetHeader;
	}
})();