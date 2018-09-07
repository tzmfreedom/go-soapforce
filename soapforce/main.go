package main

import (
	"github.com/k0kubun/pp"
	"github.com/tzmfreedom/go-soapforce"

	"os"
)

func main() {
	client := soapforce.NewClient()
	res := client.Login(os.Getenv("SALESFORCE_USERNAME"), os.Getenv("SALESFORCE_PASSWORD"))
	pp.Print(res)
	sobjects := []*soapforce.SObject{
		{
			Type: "Account",
			Extra: map[string]string{
				"Name": "Hoge",
			},
		},
	}
	sResult := client.Create(sobjects)
	pp.Print(sResult)
}
