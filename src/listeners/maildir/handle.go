// Copyright 2015 Luke Shumaker
// Copyright 2015 Zhandos Suleimenov

package maildir

import (
	md "maildir"
)

func handle(maildir md.Maildir) {
	news, err := maildir.ListNew()
	if err != nil {
		return
	}
	for _, new := range news {
		cur, err := maildir.Acknowledge(new)
		if err != nil {
			continue
		}
		// TODO: Add data about `cur` to the RDBMS, and add it
		// to the outgoing queue as nescessary.
		cur.SetInfo("foo")
	}
}
