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
    ║          ║>───╫───┬──>║ Postfix ║>─────>║ receive-email                            ║<─────────>║       ║ ║
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
    ║          ║<───╫───────[polling]────────<║ listen-twilio                            ║       ╎   ║       ║ ║
    ║          ║    ║   │                 ╎   ╚══════════════════════════════════════════╝>──┐   ╎   ║       ║ ║
    ║          ║    ║   │                 ╎                                                  │   ╎   ║       ║ ║
    ║          ║    ║   ^────────────────────────────────────────────────────────────────────┘   ╎   ║       ║ ║
    ║          ║    ║   │                 ╎                                                      ╎   ║       ║ ║
    ║          ║    ║   │                 ╎   ╔══════════════════════════════════════════╗<─────────>║       ║ ║
    ║          ║>───╫────────────────────────>║ listen-http                              ║       ╎   ║       ║ ║
    ║          ║    ║   │                 ╎   ╚══════════════════════════════════════════╝>──┐   ╎   ║       ║ ║
    ║          ║    ║   ^                 ╎                                                  │   ╎   ║       ║ ║
    ║          ║    ║   └────────────────────────────────────────────────────────────────────┘   ╎   ╚═══════╝ ║
    ║          ║    ║                     ╎                                                      ╎             ║
    ╚══════════╝    ╚══════════════════════════════════════════════════════════════════════════════════════════╝

