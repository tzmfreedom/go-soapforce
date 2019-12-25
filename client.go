package soapforce

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	DefaultApiVersion = "44.0"
	DefaultLoginUrl   = "login.salesforce.com"
)

type Client struct {
	UserInfo   *GetUserInfoResult
	ApiVersion string
	ServerUrl  string
	LoginUrl   string
	soapClient *Soap
	sessionId  string
	clientId string
	clientSecret string
}

func NewClient() *Client {
	soap := NewSoap("", true, nil)
	return &Client{
		soapClient: soap,
		ApiVersion: DefaultApiVersion,
		LoginUrl:   DefaultLoginUrl,
	}
}

func (c *Client) SetApiVersion(v string) {
	c.ApiVersion = v
	c.setLoginUrl()
}

func (c *Client) SetAccessToken(sid string) {
	c.sessionId = sid
	sessionHeader := &SessionHeader{
		SessionId: sid,
	}
	c.soapClient.AddHeader(&sessionHeader)
}

func (c *Client) SetLoginUrl(url string) {
	c.LoginUrl = url
	c.setLoginUrl()
}

func (c *Client) setLoginUrl() {
	url := fmt.Sprintf("https://%s/services/Soap/u/%s", c.LoginUrl, c.ApiVersion)
	c.soapClient.SetServerUrl(url)
}

func (c *Client) SetDebug(debug bool) {
	c.soapClient.SetDebug(debug)
}

func (c *Client) SetLogger(logger io.Writer) {
	c.soapClient.SetLogger(logger)
}

func (c *Client) SetGzip(gz bool) {
	c.soapClient.SetGzip(gz)
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
	c.ServerUrl = res.Result.ServerUrl
	c.soapClient.SetServerUrl(res.Result.ServerUrl)
	c.UserInfo = res.Result.UserInfo
	sessionHeader := &SessionHeader{
		SessionId: res.Result.SessionId,
	}
	c.soapClient.AddHeader(&sessionHeader)
	return res.Result, nil
}

func (c *Client) SetClientId(clientId string) {
	c.clientId = clientId
}

func (c *Client) SetClientSecret(clientSecret string) {
	c.clientSecret = clientSecret
}

func (c *Client) LoginWithOAuth(username, password string) error {
	params := url.Values{}
	params.Add("grant_type", "password")
	params.Add("client_id", c.clientId)
	params.Add("client_secret", c.clientSecret)
	params.Add("username", username)
	params.Add("password", password)
	resp, err := http.PostForm(fmt.Sprintf("https://%s/services/oauth2/token", c.LoginUrl), params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	tokenResponse := map[string]string{}
	err = json.Unmarshal(b, &tokenResponse)
	if err != nil {
		return err
	}

	c.sessionId = tokenResponse["access_token"]
	c.ServerUrl = tokenResponse["instance_url"]
	c.soapClient.SetServerUrl(c.ServerUrl)
	sessionHeader := &SessionHeader{
		SessionId: c.sessionId,
	}
	c.soapClient.AddHeader(&sessionHeader)
	return nil
}

func (c *Client) Refresh(refreshToken string) error {
	params := url.Values{}
	params.Add("grant_type", "refresh_token")
	params.Add("client_id", c.clientId)
	params.Add("client_secret", c.clientSecret)
	params.Add("refresh_token", refreshToken)
	resp, err := http.PostForm(fmt.Sprintf("https://%s/services/oauth2/token", c.LoginUrl), params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	tokenResponse := map[string]string{}
	err = json.Unmarshal(b, &tokenResponse)
	if err != nil {
		return err
	}

	c.sessionId = tokenResponse["access_token"]
	c.ServerUrl = tokenResponse["instance_url"]
	c.soapClient.SetServerUrl(c.ServerUrl)
	sessionHeader := &SessionHeader{
		SessionId: c.sessionId,
	}
	c.soapClient.AddHeader(&sessionHeader)
	return nil
}

func (c *Client) Logout() error {
	_, err := c.soapClient.Logout(&Logout{})
	if err != nil {
		return err
	}
	c.sessionId = ""
	c.ServerUrl = ""
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

func (c *Client) DescribeGlobal() (*DescribeGlobalResult, error) {
	res, err := c.soapClient.DescribeGlobal(&DescribeGlobal{})
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) DescribeLayout(s string, l string, ids []string) (*DescribeLayoutResultResult, error) {
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

func (c *Client) Delete(ids []string) ([]*DeleteResult, error) {
	req := &Delete{
		Ids: ids,
	}
	res, err := c.soapClient.Delete(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) Undelete(ids []string) ([]*UndeleteResult, error) {
	req := &Undelete{
		Ids: ids,
	}
	res, err := c.soapClient.Undelete(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) Retrieve(s string, ids []string, fieldList string) ([]*SObject, error) {
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

func (c *Client) SetBatchSize(size int) {
	queryOptions := &QueryOptions{
		BatchSize: int32(size),
	}
	c.soapClient.AddHeader(&queryOptions)
}

func (c *Client) SetDebuggingHeader(categories []*LogInfo) {
	debuggingHeaders := &DebuggingHeader{
		Categories: categories,
	}
	c.soapClient.AddHeader(&debuggingHeaders)
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

func (c *Client) QueryMore(ql string) (*QueryResult, error) {
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

func (c *Client) SetPassword(uid string, password string) (*SetPasswordResult, error) {
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

func (c *Client) ResetPassword(uid string) (*ResetPasswordResult, error) {
	req := &ResetPassword{
		UserId: uid,
	}
	res, err := c.soapClient.ResetPassword(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) GetUserInfo() (*GetUserInfoResult, error) {
	res, err := c.soapClient.GetUserInfo(&GetUserInfo{})
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) SendEmailMessage(ids string) (*SendEmailResult, error) {
	req := &SendEmailMessage{
		Ids: ids,
	}
	res, err := c.soapClient.SendEmailMessage(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) CompileAndTest(r *CompileAndTestRequest) (*CompileAndTestResult, error) {
	req := &CompileAndTest{
		CompileAndTestRequest: r,
	}
	res, err := c.soapClient.CompileAndTest(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) CompileClasses(scripts []string) ([]*CompileClassResult, error) {
	req := &CompileClasses{
		Scripts: scripts,
	}
	res, err := c.soapClient.CompileClasses(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) CompileTriggers(scripts []string) ([]*CompileTriggerResult, error) {
	req := &CompileTriggers{
		Scripts: scripts,
	}
	res, err := c.soapClient.CompileTriggers(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) ExecuteAnonymous(code string) (*ExecuteAnonymousResult, error) {
	req := &ExecuteAnonymous{
		String: code,
	}
	res, err := c.soapClient.ExecuteAnonymous(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) RunTests(r *RunTestsRequest) (*RunTestsResult, error) {
	req := &RunTests{
		RunTestsRequest: r,
	}
	res, err := c.soapClient.RunTests(req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) WsdlToApex(req *WsdlToApex) (*WsdlToApexResult, error) {
	res, err := c.soapClient.WsdlToApex(req)
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
