package main

import (
	"encoding/json"
	"fmt"
)

func main(){
	a:="string"

	b,err:=json.Marshal([]byte(a))
	if err!=nil{
		fmt.Println("err1:",err)
	}

	fmt.Println(b)

	var d []byte
	err=json.Unmarshal(b,&d)
	if err!=nil{
		fmt.Println("err2:",err)
	}
	fmt.Println(string(d))
}
