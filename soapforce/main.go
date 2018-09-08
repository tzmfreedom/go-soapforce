package main

import (
	"github.com/k0kubun/pp"
	"github.com/tzmfreedom/go-soapforce"

	"os"
)

var client = soapforce.NewClient()

func main() {
	client.SetDebug(true)
	res, err := client.Login(os.Getenv("SALESFORCE_USERNAME"), os.Getenv("SALESFORCE_PASSWORD"))
	if err != nil {
		panic(err)
	}
	pp.Print(res)
}

func create() {
	sobjects := []*soapforce.SObject{
		{
			Type: "Account",
			Fields: map[string]string{
				"Name": "Hoge",
			},
		},
	}
	sResult, err := client.Create(sobjects)
	if err != nil {
		panic(err)
	}
	pp.Print(sResult)
}

func update() {
	var id soapforce.ID = "001A000001WTqy6"
	sobjects := []*soapforce.SObject{
		{
			Id:   &id,
			Type: "Account",
			Fields: map[string]string{
				"Name": "popoipi!!",
			},
		},
	}
	sResult, err := client.Update(sobjects)
	if err != nil {
		panic(err)
	}
	pp.Print(sResult)
}

func upsert() {
	var id soapforce.ID = "001A000001WTqy6"
	sobjects := []*soapforce.SObject{
		{
			Id:   &id,
			Type: "Account",
			Fields: map[string]string{
				"Name": "heihei!!",
			},
		},
	}
	sResult, err := client.Upsert(sobjects, "Id")
	if err != nil {
		panic(err)
	}
	pp.Print(sResult)
}

func delete() {
	var id soapforce.ID = "001A000001WTqy6"
	ids := []*soapforce.ID{
		&id,
	}
	sResult, err := client.Delete(ids)
	if err != nil {
		panic(err)
	}
	pp.Print(sResult)

}

func undelete() {
	var id soapforce.ID = "001A000001WTqy6"
	ids := []*soapforce.ID{
		&id,
	}
	sResult, err := client.Undelete(ids)
	if err != nil {
		panic(err)
	}
	pp.Print(sResult)

}
