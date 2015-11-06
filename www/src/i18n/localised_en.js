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
			}
		}
	},
	LOGIN: {
		ERRORS: {
			"403": {
				TITLE:		"Login Failed",
				CONTENT:	"Username or password do not match."
			}
		}
	},
	ERRORS:	{
		"500": {
			TITLE:		"Something Went Wrong",
			CONTENT:	"The server has encountered an error.  " + _en.CONTACT_ADMIN
		},
		DEFAULT: {
			TITLE:		"Unknown Error",
			CONTENT:	"An unexpected response was received from the server. " + _en.CONTACT_ADMIN
		},
		CONTACT_ADMIN:	_en.CONTACT_ADMIN
	},
	USER: {
		USER:						"User",
		EMAIL_CHANGE:				"Change",
		EMAIL_DELETE:				"Delete"
		CHANGE_PASSWORD:			"Change Password",
		DELETE:						"Delete Account",
		BEGIN_EDIT:					"Begin Editing",
		STOP_EDIT:					"Stop Editing",
		FULL_NAME:					"Full Name",
		NO_FULL_NAME:				"No name set. ",
		CHANGE_NAME:				"Change Name",
		ADDRESSES:					"Contact Addresses",
	},
	GROUPS:	{
		NEW:	"New Group",
		JOIN:	"Join Group"
	},
	MEDIUMS: {
		EMAIL:		"Email",
		SMS:		"SMS",
		MMS:		"MMS"
	},
	GROUPS:			"Groups",
	NO_GROUPS:		"You are not currently in any groups.",
	OK:				"OK",
	SIGNOUT:		"Sign out"
};
