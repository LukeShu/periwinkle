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
		'angular-sortable-view',
		//periwinkle modules
		'periwinkle.i18n',
		'validation.anti-match',
		'focusOn',
		'validation.xregex',
		'periwinkle.UserService',
		'login',
		'dashboard',
		'user',
		'messages',
		'messages.message',
		'messages.thread'
	])

	.config(function($mdThemingProvider){
		$mdThemingProvider.theme('default')
		.primaryPalette('purple')
		.accentPalette('yellow');
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
		$httpProvider.defaults.headers.common['Accept'] = "application/json, text/plain;q=0.9, */*;q=0.8";
	}])

	.config(['$translateProvider', 'i18n_en', 'i18n_it', 'lang', function($translateProvider, i18n_en, i18n_it, lang) {
		$translateProvider
			.translations('en', i18n_en)
			.translations('it', i18n_it);
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
			when('/messages', {
				templateUrl:	'src/messages/messages.html',
				controller:		'MessagesController as messages'
			}).
			when('/messages/:group', {
				templateUrl:	'src/messages/thread/thread.html',
				controller:		'ThreadController as thread'
			}).
			when('/messages/:group/:message', {
				templateUrl:	'src/messages/message/message.html',
				controller:		'MessageController as message'
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
