// Copyright 2015 Richard Wisniewski
;(function() {
	'use strict';
	var _it = {
		CONTACT_ADMIN:	"Si prega di contattare un amministratore."
	};
	var localised = {
		LOGIN: {
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
					}
				}
			},
			LOGIN: {
				LOGIN:			"Accedi",
				MESSAGE:		"Non puoi accedere a questa pagina",
				ERRORS: {
					"403": {
						TITLE:		"Accesso Fallito",
						CONTENT:	"Nome utente e password non corrispondono.  "
					}
				}
			}
		},
		USER: {
			INFO: {
				USER:						"Utente",
				EMAIL_CHANGE:				"Cambia",
				EMAIL_DELETE:				"Elimina",
				CHANGE_PASSWORD:			"Cambia Password",
				DELETE:						"Elimina Account",
				BEGIN_EDIT:					"Comincia a modificare",
				STOP_EDIT:					"Finisci di modificare",
				FULL_NAME:					"Nome e Cognome",
				NO_FULL_NAME:				"Il nome non fornito. ",
				CHANGE_NAME:				"Cambia nome",
				ADDRESSES:					"Indirizzi di contatto",
				SAVING:						"Salvataggio in corso..."
			},
			GROUPS:	{
				NEW:	"Crea Gruppo",
				JOIN:	"Unisciti a un gruppo",
				GROUPS:			"Gruppi",
				NO_GROUPS:		"Non sei in nessun gruppo al momento."
			},
			NEW_GROUP:	{
				TITLE: {
					MAIN:		"Crea Gruppo",
					CREATING:	"Creazione del gruppo in corso..."
				},
				GROUP_NAME:		"Nuovo Nome del Gruppo",
				ERRORS:	{
					"409": {
						TITLE:		"Questo nome di gruppo è già usato.",
						CONTENT:	"Non puoi usare questo nome per un nuovo gruppo perché già in uso.  "
					}
				}
			},
			CHANGE_PASSWORD:	{
				TITLE: {
					MAIN:		"Cambia Password",
					CREATING:	"Cambiamento Password in corso..."
				},
				ERRORS:	{
					"409": {
						TITLE:		"Cambiamento Password Falito",
						CONTENT:	"La password inserita non è corretta.  "
					}
				},
				OLD_PASSWORD:	"Password Corrente",
				NEW_PASSWORD:	"Password Nuova",
				CONFIRM:		"Conferma Password"
			},
			NEW_ADDRESS: {
				TITLE: {
					MAIN:		"Aggiungi Indirizzo",
					CREATING:	"L'aggiunta dell'Indirizzo in corso..."
				},
				MEDIUM:		"Mezzo",
				ADDRESS:	"Indirizzo"
			},
			REDIR:	"Sarai rediretto alla pagina del tuo utente quando accedi.  "
		},
		GENERAL: {
			USERNAME_EMAIL:	"Nome utente o email",
			USERNAME:		"Nome utente",
			EMAIL:			"Email",
			PASSWORD:		"Password",
			CAPTCHA:		"CAPTCHA",
			FORM: {
				ERROR: {
					REQUIRED_FIELD:	"Questo campo è obbligatorio."
				},
				RESPONSE: {
					OK:			"OK",
					CANCEL:		"Cancella",
					CHANGE:		"Cambia",
					CREATE:		"Crea",
					ADD:		"Aggiungi"
				}
			},
			ERRORS:	{
				"500": {
					TITLE:		"Qualcosa È Sbagliato",
					CONTENT:	"C'è un errore del server.  " + _it.CONTACT_ADMIN
				},
				DEFAULT: {
					TITLE:		"Errore Sconosciuto",
					CONTENT:	"Una risposta inattesa è stata ricevuta dal server. " + _it.CONTACT_ADMIN
				},
				CONTACT_ADMIN:	_it.CONTACT_ADMIN
			},
			MEDIUMS: {
				EMAIL:		"Email",
				SMS:		"SMS",
				MMS:		"MMS"
			},
			SIGNOUT:		"Esci"
		}
	};

	angular.module('periwinkle.i18n').constant('i18n_it', localised);
})();
