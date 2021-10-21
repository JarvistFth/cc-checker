package main

import (
	"fmt"
	"time"
)

var k = "abd"

func call(s string){

	sink(s)
}

func sourceCall() string {
	 ret := time.Now().String()
	 return ret
}

func middleCall(s string) {
	sink(s)
}

func middleCallString(s string) string{
	fmt.Println(s)
	return "ok"

}

func sourceCallWithouReturn() string {
	ret := time.Now().String()
	middleCall(ret)
	return "ok"
}

func sink(s string){
	fmt.Println(s)
}

func main(){
	//t := time.Now().String()
	//call(t)

	s := sourceCall()
	s3 := middleCallString(s)
	sink(s3)

	s2 := sourceCall()
	sink(s2)
	sink(k)
	s1 := sourceCallWithouReturn()
	sink(s1)


	//i := 1
	//j := 2
	//k := 3
	//
	//a := &i
	//var b *int
	//b = &k
	//
	//a = &j
	//
	//p := &a
	//
	//q := &b
	//
	//p = q
	//
	//c := *q
	//
	//fmt.Printf("%p %p\n", p,c)


}