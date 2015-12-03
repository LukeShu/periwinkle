Installing Periwinkle
---------------------

# Compilation

First up, you have to compile the software.  There is no
pre-compilation configuration; simply run

    $ make

This will fetch any dependencies that it needs to.  Inspecting the
Makefiles will reveal what these are without too much effort (grep for
`$(NET)`); for if you are putting a build script for which network
access is not acceptable.

By default, if a fetched-dependency exists, it is left alone.  To
force the build system to try to update already-downloaded resources,
tack on `NET=FORCE`.

    $ make NET=FORCE

Compilation results in several executable programs in the `bin/`
directory within the source directory.

# How the components interact

    ╔══════════╗    ╔══════════════════════════════════════════════════════════════════════════════════════════╗
    ║          ║    ║ Your server                                                                              ║
    ║          ║    ╠══════════════════════════════════════════════════════════════════════════════════════════╣
    ║          ║    ║                     ╎                     Periwinkle                       ╎             ║
    ║          ║    ║                     ╎                                                      ╎             ║
    ║          ║    ║       ╔═════════╗   ╎   ╔══════════════════════════════════════════╗       ╎   ╔═══════╗ ║
    ║          ║>───╫───┬──>║ Postfix ║>─────>║ bin/receive-email                        ║<─────────>║       ║ ║
    ║          ║    ║   │   ╚═════════╝   ╎   ╠══════════════════════════════════════════╣       ╎   ║       ║ ║
    ║          ║    ║   │                 ╎   ║ @yourdomain.tld -> incoming mail handler ║       ╎   ║       ║ ║
    ║          ║    ║   │                 ╎   ║ @sms.gateway    -> Twilio Gateway ────────>──┐   ╎   ║       ║ ║
    ║          ║    ║   │                 ╎   ║ @mms.gateway    -> Twilio Gateway ────────>──v   ╎   ║       ║ ║
    ║          ║    ║   │                 ╎   ║ @irc.gateway    -> IRC Gateway ───────────>──v   ╎   ║       ║ ║
    ║          ║    ║   │                 ╎   ║     ...         ->        ...            ║   │   ╎   ║       ║ ║
    ║   The    ║    ║   │                 ╎   ╚══════════════════════════════════════════╝   │   ╎   ║       ║ ║
    ║ Internet ║    ║   │                 ╎                                                  │   ╎   ║       ║ ║
    ║          ║    ║   ^────────────────────────────────────────────────────────────────────┘   ╎   ║       ║ ║
    ║          ║    ║   │                 ╎                                                      ╎   ║ MySQL ║ ║
    ║          ║    ║   │                 ╎   ╔══════════════════════════════════════════╗<─────────>║       ║ ║
    ║          ║<───╫───────[polling]────────<║ bin/listen-twilio                        ║       ╎   ║       ║ ║
    ║          ║    ║   │                 ╎   ╚══════════════════════════════════════════╝>──┐   ╎   ║       ║ ║
    ║          ║    ║   │                 ╎                                  ^               │   ╎   ║       ║ ║
    ║          ║    ║   ^────────────────────────────────────────────────────│───────────────┘   ╎   ║       ║ ║
    ║          ║    ║   │                 ╎                                  v                   ╎   ║       ║ ║
    ║          ║    ║   │                 ╎   ╔══════════════════════════════════════════╗<─────────>║       ║ ║
    ║          ║>───╫────────────────────────>║ bin/listen-http                          ║       ╎   ║       ║ ║
    ║          ║    ║   │                 ╎   ╚══════════════════════════════════════════╝>──┐   ╎   ║       ║ ║
    ║          ║    ║   ^                 ╎                                                  │   ╎   ║       ║ ║
    ║          ║    ║   └────────────────────────────────────────────────────────────────────┘   ╎   ╚═══════╝ ║
    ║          ║    ║                     ╎                                                      ╎             ║
    ╚══════════╝    ╚══════════════════════════════════════════════════════════════════════════════════════════╝

With the exception that SQLite3 can be used instead of MySQL (for
development purposes), Postfix and MySQL are specifically required;
not other MTAs or RDBMSs.

Unfortunately, the standard Go language database interface doesn't
sufficiently abstract certain details of messages that the DB sends
us, and we must write code for each RDBMS that we wish to support.
The amount of code that must be written for each RDBMS is very small;
but thus far we've only done it for for SQLite3 and MySQL.

The contract with Postfix is also not abstracted such that other MTAs
may be used.  The bits of the contract that it relies on are:

 - The format of stdin for the delivery program; that is:
   A line with some information (TODO: what information) from Postfix,
   a newline, then the RFC X822-formatted message.
 - The environment variable `ORIGINAL_RECIPIENT` set to the address
   that we received the message on.
 - The interpretation of exit codes according to <sysexits.h> into
   SMTP responses.
