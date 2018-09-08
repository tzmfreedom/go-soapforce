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
	retrieve()
}

func retrieve() {
	ids := []string{ "001A000001WTqy6" }
	res, err := client.Retrieve("Account", ids, "Name, BillingAddress")
	if err != nil {
		panic(err)
	}
	pp.Print(res)
}

func describe() {
	res, err := client.DescribeGlobal()
	// res, err := client.DescribeSObject("Account")
	if err != nil {
		panic(err)
	}
	pp.Print(res)
}

func describeLayout() {
	res, err := client.DescribeLayout("Account", "", []string{})
	if err != nil {
		panic(err)
	}
	pp.Print(res)
}


func getUserInfo() {
	res, err := client.GetUserInfo()
	if err != nil {
		panic(err)
	}
	pp.Print(res)
}

func query() string {
	client.SetBatchSize(200)
	res, err := client.Query("SELECT id, Name FROM Account")
	if err != nil {
		panic(err)
	}
	pp.Print(res)
	return res.QueryLocator
}

func queryMore(ql string) string {
	res, err := client.QueryMore(ql)
	if err != nil {
		panic(err)
	}
	pp.Print(res)
	return res.QueryLocator
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
	sobjects := []*soapforce.SObject{
		{
			Id:   "001A000001WTqy6",
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
	sobjects := []*soapforce.SObject{
		{
			Id:   "001A000001WTqy6",
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
	ids := []string{
		"001A000001WSZK4",
	}
	sResult, err := client.Delete(ids)
	if err != nil {
		panic(err)
	}
	pp.Print(sResult)

}

func undelete() {
	ids := []string{
		"001A000001WTqy6",
	}
	sResult, err := client.Undelete(ids)
	if err != nil {
		panic(err)
	}
	pp.Print(sResult)

}
