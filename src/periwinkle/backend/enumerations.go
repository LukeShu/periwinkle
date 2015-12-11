// Copyright 2015 Mark Pundmann

package backend

//type PostJoin int

const (
	bounce   = 1
	moderate = 2
	accept   = 3
	no       = 1
	yes      = 1
)

var reverse = map[string]int{
	"bounce":    1,
	"moderate":  2,
	"accept":    3,
	"no":        1,
	"yes":       2,
	"public":    0,
	"confirmed": 1,
	"member":    2,
}

func Reverse(m map[string]string) []int {
	a := make([]int, len(m))
	for key, value := range m {
		a[reverse[key]] = reverse[value]
	}
	return a
}

var postjoin = map[int]string{
	1: "bounce",
	2: "moderate",
	3: "accept",
}

func PostJoin(m [3]int) map[string]string {
	a := make(map[string]string)
	a["public"] = postjoin[m[0]]
	a["confirmed"] = postjoin[m[1]]
	a["member"] = postjoin[m[2]]
	return a
}

var readexist = map[int]string{
	1: "yes",
	2: "no",
}

func ReadExist(m [2]int) map[string]string {
	a := make(map[string]string)
	a["public"] = readexist[m[0]]
	a["confirmed"] = readexist[m[1]]
	return a
}

/*
var postjoin = [...]string{
	"bounce",
	"moderate",
	"accept",
}

func (m PostJoin) String() string { return postjoin[m-1] }
*/
