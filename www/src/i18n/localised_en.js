// Copyright 2015 Richard Wisniewski
var _en = {
	CONTACT_ADMIN:	"Please contact an administrator."
};
localised.en = {
	LOGIN: {
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
			LOGIN:			"Login",
			MESSAGE:		"You are not logged in.",
			ERRORS: {
				"403": {
					TITLE:		"Login Failed",
					CONTENT:	"Username or password do not match."
				}
			}
		}
	},
	USER: {
		INFO: {
			USER:						"User",
			EMAIL_CHANGE:				"Change",
			EMAIL_DELETE:				"Delete",
			CHANGE_PASSWORD:			"Change Password",
			DELETE:						"Delete Account",
			BEGIN_EDIT:					"Begin Editing",
			STOP_EDIT:					"Stop Editing",
			FULL_NAME:					"Full Name",
			NO_FULL_NAME:				"No name set. ",
			CHANGE_NAME:				"Change Name",
			ADDRESSES:					"Contact Addresses",
			SAVING:						"Saving..."
		},
		GROUPS:	{
			NEW:	"New Group",
			JOIN:	"Join Group",
			GROUPS:			"Groups",
			NO_GROUPS:		"You are not currently in any groups."
		},
		NEW_GROUP:	{
			TITLE: {
				MAIN:		"New Group",
				CREATING:	"Creating Group...",
				FAIL:		"Fail"
			},
			GROUP_NAME:		"New Group Name"
		},
		CHANGE_PASSWORD:	{
			TITLE: {
				MAIN:		"Change Password",
				CREATING:	"Changing Password...",
				FAIL:		"Fail"
			},
			ERROR409:	{
				TITLE:		"Password Change Failed",
				CONTENT:	"You did not enter your current password correctly"
			},
			OLD_PASSWORD:	"Current Password",
			NEW_PASSWORD:	"New Password",
			CONFIRM:		"Confirm Password"
		},
		NEW_ADDRESS: {
			TITLE: {
				MAIN:		"New Address",
				CREATING:	"Adding Address",
				FAIL:		"Fail"
			},
			MEDIUM:		"Medium",
			ADDRESS:	"Address"
		},
		REDIR:	"You will be redirected back to your user once you log in. "
	},
	GENERAL: {
		USERNAME_EMAIL:	"Username or Email",
		USERNAME:		"Username",
		EMAIL:			"Email",
		PASSWORD:		"Password",
		FORM: {
			ERROR: {
				REQUIRED_FIELD:	"This field is required."
			},
			RESPONSE: {
				OK:			"OK",
				CANCEL:		"Cancel",
				CHANGE:		"Change",
				CREATE:		"Create",
				ADD:		"Add"
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
		MEDIUMS: {
			EMAIL:		"Email",
			SMS:		"SMS",
			MMS:		"MMS"
		},
		SIGNOUT:		"Sign out"
	}
};
