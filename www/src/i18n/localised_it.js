// Copyright 2015 Richard Wisniewski
var _it = {
	CONTACT_ADMIN:	"Si prega di contattare un amministratore."
};
localised.it = {
	USERNAME_EMAIL:	"Nome utente o email",
	USERNAME:		"Nome utente",
	EMAIL:			"Email",
	PASSWORD:		"Password",
	LOGIN:			"Accedi",
	REQUIRED_FIELD:	"Questo campo è obbligatorio.",
	SIGNUP: {
		SIGNUP:						"Iscriviti",
		NOT_A_USER:					"Non sei un utente?",
		ALREADY_USER:				"Sei già un utente?",
		CONFIRM_EMAIL:				"Conferma email",
		CONFIRM_EMAIL_NO_MATCH:		"I tuoi indirizzi d'email non corrispondono.",
		CONFIRM_PASSWORD:			"Conferma password",
		CONFIRM_PASSWORD_NO_MATCH:	"Le tue password non corrispondono.",
		USERNAME_INVALID:			"I nomi utente possono contenere solo lettere, numeri, _, o -.",
		EMAIL_INVALID:				"Quello non è un indirizzo d'email valido.",
		ERRORS:{
			"409": {
				TITLE:		"Utente Già Esiste",
				CONTENT:	"Nome utente o indirizzo email è già usato."
			},
			"500": {
				TITLE:		"Qualcosa È Sbagliato",
				CONTENT:	"C'è un errore del server.  " + _it.CONTACT_ADMIN
			},
		}
	},
	LOGIN: {
		ERRORS: {
			"401": {
				TITLE:		"Accesso Fallito",
				CONTENT:	"Nome utente e password non corrispondono.  "
			}
		}
	},
	ERRORS:	{
		DEFAULT: {
			TITLE:		"Errore Sconosciuto",
			CONTENT:	"Una risposta inattesa è stata ricevuta dal server. " + _it.CONTACT_ADMIN
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
	NO_GROUPS:		"Non sei in nessun gruppo al momento.",
	OK:				"OK"
};