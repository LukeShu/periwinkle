# Basic API design.

The API is accessible over the [HTTP protocol][RFC-7230].

The HTTP resources use the
[standard HTTP semantics and request types][RFC-7231]; `OPTIONS`,
`HEAD`, `GET`, `DELETE`, `POST`, `PUT`, as well as
[`PATCH`][RFC-5789].

For actions requiring authentication, there will be a session token
(see below for how to get a token).  All requests that require
authentication must set the [HTTP cookie][RFC-6265] `session_id` and
the HTTP header `X-XSRF-TOKEN` to the session token value; with the
exception that `OPTIONS`, `HEAD`, and `GET` requests don't require the
`X-XSRF-TOKEN` header to be set.

> Rationale: The double-submission of the token as both a cookie and a
> header is to protect from CSRF attacks.

If a `POST` request has a `_method` field in the submitted document,
the request shall be interpreted as a request of the type specified by
that field, instead of a `POST` request.  Similarly, if a `POST`
request has a `_xsrf_token` field in the submitted document, the
request shall be interpretted as if it had the `X-XSRF-TOKEN` HTTP
header set to that value.  Finally, if a `POST` request has a `_body`
field, then the content of just that field will be treated as the
submitted document body.

> Rationale: HTML forms may only submit `GET` and `POST` requests;
> which is silly, and requires work-arounds to allow browsers to
> emulate other request types with `POST` requests, as well as setting
> HTTP headers.  The `_body` field is an option because it isn't
> always appropriate for the root-element to be a map.

Unless otherwise specified, response documents always support the JSON
format ([RFC-7159][], [ECMA-404][]); however, other document formats
may be added in the future; to request a specific format, either
append the correct file extension to the path, or specify the MIME
type in the HTTP `Accept` header (see below for a list of MIME types
and file extension).

If a file extension is not included in the path, any path may include
a trailing "/".  A trailing "/" may therefore be used to clarify that
a "." earlier in the path was part of the base-path, not the start of
a file extension.

> Rationale: JSON is a pleasure to work with; both file extensions and
> Accept headers are the "correct" thing to do.  Plus it makes
> prototyping clients easy.

Unless otherwise specified, `POST` and `PUT` requests may be in the
following formats:
 - JSON (`Content-Type: application/json`)
 - [form-data][RFC-2388] (`Content-Type: multipart/form-data`)
 - [form-urlencoded][form-urlencoded] (`Content-Type: application/x-www-form-urlencoded`)

Unless otherwise specified, `PATCH` requests may be in the following
formats:
 - [JSON Patch][RFC-6902] (`Content-Type: application/json-patch+json`)
 - [JSON Merge Patch][RFC-7368] (`Content-Type: application/merge-patch+json`)

> Rationale: JSON is a pleasure to work with. `form-data` and
> `form-urlencoded` are also supported in order to support submitting
> requests from HTML forms.

