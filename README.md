# go-soapforce

Salesforce SOAP API Client for golang

## Install

```bash
go get github.com/tzmfreedom/go-soapforce
```

## Usage

import package and initialize client
```golang
import "github.com/tzmfreedom/go-soapforce"

client := soapforce.NewClient()
```

for sandbox
```golang
client.SetLoginUrl("test.salesforce.com")
```

set api version
```golang
client.SetApiVersion("44.0")
```

debug request/response
```golang
client.SetDebug(true)
```

Login
```golang
res, err := client.Login("username", "password")
```

Logout
```golang
res, err := client.Logout()
```

DescribeSObject
```golang
res, err := client.DescribeSObject("Account")
```

DescribeGlobal
```golang
res, err := client.DescribeGlobal()
```

DescribeLayout
```golang
recordTypeIds := []string{}
res, err := client.DescribeLayout("Account", "layout_name", recordTypeIds)
```

Create
```golang
sobjects := []*soapforce.SObject{
	{
		Type: "Account",
		Fields: map[string]string{
			"Name": "Foo",
		},
	},
}
res, err := client.Create(sobjects)
```

Update
```golang
sobjects := []*soapforce.SObject{
	{
		Id:   "001xxxxxxxxxxxxxxx",
		Type: "Account",
		Fields: map[string]string{
			"Name": "Updated Name",
		},
	},
}
sResult, err := client.Update(sobjects)
```

Upsert
```golang
sobjects := []*soapforce.SObject{
	{
		Id:   "001xxxxxxxxxxxxxxx",
		Type: "Account",
		Fields: map[string]string{
			"Name": "Upserted Name",
		},
	},
}
sResult, err := client.Upsert(sobjects, "Id")
```

Delete
```golang
ids := []string{
	"001xxxxxxxxxxxxxxx",
}
sResult, err := client.Delete(ids)
```

Undelete
```golang
ids := []string{
	"001xxxxxxxxxxxxxxx",
}
sResult, err := client.Undelete(ids)
```

Query
```golang
res, err := client.Query("SELECT id, Name FROM Account")
```

Set BatchSize
```golang
client.SetBatchSize(200)
```

QueryMore
```golang
res, err := client.Query("SELECT id FROM Account")
res, err = client.QueryMore(res.ql)
```

Retrieve
```golang
ids := []string{ "001A000001WTqy6" }
res, err := client.Retrieve("Account", ids, "Name, BillingAddress")
```

GetUserInfo
```golang
res, err := client.GetUserInfo()
```

## Contribute

Just send pull request if needed or fill an issue!

## License

The MIT License See [LICENSE](https://github.com/tzmfreedom/go-soapforce/blob/master/LICENSE) file.
