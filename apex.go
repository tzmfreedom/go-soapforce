package soapforce

import (
	"encoding/xml"
	"time"
)

// against "unused imports"
var _ time.Time
var _ xml.Name

type FlowProcessType string

const (
	FlowProcessTypeAutoLaunchedFlow FlowProcessType = "AutoLaunchedFlow"

	FlowProcessTypeFlow FlowProcessType = "Flow"

	FlowProcessTypeWorkflow FlowProcessType = "Workflow"

	FlowProcessTypeCustomEvent FlowProcessType = "CustomEvent"

	FlowProcessTypeInvocableProcess FlowProcessType = "InvocableProcess"

	FlowProcessTypeLoginFlow FlowProcessType = "LoginFlow"

	FlowProcessTypeActionPlan FlowProcessType = "ActionPlan"

	FlowProcessTypeJourneyBuilderIntegration FlowProcessType = "JourneyBuilderIntegration"

	FlowProcessTypeUserProvisioningFlow FlowProcessType = "UserProvisioningFlow"

	FlowProcessTypeSurvey FlowProcessType = "Survey"

	FlowProcessTypeForm FlowProcessType = "Form"

	FlowProcessTypeFieldServiceMobile FlowProcessType = "FieldServiceMobile"

	FlowProcessTypeOrchestrationFlow FlowProcessType = "OrchestrationFlow"

	FlowProcessTypeFieldServiceWeb FlowProcessType = "FieldServiceWeb"

	FlowProcessTypeTransactionSecurityFlow FlowProcessType = "TransactionSecurityFlow"

	FlowProcessTypeContactRequestFlow FlowProcessType = "ContactRequestFlow"
)

const (
	stringDb string = "Db"

	stringWorkflow string = "Workflow"

	stringValidation string = "Validation"

	stringCallout string = "Callout"

	stringApexcode string = "Apexcode"

	stringApexprofiling string = "Apexprofiling"

	stringVisualforce string = "Visualforce"

	stringSystem string = "System"

	stringWave string = "Wave"

	stringNba string = "Nba"

	stringAll string = "All"
)

const (
	stringLevelNone string = "None"

	stringLevelFinest string = "Finest"

	stringLevelFiner string = "Finer"

	stringLevelFine string = "Fine"

	stringLevelDebug string = "Debug"

	stringLevelInfo string = "Info"

	stringLevelWarn string = "Warn"

	stringLevelError string = "Error"
)

const (
	LogTypeNone string = "None"

	LogTypeDebugonly string = "Debugonly"

	LogTypeDb string = "Db"

	LogTypeProfiling string = "Profiling"

	LogTypeCallout string = "Callout"

	LogTypeDetail string = "Detail"
)

type CompileAndTest struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex compileAndTest"`

	CompileAndTestRequest *CompileAndTestRequest `xml:"CompileAndTestRequest,omitempty"`
}

type CompileAndTestResponse struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex compileAndTestResponse"`

	Result *CompileAndTestResult `xml:"result,omitempty"`
}

type CompileClasses struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex compileClasses"`

	Scripts []string `xml:"scripts,omitempty"`
}

type CompileClassesResponse struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex compileClassesResponse"`

	Result []*CompileClassResult `xml:"result,omitempty"`
}

type CompileTriggers struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex compileTriggers"`

	Scripts []string `xml:"scripts,omitempty"`
}

type CompileTriggersResponse struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex compileTriggersResponse"`

	Result []*CompileTriggerResult `xml:"result,omitempty"`
}

type ExecuteAnonymous struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex executeAnonymous"`

	String string `xml:"String,omitempty"`
}

type ExecuteAnonymousResponse struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex executeAnonymousResponse"`

	Result *ExecuteAnonymousResult `xml:"result,omitempty"`
}

type RunTests struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex runTests"`

	RunTestsRequest *RunTestsRequest `xml:"RunTestsRequest,omitempty"`
}

type RunTestsResponse struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex runTestsResponse"`

	Result *RunTestsResult `xml:"result,omitempty"`
}

type WsdlToApex struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex wsdlToApex"`

	Info *WsdlToApexInfo `xml:"info,omitempty"`
}

type WsdlToApexResponse struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex wsdlToApexResponse"`

	Result *WsdlToApexResult `xml:"result,omitempty"`
}

type CompileAndTestRequest struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex CompileAndTestRequest"`

	CheckOnly bool `xml:"checkOnly,omitempty"`

	Classes []string `xml:"classes,omitempty"`

	DeleteClasses []string `xml:"deleteClasses,omitempty"`

	DeleteTriggers []string `xml:"deleteTriggers,omitempty"`

	RunTestsRequest *RunTestsRequest `xml:"runTestsRequest,omitempty"`

	Triggers []string `xml:"triggers,omitempty"`
}

