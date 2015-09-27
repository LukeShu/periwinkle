(function(){
  angular
		.module('login')
		.controller('LoginController', ['$cookies', '$http', '$scope', LoginController]); 

	function LoginController($cookies, $http, $scope) {
		//gives us an anchor to the outer object from within sub objects or functions
		var self = this;
		//clears the toolbar and such so we can set it up for this view
		$scope.resetHeader();
		//adds the public fields and methods of this object to the model ($scope)
		$scope.login = self;
		//set up public fields
		self.username = '';
		self.password = '';
		self.isSignup = false;
		//prep the toolbar
		$scope.toolbar.title = 'Login';
		/*$scope.toolbar.buttons = [{
			label: "Signup",
			img_src: "assets/svg/phone.svg",
		}];
		$scope.toolbar.onclick = function(index) {
			if(index == 0) {
				$scope.login.togleSignup();
			}
		} */
		
		//check if already loged in
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
		
		//public functions
		this.login = function() {
			//http login api call
		}
		
		this.signup = function() {
			//http signup api call
		}
		
		this.togleSignup = function () {
			self.isSignup = !self.isSignup;
			if(self.isSignup) {
				$scope.toolbar.title = 'Sign Up';
			} else {
				$scope.toolbar.title = 'Login';
			}
		}
	}
})();
