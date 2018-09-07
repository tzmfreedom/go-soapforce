package main

import (
	"github.com/k0kubun/pp"
	"github.com/tzmfreedom/go-soapforce"

	"os"
)

func main() {
	client := soapforce.NewClient()
	res, err := client.Login(os.Getenv("SALESFORCE_USERNAME"), os.Getenv("SALESFORCE_PASSWORD"))
	if err != nil {
		panic(err)
	}
	pp.Print(res)
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
