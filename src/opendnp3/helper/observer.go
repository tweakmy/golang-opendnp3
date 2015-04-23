/*
Copyright (C) 2014 Jo Ee Liew liewjoee@yahoo.com

Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
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