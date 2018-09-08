package soapforce

import (
	"fmt"
	"io"
)

const (
	DefaultApiVersion = "38.0"
	DefaultLoginUrl   = "login.salesforce.com"
)

type Client struct {
	apiVersion string
	soapClient *Soap
	sessionId  string
	serverUrl  string
	loginUrl   string
}

func NewClient() *Client {
	soap := NewSoap("", true, nil)
	return &Client{
		soapClient: soap,
		apiVersion: DefaultApiVersion,
		loginUrl:   DefaultLoginUrl,
	}
}

func (c *Client) SetApiVersion(v string) {
	c.apiVersion = v
	c.setLoginUrl()
}

func (c *Client) SetLoginUrl(url string) {
	c.loginUrl = url
	c.setLoginUrl()
}

func (c *Client) setLoginUrl() {
	url := fmt.Sprintf("https://%s/services/Soap/u/%s", c.loginUrl, c.apiVersion)
	c.soapClient.SetServerUrl(url)
}

func (c *Client) SetDebug(debug bool) {
	c.soapClient.SetDebug(debug)
}

func (c *Client) SetLogger(logger io.Writer) {
	c.soapClient.SetLogger(logger)
}

func (c *Client) Login(u string, p string) (*LoginResult, error) {
	req := &Login{
		Username: u,
		Password: p,
	}
	res, err := c.soapClient.Login(req)
	if err != nil {
		return nil, err
	}
	c.sessionId = res.Result.SessionId
	c.serverUrl = res.Result.ServerUrl
	c.soapClient.SetServerUrl(res.Result.ServerUrl)
	sessionHeader := &SessionHeader{
		SessionId: res.Result.SessionId,
	}
	c.soapClient.SetHeader(&sessionHeader)
	return res.Result, nil
}

func (c *Client) Logout() error {
	_, err := c.soapClient.Logout(&Logout{})
	if err != nil {
		return err
	}
	c.sessionId = ""
	c.serverUrl = ""
	c.setLoginUrl()
	c.soapClient.ClearHeader()
	return nil
}

func (c *Client) DescribeSObject(s string) (*DescribeSObjectResult, error) {
	req := &DescribeSObject{
		SObjectType: s,
	}
	res, err := c.soapClient.DescribeSObject(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) DescribeGlobal(s string) (*DescribeGlobalResult, error) {
	res, err := c.soapClient.DescribeGlobal(&DescribeGlobal{})
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) DescribeLayout(s string, l string, ids []*ID) (*DescribeLayoutResult, error) {
	req := &DescribeLayout{
		SObjectType:   s,
		LayoutName:    l,
		RecordTypeIds: ids,
	}
	res, err := c.soapClient.DescribeLayout(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) Create(s []*SObject) ([]*SaveResult, error) {
	req := &Create{
		SObjects: s,
	}
	res, err := c.soapClient.Create(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) Update(s []*SObject) ([]*SaveResult, error) {
	req := &Update{
		SObjects: s,
	}
	res, err := c.soapClient.Update(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) Upsert(s []*SObject, key string) ([]*UpsertResult, error) {
	req := &Upsert{
		SObjects:            s,
		ExternalIDFieldName: key,
	}
	res, err := c.soapClient.Upsert(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) Merge(mergeReq []*MergeRequest) ([]*MergeResult, error) {
	req := &Merge{
		Request: mergeReq,
	}
	res, err := c.soapClient.Merge(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) Delete(ids []*ID) ([]*DeleteResult, error) {
	req := &Delete{
		Ids: ids,
	}
	res, err := c.soapClient.Delete(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) Undelete(ids []*ID) ([]*UndeleteResult, error) {
	req := &Undelete{
		Ids: ids,
	}
	res, err := c.soapClient.Undelete(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) Retrieve(s string, ids []*ID, fieldList string) ([]*SObject, error) {
	req := &Retrieve{
		SObjectType: s,
		Ids:         ids,
		FieldList:   fieldList,
	}
	res, err := c.soapClient.Retrieve(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) Query(q string) (*QueryResult, error) {
	req := &Query{
		QueryString: q,
	}
	res, err := c.soapClient.Query(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) QueryAll(q string) (*QueryResult, error) {
	req := &QueryAll{
		QueryString: q,
	}
	res, err := c.soapClient.QueryAll(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) QueryMore(ql *QueryLocator) (*QueryResult, error) {
	req := &QueryMore{
		QueryLocator: ql,
	}
	res, err := c.soapClient.QueryMore(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) Search(s string) (*SearchResult, error) {
	req := &Search{
		SearchString: s,
	}
	res, err := c.soapClient.Search(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) SetPassword(uid *ID, password string) (*SetPasswordResult, error) {
	req := &SetPassword{
		UserId:   uid,
		Password: password,
	}
	res, err := c.soapClient.SetPassword(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) ResetPassword(uid *ID) (*ResetPasswordResult, error) {
	req := &ResetPassword{
		UserId: uid,
	}
	res, err := c.soapClient.ResetPassword(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) GetUserInfo(s string) (*GetUserInfoResult, error) {
	res, err := c.soapClient.GetUserInfo(&GetUserInfo{})
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) SendEmailMessage(ids *ID) (*SendEmailResult, error) {
	req := &SendEmailMessage{
		Ids: ids,
	}
	res, err := c.soapClient.SendEmailMessage(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) SendEmail(m *Email) (*SendEmailResult, error) {
	req := &SendEmail{
		Messages: m,
	}
	res, err := c.soapClient.SendEmail(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}
