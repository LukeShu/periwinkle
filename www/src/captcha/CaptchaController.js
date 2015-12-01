// Copyright 2015 Richard Wisniewski
;(function(){
	'use strict';

	angular
	.module('captcha')
	.controller('CaptchaController', CaptchaController);

	function CaptchaController($scope, $mdDialog, $http, captcha_id) {
		var self = this;
		self.id = captcha_id;
		self.text = '';
		self.error = '';
		self.loading = false;

		self.finish = function() {
			self.loading = true;
			$http({
				method:	'POST',
				url:	'/v1/captcha/' + self.id,
				data: {
					text:	self.text
				}
			}).then(
				function success(response) {
					//store token
					$mdDialog.hide(response.data.key);
				},
				function fail(response) {
					self.error = response.data;
					self.loading = false;
				}
			);
		};
	}

})();
