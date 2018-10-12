package databases

import "errors"

const (
	Mysql = iota
)

type TypeDatabase int
var types = [...]string {"mysql",}

func (typeDatab TypeDatabase) String() string {
	return types[typeDatab]
}
func Define(name string) (TypeDatabase, error){
	for index, nameInList := range types {
		if nameInList == name {
			return TypeDatabase(index), nil
		}
	}
	return 0, errors.New("Not found mysql driver")
}