// Copyright 2015 Richard Wisniewski
var lang = "en";

try {
	lang = navigator.language.substring(0,2);
} catch(err) {}

var localised = {
	en: {},
	it: {},
	es: {}
}