If a submitted document's `Content-Type` is not one that is supported,
an HTTP 415 ("Unsupported Media Type") response is returned.  If a
`PATCH` or `PUT` request submits a document that is of a supported
type, but has fields that don't match the expected ones, then an HTTP
415 ("Unsupported Media Type") response is returned.  similarly, if a
`PATCH` request tries to access fields that don't exist, then an HTTP
415 ("Unsupported Media Type") is returned.  If the request method is
not supported for a path, an HTTP 405 ("Method Not Allowed") response
is returned.  If the request's `Accept` HTTP header, or the file
extension specifies a format that is not supported, an HTTP 406 ("Not
Acceptable") response is returned.  If the request does not have an
`Accept` HTTP header, and there are multiple possible representations
for a type, or if the `Accept` header specifies equal preference for
multiple of the possible representations, then an HTTP 300 ("Multiple
Choices") response is returned, containing a list of possibilities.

[RFC-2388]: https://tools.ietf.org/html/rfc2388
	"Returning Values from Forms: multipart/form-data"
[RFC-2616]: https://tools.ietf.org/html/rfc2616
	"RFC 2616: Hypertext Transfer Protocol -- HTTP/1.1"
[RFC-5789]: https://tools.ietf.org/html/rfc5789
	"RFC 5789: PATCH Method for HTTP"
[RFC-6265]: https://tools.ietf.org/html/rfc6265
	"RFC 6265: HTTP State Management Mechanism"
[RFC-6902]: https://tools.ietf.org/html/rfc6902
	"RFC 6902: JavaScript Object Notation (JSON) Patch"
[RFC-7159]: https://tools.ietf.org/html/rfc7159
	"RFC 7159: The JavaScript Object Notation (JSON) Data Interchange Format"
[RFC-7230]: https://tools.ietf.org/html/rfc7231
	"Hypertext Transfer Protocol (HTTP/1.1): Message Syntax and Routing"
[RFC-7231]: https://tools.ietf.org/html/rfc7231
	"Hypertext Transfer Protocol (HTTP/1.1): Semantics and Content"
[RFC-7368]: https://tools.ietf.org/html/rfc7368
	"RFC 7368: JSON Merge Patch"
[ECMA-404]: http://www.ecma-international.org/publications/files/ECMA-ST/ECMA-404.pdf
	"ECMA-404: The JSON Data Interchange Format"
[form-urlencoded]: http://www.w3.org/html/wg/drafts/html/master/semantics.html#application/x-www-form-urlencoded-encoding-algorithm
	"HTML5.1: x-www-form-urlencoded encoding Algorithm"

# File-extenson / MIME-type mapping

 * file-extension / MIME-type
 * `.json` / `application/json`
 * `.mbox` / `application/mbox`

# Paths

* `/` [`GET`]

	Returns an HTTP 302 ("Found") redirect to `/webui/`

	* `/webui/*`

		Serves resources for the fat-client web UI.

	* `/callbacks/*`

		I'm almost certain that some external API that we will
		interface with will need a callback URL that it can make
		requests to.  They will go here.  I was just skimming the
		Twilio docks, and it looks like that's not the case for
		Twilio, which really surprises me.  But I'm sure it will still
		come in to use.

	* `/s/%{identifier}` [`GET`]

		The `/s/` directory is for shortened URLs; a `GET` request to
		a valid short URL will return an HTTP 301 ("Moved
		Permanently") redirect to the appropriate full URL; otherwise
		it will return an HTTP 404 ("Not Found").

	* `/v1`

		* `/v1/session` [`POST`, `DELETE`]

			A `POST` request containing valid `login` and `password` will
			create a session; returning a HTTP 200 ("Found") containing
			`session_id`, as well as setting the `session_id` cookie.  If
			the `login` and `password` do not match a user, an HTTP 403
			("Forbidden") response is returned.

			A `DELETE` request ends the current session (if there is one),
			and returns an HTTP 204 ("No Content") response.

		* `/v1/msgs/%{msgid}` [`GET:{json,mbox}`]

			Will return an HTTP 200 ("Found") with the message having the
			specified `Message-ID`; if it is found in the message store;
			or an HTTP 404 ("Not Found") otherwise.

			TODO: It is undecided if authentication should be required to
			access messages.

		* `/v1/users` [`POST`]

			Will attempt to create a user. On success, an HTTP 201
			("Created"), with the `Location` header set to the created
			`/v1/user/%{id}`, and the response document with a list of
			URLs the resource is accessable at, differentiated by file
			extension.

			If the user can't be created, it will return an HTTP 409
			("Conflict") with a response document explaining the conflict.

			The submitted document must include "username", "email", and
			"password" fields, and may optionally include a
			"password_verification" field.

			* `/v1/users/%{id}` [`GET`, `PUT`, `PATCH`, `DELETE`]

				A `PUT` request totally replaces the user; the format is
				the same as when creating a user, except that it excludes
				any anti-spam type measures associated with original
				account creation.

				A `PATCH` request may be in either patch format, *but*
				the password can only be changed by a a JSON Patch,
				not a JSON Merge Patch.  Further, when changing the
				password, the patch must first perform a `test` to
				verify the old value of the password, then perform a
				`replace` to set the new password.  A `replace`
				without a `test` first will fail with HTTP 415
				("Unsupported Media Type").

				TODO: document the structure of the user.

		* `/v1/groups` [`POST`, `GET`]

			`POST` creates a group, `GET` lists groups that the user is
			allowed to see.

			TODO: everything else

		* `/v1/groups/%{alias}` [`GET`, `PUT`, `PATCH`, `DELETE`]

			A `PUT` request totally replaces the group; the format is the
			same as when creating a a group.

			TODO: everything else

			* `/v1/groups/%{alias}/log` [`GET`]

				Returns an HTTP 200 response containing a list of
				`Message-ID`s, if the group alias points to a valid group
				that the user is allowed to se; otherwise returns HTTP 403
				("Forbidden")

		* `/v1/captcha` [`POST`]

			The body of the POST doesn't matter.  Returns HTTP 201
			("Created") with a document body containing the ID for the
			created captcha.

			* `/v1/captcha/%{id}` [`GET:{png,wav}`, `PUT`]

				GET returns a PNG or WAV of the captcha (based on the
				`Accept:` header) with HTTP 200 ("Found"). PUT takes a
				string of what the user says the captcha said.  If the
				string matches what it says, it will return a token
				(proof of solution) with an HTTP 200 ("Found");
				otherwise it will return HTTP 403 ("Forbidden").
