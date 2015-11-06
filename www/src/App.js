// Copyright 2015 Richard Wisniewski
// Copyright 2015 Luke Shumaker
(function(){
	'use strict';

	angular.module('periwinkle', [
		'ngRoute',
		'ngMaterial',
		'ngMessages',
		'ngCookies',
		'pascalprecht.translate',
		'validation.match',
		'ngSanitize',
		//periwinkle modules
		'validation.anti-match',
		'focusOn',
		'validation.xregex',
		'periwinkle.UserService',
		'login',
		'dashboard',
		'user'
	])

	.config(function($mdThemingProvider){
		$mdThemingProvider.theme('default')
		.primaryPalette('deep-purple')
		.accentPalette('amber');
	})

	.factory('httpRequestInterceptor', ['$cookies', function ($cookies) {
		return {
			request: function (config) {
				config.headers['X-XSRF-TOKEN'] = $cookies.get('app_set_session_id');
				//debugger;
				return config;
			}
		};
	}])

	.config(['$httpProvider', function ($httpProvider) {
		$httpProvider.interceptors.push('httpRequestInterceptor');
		$httpProvider.defaults.headers.patch = {
		    'Content-Type': 'application/json;charset=utf-8'
		};
	}])

	.config(['$translateProvider', function($translateProvider) {
		$translateProvider
			.translations('en', localised.en)
			.translations('it', localised.it)
			.translations('es', localised.es);
		$translateProvider.fallbackLanguage('en');
		$translateProvider.use(lang);
		$translateProvider.useSanitizeValueStrategy('escape');
	}])

	.config(['$routeProvider', '$locationProvider',
		function($routeProvider, $locationProvider) {
			$locationProvider
			.hashPrefix('!');

			$routeProvider.
			when('/login', {
				templateUrl:	'src/login/login.html',
				controller:		'LoginController as login'
			}).
			when('/dash', {
				templateUrl:	'src/dashboard/dashboard.html',
				controller:		'DashboardController as dash'
			}).
			when('/user', {
				templateUrl:	'src/user/user.html',
				controller:		'UserController as user'
			}).
			otherwise({
				redirectTo:	'/login'
			});
		}
	])

	.filter('escapeHTML', function() {
		var div = document.createElement('div');
		return function(text) {
			div.textContent = text;
			return div.innerHTML;
		}
	});

	//user this plugin instead of the browsers regex (for unicode support)
	XRegExp.install('natives');
})();
