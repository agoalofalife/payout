package yandex

// list types deposition
const (
	TestDeps = iota
	MakeDeps
)

type TypeDeposition int
var types = [...]string {"testDeposition", "makeDeposition",}

func (typeDeps TypeDeposition) String() string {
	return types[typeDeps]
}
