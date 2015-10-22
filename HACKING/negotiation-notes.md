There's Server-driven Negotiation (RFC2616§12.1), Agent-driven
Negotiation (RFC2616§12.2), and Transport Negotiation (RFC2616§12.3).

# Server

We currently only do server-driven negotiation, but it's incomplete
and buggy.

The negotiation SHOULD be driven based on:

- Accept (content-type)
- Accept-Charset
- Accept-Encoding
- Accept-Language

but can include other things too.

We currently do not support negotiating anything but the content type.
Even if the others are non-negotiable, failing to come up with
something that matches SHOULD trigger a 406 (Not Acceptable).

Server-driven negotiation SHOULD set the `Vary:` header
(RFC2616§14.44) on responses, we currently don't.

# Agent

Let the user choose based on the lists returned in 300 (Multiple
Choices) and 406 (Not Acceptable) responses.

# Transport (caches)

The cache performs agent-driven negotiation when talking to the
upstream, and server-driven negotiation when talking to the client.

# conclusion

We should implement server-driven negotiation as a middleware that
essentially performs transport negotiation.
