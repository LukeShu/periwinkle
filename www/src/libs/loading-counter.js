(function() {
	'use strict';
	
	angular.module('loading-counter', [])
		.service('loadingService', ['$scope', LoadingCounterService]);
		
	function LoadingCounterService($scope) {
		var count = 0;
		var self = this;
		
		self.setCount = function (c) {
			count = c;
		}
		self.start = function() {
			count++;
		}
		self.finish = function() {
			count--;
			if(count <= 0) {
				//close the loader
			}
		}
	}
	
})();