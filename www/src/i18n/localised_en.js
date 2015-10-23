// Copyright 2015 Richard Wisniewski
var _en = {
	CONTACT_ADMIN:	"Please contact an administrator."
};
localised.en = {
	USERNAME_EMAIL:	"Username or Email",
	USERNAME:		"Username",
	EMAIL:			"Email",
	PASSWORD:		"Password",
	LOGIN:			"Login",
	REQUIRED_FIELD:	"This field is required.",
	SIGNUP: {
		SIGNUP:						"Sign Up",
		NOT_A_USER:					"Not a user yet?",
		ALREADY_USER:				"Aready a user?",
		CONFIRM_EMAIL:				"Confirm Email",
		CONFIRM_EMAIL_NO_MATCH:		"Your email addresses do not match.",
		CONFIRM_PASSWORD:			"Confirm password",
		CONFIRM_PASSWORD_NO_MATCH:	"Your passwords do not match",
		USERNAME_INVALID:			"Usernames may only contain letters, number, _, or -.",
		EMAIL_INVALID:				"That is not a valid email address.",
		ERRORS:{
			"409": {
				TITLE:		"User Already Exists",
				CONTENT:	"Username or Email already in use."
			},
			"500": {
				TITLE:		"Something Went Wrong",
				CONTENT:	"The server has encountered an error.  " + _en.CONTACT_ADMIN
			},
		}
	},
	LOGIN: {
		ERRORS: {
			"401": {
				TITLE:		"Login Failed",
				CONTENT:	"Username or password do not match."
			}
		}
	},
	ERRORS:	{
		DEFAULT: {
			TITLE:		"Unknown Error",
			CONTENT:	"An unexpected response was received from the server. " + _en.CONTACT_ADMIN
		},
		CONTACT_ADMIN:	_en.CONTACT_ADMIN
	},
	USER: {
		USER:						"User",
		EMAIL_CHANGE:				"Change",
		CHANGE_PASSWORD:			"Change Password",
		DELETE:						"Delete Account"
	},
	GROUPS:			"Groups",
	NO_GROUPS:		"You are not currently in any groups.",
	OK:				"OK"
};