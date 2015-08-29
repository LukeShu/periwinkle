package inotify

const (
	// Flags for the parameter of InotifyInit1().
	// These, oddly, appear to be 24-bit numbers.
	IN_CLOEXEC  uint32 = 02000000
	IN_NONBLOCK uint32 = 00004000
)

type Mask uint32

const (
	// Supported events suitable for the `mask` parameter of Inotify.AddWatch().
	IN_ACCESS        Mask = (1<< 0) // File was accessed.
	IN_MODIFY        Mask = (1<< 1) // File was modified.
	IN_ATTRIB        Mask = (1<< 2) // Metadata changed.
	IN_CLOSE_WRITE   Mask = (1<< 3) // Writtable file was closed.
	IN_CLOSE_NOWRITE Mask = (1<< 4) // Unwrittable file closed.
	IN_OPEN          Mask = (1<< 5) // File was opened.
	IN_MOVED_FROM    Mask = (1<< 6) // File was moved from X.
	IN_MOVED_TO      Mask = (1<< 7) // File was moved to Y.
	IN_CREATE        Mask = (1<< 8) // Subfile was created.
	IN_DELETE        Mask = (1<< 9) // Subfile was deleted.
	IN_DELETE_SELF   Mask = (1<<10) // Self was deleted.
	IN_MOVE_SELF     Mask = (1<<11) // Self was moved.

	// Events that appear in output without subscribing to them.
	IN_UNMOUNT       Mask = (1<<13) // Backing fs was unmounted.
	IN_Q_OVERFLOW    Mask = (1<<14) // Event queued overflowed.
	IN_IGNORED       Mask = (1<<15) // File was ignored (expect no more events).

	// Special flags that you may pass to Inotify.AddWatch()...
	// except for IN_ISDIR, which is a flag that is set on output events.
	IN_ONLYDIR       Mask = (1<<24) // Only watch the path if it is a directory.
	IN_DONT_FOLLOW   Mask = (1<<25) // Do not follow a sym link.
	IN_EXCL_UNLINK   Mask = (1<<26) // Exclude events on unlinked objects.
	IN_MASK_ADD      Mask = (1<<29) // Add to the mask of an already existing watch.
	IN_ISDIR         Mask = (1<<30) // Event occurred against dir.
	IN_ONESHOT       Mask = (1<<31) // Only send event once.

	// Convenience macros */
	IN_CLOSE      Mask = (IN_CLOSE_WRITE | IN_CLOSE_NOWRITE) // Close.
	IN_MOVE       Mask = (IN_MOVED_FROM | IN_MOVED_TO)       // Moves.
	IN_ALL_EVENTS Mask = 0x00000FFF                          // All events which a program can wait on.
)

var in_bits [32]string = [32]string{
	// mask
	/*  0 */ "IN_ACCESS",
	/*  1 */ "IN_MODIFY",
	/*  2 */ "IN_ATTRIB",
	/*  3 */ "IN_CLOSE_WRITE",
	/*  4 */ "IN_CLOSE_NOWRITE",
	/*  5 */ "IN_OPEN",
	/*  6 */ "IN_MOVED_FROM",
	/*  7 */ "IN_MOVED_TO",
	/*  8 */ "IN_CREATE",
	/*  9 */ "IN_DELETE",
	/* 10 */ "IN_DELETE_SELF",
	/* 11 */ "IN_MOVE_SELF",
	/* 12 */ "(1<<12)",
	// events sent by the kernel
	/* 13 */ "IN_UNMOUNT",
	/* 14 */ "IN_Q_OVERFLOW",
	/* 15 */ "IN_IGNORED",
	/* 16 */ "(1<<16)",
	/* 17 */ "(1<<17)",
	/* 18 */ "(1<<18)",
	/* 19 */ "(1<<19)",
	/* 20 */ "(1<<20)",
	/* 21 */ "(1<<21)",
	/* 22 */ "(1<<22)",
	/* 23 */ "(1<<23)",
	// special flags
	/* 24 */ "IN_ONLYDIR",
	/* 25 */ "IN_DONT_FOLLOW",
	/* 26 */ "IN_EXCL_UNLINK",
	/* 27 */ "(1<<27)",
	/* 28 */ "(1<<28)",
	/* 29 */ "IN_MASK_ADD",
	/* 30 */ "IN_ISDIR",
	/* 31 */ "IN_ONESHOT",
}

func (mask Mask) String() string {
	out := ""
	for i, name := range in_bits {
		if mask&(Mask(1)<<uint(i)) != 0 {
			if len(out) > 0 {
				out += "|"
			}
			out += name
		}
	}
	return out
}
