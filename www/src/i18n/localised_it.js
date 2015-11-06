// Copyright 2015 Richard Wisniewski
var _it = {
	CONTACT_ADMIN:	"Si prega di contattare un amministratore."
};
localised.it = {
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
			EMAIL_INVALID:				"That is not a valid email address.",
			ERRORS:{
				"409": {
					TITLE:		"Utente Già Esiste",
					CONTENT:	"Nome utente o indirizzo email è già usato."
				}
			}
		},
		LOGIN: {
			LOGIN:			"Accedi",
			MESSAGE:		"Non puoi accedere questo risorso",
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
			STOP_EDIT:					"Finisce a modificare",
			FULL_NAME:					"Nome e Cognome",
			NO_FULL_NAME:				"Non so che si chiama. ",
			CHANGE_NAME:				"Cambia nome",
			ADDRESSES:					"Indirizzi di contatare",
			SAVING:						"Salvando..."
		},
		GROUPS:	{
			NEW:	"Crea Gruppo",
			JOIN:	"Unisci un gruppo",
			GROUPS:			"Gruppi",
			NO_GROUPS:		"Non sei in nessun gruppo al momento."
		},
		NEW_GROUP:	{
			TITLE: {
				MAIN:		"Crea Gruppo",
				CREATING:	"Creando Group..."
			},
			GROUP_NAME:		"Come si chiamano",
			ERRORS:	{
				"409": {
					TITLE:		"Questo nome di gruppo è già usato.",
					CONTENT:	"Non puoi usare questo nome per un nuovo gruppo perché un altro gruppo lo sta già usando.  "
				}
			}
		},
		CHANGE_PASSWORD:	{
			TITLE: {
				MAIN:		"Cambia Password",
				CREATING:	"Cambiando Password..."
			},
			ERRORS:	{
				"409": {
					TITLE:		"Errore di Cambiare Password",
					CONTENT:	"Non scriva la tua password corrente corregiamente.  "
				}
			},
			OLD_PASSWORD:	"Password Corrente",
			NEW_PASSWORD:	"Password Nuova",
			CONFIRM:		"Confirma Password"
		},
		NEW_ADDRESS: {
			TITLE: {
				MAIN:		"Aggiunge Indirizzo",
				CREATING:	"Aggiungendo Indirizzo..."
			},
			MEDIUM:		"Mezzo",
			ADDRESS:	"Indirizzo"
		},
		REDIR:	"Sposterai alla pagina del suo utente quando accedi.  "
	},
	GENERAL: {
		USERNAME_EMAIL:	"Nome utente o email",
		USERNAME:		"Nome utente",
		EMAIL:			"Email",
		PASSWORD:		"Password",
		FORM: {
			ERROR: {
				REQUIRED_FIELD:	"Questo campo è obbligatorio."
			},
			RESPONSE: {
				OK:			"OK",
				CANCEL:		"Cancella",
				CHANGE:		"Cambia",
				CREATE:		"Crea",
				ADD:		"Aggiunge"
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