type RunTestsRequest struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex RunTestsRequest"`

	AllTests bool `xml:"allTests,omitempty"`

	Classes []string `xml:"classes,omitempty"`

	MaxFailedTests int32 `xml:"maxFailedTests,omitempty"`

	Namespace string `xml:"namespace,omitempty"`

	Packages []string `xml:"packages,omitempty"`

	SkipCodeCoverage bool `xml:"skipCodeCoverage,omitempty"`

	Tests []*TestsNode `xml:"tests,omitempty"`
}

type TestsNode struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex TestsNode"`

	ClassId string `xml:"classId,omitempty"`

	ClassName string `xml:"className,omitempty"`

	TestMethods []string `xml:"testMethods,omitempty"`
}

type CompileAndTestResult struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex CompileAndTestResult"`

	Classes []*CompileClassResult `xml:"classes,omitempty"`

	DeleteClasses []*DeleteApexResult `xml:"deleteClasses,omitempty"`

	DeleteTriggers []*DeleteApexResult `xml:"deleteTriggers,omitempty"`

	RunTestsResult *RunTestsResult `xml:"runTestsResult,omitempty"`

	Success bool `xml:"success,omitempty"`

	Triggers []*CompileTriggerResult `xml:"triggers,omitempty"`
}

type CompileClassResult struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex CompileClassResult"`

	BodyCrc int32 `xml:"bodyCrc,omitempty"`

	Column int32 `xml:"column,omitempty"`

	Id string `xml:"id,omitempty"`

	Line int32 `xml:"line,omitempty"`

	Name string `xml:"name,omitempty"`

	Problem string `xml:"problem,omitempty"`

	Problems []*CompileIssue `xml:"problems,omitempty"`

	Success bool `xml:"success,omitempty"`

	Warnings []*CompileIssue `xml:"warnings,omitempty"`
}

type CompileIssue struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex CompileIssue"`

	Column int32 `xml:"column,omitempty"`

	Line int32 `xml:"line,omitempty"`

	Message string `xml:"message,omitempty"`
}

type DeleteApexResult struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex DeleteApexResult"`

	Id string `xml:"id,omitempty"`

	Problem string `xml:"problem,omitempty"`

	Success bool `xml:"success,omitempty"`
}

type RunTestsResult struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex RunTestsResult"`

	ApexLogId string `xml:"apexLogId,omitempty"`

	CodeCoverage []*CodeCoverageResult `xml:"codeCoverage,omitempty"`

	CodeCoverageWarnings []*CodeCoverageWarning `xml:"codeCoverageWarnings,omitempty"`

	Failures []*RunTestFailure `xml:"failures,omitempty"`

	FlowCoverage []*FlowCoverageResult `xml:"flowCoverage,omitempty"`

	FlowCoverageWarnings []*FlowCoverageWarning `xml:"flowCoverageWarnings,omitempty"`

	NumFailures int32 `xml:"numFailures,omitempty"`

	NumTestsRun int32 `xml:"numTestsRun,omitempty"`

	Successes []*RunTestSuccess `xml:"successes,omitempty"`

	TotalTime float64 `xml:"totalTime,omitempty"`
}

type CodeCoverageResult struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex CodeCoverageResult"`

	Id string `xml:"id,omitempty"`

	LocationsNotCovered []*CodeLocation `xml:"locationsNotCovered,omitempty"`

	Name string `xml:"name,omitempty"`

	Namespace string `xml:"namespace,omitempty"`

	NumLocations int32 `xml:"numLocations,omitempty"`

	NumLocationsNotCovered int32 `xml:"numLocationsNotCovered,omitempty"`

	Type_ string `xml:"type,omitempty"`
}

type CodeLocation struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex CodeLocation"`

	Column int32 `xml:"column,omitempty"`

	Line int32 `xml:"line,omitempty"`

	NumExecutions int32 `xml:"numExecutions,omitempty"`

	Time float64 `xml:"time,omitempty"`
}

type CodeCoverageWarning struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex CodeCoverageWarning"`

	Id string `xml:"id,omitempty"`

	Message string `xml:"message,omitempty"`

	Name string `xml:"name,omitempty"`

	Namespace string `xml:"namespace,omitempty"`
}

