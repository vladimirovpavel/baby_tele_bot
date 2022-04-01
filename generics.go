package main

import (
	"fmt"

	"github.com/cheekybits/genny/generic"
)

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "MyType=string,int"

type MyType generic.Type

func PrintMyType(value MyType) {
	fmt.Printf("%#v", value)
}

func callPrintMyType() {
	var v MyType
	PrintMyType(v)
}
