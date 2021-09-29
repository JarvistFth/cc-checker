package main

import (
	"fmt"
	"time"
)

func call(s string){

	sink(s)
}

func sink(s string){
	fmt.Println(s)
}

func main(){
	t := time.Now().String()
	call(t)



	i := 1
	j := 2
	k := 3

	a := &i
	var b *int
	b = &k

	a = &j

	p := &a

	q := &b

	p = q

	c := *q

	fmt.Printf("%p %p\n", p,c)


}