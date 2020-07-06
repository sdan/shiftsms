package main

import (
	"encoding/json"
	"fmt"
)

type JSONResponse struct {
	Value1 string `json:"key1"`
	Value2 string `json:"key2"`
	Nested Nested `json:"nested"`
}

type Nested struct {
	NestValue1 string `json:"nestkey1"`
	NestValue2 string `json:"nestkey2"`
}

func main() {

	nested := Nested{
		NestValue1: "nest value 1",
		NestValue2: "nest value 2",
	}

	jsonResponse := JSONResponse{
		Value1: "value 1",
		Value2: "value 2",
		Nested: nested,
	}

	// Try uncommenting the section below and commenting out lines 21-30 the result will be the same meaning you can declare inline

	// jsonResponse := JSONResponse{
	// 	Value1: "value 1",
	// 	Value2: "value 2",
	// 	Nested: Nested{
	// 		NestValue1: "nest value 1",
	// 		NestValue2: "nest value 2",
	// 	},
	// }

	fmt.Printf("The struct returned before marshalling\n\n")
	fmt.Printf("%+v\n\n\n\n", jsonResponse)

	// The MarshalIndent function only serves to pretty print, json.Marshal() is what would normally be used
	byteArray, err := json.MarshalIndent(jsonResponse, "", "  ")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("The JSON response returned when the struct is marshalled\n\n")
	fmt.Println(string(byteArray))
}
