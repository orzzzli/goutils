package path

import (
	"fmt"
	"testing"
)

func TestGetFileName(t *testing.T) {
	res := GetFileName("/Users/Orz/Documents/Github/goutils/go.mod")
	fmt.Println(res)
	res = GetFileName("/Users/Orz/Documents/Github/goutils/")
	fmt.Println(res)
}

func TestGetPath(t *testing.T) {
	res := GetPath("/Users/Orz/Documents/Github/goutils/go.mod")
	fmt.Println(res)
	res = GetPath("/Users/Orz/Documents/Github/goutils/")
	fmt.Println(res)
}

func TestGetPathDiff(t *testing.T) {
	res := GetPathDiff("/Users/Orz/Documents/Github/goutils/go.mod","/Users/Orz/Documents/Github",false)
	fmt.Println(res)
	res = GetPathDiff("/Users/Orz/Documents/Github/goutils/go.mod","/Users/Orz/Documents/Github",true)
	fmt.Println(res)
}

func TestGetAllFiles(t *testing.T) {
	var temp []string
	err := GetAllFiles("/Users/Orz/Documents/Github/goutils",&temp,true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(temp)
	var temp2 []string
	err = GetAllFiles("/Users/Orz/Documents/Github/goutils",&temp2,false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(temp2)
}
