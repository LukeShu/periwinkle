(function(){
  angular
		.module('login')
		.controller('LoginController', ['$cookies', '$http', '$scope', LoginController]); 

	function LoginController($cookies, $http, $scope) {
		this.username = '';
		this.password = '';
		var sessionID = $cookies.get("sessionID");
		if(sessionID !== undefined) {
			//show loading spinner
			//validate sessionID
			//if valid take down spinner
			//if invalid (or error)
				//show login screen
				//clear cookie
				//show warning bar
		}
		
		function login() {
			//http login api call
		}
	}
})();
