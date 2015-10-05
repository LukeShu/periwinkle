package cfg

import (
	"maildir"
	"net/http"
)

const IncomingMail maildir.Maildir = "/srv/periwinkle/Maildir"
const WebUiDir http.Dir = "./www"
const WebAddr string = ":8080"
