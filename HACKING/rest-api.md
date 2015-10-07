# Basic API design.

There are a number of paths that respond to [HTTP][RFC-2616] `POST`,
`PUT`, [`PATCH`][RFC-5789], `DELETE`, and `GET` requests, and may
return a document in response.

For actions requiring authentication, there will be a session token
(see below for how to get a token).  All requests that require
authentication must submit an HTTP cookie `session_id` set to the
token value.  Further, for requests that have a document that is
submitted (i.e., `POST`, `PUT`, and `PATCH`) the document must include
a `session_id` field that is also set to the token value.

> Rationale: The double-submission of the token for requests with a
> body is to protect from CSRF attacks.  This cannot be done for `GET`
> or `DELETE` requests, which is OK.  It is OK for `GET` requests
> because if the API is implemented correctly, they are not vulnerable
> to CSRF attacks.  It is OK for `DELETE` requests, because modern web
> browsers will not perform cross-domain `DELETE` requests unless
> explicitly enabled by the server using the
> `Access-Control-Allow-Origin` HTTP header.

If a `POST` request has a `_method` field in the submitted document,
the request shall be interpreted as a request of the type specified by
that field, instead of a `POST` request.  Even if the specified
`_method` does not have a document that is submitted (e.g., `DELETE`),
because it started as a `POST` request, the submitted document it must
still contain the `session_id` attribute.

> Rationale: HTML forms may only submit `GET` and `POST` requests;
> which is silly, and requires work-arounds to allow browsers to
> emulate other request types with `POST` requests.

Unless otherwise specified, response documents are always in JSON
([RFC-7159][], [ECMA-404][]); however, other document formats may be
added in the future; to request a specific format, either append the
correct file extension to the path, or specify the MIME type in the
HTTP `Accept` header (see below for a list of MIME types and file
extension).

If a file extension is not included in the path, any path may include
a trailing "/".  A trailing "/" may therefore be used to clarify that
a "." earlier in the path was part of the base-path, not the start of
a file extension.

> Rationale: JSON is a pleasure to work with; both file extensions and
> Accept headers are the "correct" thing to do.  Plus it makes
> prototyping clients easy.

For requests in which information is submitted to the server (that is,
everything but `GET` requests), the document may be submitted in
either JSON format or [form-data][RFC-2388] format; as specified by
the HTTP `Content-Type` header (with values of `application/json` and
`multipart/form-data` respectively).

> Rationale: JSON is a pleasure to work with. `form-data` is also
> supported in order to support submitting requests from HTML forms.

If a submitted document's `Content-Type` is not one that is supported,
an HTTP 415 ("Unsupported Media Type") response is returned.  If a
submitted document is of a supported type, but the document structure
does not match the expected format, then an HTTP 400 ("Bad Request")
response is returned.  If the request method is not supported for a
path, an HTTP 405 ("Method Not Allowed") response is returned.  If the
request's `Accept` HTTP header, or the file extension specifies a
format that is not supported, an HTTP 406 ("Not Acceptable") response
is returned.

[RFC-2388]: https://tools.ietf.org/html/rfc2388
	"Returning Values from Forms: multipart/form-data"
[RFC-2616]: https://tools.ietf.org/html/rfc2616
	"RFC 2616: Hypertext Transfer Protocol -- HTTP/1.1"
[RFC-5789]: https://tools.ietf.org/html/rfc5789
	"RFC 5789: PATCH Method for HTTP"
[RFC-7159]: https://tools.ietf.org/html/rfc7159
	"RFC 7159: The JavaScript Object Notation (JSON) Data Interchange Format"
[ECMA-404]: http://www.ecma-international.org/publications/files/ECMA-ST/ECMA-404.pdf
	"ECMA-404: The JSON Data Interchange Format"

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

	* `/session` [`POST`, `DELETE`]

		A `POST` request containing valid `login` and `password` will
		create a session; returning a HTTP 200 ("Found") containing
		`session_id`, as well as setting the `session_id` cookie.  If
		the `login` and `password` do not match a user, an HTTP 401
		("Unauthorized") response is returned.

		A `DELETE` request ends the current session (if there is one),
		and returns an HTTP 204 ("No Content") response.

	* `/s/%{identifier}` [`GET`]

		The `/s/` directory is for shortened URLs; a `GET` request to
		a valid short URL will return an HTTP 301 ("Moved
		Permanently") redirect to the appropriate full URL; otherwise
		it will return an HTTP 404 ("Not Found").

	* `/msgs/%{msgid}` [`GET:{json,mbox}`]

		Will return an HTTP 200 ("Found") with the message having the
		specified `Message-ID`; if it is found in the message store;
		or an HTTP 404 ("Not Found") otherwise.

		TODO: It is undecided if authentication should be required to
		access messages.

	* `/users` [`POST`]

		Will attempt to create a user. On success, an HTTP 201
		("Created"), with the `Location` header set to the created
		`/user/%{alias}`, and the response document with a list of
		URLs the resource is accessable at, differentiated by file
		extension.

		If the user can't be created, it will return an HTTP 409
		("Conflict") with a response document explaining the conflict.

		The submitted document must include "username", "email", and
		"password" fields, and may optionally include a
		"password_verification" field.

		* `/users/%{alias}` [`GET`, `PUT`, `PATCH`, `DELETE`]

			A `PUT` request totally replaces the user; the format is
			the same as when creating a user, except that it excludes
			any anti-spam type measures associated with original
			account creation.

			TODO: everything else

	* `/groups` [`POST`, `GET`]

		`POST` creates a group, `GET` lists groups that the user is
		allowed to see.

		TODO: everything else

	* `/groups/%{alias}` [`GET`, `PUT`, `PATCH`, `DELETE`]

		A `PUT` request totally replaces the group; the format is the
		same as when creating a a group.

		TODO: everything else

		* `/groups/%{alias}/log` [`GET`]

			Returns an HTTP 200 response containing a list of
			`Message-ID`s, if the group alias points to a valid group
			that the user is allowed to se; otherwise returns HTTP 401
			("Unauthorized")
