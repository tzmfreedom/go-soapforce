# go-soapforce

Salesforce SOAP API Client for golang

## Install

```bash
go get github.com/tzmfreedom/go-soapforce
```

## Usage

Initialize Client
```
client := soapforce.NewClient()
```

Login
```
res, err := client.Login("username", "password")
```

Create
```golang
sobjects := []*soapforce.SObject{
	{
		Type: "Account",
		Fields: map[string]string{
			"Name": "Hoge",
		},
	},
}
res, err := client.Create(sobjects)
```

## Contribute

Just send pull request if needed or fill an issue!

## License

The MIT License See [LICENSE](https://github.com/tzmfreedom/yasd/blob/master/LICENSE) file.
