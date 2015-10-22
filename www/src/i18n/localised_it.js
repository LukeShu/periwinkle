// Copyright 2015 Richard Wisniewski
var _it = {
	CONTACT_ADMIN:	"Contatta un amministratore."
};
localised.it = {
	USERNAME_EMAIL:	"Nome utente o email",
	USERNAME:		"Nome utente",
	EMAIL:			"Email",
	PASSWORD:		"Password",
	LOGIN:			"Accedi",
	REQUIRED_FIELD:	"Questo campo è necessario.",
	SIGNUP: {
		SIGNUP:						"Inscriviti",
		NOT_A_USER:					"Non sei un utente?",
		ALREADY_USER:				"Sei già un utente?",
		CONFIRM_EMAIL:				"Conferma email",
		CONFIRM_EMAIL_NO_MATCH:		"I tuoi indirizzi d'email non sono simili.",
		CONFIRM_PASSWORD:			"Conferma password",
		CONFIRM_PASSWORD_NO_MATCH:	"I tuoi password non sono simili.",
		USERNAME_INVALID:			"I nomi untente possono contenere lettere, numeri, _, o -.",
		EMAIL_INVALID:				"Quello non è un indirizzo d'email valido.",
		ERRORS:{
			"409": {
				TITLE:		"Utente Già Esiste",
				CONTENT:	"Nome utente o indirizzo email è già usato."
			},
			"500": {
				TITLE:		"Qualcosa Sbaglia",
				CONTENT:	"Il server inconta un errore.  " + _it.CONTACT_ADMIN
			},
		}
	},
	LOGIN: {
		ERRORS: {
			"401": {
				TITLE:		"Accesso Falito",
				CONTENT:	"Nome utente e password non si correspondono.  "
			}
		}
	},
	ERRORS:	{
		DEFAULT: {
			TITLE:		"Errore Sconosciuto",
			CONTENT:	"Ricevuto una risposta inattesa dal server. " + _it.CONTACT_ADMIN
		},
		CONTACT_ADMIN:	_it.CONTACT_ADMIN
	},
	USER: {
		USER:						"Utente",
		EMAIL_CHANGE:				"Cambia",
		CHANGE_PASSWORD:			"Cambia Password",
		DELETE:						"Elimina account"
	},
	GROUPS:			"Gruppi",
	NO_GROUPS:		"Non hai nessuno gruppi adesso."
};