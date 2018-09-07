package soapforce

const DefaultApiVersion = "38.0"

type Client struct {
	apiVersion string
	soapClient *Soap
	sessionId  string
	serverUrl  string
}

func NewClient() *Client {
	soap := NewSoap("", true, nil)
	return &Client{
		soapClient: soap,
		apiVersion: DefaultApiVersion,
	}
}

func (c *Client) Login(u string, p string) *LoginResult {
	req := &Login{
		Username: u,
		Password: p,
	}
	res, err := c.soapClient.Login(req)
	if err != nil {
		panic(err)
	}
	c.sessionId = res.Result.SessionId
	c.serverUrl = res.Result.ServerUrl
	c.soapClient.SetServerUrl(res.Result.ServerUrl)
	sessionHeader := &SessionHeader {
		SessionId: res.Result.SessionId,
	}
	c.soapClient.SetHeader(&sessionHeader)
	return res.Result
}

func (c *Client) Create(s []*SObject) []*SaveResult {
	req := &Create{
		SObjects: s,
	}
	res, err := c.soapClient.Create(req)
	if err != nil {
		panic(err)
	}
	return res.Result
}
