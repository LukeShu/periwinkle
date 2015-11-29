// +build ignore

package xgettext

func gettext(str string) {}

func ngettext(str string) {}

func fn() {
	/* This is the first comment.  */
	gettext("foo")

	/* This is the second comment: not extracted  */
	gettext(
		"bar")

	gettext(
		/* This is the third comment.  */
		"baz")

	ngettext("FOO" + ("BAR" + "BAZ"))
}
