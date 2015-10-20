// Copyright 2015 Richard Wisniewski
(function() {
	'use strict';

	angular.module('validation.xregex', []);

	angular.module('validation.xregex').directive('xregPattern', patternTest);

	function patternTest ($parse) {
		return {
		    require: '?ngModel',
		    restrict: 'A',
		    link: function(scope, elem, attrs, ctrl) {

				ctrl.$validators.xregPattern = function(){
					var regex = XRegExp(attrs.xregPattern)
					var value = regex.test(ctrl.$viewValue) === true;
					return value;
				};
			}
		};
	}
})();
