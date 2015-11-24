// Copyright 2015 Mark Pundmann
package backend

type Existence int

const (
	public       = 1
	confirmed    = 2
	member       = 3
	moderator    = 3
	auto         = 1
	confirmation = 2
)

var reverse = map[string]int{
	"public":       1,
	"confirmed":    2,
	"member":       3,
	"moderator":    3,
	"auto":         1,
	"confirmation": 2,
}

func Reverse(m string) int { return reverse[m] }

var existence = [...]string{
	"public",
	"confirmed",
	"member",
}

func (m Existence) String() string { return existence[m-1] }

type Read int

var read = [...]string{
	"public",
	"confirmed",
	"member",
}

func (m Read) String() string { return read[m-1] }

type Post int

var post = [...]string{
	"public",
	"confirmed",
	"moderator",
}

func (m Post) String() string { return post[m-1] }

type Join int

var join = [...]string{
	"auto",
	"confirmation",
}

func (m Join) String() string { return join[m-1] }
