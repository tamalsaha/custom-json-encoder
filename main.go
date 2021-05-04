package main

import (
	"encoding/json"
	"fmt"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
)

type Instance struct {
	Label string `json:"label,omitempty" tf:"Label"`
	Alert Alert  `json:"alert,omitempty" tf:"Alert"`
}

type Alert struct {
	Name string `json:"name,omitempty" tf:"Name"`
}

type instanceCodec struct {
	jsonit jsoniter.API
}

func (instanceCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return (*Alert)(ptr) == nil
}

func (ic instanceCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	alert := (*Alert)(ptr)
	var alerts []Alert
	if alert != nil {
		alerts = []Alert{*alert}
	}

	jsoniter.RegisterTypeEncoder(reflect2.TypeOf(Alert{}).String(), nil)
	byt, _ := ic.jsonit.Marshal(alerts)
	jsoniter.RegisterTypeEncoder(reflect2.TypeOf(Alert{}).String(), ic)

	stream.Write(byt)
}

func (ic instanceCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	switch iter.WhatIsNext() {
	case jsoniter.NilValue:
		iter.Skip()
		*(*Alert)(ptr) = Alert{}
		return
	case jsoniter.ArrayValue:
		alertsByte := iter.SkipAndReturnBytes()
		if len(alertsByte) > 0 {
			var alerts []Alert
			jsoniter.RegisterTypeDecoder(reflect2.TypeOf(Alert{}).String(), nil)
			ic.jsonit.Unmarshal(alertsByte, &alerts)
			jsoniter.RegisterTypeDecoder(reflect2.TypeOf(Alert{}).String(), ic)
			if len(alerts) > 0 {
				*(*Alert)(ptr) = alerts[0]
			} else {
				*(*Alert)(ptr) = Alert{}
			}
		} else {
			*(*Alert)(ptr) = Alert{}
		}
	default:
		iter.ReportError("decode Alert", "unexpected JSON type")
	}
}

func main() {
	jsonit := jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		TagKey:                 "tf",
	}.Froze()
	jsoniter.RegisterTypeDecoder(reflect2.TypeOf(Alert{}).String(), instanceCodec{
		jsonit: jsonit,
	})
	jsoniter.RegisterTypeEncoder(reflect2.TypeOf(Alert{}).String(), instanceCodec{
		jsonit: jsonit,
	})

	instance := Instance{
		Label: "test",
		Alert: Alert{
			Name: "test2",
		},
	}

	jsonByt, err := json.Marshal(instance)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonByt))

	jsonItByt, err := jsonit.Marshal(instance)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonItByt))

	jsonInstance := Instance{}
	err = json.Unmarshal(jsonByt, &jsonInstance)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", jsonInstance)

	jsonItInstance := Instance{}
	err = jsonit.Unmarshal(jsonItByt, &jsonItInstance)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", jsonItInstance)
}

