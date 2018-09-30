# go-comments

go-comments is a tool for adding templates to `should-commented` variable names, structure definitions, and function names

# example

If you write the following code

```
package main

import "fmt"

type Aoge int

const (
	Hoge = "Hoge"
	hoge = "hoge"
	Hage = "Hage"
)

const Fugo = "Fugo"

type Fuga struct {
	Name string
}

type fuga struct {
	name string
}

func main() {
	fmt.Println("vim-go")
}

func Foo() {

}

func bar() {

}

```

go-comments output following code

```
package main

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

}
```
