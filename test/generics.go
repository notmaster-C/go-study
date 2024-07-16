package test

import (
	"fmt"
)

type intA int
type MyInt interface {
	~int | ~int8 | int16 | int32 | int64
}
type MySlice[T int | float64 | float32] []T

type MyMap[KEY int | string, VALUE float32 | float64] map[KEY]VALUE

func TestGenerics() {
	var a MySlice[float32] = []float32{1, 2, 3}
	fmt.Println(a, "Type: %T", a)
	var m1 MyMap[string, float64] = map[string]float64{
		"go":   9.9,
		"java": 9.8,
	}
	fmt.Println(m1)

	//无意义泛型
	type W[T int | string] int
	var w W[int] = 1
	var w2 W[string] = 2
	fmt.Println(w, w2)
	// 下面这样写 就会报错
	// var w3 W[string] = "123"

	var s MySlice[int] = []int{1, 2, 3, 4}
	fmt.Println(s.Sum())
	// 可推断就不写类型 int int8 int64这种推断有误的需要写
	fmt.Println(Add[int](1, 2))
	fmt.Println(Add("1", "2"))

	fmt.Println(GetMaxNum(1, 2))

	var b1 intA = 1
	var b2 intA = 2
	fmt.Println(GetMaxNum(b1, b2))

	fmt.Println(GetMaxNum(b1, 2))
}
func printArray[T any](arr []T) {
	for _, v := range arr {
		fmt.Println(v)
	}
}

func Add[T int | float64 | string](a T, b T) T {
	return a + b
}

func (s MySlice[T]) Sum() (sum T) {
	for _, v := range s {
		sum += v
	}
	return
}

func GetMaxNum[T MyInt](a T, b T) T {
	if a > b {
		return a
	}
	return b
}
