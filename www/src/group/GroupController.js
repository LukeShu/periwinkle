// Copyright 2015 Richard Wisniewski
;(function(){
	'use strict';

	angular
	.module('group')
	.controller('GroupController', ['$scope', '$http', '$routeParams', 'UserService', GroupController]);

	function GroupController($scope, $http, $routeParams, userService) {
		var self = $scope.subs = this;

		$scope.permissions = self.groupname = $routeParams.group;
		$scope.title = self.groupname;
		$scope.loading = false;

		self.permissions = {
			post : {
				public: 'bounce',
				confirmed: 'moderate',
				member: 'accept'
			},
			join : {
				public: 'bounce',
				confirmed: 'moderate',
				member: 'accept'
			},
			read: {
				public: 'no',
				confirmed: 'no'
			},
			exists: {
				public: 'yes',
				confirmed: 'yes'
			},
		};
		self.permissions_status = {
			loading: false,
			load:	function() {
				self.permissions_status.loading = true;
				$http({
					method:	'GET',
					url:	'/v1/groups/' + self.groupname
				}).then(
					function success(response) {
						self.permissions.exists = response.data.existence;
						self.permissions.read = response.data.read;
						self.permissions.join = response.data.join;
						self.permissions.post = response.data.post;
						self.permissions_status.loading = false;
					},
					function fail(response) {
						//show error to user
						self.permissions_status.loading = false;
						var status_code = response.status;
						var reason = response.data;
						//show alert
						switch(status_code){
							case 500:
								$scope.showError('GENERAL.ERRORS.500.TITLE', 'GENERAL.ERRORS.500.CONTENT', reason, 'body', 'body');
								break;
							default:
								$scope.showError('GENERAL.ERRORS.DEFAULT.TITLE', 'GENERAL.ERRORS.DEFAULT.CONTENT', reason, 'body', 'body');
						}
					}
				)
			},
			submit:	function() {
				debugger;
				self.permissions_status.loading = true;
				$http({
					method: 'PATCH',
					url: '/v1/groups/' + self.groupname,
					headers: {
						'Content-Type': 'application/json-patch+json'
					},
					data: [
						{
							'op':		'replace',
							'path':		'/existence',
							'value':	self.permissions.exists
						},
						{
							'op':		'replace',
							'path':		'/read',
							'value':	self.permissions.read
						},
						{
							'op':		'replace',
							'path':		'/post',
							'value':	self.permissions.post
						},
						{
							'op':		'replace',
							'path':		'/join',
							'value':	self.permissions.join
						}
					]
				}).then(
					function success (response) {
					},
					function fail (response) {
					}
				);
			}
		};
		self.memmbers = [];
		self.members_status = {
			loading: false,
			load: function() {

			},
			'new': function() {

			}
		};
		self.addresses = {
			email:	[],
			sms:	[],
			mms:	[]
		};
		self.addresses_status = {
			loading: false,
			load: function() {
				self.addresses_status.loading = true;
				debugger;
				$http({
					method: 'GET',
					url: '/v1/users/' + userService.user_id,
					headers: {
						'Content-Type': 'application/json'
					}
				}).then(
					function success(response) {
						self.addresses.email = [];
						self.addresses.sms = [];
						self.addresses.mms = [];
						if(response.data.addresses != null && response.data.addresses.length > 0) {
							response.data.addresses.sort(function(a, b) {
								return a.sort_order - b.sort_order;
							});
							var i;
							for (i in response.data.addresses) {
								self.addresses[response.data.addresses[i].medium].push({
									address:		response.data.addresses[i].address,
									is:			 	false
								});
							}
							//TODO: get subscription list
						}
						self.addresses_status.loading = false;
					},
					function fail(response) {
						//show error to user
						self.info.status.loading = false;
						var status_code = response.status;
						var reason = response.data;
						//show alert
						switch(status_code){
							case 500:
								$scope.showError('GENERAL.ERRORS.500.TITLE', 'GENERAL.ERRORS.500.CONTENT', reason, 'body', 'body');
								break;
							default:
								$scope.showError('GENERAL.ERRORS.DEFAULT.TITLE', 'GENERAL.ERRORS.DEFAULT.CONTENT', reason, 'body', 'body');
						}
					}
				);
			},
			submit:	function() {

			}
		};

		self.load = function() {
			$scope.loading = true;
			userService.validate(
				function success() {
					$scope.loading = false;
					self.permissions_status.load();
					self.members_status.load();
					self.addresses_status.load();
				},
				function fail(status) {
					debugger;
					var status_code = response.status;
					var reason = response.data;
					//show alert
					switch(status_code){
						case 500:
							$scope.showError('GENERAL.ERRORS.500.TITLE', 'GENERAL.ERRORS.500.CONTENT', reason, 'body', 'body');
							break;
						default:
							$scope.showError('GENERAL.ERRORS.DEFAULT.TITLE', 'GENERAL.ERRORS.DEFAULT.CONTENT', reason, 'body', 'body');
					}
				},
				function noSession_cb() {
					userService.loginRedir.has = true;
					userService.loginRedir.path = $location.path();
					userService.loginRedir.message = 'USER.REDIR';
					$location.path('/login');
				}
			);
		};
		self.load();
	}

	//from internet
	function deepCompare () {
	  var i, l, leftChain, rightChain;

	  function compare2Objects (x, y) {
	    var p;

	    // remember that NaN === NaN returns false
	    // and isNaN(undefined) returns true
	    if (isNaN(x) && isNaN(y) && typeof x === 'number' && typeof y === 'number') {
	         return true;
	    }

	    // Compare primitives and functions.
	    // Check if both arguments link to the same object.
	    // Especially useful on step when comparing prototypes
	    if (x === y) {
	        return true;
	    }

	    // Works in case when functions are created in constructor.
	    // Comparing dates is a common scenario. Another built-ins?
	    // We can even handle functions passed across iframes
	    if ((typeof x === 'function' && typeof y === 'function') ||
	       (x instanceof Date && y instanceof Date) ||
	       (x instanceof RegExp && y instanceof RegExp) ||
	       (x instanceof String && y instanceof String) ||
	       (x instanceof Number && y instanceof Number)) {
	        return x.toString() === y.toString();
	    }

	    // At last checking prototypes as good a we can
	    if (!(x instanceof Object && y instanceof Object)) {
	        return false;
	    }

	    if (x.isPrototypeOf(y) || y.isPrototypeOf(x)) {
	        return false;
	    }

	    if (x.constructor !== y.constructor) {
	        return false;
	    }

	    if (x.prototype !== y.prototype) {
	        return false;
	    }

	    // Check for infinitive linking loops
	    if (leftChain.indexOf(x) > -1 || rightChain.indexOf(y) > -1) {
	         return false;
	    }

	    // Quick checking of one object beeing a subset of another.
	    // todo: cache the structure of arguments[0] for performance
	    for (p in y) {
	        if (y.hasOwnProperty(p) !== x.hasOwnProperty(p)) {
	            return false;
	        }
	        else if (typeof y[p] !== typeof x[p]) {
	            return false;
	        }
	    }

	    for (p in x) {
	        if (y.hasOwnProperty(p) !== x.hasOwnProperty(p)) {
	            return false;
	        }
	        else if (typeof y[p] !== typeof x[p]) {
	            return false;
	        }

	        switch (typeof (x[p])) {
	            case 'object':
	            case 'function':

	                leftChain.push(x);
	                rightChain.push(y);

	                if (!compare2Objects (x[p], y[p])) {
	                    return false;
	                }

	                leftChain.pop();
	                rightChain.pop();
	                break;

	            default:
	                if (x[p] !== y[p]) {
	                    return false;
	                }
	                break;
	        }
	    }

	    return true;
	  }

	  if (arguments.length < 1) {
	    return true; //Die silently? Don't know how to handle such case, please help...
	    // throw "Need two or more arguments to compare";
	  }

	  for (i = 1, l = arguments.length; i < l; i++) {

	      leftChain = []; //Todo: this can be cached
	      rightChain = [];

	      if (!compare2Objects(arguments[0], arguments[i])) {
	          return false;
	      }
	  }

	  return true;
	}
})();
