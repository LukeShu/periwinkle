(function(){
	'use strict';

	angular
		.module('user')
		.controller('UserController', ['$cookies', '$http', '$scope', '$interval', UserController]); 

	function UserController($cookies, $http, $scope, $interval) {
		//gives us an anchor to the outer object from within sub objects or functions
		var self = this;
		//clears the toolbar and such so we can set it up for this view
		$scope.resetHeader();
		//set up public fields
		self.username = '';
		self.password = '';
		self.email = '';
		
		self.reload = function() {
			
		};
		self.reload();
	}
})();
