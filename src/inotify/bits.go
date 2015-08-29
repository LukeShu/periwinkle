package inotify

/* Flags for the parameter of inotify_init1.  */
const (
	/* These are, oddly, 24-bit numbers */
	IN_CLOEXEC  uint32 = 02000000
	IN_NONBLOCK uint32 = 00004000
)

type Mask uint32

const (
	/* Supported events suitable for MASK parameter of INOTIFY_ADD_WATCH.  */
	IN_ACCESS        Mask = 0x00000001 /* File was accessed.  */
	IN_MODIFY        Mask = 0x00000002 /* File was modified.  */
	IN_ATTRIB        Mask = 0x00000004 /* Metadata changed.  */
	IN_CLOSE_WRITE   Mask = 0x00000008 /* Writtable file was closed.  */
	IN_CLOSE_NOWRITE Mask = 0x00000010 /* Unwrittable file closed.  */
	IN_OPEN          Mask = 0x00000020 /* File was opened.  */
	IN_MOVED_FROM    Mask = 0x00000040 /* File was moved from X.  */
	IN_MOVED_TO      Mask = 0x00000080 /* File was moved to Y.  */
	IN_CREATE        Mask = 0x00000100 /* Subfile was created.  */
	IN_DELETE        Mask = 0x00000200 /* Subfile was deleted.  */
	IN_DELETE_SELF   Mask = 0x00000400 /* Self was deleted.  */
	IN_MOVE_SELF     Mask = 0x00000800 /* Self was moved.  */

	/* Events sent by the kernel.  */
	IN_UNMOUNT       Mask = 0x00002000 /* Backing fs was unmounted.  */
	IN_Q_OVERFLOW    Mask = 0x00004000 /* Event queued overflowed.  */
	IN_IGNORED       Mask = 0x00008000 /* File was ignored.  */

	/* Special flags.  */
	IN_ONLYDIR       Mask = 0x01000000 /* Only watch the path if it is a directory.  */
	IN_DONT_FOLLOW   Mask = 0x02000000 /* Do not follow a sym link.  */
	IN_EXCL_UNLINK   Mask = 0x04000000 /* Exclude events on unlinked objects.  */
	IN_MASK_ADD      Mask = 0x20000000 /* Add to the mask of an already existing watch.  */
	IN_ISDIR         Mask = 0x40000000 /* Event occurred against dir.  */
	IN_ONESHOT       Mask = 0x80000000 /* Only send event once.  */

	/* Convenience macros */
	IN_CLOSE      Mask = (IN_CLOSE_WRITE | IN_CLOSE_NOWRITE) /* Close.  */
	IN_MOVE       Mask = (IN_MOVED_FROM | IN_MOVED_TO)       /* Moves.  */
	IN_ALL_EVENTS Mask = 0x00000FFF                          /* All events which a program can wait on.  */

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
