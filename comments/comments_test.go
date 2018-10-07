package comments_test

import (
	"reflect"
	"testing"

	"github.com/gainings/go-comments/comments"
)

func TestProcess(t *testing.T) {
	got, err := comments.Process("testdata/sample1.go", nil)
	if err != nil {
		t.Errorf("Process failed becaus of %#v", err)
	}
	want := `package main

import "fmt"

//Aoge is TODO: need to enter a comment
type Aoge int

const (
	//Hoge is TODO: need to enter a comment
	Hoge	= "Hoge"
	hoge	= "hoge"
	//Hage is TODO: need to enter a comment
	Hage	= "Hage"
)

//Fugo is TODO: need to enter a comment
const Fugo = "Fugo"

//Fuga is TODO: need to enter a comment
type Fuga struct {
	Name string
}

type fuga struct {
	name string
}

func main() {
	fmt.Println("vim-go")
}

//Foo is TODO: need to enter a comment
func Foo() {

}

func bar() {

}`
	if reflect.DeepEqual(want, string(got)) {
		t.Errorf("want %s \n got %s", want, string(got))
	}
}
