var periwinkleApp = angular.module('periwinkleApp', [
	'ngRoute',
	'periwinkleControllers'
]);

periwinkleApp.config(['$routeProvider',
	function($routeProvider) {
		$routeProvider.
		when('/login', {
			templateUrl:	'templates/login.html',
			controller:		'LoginController'
		}).
		when('/dashboard', {
			templateUrl:	'templates/dashboard.html',
			controller:		'DashboardController'
		}).
		when('/messaages/:groupid', {
			templateURL:	'templates/messages.html',
			controller:		'MessaagesController'
		}).
		when('/settings', {
			templateUrl:	'templates/settings.html',
			controller:		'SettingsController'
		}).
		otherwise({
			redirectTo:	'/login'
		})
	}
]});