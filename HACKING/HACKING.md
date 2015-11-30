Try to remember to make sure that `make check` passes before you
commit.

Avoid using the `error` type, use `locale.Error` instead.  To convert
an `error` from an external library to a `locale.Error`, use
`locale.UntranslatedError(err)`.

Avoid using `panic()` unless you know what you are doing.  In the
`periwinkle/backend` package, these alternatives exist:

 * `dbError()`: if you encountered a database error.
 * `programmerError()`: if you received an argument that was invalid,
   or some other condition meaning that you or another committer
   screwed up.

Avoid using `log` or `fmt.Print*`, use `periwinkle.Logf` or
`periwinkle.LogErr` instead.

Avoid using the `os` package.  If you need something from the outside
system, you should have it passed to you as an argument.  Obviously,
there need to be exceptions for the top-level code that will call you,
passing in the environment stuff!  Right now, those exceptions are:
 - `*/cmd/*` (anything with `package main`)
 - `periwinkle/cfg/parse.go`
 - `periwinkle/util.go`
 - `postfixpipe/`
 - `maildir/`
 - `locale/gettext/`

Avoid doing any raw GORM/SQL outside of `periwinkle/backend`.  If you
need to do something that you need raw GORM/SQL for, add a function or
method to do it to the `backend` package.

Avoid doing input validation in Go, do it in SQL instead.  There are
exceptions to this, but I can't think of a better heuristic than "I
know it when I see it."

DO NOT do a database call inside of a `for` loop.  Figure out how to
hoist the DB call to outside of the loop.  (Obvious exception: the
main program loop, if there is one).

To add a table to the database in `periwinkle/backend`, define the
struct like the existing tables, write a `dbSchema` method, and add it
to the list in `src/periwinkle/backend/tables.go`.  Optionally, you
may define a `dbSeed` method that sets up initial data into the table.
