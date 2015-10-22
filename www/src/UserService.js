// Copyright 2015 Richard Wisniewski
(function(){
	'use strict';

	angular
		.module('periwinkle.UserService', [])
		.service('UserService', ['$scope', '$http', '$translate', UserService]);
		
	function UserService ($scope, $http, $translate) {
		self = this;
		self.username = '';
		self.session_id = '';
		self.language = '';
	}
})();