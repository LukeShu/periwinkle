// Copyright 2015 Richard Wisniewski
;(function() {
	'use strict';
	var lang = "en";

	try {
		lang = navigator.language.substring(0,2);
	} catch(err) {}

	angular.module( 'periwinkle.i18n', [])
	.constant('lang', lang);
})();
