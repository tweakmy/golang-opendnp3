package helper

import (
	//"fmt"
)

//type Publisher interface {
//    Publish(value interface{})
//}

type Observer interface {
    Notify(value interface{})
}

type ObserverFunc func(value interface{})

func (fn ObserverFunc) Notify(value interface{}){
    fn(value)
}

type Observable []Observer

func (observers *Observable) AddObserver(a Observer){
    *observers = append(*observers, a)
}

func (observers *Observable) DelObserver(a Observer){
    //*observers = append(*observers, a)
	for i,key := range *observers {
		if key == a {
			println(i)
		}
	}
}

func (observers Observable) Publish(value interface{}){
    for _, obs := range observers {
        obs.Notify(value)
    }
}

//type Field struct {
//	Value int64
//	Observable
//}

//func (f *Field) Set(v int64){
//	f.Value = v
//	f.Publish(v)
//}
//
//
//func Listen(value interface{}){
//	fmt.Printf("new value 1: %v\n", value)
//}

//func Listen2(value interface{}){
//	fmt.Printf("new value 2: %v\n", value)
//}

//func main() {
//	v := &Field{}
//	v.AddObserver(ObserverFunc(Listen))
//	v.AddObserver(ObserverFunc(Listen2))
//	v.Set(105)
//	
//	fmt.Println("Hello, playground")
//}