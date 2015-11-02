// Copyright 2015 Richard Wisniewski
var _en = {
	CONTACT_ADMIN:	"Please contact an administrator."“请联系管理人员”
};
localised.en = {
	USERNAME_EMAIL:	“用户名电子邮件地址”,
	USERNAME:		“用户名”,
	EMAIL:			“电子邮件地址”,
	PASSWORD:		“密码”,
	LOGIN:			“登录”,
	REQUIRED_FIELD:	“此处必填”,
	SIGNUP: {
		SIGNUP:						“注册”,
		NOT_A_USER:					“还不是用户？”,
		ALREADY_USER:				“已经是用户？”,
		CONFIRM_EMAIL:				“确认邮箱”,
		CONFIRM_EMAIL_NO_MATCH:		“您的邮件地址不符”,
		CONFIRM_PASSWORD:			“确认密码”,
		CONFIRM_PASSWORD_NO_MATCH:	"Your passwords do not match",“您的密码不符”
		USERNAME_INVALID:			"Usernames may only contain letters, number, _, or -.",“用户名只能包含字母，数字，"_",或者"-"."
		EMAIL_INVALID:				"That is not a valid email address.","这不是一个有效的电子邮件地址."
		ERRORS:{
			"409": {
				TITLE:		"User Already Exists","该用户已经存在"
				CONTENT:	"Username or Email already in use.""用户名或者电子邮件已经被使用”
			}
		}
	},
	LOGIN: {
		ERRORS: {
			"403": {
				TITLE:		"Login Failed","登陆失败”
				CONTENT:	"用户名或者密码不匹配"
			}
		}
	},
	ERRORS:	{
		"500": {
			TITLE:		"Something Went Wrong","出现了一个错误"
			CONTENT:	"The server has encountered an error.  " + _en.CONTACT_ADMIN,"服务器遇到一个问题"
		},
		DEFAULT: {
			TITLE:		"Unknown Error",“不明错误”
			CONTENT:	"An unexpected response was received from the server. " + _en.CONTACT_ADMIN,"服务器给出了一个非预期的反应"
		},
		CONTACT_ADMIN:	_en.CONTACT_ADMIN
	},
	USER: {
		USER:						"User","用户"
		EMAIL_CHANGE:				"Change","变更"
		CHANGE_PASSWORD:			"Change Password","变更密码"
		DELETE:						"Delete Account","删除账户"
	},
	GROUPS:			"Groups","群"
	NO_GROUPS:		"You are not currently in any groups.","您现在不在任何一个群中"
	OK:				"OK","好的"
	SIGNOUT:		"Sign out","登出"
};
