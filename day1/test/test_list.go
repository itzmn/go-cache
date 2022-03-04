package main

import (
	"container/list"
	"fmt"
)

func testList() {
	l := list.New()

	l.PushBack(2)
	l.PushBack(3)
	l.PushBack(4)

	back := l.Back()
	l.MoveToFront(back)
	//fmt.Println(back.Value)

	for front := l.Front(); front != nil; front = front.Next() {
		fmt.Printf("front value:%v\n", front.Value)
	}
}

func main() {
	testList()

}
