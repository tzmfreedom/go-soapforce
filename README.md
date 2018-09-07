# go-soapforce

Salesforce SOAP API Client for golang

## Install

```bash
go get github.com/tzmfreedom/go-soapforce
```

## Usage

Initialize Client
```golang
client := soapforce.NewClient()
```

for sandbox
```golang
client.SetLoginUrl("test.salesforce.com")
```

set api version
```golang
client.SetApiVersion("38.0")
```

Login
```golang
res, err := client.Login("username", "password")
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
```

## Contribute

Just send pull request if needed or fill an issue!

## License

The MIT License See [LICENSE](https://github.com/tzmfreedom/go-soapforce/blob/master/LICENSE) file.
