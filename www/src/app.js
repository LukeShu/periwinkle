var periwinkleApp = angular.module('periwinkleApp', [
	'ngRoute',
	'periwinkleControllers',
	'ngMaterial',
	'ngCookies'
]);

periwinkleApp.config(function($mdThemingProvider){
	$mdThemingProvider.theme('default')
		.primaryPalette('deep-purple')
		.accentPalette('cyan');
});

periwinkleApp.config(['$routeProvider',
	function($routeProvider) {
		$routeProvider.
		when('/login', {
			templateUrl:	'src/login/login.html',
			controller:		'LoginController'
		}).
		when('/dashboard', {
			templateUrl:	'src/dashboard/dashboard.html',
			controller:		'DashboardController'
		}).
		when('/messaages/:groupid', {
			templateURL:	'src/messages/messages.html',
			controller:		'MessaagesController'
		}).
		when('/settings', {
			templateUrl:	'src/settings/settings.html',
			controller:		'SettingsController'
		}).
		otherwise({
			//redirectTo:	'/login'
		})
	}
]);