type RunTestFailure struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex RunTestFailure"`

	Id string `xml:"id,omitempty"`

	Message string `xml:"message,omitempty"`

	MethodName string `xml:"methodName,omitempty"`

	Name string `xml:"name,omitempty"`

	Namespace string `xml:"namespace,omitempty"`

	SeeAllData bool `xml:"seeAllData,omitempty"`

	StackTrace string `xml:"stackTrace,omitempty"`

	Time float64 `xml:"time,omitempty"`

	Type_ string `xml:"type,omitempty"`
}

type FlowCoverageResult struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex FlowCoverageResult"`

	ElementsNotCovered []string `xml:"elementsNotCovered,omitempty"`

	FlowId string `xml:"flowId,omitempty"`

	FlowName string `xml:"flowName,omitempty"`

	FlowNamespace string `xml:"flowNamespace,omitempty"`

	NumElements int32 `xml:"numElements,omitempty"`

	NumElementsNotCovered int32 `xml:"numElementsNotCovered,omitempty"`

	ProcessType *FlowProcessType `xml:"processType,omitempty"`
}

type FlowCoverageWarning struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex FlowCoverageWarning"`

	FlowId string `xml:"flowId,omitempty"`

	FlowName string `xml:"flowName,omitempty"`

	FlowNamespace string `xml:"flowNamespace,omitempty"`

	Message string `xml:"message,omitempty"`
}

type RunTestSuccess struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex RunTestSuccess"`

	Id string `xml:"id,omitempty"`

	MethodName string `xml:"methodName,omitempty"`

	Name string `xml:"name,omitempty"`

	Namespace string `xml:"namespace,omitempty"`

	SeeAllData bool `xml:"seeAllData,omitempty"`

	Time float64 `xml:"time,omitempty"`
}

type CompileTriggerResult struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex CompileTriggerResult"`

	BodyCrc int32 `xml:"bodyCrc,omitempty"`

	Column int32 `xml:"column,omitempty"`

	Id string `xml:"id,omitempty"`

	Line int32 `xml:"line,omitempty"`

	Name string `xml:"name,omitempty"`

	Problem string `xml:"problem,omitempty"`

	Problems []*CompileIssue `xml:"problems,omitempty"`

	Success bool `xml:"success,omitempty"`

	Warnings []*CompileIssue `xml:"warnings,omitempty"`
}

type ExecuteAnonymousResult struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex ExecuteAnonymousResult"`

	Column int32 `xml:"column,omitempty"`

	CompileProblem string `xml:"compileProblem,omitempty"`

	Compiled bool `xml:"compiled,omitempty"`

	ExceptionMessage string `xml:"exceptionMessage,omitempty"`

	ExceptionStackTrace string `xml:"exceptionStackTrace,omitempty"`

	Line int32 `xml:"line,omitempty"`

	Success bool `xml:"success,omitempty"`
}

type WsdlToApexInfo struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex WsdlToApexInfo"`

	Mapping []*NamespacePackagePair `xml:"mapping,omitempty"`

	Wsdl string `xml:"wsdl,omitempty"`
}

type NamespacePackagePair struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex NamespacePackagePair"`

	Namespace string `xml:"namespace,omitempty"`

	PackageName string `xml:"packageName,omitempty"`
}

type WsdlToApexResult struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2006/08/apex WsdlToApexResult"`

	ApexScripts []string `xml:"apexScripts,omitempty"`

	Errors []string `xml:"errors,omitempty"`

	Success bool `xml:"success,omitempty"`
}

/* Compile one or more Apex Classes, Triggers, and run tests. */
func (s *Soap) CompileAndTest(request *CompileAndTest) (*CompileAndTestResponse, error) {
	response := new(CompileAndTestResponse)
	err := s.client.Call(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/* Compile one or more Apex Classes. */
func (s *Soap) CompileClasses(request *CompileClasses) (*CompileClassesResponse, error) {
	response := new(CompileClassesResponse)
	err := s.client.Call(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/* Compile Apex Trigger code blocks. */
func (s *Soap) CompileTriggers(request *CompileTriggers) (*CompileTriggersResponse, error) {
	response := new(CompileTriggersResponse)
	err := s.client.Call(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/* Execute an anonymous Apex code block */
func (s *Soap) ExecuteAnonymous(request *ExecuteAnonymous) (*ExecuteAnonymousResponse, error) {
	response := new(ExecuteAnonymousResponse)
	err := s.client.Call(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/* Execute test methods */
func (s *Soap) RunTests(request *RunTests) (*RunTestsResponse, error) {
	response := new(RunTestsResponse)
	err := s.client.Call(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/* Generate Apex packages from WSDL for web s callouts */
func (s *Soap) WsdlToApex(request *WsdlToApex) (*WsdlToApexResponse, error) {
	response := new(WsdlToApexResponse)
	err := s.client.Call(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
