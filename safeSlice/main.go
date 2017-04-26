package main

import (
	"fmt"

	"github.com/lempiy/Goroutine_practice/safeSlice/slice"
)

func main() {
	mySuperSlice := safeSlice.New()
	mySuperSlice.Append(122)
	mySuperSlice.Append("Hello Go!")
	fmt.Printf("%v %v\n", (mySuperSlice.At(0)).(int), (mySuperSlice.At(1)).(string))
	fmt.Printf("len - %d\n", mySuperSlice.Len())
	mySuperSlice.Update(0, func(value interface{}) interface{} {
		value = (value).(int) * 2
		return value
	})
	fmt.Printf("%v\n", (mySuperSlice.At(0)).(int))
	mySuperSlice.Delete(0)
	fmt.Printf("len - %d, first - %v\n", mySuperSlice.Len(), (mySuperSlice.At(0)).(string))
	slice := mySuperSlice.Close()
	fmt.Printf("closed slice - %v\n", slice)
}
