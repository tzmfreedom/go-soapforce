package soapforce

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/xml"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"
)

// against "unused imports"
var _ time.Time
var _ xml.Name

type SObject struct {
	Type string `xml:"type,omitempty"`

	FieldsToNull []string `xml:"fieldsToNull,omitempty"`

	Id string `xml:"Id,omitempty"`

	Fields map[string]interface{}
}

func (s *SObject) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "sObjects"
	start.Name.Space = "urn:sobject.partner.soap.sforce.com"
	e.EncodeToken(start)
	e.EncodeElement(s.Type, xml.StartElement{Name: xml.Name{Local: "type"}})
	if s.Id != "" {
		e.EncodeElement(s.Id, xml.StartElement{Name: xml.Name{Local: "Id"}})
	}
	if s.FieldsToNull != nil {
		for _, v := range s.FieldsToNull {
			e.EncodeElement(v, xml.StartElement{Name: xml.Name{Local: "fieldsToNull"}})
		}
	}
	for k, v := range s.Fields {
		if _, ok := v.(string); ok {
			e.EncodeElement(v, xml.StartElement{Name: xml.Name{Local: k}})
		} else if obj, ok := v.(map[string]string); ok {
			ref := xml.StartElement{Name: xml.Name{Local: k}}
			e.EncodeToken(ref)
			for innerKey, innerValue := range obj {
				e.EncodeElement(innerValue, xml.StartElement{Name: xml.Name{Local: innerKey}})
			}
			e.EncodeToken(ref.End())
		}
	}
	e.EncodeToken(start.End())
	return nil
}

func (s *SObject) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return decodeSObject(d, s, "")
}

func decodeSObject(d *xml.Decoder, s *SObject, objectName string) error {
	s.Fields = make(map[string]interface{})
	for {
		token, err := d.Token()
		if token == nil {
			break
		}
		if err != nil {
			return err
		}
		if t, ok := token.(xml.EndElement); ok {
			if t.Name.Local == objectName {
				break
			}
		}
		if t, ok := token.(xml.StartElement); ok {
			if t.Name.Local == "Id" || t.Name.Local == "type" {
				var v string
				if err := d.DecodeElement(&v, &t); err != nil {
					return err
				}
				if t.Name.Local == "Id" {
					s.Id = v
				} else if t.Name.Local == "type" {
					s.Type = v
				}
			} else {
				if len(t.Attr) > 0 {
					a := t.Attr[0]
					if a.Name.Local == "type" && a.Value == "sf:sObject" {
						v := &SObject{}
						err := decodeSObject(d, v, t.Name.Local)
						if err != nil {
							return err
						}
						s.Fields[t.Name.Local] = v
					} else if a.Name.Local == "type" && a.Value == "QueryResult" {
						v := &QueryResult{}
						if err := d.DecodeElement(v, &t); err != nil {
							return err
						}
						s.Fields[t.Name.Local] = v
					} else {
						var v string
						if err := d.DecodeElement(&v, &t); err != nil {
							return err
						}
						s.Fields[t.Name.Local] = v
					}
				} else {
					var v string
					if err := d.DecodeElement(&v, &t); err != nil {
						return err
					}
					s.Fields[t.Name.Local] = v
				}
			}
		}
	}
	return nil
}

type StatusCode string

const (
	StatusCodeALL_OR_NONE_OPERATION_ROLLED_BACK StatusCode = "ALL_OR_NONE_OPERATION_ROLLED_BACK"

	StatusCodeALREADY_IN_PROCESS StatusCode = "ALREADY_IN_PROCESS"

	StatusCodeAPEX_DATA_ACCESS_RESTRICTION StatusCode = "APEX_DATA_ACCESS_RESTRICTION"

	StatusCodeASSIGNEE_TYPE_REQUIRED StatusCode = "ASSIGNEE_TYPE_REQUIRED"

	StatusCodeAURA_COMPILE_ERROR StatusCode = "AURA_COMPILE_ERROR"

	StatusCodeBAD_CUSTOM_ENTITY_PARENT_DOMAIN StatusCode = "BAD_CUSTOM_ENTITY_PARENT_DOMAIN"

	StatusCodeBCC_NOT_ALLOWED_IF_BCC_COMPLIANCE_ENABLED StatusCode = "BCC_NOT_ALLOWED_IF_BCC_COMPLIANCE_ENABLED"

	StatusCodeCANNOT_CASCADE_PRODUCT_ACTIVE StatusCode = "CANNOT_CASCADE_PRODUCT_ACTIVE"

	StatusCodeCANNOT_CHANGE_FIELD_TYPE_OF_APEX_REFERENCED_FIELD StatusCode = "CANNOT_CHANGE_FIELD_TYPE_OF_APEX_REFERENCED_FIELD"

	StatusCodeCANNOT_CHANGE_FIELD_TYPE_OF_REFERENCED_FIELD StatusCode = "CANNOT_CHANGE_FIELD_TYPE_OF_REFERENCED_FIELD"

	StatusCodeCANNOT_CREATE_ANOTHER_MANAGED_PACKAGE StatusCode = "CANNOT_CREATE_ANOTHER_MANAGED_PACKAGE"

	StatusCodeCANNOT_DEACTIVATE_DIVISION StatusCode = "CANNOT_DEACTIVATE_DIVISION"

	StatusCodeCANNOT_DELETE_GLOBAL_ACTION_LIST StatusCode = "CANNOT_DELETE_GLOBAL_ACTION_LIST"

	StatusCodeCANNOT_DELETE_LAST_DATED_CONVERSION_RATE StatusCode = "CANNOT_DELETE_LAST_DATED_CONVERSION_RATE"

	StatusCodeCANNOT_DELETE_MANAGED_OBJECT StatusCode = "CANNOT_DELETE_MANAGED_OBJECT"

	StatusCodeCANNOT_DISABLE_LAST_ADMIN StatusCode = "CANNOT_DISABLE_LAST_ADMIN"

	StatusCodeCANNOT_ENABLE_IP_RESTRICT_REQUESTS StatusCode = "CANNOT_ENABLE_IP_RESTRICT_REQUESTS"

	StatusCodeCANNOT_EXECUTE_FLOW_TRIGGER StatusCode = "CANNOT_EXECUTE_FLOW_TRIGGER"

	StatusCodeCANNOT_FREEZE_SELF StatusCode = "CANNOT_FREEZE_SELF"

	StatusCodeCANNOT_INSERT_UPDATE_ACTIVATE_ENTITY StatusCode = "CANNOT_INSERT_UPDATE_ACTIVATE_ENTITY"

	StatusCodeCANNOT_MODIFY_MANAGED_OBJECT StatusCode = "CANNOT_MODIFY_MANAGED_OBJECT"

	StatusCodeCANNOT_PASSWORD_LOCKOUT StatusCode = "CANNOT_PASSWORD_LOCKOUT"

	StatusCodeCANNOT_POST_TO_ARCHIVED_GROUP StatusCode = "CANNOT_POST_TO_ARCHIVED_GROUP"

	StatusCodeCANNOT_RENAME_APEX_REFERENCED_FIELD StatusCode = "CANNOT_RENAME_APEX_REFERENCED_FIELD"

	StatusCodeCANNOT_RENAME_APEX_REFERENCED_OBJECT StatusCode = "CANNOT_RENAME_APEX_REFERENCED_OBJECT"

	StatusCodeCANNOT_RENAME_REFERENCED_FIELD StatusCode = "CANNOT_RENAME_REFERENCED_FIELD"

	StatusCodeCANNOT_RENAME_REFERENCED_OBJECT StatusCode = "CANNOT_RENAME_REFERENCED_OBJECT"

	StatusCodeCANNOT_REPARENT_RECORD StatusCode = "CANNOT_REPARENT_RECORD"

	StatusCodeCANNOT_UPDATE_CONVERTED_LEAD StatusCode = "CANNOT_UPDATE_CONVERTED_LEAD"

	StatusCodeCANT_DISABLE_CORP_CURRENCY StatusCode = "CANT_DISABLE_CORP_CURRENCY"

	StatusCodeCANT_UNSET_CORP_CURRENCY StatusCode = "CANT_UNSET_CORP_CURRENCY"

	StatusCodeCHILD_SHARE_FAILS_PARENT StatusCode = "CHILD_SHARE_FAILS_PARENT"

	StatusCodeCIRCULAR_DEPENDENCY StatusCode = "CIRCULAR_DEPENDENCY"

	StatusCodeCLEAN_SERVICE_ERROR StatusCode = "CLEAN_SERVICE_ERROR"

	StatusCodeCOLLISION_DETECTED StatusCode = "COLLISION_DETECTED"

	StatusCodeCOMMUNITY_NOT_ACCESSIBLE StatusCode = "COMMUNITY_NOT_ACCESSIBLE"

	StatusCodeCONFLICTING_ENVIRONMENT_HUB_MEMBER StatusCode = "CONFLICTING_ENVIRONMENT_HUB_MEMBER"

	StatusCodeCONFLICTING_SSO_USER_MAPPING StatusCode = "CONFLICTING_SSO_USER_MAPPING"

	StatusCodeCUSTOM_APEX_ERROR StatusCode = "CUSTOM_APEX_ERROR"

	StatusCodeCUSTOM_CLOB_FIELD_LIMIT_EXCEEDED StatusCode = "CUSTOM_CLOB_FIELD_LIMIT_EXCEEDED"

	StatusCodeCUSTOM_ENTITY_OR_FIELD_LIMIT StatusCode = "CUSTOM_ENTITY_OR_FIELD_LIMIT"

	StatusCodeCUSTOM_FIELD_INDEX_LIMIT_EXCEEDED StatusCode = "CUSTOM_FIELD_INDEX_LIMIT_EXCEEDED"

	StatusCodeCUSTOM_INDEX_EXISTS StatusCode = "CUSTOM_INDEX_EXISTS"

	StatusCodeCUSTOM_LINK_LIMIT_EXCEEDED StatusCode = "CUSTOM_LINK_LIMIT_EXCEEDED"

	StatusCodeCUSTOM_METADATA_LIMIT_EXCEEDED StatusCode = "CUSTOM_METADATA_LIMIT_EXCEEDED"

	StatusCodeCUSTOM_SETTINGS_LIMIT_EXCEEDED StatusCode = "CUSTOM_SETTINGS_LIMIT_EXCEEDED"

	StatusCodeCUSTOM_TAB_LIMIT_EXCEEDED StatusCode = "CUSTOM_TAB_LIMIT_EXCEEDED"

	StatusCodeDATACLOUDADDRESS_NO_RECORDS_FOUND StatusCode = "DATACLOUDADDRESS_NO_RECORDS_FOUND"

	StatusCodeDATACLOUDADDRESS_PROCESSING_ERROR StatusCode = "DATACLOUDADDRESS_PROCESSING_ERROR"

	StatusCodeDATACLOUDADDRESS_SERVER_ERROR StatusCode = "DATACLOUDADDRESS_SERVER_ERROR"

	StatusCodeDELETE_FAILED StatusCode = "DELETE_FAILED"

	StatusCodeDELETE_OPERATION_TOO_LARGE StatusCode = "DELETE_OPERATION_TOO_LARGE"

	StatusCodeDELETE_REQUIRED_ON_CASCADE StatusCode = "DELETE_REQUIRED_ON_CASCADE"

	StatusCodeDEPENDENCY_EXISTS StatusCode = "DEPENDENCY_EXISTS"

	StatusCodeDUPLICATES_DETECTED StatusCode = "DUPLICATES_DETECTED"

	StatusCodeDUPLICATE_CASE_SOLUTION StatusCode = "DUPLICATE_CASE_SOLUTION"

	StatusCodeDUPLICATE_COMM_NICKNAME StatusCode = "DUPLICATE_COMM_NICKNAME"

	StatusCodeDUPLICATE_CUSTOM_ENTITY_DEFINITION StatusCode = "DUPLICATE_CUSTOM_ENTITY_DEFINITION"

	StatusCodeDUPLICATE_CUSTOM_TAB_MOTIF StatusCode = "DUPLICATE_CUSTOM_TAB_MOTIF"

	StatusCodeDUPLICATE_DEVELOPER_NAME StatusCode = "DUPLICATE_DEVELOPER_NAME"

	StatusCodeDUPLICATE_EXTERNAL_ID StatusCode = "DUPLICATE_EXTERNAL_ID"

	StatusCodeDUPLICATE_MASTER_LABEL StatusCode = "DUPLICATE_MASTER_LABEL"

	StatusCodeDUPLICATE_SENDER_DISPLAY_NAME StatusCode = "DUPLICATE_SENDER_DISPLAY_NAME"

	StatusCodeDUPLICATE_USERNAME StatusCode = "DUPLICATE_USERNAME"

	StatusCodeDUPLICATE_VALUE StatusCode = "DUPLICATE_VALUE"

	StatusCodeEMAIL_ADDRESS_BOUNCED StatusCode = "EMAIL_ADDRESS_BOUNCED"

	StatusCodeEMAIL_EXTERNAL_TRANSPORT_CONNECTION_ERROR StatusCode = "EMAIL_EXTERNAL_TRANSPORT_CONNECTION_ERROR"

	StatusCodeEMAIL_EXTERNAL_TRANSPORT_TOKEN_ERROR StatusCode = "EMAIL_EXTERNAL_TRANSPORT_TOKEN_ERROR"

	StatusCodeEMAIL_EXTERNAL_TRANSPORT_TOO_MANY_REQUESTS_ERROR StatusCode = "EMAIL_EXTERNAL_TRANSPORT_TOO_MANY_REQUESTS_ERROR"

	StatusCodeEMAIL_EXTERNAL_TRANSPORT_UNKNOWN_ERROR StatusCode = "EMAIL_EXTERNAL_TRANSPORT_UNKNOWN_ERROR"

	StatusCodeEMAIL_NOT_PROCESSED_DUE_TO_PRIOR_ERROR StatusCode = "EMAIL_NOT_PROCESSED_DUE_TO_PRIOR_ERROR"

	StatusCodeEMAIL_OPTED_OUT StatusCode = "EMAIL_OPTED_OUT"

	StatusCodeEMAIL_TEMPLATE_FORMULA_ERROR StatusCode = "EMAIL_TEMPLATE_FORMULA_ERROR"

	StatusCodeEMAIL_TEMPLATE_MERGEFIELD_ACCESS_ERROR StatusCode = "EMAIL_TEMPLATE_MERGEFIELD_ACCESS_ERROR"

	StatusCodeEMAIL_TEMPLATE_MERGEFIELD_ERROR StatusCode = "EMAIL_TEMPLATE_MERGEFIELD_ERROR"

	StatusCodeEMAIL_TEMPLATE_MERGEFIELD_VALUE_ERROR StatusCode = "EMAIL_TEMPLATE_MERGEFIELD_VALUE_ERROR"

	StatusCodeEMAIL_TEMPLATE_PROCESSING_ERROR StatusCode = "EMAIL_TEMPLATE_PROCESSING_ERROR"

	StatusCodeEMPTY_SCONTROL_FILE_NAME StatusCode = "EMPTY_SCONTROL_FILE_NAME"

	StatusCodeENTITY_FAILED_IFLASTMODIFIED_ON_UPDATE StatusCode = "ENTITY_FAILED_IFLASTMODIFIED_ON_UPDATE"

	StatusCodeENTITY_IS_ARCHIVED StatusCode = "ENTITY_IS_ARCHIVED"

	StatusCodeENTITY_IS_DELETED StatusCode = "ENTITY_IS_DELETED"

	StatusCodeENTITY_IS_LOCKED StatusCode = "ENTITY_IS_LOCKED"

	StatusCodeENTITY_SAVE_ERROR StatusCode = "ENTITY_SAVE_ERROR"

	StatusCodeENTITY_SAVE_VALIDATION_ERROR StatusCode = "ENTITY_SAVE_VALIDATION_ERROR"

	StatusCodeENVIRONMENT_HUB_MEMBERSHIP_CONFLICT StatusCode = "ENVIRONMENT_HUB_MEMBERSHIP_CONFLICT"

	StatusCodeENVIRONMENT_HUB_MEMBERSHIP_ERROR_JOINING_HUB StatusCode = "ENVIRONMENT_HUB_MEMBERSHIP_ERROR_JOINING_HUB"

	StatusCodeENVIRONMENT_HUB_MEMBERSHIP_USER_ALREADY_IN_HUB StatusCode = "ENVIRONMENT_HUB_MEMBERSHIP_USER_ALREADY_IN_HUB"

	StatusCodeENVIRONMENT_HUB_MEMBERSHIP_USER_NOT_ORG_ADMIN StatusCode = "ENVIRONMENT_HUB_MEMBERSHIP_USER_NOT_ORG_ADMIN"

	StatusCodeERROR_IN_MAILER StatusCode = "ERROR_IN_MAILER"

	StatusCodeEXCHANGE_WEB_SERVICES_URL_INVALID StatusCode = "EXCHANGE_WEB_SERVICES_URL_INVALID"

	StatusCodeFAILED_ACTIVATION StatusCode = "FAILED_ACTIVATION"

	StatusCodeFIELD_CUSTOM_VALIDATION_EXCEPTION StatusCode = "FIELD_CUSTOM_VALIDATION_EXCEPTION"

	StatusCodeFIELD_FILTER_VALIDATION_EXCEPTION StatusCode = "FIELD_FILTER_VALIDATION_EXCEPTION"

	StatusCodeFIELD_INTEGRITY_EXCEPTION StatusCode = "FIELD_INTEGRITY_EXCEPTION"

	StatusCodeFIELD_KEYWORD_LIST_MATCH_LIMIT StatusCode = "FIELD_KEYWORD_LIST_MATCH_LIMIT"

	StatusCodeFIELD_MAPPING_ERROR StatusCode = "FIELD_MAPPING_ERROR"

	StatusCodeFIELD_MODERATION_RULE_BLOCK StatusCode = "FIELD_MODERATION_RULE_BLOCK"

	StatusCodeFIELD_NOT_UPDATABLE StatusCode = "FIELD_NOT_UPDATABLE"

	StatusCodeFILE_EXTENSION_NOT_ALLOWED StatusCode = "FILE_EXTENSION_NOT_ALLOWED"

	StatusCodeFILE_SIZE_LIMIT_EXCEEDED StatusCode = "FILE_SIZE_LIMIT_EXCEEDED"

	StatusCodeFILTERED_LOOKUP_LIMIT_EXCEEDED StatusCode = "FILTERED_LOOKUP_LIMIT_EXCEEDED"

	StatusCodeFIND_DUPLICATES_ERROR StatusCode = "FIND_DUPLICATES_ERROR"

	StatusCodeFUNCTIONALITY_NOT_ENABLED StatusCode = "FUNCTIONALITY_NOT_ENABLED"

	StatusCodeHAS_PUBLIC_REFERENCES StatusCode = "HAS_PUBLIC_REFERENCES"

	StatusCodeHTML_FILE_UPLOAD_NOT_ALLOWED StatusCode = "HTML_FILE_UPLOAD_NOT_ALLOWED"

	StatusCodeIMAGE_TOO_LARGE StatusCode = "IMAGE_TOO_LARGE"

	StatusCodeINACTIVE_OWNER_OR_USER StatusCode = "INACTIVE_OWNER_OR_USER"

	StatusCodeINACTIVE_RULE_ERROR StatusCode = "INACTIVE_RULE_ERROR"

	StatusCodeINSERT_UPDATE_DELETE_NOT_ALLOWED_DURING_MAINTENANCE StatusCode = "INSERT_UPDATE_DELETE_NOT_ALLOWED_DURING_MAINTENANCE"

	StatusCodeINSUFFICIENT_ACCESS_ON_CROSS_REFERENCE_ENTITY StatusCode = "INSUFFICIENT_ACCESS_ON_CROSS_REFERENCE_ENTITY"

	StatusCodeINSUFFICIENT_ACCESS_OR_READONLY StatusCode = "INSUFFICIENT_ACCESS_OR_READONLY"

	StatusCodeINSUFFICIENT_ACCESS_TO_INSIGHTSEXTERNALDATA StatusCode = "INSUFFICIENT_ACCESS_TO_INSIGHTSEXTERNALDATA"

	StatusCodeINSUFFICIENT_CREDITS StatusCode = "INSUFFICIENT_CREDITS"

	StatusCodeINVALID_ACCESS_LEVEL StatusCode = "INVALID_ACCESS_LEVEL"

	StatusCodeINVALID_ARGUMENT_TYPE StatusCode = "INVALID_ARGUMENT_TYPE"

	StatusCodeINVALID_ASSIGNEE_TYPE StatusCode = "INVALID_ASSIGNEE_TYPE"

	StatusCodeINVALID_ASSIGNMENT_RULE StatusCode = "INVALID_ASSIGNMENT_RULE"

	StatusCodeINVALID_BATCH_OPERATION StatusCode = "INVALID_BATCH_OPERATION"

	StatusCodeINVALID_CONTENT_TYPE StatusCode = "INVALID_CONTENT_TYPE"

	StatusCodeINVALID_CREDIT_CARD_INFO StatusCode = "INVALID_CREDIT_CARD_INFO"

	StatusCodeINVALID_CROSS_REFERENCE_KEY StatusCode = "INVALID_CROSS_REFERENCE_KEY"

	StatusCodeINVALID_CROSS_REFERENCE_TYPE_FOR_FIELD StatusCode = "INVALID_CROSS_REFERENCE_TYPE_FOR_FIELD"

	StatusCodeINVALID_CURRENCY_CONV_RATE StatusCode = "INVALID_CURRENCY_CONV_RATE"

	StatusCodeINVALID_CURRENCY_CORP_RATE StatusCode = "INVALID_CURRENCY_CORP_RATE"

	StatusCodeINVALID_CURRENCY_ISO StatusCode = "INVALID_CURRENCY_ISO"

	StatusCodeINVALID_DATA_CATEGORY_GROUP_REFERENCE StatusCode = "INVALID_DATA_CATEGORY_GROUP_REFERENCE"

	StatusCodeINVALID_DATA_URI StatusCode = "INVALID_DATA_URI"

	StatusCodeINVALID_EMAIL_ADDRESS StatusCode = "INVALID_EMAIL_ADDRESS"

	StatusCodeINVALID_EMPTY_KEY_OWNER StatusCode = "INVALID_EMPTY_KEY_OWNER"

	StatusCodeINVALID_ENTITY_FOR_MATCH_ENGINE_ERROR StatusCode = "INVALID_ENTITY_FOR_MATCH_ENGINE_ERROR"

	StatusCodeINVALID_ENTITY_FOR_MATCH_OPERATION_ERROR StatusCode = "INVALID_ENTITY_FOR_MATCH_OPERATION_ERROR"

	StatusCodeINVALID_ENTITY_FOR_UPSERT StatusCode = "INVALID_ENTITY_FOR_UPSERT"

	StatusCodeINVALID_ENVIRONMENT_HUB_MEMBER StatusCode = "INVALID_ENVIRONMENT_HUB_MEMBER"

	StatusCodeINVALID_EVENT_DELIVERY StatusCode = "INVALID_EVENT_DELIVERY"

	StatusCodeINVALID_EVENT_SUBSCRIPTION StatusCode = "INVALID_EVENT_SUBSCRIPTION"

	StatusCodeINVALID_FIELD StatusCode = "INVALID_FIELD"

	StatusCodeINVALID_FIELD_FOR_INSERT_UPDATE StatusCode = "INVALID_FIELD_FOR_INSERT_UPDATE"

	StatusCodeINVALID_FIELD_WHEN_USING_TEMPLATE StatusCode = "INVALID_FIELD_WHEN_USING_TEMPLATE"

	StatusCodeINVALID_FILTER_ACTION StatusCode = "INVALID_FILTER_ACTION"

	StatusCodeINVALID_GOOGLE_DOCS_URL StatusCode = "INVALID_GOOGLE_DOCS_URL"

	StatusCodeINVALID_ID_FIELD StatusCode = "INVALID_ID_FIELD"

	StatusCodeINVALID_INET_ADDRESS StatusCode = "INVALID_INET_ADDRESS"

	StatusCodeINVALID_INPUT StatusCode = "INVALID_INPUT"

	StatusCodeINVALID_LINEITEM_CLONE_STATE StatusCode = "INVALID_LINEITEM_CLONE_STATE"

	StatusCodeINVALID_MARKUP StatusCode = "INVALID_MARKUP"

	StatusCodeINVALID_MASTER_OR_TRANSLATED_SOLUTION StatusCode = "INVALID_MASTER_OR_TRANSLATED_SOLUTION"

	StatusCodeINVALID_MESSAGE_ID_REFERENCE StatusCode = "INVALID_MESSAGE_ID_REFERENCE"

	StatusCodeINVALID_NAMESPACE_PREFIX StatusCode = "INVALID_NAMESPACE_PREFIX"

	StatusCodeINVALID_OAUTH_URL StatusCode = "INVALID_OAUTH_URL"

	StatusCodeINVALID_OPERATION StatusCode = "INVALID_OPERATION"

	StatusCodeINVALID_OPERATOR StatusCode = "INVALID_OPERATOR"

	StatusCodeINVALID_OR_NULL_FOR_RESTRICTED_PICKLIST StatusCode = "INVALID_OR_NULL_FOR_RESTRICTED_PICKLIST"

	StatusCodeINVALID_OWNER StatusCode = "INVALID_OWNER"

	StatusCodeINVALID_PACKAGE_LICENSE StatusCode = "INVALID_PACKAGE_LICENSE"

	StatusCodeINVALID_PACKAGE_VERSION StatusCode = "INVALID_PACKAGE_VERSION"

	StatusCodeINVALID_PARTNER_NETWORK_STATUS StatusCode = "INVALID_PARTNER_NETWORK_STATUS"

	StatusCodeINVALID_PERSON_ACCOUNT_OPERATION StatusCode = "INVALID_PERSON_ACCOUNT_OPERATION"

	StatusCodeINVALID_QUERY_LOCATOR StatusCode = "INVALID_QUERY_LOCATOR"

	StatusCodeINVALID_READ_ONLY_USER_DML StatusCode = "INVALID_READ_ONLY_USER_DML"

	StatusCodeINVALID_RUNTIME_VALUE StatusCode = "INVALID_RUNTIME_VALUE"

	StatusCodeINVALID_SAVE_AS_ACTIVITY_FLAG StatusCode = "INVALID_SAVE_AS_ACTIVITY_FLAG"

	StatusCodeINVALID_SESSION_ID StatusCode = "INVALID_SESSION_ID"

	StatusCodeINVALID_SETUP_OWNER StatusCode = "INVALID_SETUP_OWNER"

	StatusCodeINVALID_SIGNUP_COUNTRY StatusCode = "INVALID_SIGNUP_COUNTRY"

	StatusCodeINVALID_SIGNUP_OPTION StatusCode = "INVALID_SIGNUP_OPTION"

	StatusCodeINVALID_SITE_DELETE_EXCEPTION StatusCode = "INVALID_SITE_DELETE_EXCEPTION"

	StatusCodeINVALID_SITE_FILE_IMPORTED_EXCEPTION StatusCode = "INVALID_SITE_FILE_IMPORTED_EXCEPTION"

	StatusCodeINVALID_SITE_FILE_TYPE_EXCEPTION StatusCode = "INVALID_SITE_FILE_TYPE_EXCEPTION"

	StatusCodeINVALID_STATUS StatusCode = "INVALID_STATUS"

	StatusCodeINVALID_SUBDOMAIN StatusCode = "INVALID_SUBDOMAIN"

	StatusCodeINVALID_TYPE StatusCode = "INVALID_TYPE"

	StatusCodeINVALID_TYPE_FOR_OPERATION StatusCode = "INVALID_TYPE_FOR_OPERATION"

	StatusCodeINVALID_TYPE_ON_FIELD_IN_RECORD StatusCode = "INVALID_TYPE_ON_FIELD_IN_RECORD"

	StatusCodeINVALID_USERID StatusCode = "INVALID_USERID"

	StatusCodeIP_RANGE_LIMIT_EXCEEDED StatusCode = "IP_RANGE_LIMIT_EXCEEDED"

	StatusCodeJIGSAW_IMPORT_LIMIT_EXCEEDED StatusCode = "JIGSAW_IMPORT_LIMIT_EXCEEDED"

	StatusCodeLICENSE_LIMIT_EXCEEDED StatusCode = "LICENSE_LIMIT_EXCEEDED"

	StatusCodeLIGHT_PORTAL_USER_EXCEPTION StatusCode = "LIGHT_PORTAL_USER_EXCEPTION"

	StatusCodeLIMIT_EXCEEDED StatusCode = "LIMIT_EXCEEDED"

	StatusCodeMALFORMED_ID StatusCode = "MALFORMED_ID"

	StatusCodeMANAGER_NOT_DEFINED StatusCode = "MANAGER_NOT_DEFINED"

	StatusCodeMASSMAIL_RETRY_LIMIT_EXCEEDED StatusCode = "MASSMAIL_RETRY_LIMIT_EXCEEDED"

	StatusCodeMASS_MAIL_LIMIT_EXCEEDED StatusCode = "MASS_MAIL_LIMIT_EXCEEDED"

	StatusCodeMATCH_DEFINITION_ERROR StatusCode = "MATCH_DEFINITION_ERROR"

	StatusCodeMATCH_OPERATION_ERROR StatusCode = "MATCH_OPERATION_ERROR"

	StatusCodeMATCH_OPERATION_INVALID_ENGINE_ERROR StatusCode = "MATCH_OPERATION_INVALID_ENGINE_ERROR"

	StatusCodeMATCH_OPERATION_INVALID_RULE_ERROR StatusCode = "MATCH_OPERATION_INVALID_RULE_ERROR"

	StatusCodeMATCH_OPERATION_MISSING_ENGINE_ERROR StatusCode = "MATCH_OPERATION_MISSING_ENGINE_ERROR"

	StatusCodeMATCH_OPERATION_MISSING_OBJECT_TYPE_ERROR StatusCode = "MATCH_OPERATION_MISSING_OBJECT_TYPE_ERROR"

	StatusCodeMATCH_OPERATION_MISSING_OPTIONS_ERROR StatusCode = "MATCH_OPERATION_MISSING_OPTIONS_ERROR"

	StatusCodeMATCH_OPERATION_MISSING_RULE_ERROR StatusCode = "MATCH_OPERATION_MISSING_RULE_ERROR"

	StatusCodeMATCH_OPERATION_UNKNOWN_RULE_ERROR StatusCode = "MATCH_OPERATION_UNKNOWN_RULE_ERROR"

	StatusCodeMATCH_OPERATION_UNSUPPORTED_VERSION_ERROR StatusCode = "MATCH_OPERATION_UNSUPPORTED_VERSION_ERROR"

	StatusCodeMATCH_PRECONDITION_FAILED StatusCode = "MATCH_PRECONDITION_FAILED"

	StatusCodeMATCH_RUNTIME_ERROR StatusCode = "MATCH_RUNTIME_ERROR"

	StatusCodeMATCH_SERVICE_ERROR StatusCode = "MATCH_SERVICE_ERROR"

	StatusCodeMATCH_SERVICE_TIMED_OUT StatusCode = "MATCH_SERVICE_TIMED_OUT"

	StatusCodeMATCH_SERVICE_UNAVAILABLE_ERROR StatusCode = "MATCH_SERVICE_UNAVAILABLE_ERROR"

	StatusCodeMAXIMUM_CCEMAILS_EXCEEDED StatusCode = "MAXIMUM_CCEMAILS_EXCEEDED"

	StatusCodeMAXIMUM_DASHBOARD_COMPONENTS_EXCEEDED StatusCode = "MAXIMUM_DASHBOARD_COMPONENTS_EXCEEDED"

	StatusCodeMAXIMUM_HIERARCHY_CHILDREN_REACHED StatusCode = "MAXIMUM_HIERARCHY_CHILDREN_REACHED"

	StatusCodeMAXIMUM_HIERARCHY_LEVELS_REACHED StatusCode = "MAXIMUM_HIERARCHY_LEVELS_REACHED"

	StatusCodeMAXIMUM_HIERARCHY_TREE_SIZE_REACHED StatusCode = "MAXIMUM_HIERARCHY_TREE_SIZE_REACHED"

	StatusCodeMAXIMUM_SIZE_OF_ATTACHMENT StatusCode = "MAXIMUM_SIZE_OF_ATTACHMENT"

	StatusCodeMAXIMUM_SIZE_OF_DOCUMENT StatusCode = "MAXIMUM_SIZE_OF_DOCUMENT"

	StatusCodeMAX_ACTIONS_PER_RULE_EXCEEDED StatusCode = "MAX_ACTIONS_PER_RULE_EXCEEDED"

	StatusCodeMAX_ACTIVE_RULES_EXCEEDED StatusCode = "MAX_ACTIVE_RULES_EXCEEDED"

	StatusCodeMAX_APPROVAL_STEPS_EXCEEDED StatusCode = "MAX_APPROVAL_STEPS_EXCEEDED"

	StatusCodeMAX_DEPTH_IN_FLOW_EXECUTION StatusCode = "MAX_DEPTH_IN_FLOW_EXECUTION"

	StatusCodeMAX_FORMULAS_PER_RULE_EXCEEDED StatusCode = "MAX_FORMULAS_PER_RULE_EXCEEDED"

	StatusCodeMAX_RULES_EXCEEDED StatusCode = "MAX_RULES_EXCEEDED"

	StatusCodeMAX_RULE_ENTRIES_EXCEEDED StatusCode = "MAX_RULE_ENTRIES_EXCEEDED"

	StatusCodeMAX_TASK_DESCRIPTION_EXCEEEDED StatusCode = "MAX_TASK_DESCRIPTION_EXCEEEDED"

	StatusCodeMAX_TM_RULES_EXCEEDED StatusCode = "MAX_TM_RULES_EXCEEDED"

	StatusCodeMAX_TM_RULE_ITEMS_EXCEEDED StatusCode = "MAX_TM_RULE_ITEMS_EXCEEDED"

	StatusCodeMERGE_FAILED StatusCode = "MERGE_FAILED"

	StatusCodeMETADATA_FIELD_UPDATE_ERROR StatusCode = "METADATA_FIELD_UPDATE_ERROR"

	StatusCodeMISSING_ARGUMENT StatusCode = "MISSING_ARGUMENT"

	StatusCodeMISSING_RECORD StatusCode = "MISSING_RECORD"

	StatusCodeMIXED_DML_OPERATION StatusCode = "MIXED_DML_OPERATION"

	StatusCodeNONUNIQUE_SHIPPING_ADDRESS StatusCode = "NONUNIQUE_SHIPPING_ADDRESS"

	StatusCodeNO_APPLICABLE_PROCESS StatusCode = "NO_APPLICABLE_PROCESS"

	StatusCodeNO_ATTACHMENT_PERMISSION StatusCode = "NO_ATTACHMENT_PERMISSION"

	StatusCodeNO_INACTIVE_DIVISION_MEMBERS StatusCode = "NO_INACTIVE_DIVISION_MEMBERS"

	StatusCodeNO_MASS_MAIL_PERMISSION StatusCode = "NO_MASS_MAIL_PERMISSION"

	StatusCodeNO_PARTNER_PERMISSION StatusCode = "NO_PARTNER_PERMISSION"

	StatusCodeNO_SUCH_USER_EXISTS StatusCode = "NO_SUCH_USER_EXISTS"

	StatusCodeNUMBER_OUTSIDE_VALID_RANGE StatusCode = "NUMBER_OUTSIDE_VALID_RANGE"

	StatusCodeNUM_HISTORY_FIELDS_BY_SOBJECT_EXCEEDED StatusCode = "NUM_HISTORY_FIELDS_BY_SOBJECT_EXCEEDED"

	StatusCodeOPTED_OUT_OF_MASS_MAIL StatusCode = "OPTED_OUT_OF_MASS_MAIL"

	StatusCodeOP_WITH_INVALID_USER_TYPE_EXCEPTION StatusCode = "OP_WITH_INVALID_USER_TYPE_EXCEPTION"

	StatusCodePACKAGE_LICENSE_REQUIRED StatusCode = "PACKAGE_LICENSE_REQUIRED"

	StatusCodePACKAGING_API_INSTALL_FAILED StatusCode = "PACKAGING_API_INSTALL_FAILED"

	StatusCodePACKAGING_API_UNINSTALL_FAILED StatusCode = "PACKAGING_API_UNINSTALL_FAILED"

	StatusCodePALI_INVALID_ACTION_ID StatusCode = "PALI_INVALID_ACTION_ID"

	StatusCodePALI_INVALID_ACTION_NAME StatusCode = "PALI_INVALID_ACTION_NAME"

	StatusCodePALI_INVALID_ACTION_TYPE StatusCode = "PALI_INVALID_ACTION_TYPE"

	StatusCodePAL_INVALID_ASSISTANT_RECOMMENDATION_TYPE_ID StatusCode = "PAL_INVALID_ASSISTANT_RECOMMENDATION_TYPE_ID"

	StatusCodePAL_INVALID_ENTITY_ID StatusCode = "PAL_INVALID_ENTITY_ID"

	StatusCodePAL_INVALID_FLEXIPAGE_ID StatusCode = "PAL_INVALID_FLEXIPAGE_ID"

	StatusCodePAL_INVALID_LAYOUT_ID StatusCode = "PAL_INVALID_LAYOUT_ID"

	StatusCodePAL_INVALID_PARAMETERS StatusCode = "PAL_INVALID_PARAMETERS"

	StatusCodePA_API_EXCEPTION StatusCode = "PA_API_EXCEPTION"

	StatusCodePA_AXIS_FAULT StatusCode = "PA_AXIS_FAULT"

	StatusCodePA_INVALID_ID_EXCEPTION StatusCode = "PA_INVALID_ID_EXCEPTION"

	StatusCodePA_NO_ACCESS_EXCEPTION StatusCode = "PA_NO_ACCESS_EXCEPTION"

	StatusCodePA_NO_DATA_FOUND_EXCEPTION StatusCode = "PA_NO_DATA_FOUND_EXCEPTION"

	StatusCodePA_URI_SYNTAX_EXCEPTION StatusCode = "PA_URI_SYNTAX_EXCEPTION"

	StatusCodePA_VISIBLE_ACTIONS_FILTER_ORDERING_EXCEPTION StatusCode = "PA_VISIBLE_ACTIONS_FILTER_ORDERING_EXCEPTION"

	StatusCodePORTAL_NO_ACCESS StatusCode = "PORTAL_NO_ACCESS"

	StatusCodePORTAL_USER_ALREADY_EXISTS_FOR_CONTACT StatusCode = "PORTAL_USER_ALREADY_EXISTS_FOR_CONTACT"

	StatusCodePORTAL_USER_CREATION_RESTRICTED_WITH_ENCRYPTION StatusCode = "PORTAL_USER_CREATION_RESTRICTED_WITH_ENCRYPTION"

	StatusCodePRIVATE_CONTACT_ON_ASSET StatusCode = "PRIVATE_CONTACT_ON_ASSET"

	StatusCodePROCESSING_HALTED StatusCode = "PROCESSING_HALTED"

	StatusCodeQA_INVALID_CREATE_FEED_ITEM StatusCode = "QA_INVALID_CREATE_FEED_ITEM"

	StatusCodeQA_INVALID_SUCCESS_MESSAGE StatusCode = "QA_INVALID_SUCCESS_MESSAGE"

	StatusCodeQUERY_TIMEOUT StatusCode = "QUERY_TIMEOUT"

	StatusCodeQUICK_ACTION_LIST_ITEM_NOT_ALLOWED StatusCode = "QUICK_ACTION_LIST_ITEM_NOT_ALLOWED"

	StatusCodeQUICK_ACTION_LIST_NOT_ALLOWED StatusCode = "QUICK_ACTION_LIST_NOT_ALLOWED"

	StatusCodeRECORD_IN_USE_BY_WORKFLOW StatusCode = "RECORD_IN_USE_BY_WORKFLOW"

	StatusCodeREL_FIELD_BAD_ACCESSIBILITY StatusCode = "REL_FIELD_BAD_ACCESSIBILITY"

	StatusCodeREPUTATION_MINIMUM_NUMBER_NOT_REACHED StatusCode = "REPUTATION_MINIMUM_NUMBER_NOT_REACHED"

	StatusCodeREQUEST_RUNNING_TOO_LONG StatusCode = "REQUEST_RUNNING_TOO_LONG"

	StatusCodeREQUIRED_FEATURE_MISSING StatusCode = "REQUIRED_FEATURE_MISSING"

	StatusCodeREQUIRED_FIELD_MISSING StatusCode = "REQUIRED_FIELD_MISSING"

	StatusCodeRETRIEVE_EXCHANGE_ATTACHMENT_FAILED StatusCode = "RETRIEVE_EXCHANGE_ATTACHMENT_FAILED"

	StatusCodeRETRIEVE_EXCHANGE_EMAIL_FAILED StatusCode = "RETRIEVE_EXCHANGE_EMAIL_FAILED"

	StatusCodeRETRIEVE_EXCHANGE_EVENT_FAILED StatusCode = "RETRIEVE_EXCHANGE_EVENT_FAILED"

	StatusCodeSALESFORCE_INBOX_TRANSPORT_CONNECTION_ERROR StatusCode = "SALESFORCE_INBOX_TRANSPORT_CONNECTION_ERROR"

	StatusCodeSALESFORCE_INBOX_TRANSPORT_TOKEN_ERROR StatusCode = "SALESFORCE_INBOX_TRANSPORT_TOKEN_ERROR"

	StatusCodeSALESFORCE_INBOX_TRANSPORT_UNKNOWN_ERROR StatusCode = "SALESFORCE_INBOX_TRANSPORT_UNKNOWN_ERROR"

	StatusCodeSELF_REFERENCE_FROM_FLOW StatusCode = "SELF_REFERENCE_FROM_FLOW"

	StatusCodeSELF_REFERENCE_FROM_TRIGGER StatusCode = "SELF_REFERENCE_FROM_TRIGGER"

	StatusCodeSHARE_NEEDED_FOR_CHILD_OWNER StatusCode = "SHARE_NEEDED_FOR_CHILD_OWNER"

	StatusCodeSINGLE_EMAIL_LIMIT_EXCEEDED StatusCode = "SINGLE_EMAIL_LIMIT_EXCEEDED"

	StatusCodeSOCIAL_ACCOUNT_NOT_FOUND StatusCode = "SOCIAL_ACCOUNT_NOT_FOUND"

	StatusCodeSOCIAL_ACTION_INVALID StatusCode = "SOCIAL_ACTION_INVALID"

	StatusCodeSOCIAL_POST_INVALID StatusCode = "SOCIAL_POST_INVALID"

	StatusCodeSOCIAL_POST_NOT_FOUND StatusCode = "SOCIAL_POST_NOT_FOUND"

	StatusCodeSTANDARD_PRICE_NOT_DEFINED StatusCode = "STANDARD_PRICE_NOT_DEFINED"

	StatusCodeSTORAGE_LIMIT_EXCEEDED StatusCode = "STORAGE_LIMIT_EXCEEDED"

	StatusCodeSTRING_TOO_LONG StatusCode = "STRING_TOO_LONG"

	StatusCodeSUBDOMAIN_IN_USE StatusCode = "SUBDOMAIN_IN_USE"

	StatusCodeTABSET_LIMIT_EXCEEDED StatusCode = "TABSET_LIMIT_EXCEEDED"

	StatusCodeTEMPLATE_NOT_ACTIVE StatusCode = "TEMPLATE_NOT_ACTIVE"

	StatusCodeTEMPLATE_NOT_FOUND StatusCode = "TEMPLATE_NOT_FOUND"

	StatusCodeTERRITORY_REALIGN_IN_PROGRESS StatusCode = "TERRITORY_REALIGN_IN_PROGRESS"

	StatusCodeTEXT_DATA_OUTSIDE_SUPPORTED_CHARSET StatusCode = "TEXT_DATA_OUTSIDE_SUPPORTED_CHARSET"

	StatusCodeTOO_MANY_APEX_REQUESTS StatusCode = "TOO_MANY_APEX_REQUESTS"

	StatusCodeTOO_MANY_ENUM_VALUE StatusCode = "TOO_MANY_ENUM_VALUE"

	StatusCodeTOO_MANY_POSSIBLE_USERS_EXIST StatusCode = "TOO_MANY_POSSIBLE_USERS_EXIST"

	StatusCodeTRANSFER_REQUIRES_READ StatusCode = "TRANSFER_REQUIRES_READ"

	StatusCodeUNABLE_TO_LOCK_ROW StatusCode = "UNABLE_TO_LOCK_ROW"

	StatusCodeUNAVAILABLE_RECORDTYPE_EXCEPTION StatusCode = "UNAVAILABLE_RECORDTYPE_EXCEPTION"

	StatusCodeUNAVAILABLE_REF StatusCode = "UNAVAILABLE_REF"

	StatusCodeUNDELETE_FAILED StatusCode = "UNDELETE_FAILED"

	StatusCodeUNKNOWN_EXCEPTION StatusCode = "UNKNOWN_EXCEPTION"

	StatusCodeUNSAFE_HTML_CONTENT StatusCode = "UNSAFE_HTML_CONTENT"

	StatusCodeUNSPECIFIED_EMAIL_ADDRESS StatusCode = "UNSPECIFIED_EMAIL_ADDRESS"

	StatusCodeUNSUPPORTED_APEX_TRIGGER_OPERATON StatusCode = "UNSUPPORTED_APEX_TRIGGER_OPERATON"

	StatusCodeUNVERIFIED_SENDER_ADDRESS StatusCode = "UNVERIFIED_SENDER_ADDRESS"

	StatusCodeUSER_OWNS_PORTAL_ACCOUNT_EXCEPTION StatusCode = "USER_OWNS_PORTAL_ACCOUNT_EXCEPTION"

	StatusCodeUSER_WITH_APEX_SHARES_EXCEPTION StatusCode = "USER_WITH_APEX_SHARES_EXCEPTION"

	StatusCodeVF_COMPILE_ERROR StatusCode = "VF_COMPILE_ERROR"

	StatusCodeWEBLINK_SIZE_LIMIT_EXCEEDED StatusCode = "WEBLINK_SIZE_LIMIT_EXCEEDED"

	StatusCodeWEBLINK_URL_INVALID StatusCode = "WEBLINK_URL_INVALID"

	StatusCodeWRONG_CONTROLLER_TYPE StatusCode = "WRONG_CONTROLLER_TYPE"

	StatusCodeXCLEAN_UNEXPECTED_ERROR StatusCode = "XCLEAN_UNEXPECTED_ERROR"
)

type ExtendedErrorCode string

type ShareAccessLevel string

const (
	ShareAccessLevelRead ShareAccessLevel = "Read"

	ShareAccessLevelEdit ShareAccessLevel = "Edit"

	ShareAccessLevelAll ShareAccessLevel = "All"
)

type FieldType string

const (
	FieldTypeString FieldType = "string"

	FieldTypePicklist FieldType = "picklist"

	FieldTypeMultipicklist FieldType = "multipicklist"

	FieldTypeCombobox FieldType = "combobox"

	FieldTypeReference FieldType = "reference"

	FieldTypeBase64 FieldType = "base64"

	FieldTypeBoolean FieldType = "boolean"

	FieldTypeCurrency FieldType = "currency"

	FieldTypeTextarea FieldType = "textarea"

	FieldTypeInt FieldType = "int"

	FieldTypeDouble FieldType = "double"

	FieldTypePercent FieldType = "percent"

	FieldTypePhone FieldType = "phone"

	FieldTypeId FieldType = "id"

	FieldTypeDate FieldType = "date"

	FieldTypeDatetime FieldType = "datetime"

	FieldTypeTime FieldType = "time"

	FieldTypeUrl FieldType = "url"

	FieldTypeEmail FieldType = "email"

	FieldTypeEncryptedstring FieldType = "encryptedstring"

	FieldTypeDatacategorygroupreference FieldType = "datacategorygroupreference"

	FieldTypeLocation FieldType = "location"

	FieldTypeAddress FieldType = "address"

	FieldTypeAnyType FieldType = "anyType"

	FieldTypeComplexvalue FieldType = "complexvalue"
)

type SoapType string

const (
	SoapTypeTnsID SoapType = "tns:ID"

	SoapTypeXsdbase64Binary SoapType = "xsd:base64Binary"

	SoapTypeXsdboolean SoapType = "xsd:boolean"

	SoapTypeXsddouble SoapType = "xsd:double"

	SoapTypeXsdint SoapType = "xsd:int"

	SoapTypeXsdstring SoapType = "xsd:string"

	SoapTypeXsddate SoapType = "xsd:date"

	SoapTypeXsddateTime SoapType = "xsd:dateTime"

	SoapTypeXsdtime SoapType = "xsd:time"

	SoapTypeTnslocation SoapType = "tns:location"

	SoapTypeTnsaddress SoapType = "tns:address"

	SoapTypeXsdanyType SoapType = "xsd:anyType"

	SoapTypeUrnRelationshipReferenceTo SoapType = "urn:RelationshipReferenceTo"

	SoapTypeUrnJunctionIdListNames SoapType = "urn:JunctionIdListNames"

	SoapTypeUrnSearchLayoutFieldsDisplayed SoapType = "urn:SearchLayoutFieldsDisplayed"

	SoapTypeUrnSearchLayoutField SoapType = "urn:SearchLayoutField"

	SoapTypeUrnSearchLayoutButtonsDisplayed SoapType = "urn:SearchLayoutButtonsDisplayed"

	SoapTypeUrnSearchLayoutButton SoapType = "urn:SearchLayoutButton"

	SoapTypeUrnRecordTypesSupported SoapType = "urn:RecordTypesSupported"
)

type DifferenceType string

const (
	DifferenceTypeDIFFERENT DifferenceType = "DIFFERENT"

	DifferenceTypeNULL DifferenceType = "NULL"

	DifferenceTypeSAME DifferenceType = "SAME"

	DifferenceTypeSIMILAR DifferenceType = "SIMILAR"
)

type Article string

const (
	ArticleNone Article = "None"

	ArticleIndefinite Article = "Indefinite"

	ArticleDefinite Article = "Definite"
)

type CaseType string

const (
	CaseTypeNominative CaseType = "Nominative"

	CaseTypeAccusative CaseType = "Accusative"

	CaseTypeGenitive CaseType = "Genitive"

	CaseTypeDative CaseType = "Dative"

	CaseTypeInessive CaseType = "Inessive"

	CaseTypeElative CaseType = "Elative"

	CaseTypeIllative CaseType = "Illative"

	CaseTypeAdessive CaseType = "Adessive"

	CaseTypeAblative CaseType = "Ablative"

	CaseTypeAllative CaseType = "Allative"

	CaseTypeEssive CaseType = "Essive"

	CaseTypeTranslative CaseType = "Translative"

	CaseTypePartitive CaseType = "Partitive"

	CaseTypeObjective CaseType = "Objective"

	CaseTypeSubjective CaseType = "Subjective"

	CaseTypeInstrumental CaseType = "Instrumental"

	CaseTypePrepositional CaseType = "Prepositional"

	CaseTypeLocative CaseType = "Locative"

	CaseTypeVocative CaseType = "Vocative"

	CaseTypeSublative CaseType = "Sublative"

	CaseTypeSuperessive CaseType = "Superessive"

	CaseTypeDelative CaseType = "Delative"

	CaseTypeCausalfinal CaseType = "Causalfinal"

	CaseTypeEssiveformal CaseType = "Essiveformal"

	CaseTypeTermanative CaseType = "Termanative"

	CaseTypeDistributive CaseType = "Distributive"

	CaseTypeErgative CaseType = "Ergative"

	CaseTypeAdverbial CaseType = "Adverbial"

	CaseTypeAbessive CaseType = "Abessive"

	CaseTypeComitative CaseType = "Comitative"
)

type Gender string

const (
	GenderNeuter Gender = "Neuter"

	GenderMasculine Gender = "Masculine"

	GenderFeminine Gender = "Feminine"

	GenderAnimateMasculine Gender = "AnimateMasculine"
)

type GrammaticalNumber string

const (
	GrammaticalNumberSingular GrammaticalNumber = "Singular"

	GrammaticalNumberPlural GrammaticalNumber = "Plural"
)

type Possessive string

const (
	PossessiveNone Possessive = "None"

	PossessiveFirst Possessive = "First"

	PossessiveSecond Possessive = "Second"
)

type StartsWith string

const (
	StartsWithConsonant StartsWith = "Consonant"

	StartsWithVowel StartsWith = "Vowel"

	StartsWithSpecial StartsWith = "Special"
)

type ComponentInstancePropertyTypeEnum string

const (
	ComponentInstancePropertyTypeEnumDecorator ComponentInstancePropertyTypeEnum = "decorator"
)

type FlexipageContextTypeEnum string

const (
	FlexipageContextTypeEnumENTITYNAME FlexipageContextTypeEnum = "ENTITYNAME"
)

type FeedLayoutFilterType string

const (
	FeedLayoutFilterTypeAllUpdates FeedLayoutFilterType = "AllUpdates"

	FeedLayoutFilterTypeFeedItemType FeedLayoutFilterType = "FeedItemType"

	FeedLayoutFilterTypeCustom FeedLayoutFilterType = "Custom"
)

type TabOrderType string

const (
	TabOrderTypeLeftToRight TabOrderType = "LeftToRight"

	TabOrderTypeTopToBottom TabOrderType = "TopToBottom"
)

type WebLinkWindowType string

const (
	WebLinkWindowTypeNewWindow WebLinkWindowType = "newWindow"

	WebLinkWindowTypeSidebar WebLinkWindowType = "sidebar"

	WebLinkWindowTypeNoSidebar WebLinkWindowType = "noSidebar"

	WebLinkWindowTypeReplace WebLinkWindowType = "replace"

	WebLinkWindowTypeOnClickJavaScript WebLinkWindowType = "onClickJavaScript"
)

type WebLinkPosition string

const (
	WebLinkPositionFullScreen WebLinkPosition = "fullScreen"

	WebLinkPositionNone WebLinkPosition = "none"

	WebLinkPositionTopLeft WebLinkPosition = "topLeft"
)

type WebLinkType string

const (
	WebLinkTypeUrl WebLinkType = "url"

	WebLinkTypeSControl WebLinkType = "sControl"

	WebLinkTypeJavascript WebLinkType = "javascript"

	WebLinkTypePage WebLinkType = "page"

	WebLinkTypeFlow WebLinkType = "flow"
)

type LayoutComponentType string

const (
	LayoutComponentTypeReportChart LayoutComponentType = "ReportChart"

	LayoutComponentTypeField LayoutComponentType = "Field"

	LayoutComponentTypeSeparator LayoutComponentType = "Separator"

	LayoutComponentTypeSControl LayoutComponentType = "SControl"

	LayoutComponentTypeEmptySpace LayoutComponentType = "EmptySpace"

	LayoutComponentTypeVisualforcePage LayoutComponentType = "VisualforcePage"

	LayoutComponentTypeExpandedLookup LayoutComponentType = "ExpandedLookup"

	LayoutComponentTypeAuraComponent LayoutComponentType = "AuraComponent"

	LayoutComponentTypeCanvas LayoutComponentType = "Canvas"

	LayoutComponentTypeCustomLink LayoutComponentType = "CustomLink"

	LayoutComponentTypeAnalyticsCloud LayoutComponentType = "AnalyticsCloud"
)

type ReportChartSize string

const (
	ReportChartSizeSMALL ReportChartSize = "SMALL"

	ReportChartSizeMEDIUM ReportChartSize = "MEDIUM"

	ReportChartSizeLARGE ReportChartSize = "LARGE"
)

type EmailPriority string

const (
	EmailPriorityHighest EmailPriority = "Highest"

	EmailPriorityHigh EmailPriority = "High"

	EmailPriorityNormal EmailPriority = "Normal"

	EmailPriorityLow EmailPriority = "Low"

	EmailPriorityLowest EmailPriority = "Lowest"
)

type SendEmailOptOutPolicy string

const (
	SendEmailOptOutPolicySEND SendEmailOptOutPolicy = "SEND"

	SendEmailOptOutPolicyFILTER SendEmailOptOutPolicy = "FILTER"

	SendEmailOptOutPolicyREJECT SendEmailOptOutPolicy = "REJECT"
)

type OrderByDirection string

const (
	OrderByDirectionAscending OrderByDirection = "ascending"

	OrderByDirectionDescending OrderByDirection = "descending"
)

type OrderByNullsPosition string

const (
	OrderByNullsPositionFirst OrderByNullsPosition = "first"

	OrderByNullsPositionLast OrderByNullsPosition = "last"
)

type SoqlOperator string

const (
	SoqlOperatorEquals SoqlOperator = "equals"

	SoqlOperatorExcludes SoqlOperator = "excludes"

	SoqlOperatorGreaterThan SoqlOperator = "greaterThan"

	SoqlOperatorGreaterThanOrEqualTo SoqlOperator = "greaterThanOrEqualTo"

	SoqlOperatorIn SoqlOperator = "in"

	SoqlOperatorIncludes SoqlOperator = "includes"

	SoqlOperatorLessThan SoqlOperator = "lessThan"

	SoqlOperatorLessThanOrEqualTo SoqlOperator = "lessThanOrEqualTo"

	SoqlOperatorLike SoqlOperator = "like"

	SoqlOperatorNotEquals SoqlOperator = "notEquals"

	SoqlOperatorNotIn SoqlOperator = "notIn"

	SoqlOperatorWithin SoqlOperator = "within"
)

type SoqlConjunction string

const (
	SoqlConjunctionAnd SoqlConjunction = "and"

	SoqlConjunctionOr SoqlConjunction = "or"
)

type AppMenuType string

const (
	AppMenuTypeAppSwitcher AppMenuType = "AppSwitcher"

	AppMenuTypeSalesforce1 AppMenuType = "Salesforce1"

	AppMenuTypeNetworkTabs AppMenuType = "NetworkTabs"
)

type ListViewIsSoqlCompatible string

const (
	ListViewIsSoqlCompatibleTRUE ListViewIsSoqlCompatible = "TRUE"

	ListViewIsSoqlCompatibleFALSE ListViewIsSoqlCompatible = "FALSE"

	ListViewIsSoqlCompatibleALL ListViewIsSoqlCompatible = "ALL"
)

const (
	DebugLevelNone string = "None"

	DebugLevelDebugOnly string = "DebugOnly"

	DebugLevelDb string = "Db"

	DebugLevelProfiling string = "Profiling"

	DebugLevelCallout string = "Callout"

	DebugLevelDetail string = "Detail"
)

const (
	LogCategoryDb string = "Db"

	LogCategoryWorkflow string = "Workflow"

	LogCategoryValidation string = "Validation"

	LogCategoryCallout string = "Callout"

	LogCategoryApex_code string = "Apex_code"

	LogCategoryApex_profiling string = "Apex_profiling"

	LogCategoryVisualforce string = "Visualforce"

	LogCategorySystem string = "System"

	LogCategoryAll string = "All"
)

const (
	LogCategoryLevelNone string = "None"

	LogCategoryLevelFinest string = "Finest"

	LogCategoryLevelFiner string = "Finer"

	LogCategoryLevelFine string = "Fine"

	LogCategoryLevelDebug string = "Debug"

	LogCategoryLevelInfo string = "Info"

	LogCategoryLevelWarn string = "Warn"

	LogCategoryLevelError string = "Error"
)

type OwnerChangeOptionType string

const (
	OwnerChangeOptionTypeEnforceNewOwnerHasReadAccess OwnerChangeOptionType = "EnforceNewOwnerHasReadAccess"

	OwnerChangeOptionTypeTransferOpenActivities OwnerChangeOptionType = "TransferOpenActivities"

	OwnerChangeOptionTypeTransferNotesAndAttachments OwnerChangeOptionType = "TransferNotesAndAttachments"

	OwnerChangeOptionTypeTransferOthersOpenOpportunities OwnerChangeOptionType = "TransferOthersOpenOpportunities"

	OwnerChangeOptionTypeTransferOwnedOpenOpportunities OwnerChangeOptionType = "TransferOwnedOpenOpportunities"

	OwnerChangeOptionTypeTransferContracts OwnerChangeOptionType = "TransferContracts"

	OwnerChangeOptionTypeTransferOrders OwnerChangeOptionType = "TransferOrders"

	OwnerChangeOptionTypeTransferContacts OwnerChangeOptionType = "TransferContacts"
)

type FindDuplicates struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com findDuplicates"`

	SObjects []*SObject `xml:"sObjects,omitempty"`
}

type FindDuplicatesResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com findDuplicatesResponse"`

	Result []*FindDuplicatesResult `xml:"result,omitempty"`
}

type Login struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com login"`

	Username string `xml:"username,omitempty"`

	Password string `xml:"password,omitempty"`
}

type LoginResponse struct {
	Result *LoginResult `xml:"result,omitempty"`
}

type DescribeSObject struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeSObject"`

	SObjectType string `xml:"sObjectType,omitempty"`
}

type DescribeSObjectResponse struct {
	Result *DescribeSObjectResult `xml:"result,omitempty"`
}

type DescribeSObjects struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeSObjects"`

	SObjectType string `xml:"sObjectType,omitempty"`
}

type DescribeSObjectsResponse struct {
	Result *DescribeSObjectResult `xml:"result,omitempty"`
}

type DescribeGlobal struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeGlobal"`
}

type DescribeGlobalResponse struct {
	Result *DescribeGlobalResult `xml:"result,omitempty"`
}

type DescribeGlobalTheme struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeGlobalTheme"`
}

type DescribeGlobalThemeResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeGlobalThemeResponse"`

	Result *DescribeGlobalTheme `xml:"result,omitempty"`
}

type DescribeTheme struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeTheme"`

	SobjectType []string `xml:"sobjectType,omitempty"`
}

type DescribeThemeResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeThemeResponse"`

	Result *DescribeThemeResult `xml:"result,omitempty"`
}

type DescribeDataCategoryGroups struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeDataCategoryGroups"`

	SObjectType string `xml:"sObjectType,omitempty"`
}

type DescribeDataCategoryGroupsResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeDataCategoryGroupsResponse"`

	Result *DescribeDataCategoryGroupResult `xml:"result,omitempty"`
}

type DescribeDataCategoryGroupStructures struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeDataCategoryGroupStructures"`

	Pairs *DataCategoryGroupSobjectTypePair `xml:"pairs,omitempty"`

	TopCategoriesOnly bool `xml:"topCategoriesOnly,omitempty"`
}

type DescribeDataCategoryGroupStructuresResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeDataCategoryGroupStructuresResponse"`

	Result *DescribeDataCategoryGroupStructureResult `xml:"result,omitempty"`
}

type DescribeKnowledgeSettings struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeKnowledgeSettings"`
}

type DescribeKnowledgeSettingsResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeKnowledgeSettingsResponse"`

	Result *KnowledgeSettings `xml:"result,omitempty"`
}

type DescribeFlexiPages struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeFlexiPages"`

	FlexiPages []string `xml:"flexiPages,omitempty"`

	Contexts []*FlexipageContext `xml:"contexts,omitempty"`
}

type DescribeFlexiPagesResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeFlexiPagesResponse"`

	Result []*DescribeFlexiPageResult `xml:"result,omitempty"`
}

type DescribeAppMenu struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeAppMenu"`

	AppMenuType *AppMenuType `xml:"appMenuType,omitempty"`

	NetworkId string `xml:"networkId,omitempty"`
}

type DescribeAppMenuResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeAppMenuResponse"`

	Result *DescribeAppMenuResult `xml:"result,omitempty"`
}

type DescribeLayout struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeLayout"`

	SObjectType string `xml:"sObjectType,omitempty"`

	LayoutName string `xml:"layoutName,omitempty"`

	RecordTypeIds []string `xml:"recordTypeIds,omitempty"`
}

type DescribeLayoutResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeLayoutResponse"`

	Result *DescribeLayoutResultResult `xml:"result,omitempty"`
}

type DescribeCompactLayouts struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeCompactLayouts"`

	SObjectType string `xml:"sObjectType,omitempty"`

	RecordTypeIds []string `xml:"recordTypeIds,omitempty"`
}

type DescribeCompactLayoutsResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeCompactLayoutsResponse"`

	Result *DescribeCompactLayoutsResult `xml:"result,omitempty"`
}

type DescribePrimaryCompactLayouts struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describePrimaryCompactLayouts"`

	SObjectTypes []string `xml:"sObjectTypes,omitempty"`
}

type DescribePrimaryCompactLayoutsResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describePrimaryCompactLayoutsResponse"`

	Result []*DescribeCompactLayout `xml:"result,omitempty"`
}

type DescribePathAssistants struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describePathAssistants"`

	SObjectType string `xml:"sObjectType,omitempty"`

	PicklistValue string `xml:"picklistValue,omitempty"`

	RecordTypeIds []string `xml:"recordTypeIds,omitempty"`
}

type DescribePathAssistantsResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describePathAssistantsResponse"`

	Result *DescribePathAssistantsResult `xml:"result,omitempty"`
}

type DescribeApprovalLayoutParameter struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeApprovalLayout"`

	SObjectType string `xml:"sObjectType,omitempty"`

	ApprovalProcessNames []string `xml:"approvalProcessNames,omitempty"`
}

type DescribeApprovalLayoutResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeApprovalLayoutResponse"`

	Result *DescribeApprovalLayoutResult `xml:"result,omitempty"`
}

type DescribeSoftphoneLayout struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeSoftphoneLayout"`
}

type DescribeSoftphoneLayoutResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeSoftphoneLayoutResponse"`

	Result *DescribeSoftphoneLayoutResult `xml:"result,omitempty"`
}

type DescribeSoqlListViews struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeSoqlListViews"`

	Request *DescribeSoqlListViewsRequest `xml:"request,omitempty"`
}

type DescribeSoqlListViewsResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeSoqlListViewsResponse"`

	Result *DescribeSoqlListViewResult `xml:"result,omitempty"`
}

type ExecuteListView struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com executeListView"`

	Request *ExecuteListViewRequest `xml:"request,omitempty"`
}

type ExecuteListViewResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com executeListViewResponse"`

	Result *ExecuteListViewResult `xml:"result,omitempty"`
}

type DescribeSObjectListViews struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeSObjectListViews"`

	SObjectType string `xml:"sObjectType,omitempty"`

	RecentsOnly bool `xml:"recentsOnly,omitempty"`

	IsSoqlCompatible *ListViewIsSoqlCompatible `xml:"isSoqlCompatible,omitempty"`

	Limit int32 `xml:"limit,omitempty"`

	Offset int32 `xml:"offset,omitempty"`
}

type DescribeSObjectListViewsResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeSObjectListViewsResponse"`

	Result *DescribeSoqlListViewResult `xml:"result,omitempty"`
}

type DescribeSearchLayouts struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeSearchLayouts"`

	SObjectType []string `xml:"sObjectType,omitempty"`
}

type DescribeSearchLayoutsResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeSearchLayoutsResponse"`

	Result []*DescribeSearchLayoutResult `xml:"result,omitempty"`
}

type DescribeSearchScopeOrder struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeSearchScopeOrder"`
}

type DescribeSearchScopeOrderResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeSearchScopeOrderResponse"`

	Result []*DescribeSearchScopeOrderResult `xml:"result,omitempty"`
}

type DescribeSearchableEntities struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeSearchableEntities"`

	IncludeOnlyEntitiesWithTabs bool `xml:"includeOnlyEntitiesWithTabs,omitempty"`
}

type DescribeSearchableEntitiesResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeSearchableEntitiesResponse"`

	Result []*DescribeSearchableEntityResult `xml:"result,omitempty"`
}

type DescribeTabs struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeTabs"`
}

type DescribeTabsResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeTabsResponse"`

	Result []*DescribeTabSetResult `xml:"result,omitempty"`
}

type DescribeAllTabs struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeAllTabs"`
}

type DescribeAllTabsResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeAllTabsResponse"`

	Result []*DescribeTab `xml:"result,omitempty"`
}

type DescribeNouns struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeNouns"`

	Nouns string `xml:"nouns,omitempty"`

	OnlyRenamed bool `xml:"onlyRenamed,omitempty"`

	IncludeFields bool `xml:"includeFields,omitempty"`
}

type DescribeNounsResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeNounsResponse"`

	Result []*DescribeNounResult `xml:"result,omitempty"`
}

type Create struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com create"`

	SObjects []*SObject `xml:"sObjects,omitempty"`
}

type CreateResponse struct {
	Result []*SaveResult `xml:"result,omitempty"`
}

type SendEmail struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com sendEmail"`

	Messages *Email `xml:"messages,omitempty"`
}

type SendEmailResponse struct {
	Result *SendEmailResult `xml:"result,omitempty"`
}

type RenderEmailTemplate struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com renderEmailTemplate"`

	RenderRequests *RenderEmailTemplateRequest `xml:"renderRequests,omitempty"`
}

type RenderEmailTemplateResponse struct {
	Result *RenderEmailTemplateResult `xml:"result,omitempty"`
}

type SendEmailMessage struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com sendEmailMessage"`

	Ids string `xml:"ids,omitempty"`
}

type SendEmailMessageResponse struct {
	Result *SendEmailResult `xml:"result,omitempty"`
}

type Update struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com update"`

	SObjects []*SObject `xml:"sObjects,omitempty"`
}

type UpdateResponse struct {
	Result []*SaveResult `xml:"result,omitempty"`
}

type Upsert struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com upsert"`

	ExternalIDFieldName string `xml:"externalIDFieldName,omitempty"`

	SObjects []*SObject `xml:"sObjects,omitempty"`
}

type UpsertResponse struct {
	Result []*UpsertResult `xml:"result,omitempty"`
}

type Merge struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com merge"`

	Request []*MergeRequest `xml:"request,omitempty"`
}

type MergeResponse struct {
	Result []*MergeResult `xml:"result,omitempty"`
}

type Delete struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com delete"`

	Ids []string `xml:"ids,omitempty"`
}

type DeleteResponse struct {
	Result []*DeleteResult `xml:"result,omitempty"`
}

type Undelete struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com undelete"`

	Ids []string `xml:"ids,omitempty"`
}

type UndeleteResponse struct {
	Result []*UndeleteResult `xml:"result,omitempty"`
}

type EmptyRecycleBin struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com emptyRecycleBin"`

	Ids []string `xml:"ids,omitempty"`
}

type EmptyRecycleBinResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com emptyRecycleBinResponse"`

	Result []*EmptyRecycleBinResult `xml:"result,omitempty"`
}

type Process struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com process"`

	Actions []*ProcessRequest `xml:"actions,omitempty"`
}

type ProcessResponse struct {
	Result []*ProcessResult `xml:"result,omitempty"`
}

type PerformQuickActions struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com performQuickActions"`

	QuickActions []*PerformQuickActionRequest `xml:"quickActions,omitempty"`
}

type PerformQuickActionsResponse struct {
	Result []*PerformQuickActionResult `xml:"result,omitempty"`
}

type RetrieveQuickActionTemplates struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com retrieveQuickActionTemplates"`

	QuickActionNames []string `xml:"quickActionNames,omitempty"`

	ContextId string `xml:"contextId,omitempty"`
}

type RetrieveQuickActionTemplatesResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com retrieveQuickActionTemplatesResponse"`

	Result []*QuickActionTemplateResult `xml:"result,omitempty"`
}

type DescribeQuickActions struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeQuickActions"`

	QuickActions []string `xml:"quickActions,omitempty"`
}

type DescribeQuickActionsResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeQuickActionsResponse"`

	Result []*DescribeQuickActionResult `xml:"result,omitempty"`
}

type DescribeAvailableQuickActions struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeAvailableQuickActions"`

	ContextType string `xml:"contextType,omitempty"`
}

type DescribeAvailableQuickActionsResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeAvailableQuickActionsResponse"`

	Result []*DescribeAvailableQuickActionResult `xml:"result,omitempty"`
}

type DescribeVisualForce struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeVisualForce"`

	IncludeAllDetails bool `xml:"includeAllDetails,omitempty"`

	NamespacePrefix string `xml:"namespacePrefix,omitempty"`
}

type DescribeVisualForceResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com describeVisualForceResponse"`

	Result *DescribeVisualForceResult `xml:"result,omitempty"`
}

type Retrieve struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com retrieve"`

	FieldList string `xml:"fieldList,omitempty"`

	SObjectType string `xml:"sObjectType,omitempty"`

	Ids []string `xml:"ids,omitempty"`
}

type RetrieveResponse struct {
	Result []*SObject `xml:"result,omitempty"`
}

type ConvertLead struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com convertLead"`

	LeadConverts []*LeadConvert `xml:"leadConverts,omitempty"`
}

type ConvertLeadResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com convertLeadResponse"`

	Result []*LeadConvertResult `xml:"result,omitempty"`
}

type GetUpdated struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com getUpdated"`

	SObjectType string `xml:"sObjectType,omitempty"`

	StartDate time.Time `xml:"startDate,omitempty"`

	EndDate time.Time `xml:"endDate,omitempty"`
}

type GetUpdatedResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com getUpdatedResponse"`

	Result *GetUpdatedResult `xml:"result,omitempty"`
}

type GetDeleted struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com getDeleted"`

	SObjectType string `xml:"sObjectType,omitempty"`

	StartDate time.Time `xml:"startDate,omitempty"`

	EndDate time.Time `xml:"endDate,omitempty"`
}

type GetDeletedResponse struct {
	Result *GetDeletedResult `xml:"result,omitempty"`
}

type Logout struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com logout"`
}

type LogoutResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com logoutResponse"`
}

type InvalidateSessions struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com invalidateSessions"`

	SessionIds []string `xml:"sessionIds,omitempty"`
}

type InvalidateSessionsResponse struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com invalidateSessionsResponse"`

	Result []*InvalidateSessionsResult `xml:"result,omitempty"`
}

type Query struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com query"`

	QueryString string `xml:"queryString,omitempty"`
}

type QueryResponse struct {
	Result *QueryResult `xml:"result,omitempty"`
}

type QueryAll struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com queryAll"`

	QueryString string `xml:"queryString,omitempty"`
}

type QueryAllResponse struct {
	Result *QueryResult `xml:"result,omitempty"`
}

type QueryMore struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com queryMore"`

	QueryLocator string `xml:"queryLocator,omitempty"`
}

type QueryMoreResponse struct {
	Result *QueryResult `xml:"result,omitempty"`
}

type Search struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com search"`

	SearchString string `xml:"searchString,omitempty"`
}

type SearchResponse struct {
	Result *SearchResult `xml:"result,omitempty"`
}

type GetServerTimestamp struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com getServerTimestamp"`
}

type GetServerTimestampResponse struct {
	Result *GetServerTimestampResult `xml:"result,omitempty"`
}

type SetPassword struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com setPassword"`

	UserId string `xml:"userId,omitempty"`

	Password string `xml:"password,omitempty"`
}

type SetPasswordResponse struct {
	Result *SetPasswordResult `xml:"result,omitempty"`
}

type ResetPassword struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com resetPassword"`

	UserId string `xml:"userId,omitempty"`
}

type ResetPasswordResponse struct {
	Result *ResetPasswordResult `xml:"result,omitempty"`
}

type GetUserInfo struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com getUserInfo"`
}

type GetUserInfoResponse struct {
	Result *GetUserInfoResult `xml:"result,omitempty"`
}

type SessionHeader struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com SessionHeader"`

	SessionId string `xml:"sessionId,omitempty"`
}

type LoginScopeHeader struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com LoginScopeHeader"`

	OrganizationId string `xml:"organizationId,omitempty"`

	PortalId string `xml:"portalId,omitempty"`
}

type CallOptions struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com CallOptions"`

	Client string `xml:"client,omitempty"`

	DefaultNamespace string `xml:"defaultNamespace,omitempty"`
}

type QueryOptions struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com QueryOptions"`

	BatchSize int32 `xml:"batchSize,omitempty"`
}

type DebuggingHeader struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DebuggingHeader"`

	Categories []*LogInfo `xml:"categories,omitempty"`

	DebugLevel string `xml:"debugLevel,omitempty"`
}

type DebuggingInfo struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DebuggingInfo"`

	DebugLog string `xml:"debugLog,omitempty"`
}

type PackageVersionHeader struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com PackageVersionHeader"`

	PackageVersions []*PackageVersion `xml:"packageVersions,omitempty"`
}

type AllowFieldTruncationHeader struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com AllowFieldTruncationHeader"`

	AllowFieldTruncation bool `xml:"allowFieldTruncation,omitempty"`
}

type DisableFeedTrackingHeader struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DisableFeedTrackingHeader"`

	DisableFeedTracking bool `xml:"disableFeedTracking,omitempty"`
}

type StreamingEnabledHeader struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com StreamingEnabledHeader"`

	StreamingEnabled bool `xml:"streamingEnabled,omitempty"`
}

type AllOrNoneHeader struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com AllOrNoneHeader"`

	AllOrNone bool `xml:"allOrNone,omitempty"`
}

type DuplicateRuleHeader struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DuplicateRuleHeader"`

	AllowSave bool `xml:"allowSave,omitempty"`

	IncludeRecordDetails bool `xml:"includeRecordDetails,omitempty"`

	RunAsCurrentUser bool `xml:"runAsCurrentUser,omitempty"`
}

type LimitInfoHeader struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com LimitInfoHeader"`

	LimitInfo *LimitInfo `xml:"limitInfo,omitempty"`
}

type MruHeader struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com MruHeader"`

	UpdateMru bool `xml:"updateMru,omitempty"`
}

type EmailHeader struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com EmailHeader"`

	TriggerAutoResponseEmail bool `xml:"triggerAutoResponseEmail,omitempty"`

	TriggerOtherEmail bool `xml:"triggerOtherEmail,omitempty"`

	TriggerUserEmail bool `xml:"triggerUserEmail,omitempty"`
}

type AssignmentRuleHeader struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com AssignmentRuleHeader"`

	AssignmentRuleId string `xml:"assignmentRuleId,omitempty"`

	UseDefaultRule bool `xml:"useDefaultRule,omitempty"`
}

type UserTerritoryDeleteHeader struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com UserTerritoryDeleteHeader"`

	TransferToUserId string `xml:"transferToUserId,omitempty"`
}

type LocaleOptions struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com LocaleOptions"`

	Language string `xml:"language,omitempty"`

	LocalizeErrors bool `xml:"localizeErrors,omitempty"`
}

type OwnerChangeOptions struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com OwnerChangeOptions"`

	Options []*OwnerChangeOption `xml:"options,omitempty"`
}

type Address struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com address"`

	*Location

	City string `xml:"city,omitempty"`

	Country string `xml:"country,omitempty"`

	CountryCode string `xml:"countryCode,omitempty"`

	GeocodeAccuracy string `xml:"geocodeAccuracy,omitempty"`

	PostalCode string `xml:"postalCode,omitempty"`

	State string `xml:"state,omitempty"`

	StateCode string `xml:"stateCode,omitempty"`

	Street string `xml:"street,omitempty"`
}

type Location struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com location"`

	Latitude float64 `xml:"latitude,omitempty"`

	Longitude float64 `xml:"longitude,omitempty"`
}

type QueryResult struct {
	Done bool `xml:"done,omitempty"`

	QueryLocator string `xml:"queryLocator,omitempty"`

	Records []*SObject `xml:"records,omitempty"`

	Size int32 `xml:"size,omitempty"`
}

type SearchResult struct {
	QueryId string `xml:"queryId,omitempty"`

	SearchRecords []*SearchRecord `xml:"searchRecords,omitempty"`

	SearchResultsMetadata *SearchResultsMetadata `xml:"searchResultsMetadata,omitempty"`
}

type SearchRecord struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com SearchRecord"`

	Record *SObject `xml:"record,omitempty"`

	Snippet *SearchSnippet `xml:"snippet,omitempty"`
}

type SearchSnippet struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com SearchSnippet"`

	Text string `xml:"text,omitempty"`

	WholeFields []*NameValuePair `xml:"wholeFields,omitempty"`
}

type SearchResultsMetadata struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com SearchResultsMetadata"`

	EntityLabelMetadata []*LabelsSearchMetadata `xml:"entityLabelMetadata,omitempty"`
}

type LabelsSearchMetadata struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com LabelsSearchMetadata"`

	EntityFieldLabels []*NameValuePair `xml:"entityFieldLabels,omitempty"`

	EntityName string `xml:"entityName,omitempty"`
}

type RelationshipReferenceTo struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com RelationshipReferenceTo"`

	ReferenceTo []string `xml:"referenceTo,omitempty"`
}

type RecordTypesSupported struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com RecordTypesSupported"`

	RecordTypeInfos []*RecordTypeInfo `xml:"recordTypeInfos,omitempty"`
}

type JunctionIdListNames struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com JunctionIdListNames"`

	Names []string `xml:"names,omitempty"`
}

type SearchLayoutButtonsDisplayed struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com SearchLayoutButtonsDisplayed"`

	Applicable bool `xml:"applicable,omitempty"`

	Buttons []*SearchLayoutButton `xml:"buttons,omitempty"`
}

type SearchLayoutButton struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com SearchLayoutButton"`

	ApiName string `xml:"apiName,omitempty"`

	Label string `xml:"label,omitempty"`
}

type SearchLayoutFieldsDisplayed struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com SearchLayoutFieldsDisplayed"`

	Applicable bool `xml:"applicable,omitempty"`

	Fields []*SearchLayoutField `xml:"fields,omitempty"`
}

type SearchLayoutField struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com SearchLayoutField"`

	ApiName string `xml:"apiName,omitempty"`

	Label string `xml:"label,omitempty"`

	Sortable bool `xml:"sortable,omitempty"`
}

type NameValuePair struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com NameValuePair"`

	Name string `xml:"name,omitempty"`

	Value string `xml:"value,omitempty"`
}

type NameObjectValuePair struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com NameObjectValuePair"`

	Name string `xml:"name,omitempty"`

	Value []interface{} `xml:"value,omitempty"`
}

type GetUpdatedResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com GetUpdatedResult"`

	Ids []string `xml:"ids,omitempty"`

	LatestDateCovered time.Time `xml:"latestDateCovered,omitempty"`
}

type GetDeletedResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com GetDeletedResult"`

	DeletedRecords []*DeletedRecord `xml:"deletedRecords,omitempty"`

	EarliestDateAvailable time.Time `xml:"earliestDateAvailable,omitempty"`

	LatestDateCovered time.Time `xml:"latestDateCovered,omitempty"`
}

type DeletedRecord struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DeletedRecord"`

	DeletedDate time.Time `xml:"deletedDate,omitempty"`

	Id string `xml:"id,omitempty"`
}

type GetServerTimestampResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com GetServerTimestampResult"`

	Timestamp time.Time `xml:"timestamp,omitempty"`
}

type InvalidateSessionsResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com InvalidateSessionsResult"`

	Errors []*Error `xml:"errors,omitempty"`

	Success bool `xml:"success,omitempty"`
}

type SetPasswordResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com SetPasswordResult"`
}

type ResetPasswordResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com ResetPasswordResult"`

	Password string `xml:"password,omitempty"`
}

type GetUserInfoResult struct {
	AccessibilityMode bool `xml:"accessibilityMode,omitempty"`

	CurrencySymbol string `xml:"currencySymbol,omitempty"`

	OrgAttachmentFileSizeLimit int32 `xml:"orgAttachmentFileSizeLimit,omitempty"`

	OrgDefaultCurrencyIsoCode string `xml:"orgDefaultCurrencyIsoCode,omitempty"`

	OrgDefaultCurrencyLocale string `xml:"orgDefaultCurrencyLocale,omitempty"`

	OrgDisallowHtmlAttachments bool `xml:"orgDisallowHtmlAttachments,omitempty"`

	OrgHasPersonAccounts bool `xml:"orgHasPersonAccounts,omitempty"`

	OrganizationId string `xml:"organizationId,omitempty"`

	OrganizationMultiCurrency bool `xml:"organizationMultiCurrency,omitempty"`

	OrganizationName string `xml:"organizationName,omitempty"`

	ProfileId string `xml:"profileId,omitempty"`

	RoleId string `xml:"roleId,omitempty"`

	SessionSecondsValid int32 `xml:"sessionSecondsValid,omitempty"`

	UserDefaultCurrencyIsoCode string `xml:"userDefaultCurrencyIsoCode,omitempty"`

	UserEmail string `xml:"userEmail,omitempty"`

	UserFullName string `xml:"userFullName,omitempty"`

	UserId string `xml:"userId,omitempty"`

	UserLanguage string `xml:"userLanguage,omitempty"`

	UserLocale string `xml:"userLocale,omitempty"`

	UserName string `xml:"userName,omitempty"`

	UserTimeZone string `xml:"userTimeZone,omitempty"`

	UserType string `xml:"userType,omitempty"`

	UserUiSkin string `xml:"userUiSkin,omitempty"`
}

type LoginResult struct {
	MetadataServerUrl string `xml:"metadataServerUrl,omitempty"`

	PasswordExpired bool `xml:"passwordExpired,omitempty"`

	Sandbox bool `xml:"sandbox,omitempty"`

	ServerUrl string `xml:"serverUrl,omitempty"`

	SessionId string `xml:"sessionId,omitempty"`

	UserId string `xml:"userId,omitempty"`

	UserInfo *GetUserInfoResult `xml:"userInfo,omitempty"`
}

type ExtendedErrorDetails struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com ExtendedErrorDetails"`

	ExtendedErrorCode *ExtendedErrorCode `xml:"extendedErrorCode,omitempty"`
}

type Error struct {
	ExtendedErrorDetails []*ExtendedErrorDetails `xml:"extendedErrorDetails,omitempty"`

	Fields []string `xml:"fields,omitempty"`

	Message string `xml:"message,omitempty"`

	StatusCode *StatusCode `xml:"statusCode,omitempty"`
}

type SendEmailError struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com SendEmailError"`

	Fields []string `xml:"fields,omitempty"`

	Message string `xml:"message,omitempty"`

	StatusCode *StatusCode `xml:"statusCode,omitempty"`

	TargetObjectId string `xml:"targetObjectId,omitempty"`
}

type SaveResult struct {
	Errors []*Error `xml:"errors,omitempty"`

	Id string `xml:"id,omitempty"`

	Success bool `xml:"success,omitempty"`
}

type RenderEmailTemplateError struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com RenderEmailTemplateError"`

	FieldName string `xml:"fieldName,omitempty"`

	Message string `xml:"message,omitempty"`

	Offset int32 `xml:"offset,omitempty"`

	StatusCode *StatusCode `xml:"statusCode,omitempty"`
}

type UpsertResult struct {
	Created bool `xml:"created,omitempty"`

	Errors []*Error `xml:"errors,omitempty"`

	Id string `xml:"id,omitempty"`

	Success bool `xml:"success,omitempty"`
}

type PerformQuickActionResult struct {
	ContextId string `xml:"contextId,omitempty"`

	Created bool `xml:"created,omitempty"`

	Errors []*Error `xml:"errors,omitempty"`

	FeedItemIds []string `xml:"feedItemIds,omitempty"`

	Ids []string `xml:"ids,omitempty"`

	Success bool `xml:"success,omitempty"`

	SuccessMessage string `xml:"successMessage,omitempty"`
}

type QuickActionTemplateResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com QuickActionTemplateResult"`

	DefaultValueFormulas *SObject `xml:"defaultValueFormulas,omitempty"`

	DefaultValues *SObject `xml:"defaultValues,omitempty"`

	Errors []*Error `xml:"errors,omitempty"`

	Success bool `xml:"success,omitempty"`
}

type MergeRequest struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com MergeRequest"`

	MasterRecord *SObject `xml:"masterRecord,omitempty"`

	RecordToMergeIds []string `xml:"recordToMergeIds,omitempty"`
}

type MergeResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com MergeResult"`

	Errors []*Error `xml:"errors,omitempty"`

	Id string `xml:"id,omitempty"`

	MergedRecordIds []string `xml:"mergedRecordIds,omitempty"`

	Success bool `xml:"success,omitempty"`

	UpdatedRelatedIds []string `xml:"updatedRelatedIds,omitempty"`
}

type ProcessRequest struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com ProcessRequest"`

	Comments string `xml:"comments,omitempty"`

	NextApproverIds []string `xml:"nextApproverIds,omitempty"`
}

type ProcessSubmitRequest struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com ProcessSubmitRequest"`

	*ProcessRequest

	ObjectId string `xml:"objectId,omitempty"`

	SubmitterId string `xml:"submitterId,omitempty"`

	ProcessDefinitionNameOrId string `xml:"processDefinitionNameOrId,omitempty"`

	SkipEntryCriteria bool `xml:"skipEntryCriteria,omitempty"`
}

type ProcessWorkitemRequest struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com ProcessWorkitemRequest"`

	*ProcessRequest

	Action string `xml:"action,omitempty"`

	WorkitemId string `xml:"workitemId,omitempty"`
}

type PerformQuickActionRequest struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com PerformQuickActionRequest"`

	ContextId string `xml:"contextId,omitempty"`

	QuickActionName string `xml:"quickActionName,omitempty"`

	Records []*SObject `xml:"records,omitempty"`
}

type DescribeAvailableQuickActionResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeAvailableQuickActionResult"`

	ActionEnumOrId string `xml:"actionEnumOrId,omitempty"`

	Label string `xml:"label,omitempty"`

	Name string `xml:"name,omitempty"`

	Type_ string `xml:"type,omitempty"`
}

type DescribeQuickActionResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeQuickActionResult"`

	AccessLevelRequired *ShareAccessLevel `xml:"accessLevelRequired,omitempty"`

	ActionEnumOrId string `xml:"actionEnumOrId,omitempty"`

	CanvasApplicationId string `xml:"canvasApplicationId,omitempty"`

	CanvasApplicationName string `xml:"canvasApplicationName,omitempty"`

	Colors []*DescribeColor `xml:"colors,omitempty"`

	ContextSobjectType string `xml:"contextSobjectType,omitempty"`

	DefaultValues []*DescribeQuickActionDefaultValue `xml:"defaultValues,omitempty"`

	Height int32 `xml:"height,omitempty"`

	IconName string `xml:"iconName,omitempty"`

	IconUrl string `xml:"iconUrl,omitempty"`

	Icons []*DescribeIcon `xml:"icons,omitempty"`

	Label string `xml:"label,omitempty"`

	Layout *DescribeLayoutSection `xml:"layout,omitempty"`

	LightningComponentBundleId string `xml:"lightningComponentBundleId,omitempty"`

	LightningComponentBundleName string `xml:"lightningComponentBundleName,omitempty"`

	LightningComponentQualifiedName string `xml:"lightningComponentQualifiedName,omitempty"`

	MiniIconUrl string `xml:"miniIconUrl,omitempty"`

	Name string `xml:"name,omitempty"`

	ShowQuickActionLcHeader bool `xml:"showQuickActionLcHeader,omitempty"`

	ShowQuickActionVfHeader bool `xml:"showQuickActionVfHeader,omitempty"`

	TargetParentField string `xml:"targetParentField,omitempty"`

	TargetRecordTypeId string `xml:"targetRecordTypeId,omitempty"`

	TargetSobjectType string `xml:"targetSobjectType,omitempty"`

	Type_ string `xml:"type,omitempty"`

	VisualforcePageName string `xml:"visualforcePageName,omitempty"`

	VisualforcePageUrl string `xml:"visualforcePageUrl,omitempty"`

	Width int32 `xml:"width,omitempty"`
}

type DescribeQuickActionDefaultValue struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeQuickActionDefaultValue"`

	DefaultValue string `xml:"defaultValue,omitempty"`

	Field string `xml:"field,omitempty"`
}

type DescribeVisualForceResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeVisualForceResult"`

	Domain string `xml:"domain,omitempty"`
}

type ProcessResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com ProcessResult"`

	ActorIds []string `xml:"actorIds,omitempty"`

	EntityId string `xml:"entityId,omitempty"`

	Errors []*Error `xml:"errors,omitempty"`

	InstanceId string `xml:"instanceId,omitempty"`

	InstanceStatus string `xml:"instanceStatus,omitempty"`

	NewWorkitemIds []string `xml:"newWorkitemIds,omitempty"`

	Success bool `xml:"success,omitempty"`
}

type DeleteResult struct {
	Errors []*Error `xml:"errors,omitempty"`

	Id string `xml:"id,omitempty"`

	Success bool `xml:"success,omitempty"`
}

type UndeleteResult struct {
	Errors []*Error `xml:"errors,omitempty"`

	Id string `xml:"id,omitempty"`

	Success bool `xml:"success,omitempty"`
}

type EmptyRecycleBinResult struct {
	Errors []*Error `xml:"errors,omitempty"`

	Id string `xml:"id,omitempty"`

	Success bool `xml:"success,omitempty"`
}

type LeadConvert struct {
	AccountId string `xml:"accountId,omitempty"`

	ContactId string `xml:"contactId,omitempty"`

	ConvertedStatus string `xml:"convertedStatus,omitempty"`

	DoNotCreateOpportunity bool `xml:"doNotCreateOpportunity,omitempty"`

	LeadId string `xml:"leadId,omitempty"`

	OpportunityName string `xml:"opportunityName,omitempty"`

	OverwriteLeadSource bool `xml:"overwriteLeadSource,omitempty"`

	OwnerId string `xml:"ownerId,omitempty"`

	SendNotificationEmail bool `xml:"sendNotificationEmail,omitempty"`
}

type LeadConvertResult struct {
	AccountId string `xml:"accountId,omitempty"`

	ContactId string `xml:"contactId,omitempty"`

	Errors []*Error `xml:"errors,omitempty"`

	LeadId string `xml:"leadId,omitempty"`

	OpportunityId string `xml:"opportunityId,omitempty"`

	Success bool `xml:"success,omitempty"`
}

type DescribeSObjectResult struct {
	ActionOverrides []*ActionOverride `xml:"actionOverrides,omitempty"`

	Activateable bool `xml:"activateable,omitempty"`

	ChildRelationships []*ChildRelationship `xml:"childRelationships,omitempty"`

	CompactLayoutable bool `xml:"compactLayoutable,omitempty"`

	Createable bool `xml:"createable,omitempty"`

	Custom bool `xml:"custom,omitempty"`

	CustomSetting bool `xml:"customSetting,omitempty"`

	Deletable bool `xml:"deletable,omitempty"`

	DeprecatedAndHidden bool `xml:"deprecatedAndHidden,omitempty"`

	FeedEnabled bool `xml:"feedEnabled,omitempty"`

	Fields []*Field `xml:"fields,omitempty"`

	HasSubtypes bool `xml:"hasSubtypes,omitempty"`

	IdEnabled bool `xml:"idEnabled,omitempty"`

	KeyPrefix string `xml:"keyPrefix,omitempty"`

	Label string `xml:"label,omitempty"`

	LabelPlural string `xml:"labelPlural,omitempty"`

	Layoutable bool `xml:"layoutable,omitempty"`

	Mergeable bool `xml:"mergeable,omitempty"`

	MruEnabled bool `xml:"mruEnabled,omitempty"`

	Name string `xml:"name,omitempty"`

	NamedLayoutInfos []*NamedLayoutInfo `xml:"namedLayoutInfos,omitempty"`

	NetworkScopeFieldName string `xml:"networkScopeFieldName,omitempty"`

	Queryable bool `xml:"queryable,omitempty"`

	RecordTypeInfos []*RecordTypeInfo `xml:"recordTypeInfos,omitempty"`

	Replicateable bool `xml:"replicateable,omitempty"`

	Retrieveable bool `xml:"retrieveable,omitempty"`

	SearchLayoutable bool `xml:"searchLayoutable,omitempty"`

	Searchable bool `xml:"searchable,omitempty"`

	SupportedScopes []*ScopeInfo `xml:"supportedScopes,omitempty"`

	Triggerable bool `xml:"triggerable,omitempty"`

	Undeletable bool `xml:"undeletable,omitempty"`

	Updateable bool `xml:"updateable,omitempty"`

	UrlDetail string `xml:"urlDetail,omitempty"`

	UrlEdit string `xml:"urlEdit,omitempty"`

	UrlNew string `xml:"urlNew,omitempty"`
}

type DescribeGlobalSObjectResult struct {
	Activateable bool `xml:"activateable,omitempty"`

	Createable bool `xml:"createable,omitempty"`

	Custom bool `xml:"custom,omitempty"`

	CustomSetting bool `xml:"customSetting,omitempty"`

	Deletable bool `xml:"deletable,omitempty"`

	DeprecatedAndHidden bool `xml:"deprecatedAndHidden,omitempty"`

	FeedEnabled bool `xml:"feedEnabled,omitempty"`

	HasSubtypes bool `xml:"hasSubtypes,omitempty"`

	IdEnabled bool `xml:"idEnabled,omitempty"`

	KeyPrefix string `xml:"keyPrefix,omitempty"`

	Label string `xml:"label,omitempty"`

	LabelPlural string `xml:"labelPlural,omitempty"`

	Layoutable bool `xml:"layoutable,omitempty"`

	Mergeable bool `xml:"mergeable,omitempty"`

	MruEnabled bool `xml:"mruEnabled,omitempty"`

	Name string `xml:"name,omitempty"`

	Queryable bool `xml:"queryable,omitempty"`

	Replicateable bool `xml:"replicateable,omitempty"`

	Retrieveable bool `xml:"retrieveable,omitempty"`

	Searchable bool `xml:"searchable,omitempty"`

	Triggerable bool `xml:"triggerable,omitempty"`

	Undeletable bool `xml:"undeletable,omitempty"`

	Updateable bool `xml:"updateable,omitempty"`
}

type ChildRelationship struct {
	CascadeDelete bool `xml:"cascadeDelete,omitempty"`

	ChildSObject string `xml:"childSObject,omitempty"`

	DeprecatedAndHidden bool `xml:"deprecatedAndHidden,omitempty"`

	Field string `xml:"field,omitempty"`

	JunctionIdListNames []string `xml:"junctionIdListNames,omitempty"`

	JunctionReferenceTo []string `xml:"junctionReferenceTo,omitempty"`

	RelationshipName string `xml:"relationshipName,omitempty"`

	RestrictedDelete bool `xml:"restrictedDelete,omitempty"`
}

type DescribeGlobalResult struct {
	Encoding string `xml:"encoding,omitempty"`

	MaxBatchSize int32 `xml:"maxBatchSize,omitempty"`

	Sobjects []*DescribeGlobalSObjectResult `xml:"sobjects,omitempty"`
}

type DescribeGlobalThemeResult struct {
	Global *DescribeGlobalResult `xml:"global,omitempty"`

	Theme *DescribeThemeResult `xml:"theme,omitempty"`
}

type ScopeInfo struct {
	Label string `xml:"label,omitempty"`

	Name string `xml:"name,omitempty"`
}

type FilteredLookupInfo struct {
	ControllingFields []string `xml:"controllingFields,omitempty"`

	Dependent bool `xml:"dependent,omitempty"`

	OptionalFilter bool `xml:"optionalFilter,omitempty"`
}

type Field struct {
	Aggregatable bool `xml:"aggregatable,omitempty"`

	AutoNumber bool `xml:"autoNumber,omitempty"`

	ByteLength int32 `xml:"byteLength,omitempty"`

	Calculated bool `xml:"calculated,omitempty"`

	CalculatedFormula string `xml:"calculatedFormula,omitempty"`

	CascadeDelete bool `xml:"cascadeDelete,omitempty"`

	CaseSensitive bool `xml:"caseSensitive,omitempty"`

	CompoundFieldName string `xml:"compoundFieldName,omitempty"`

	ControllerName string `xml:"controllerName,omitempty"`

	Createable bool `xml:"createable,omitempty"`

	Custom bool `xml:"custom,omitempty"`

	DefaultValue interface{} `xml:"defaultValue,omitempty"`

	DefaultValueFormula string `xml:"defaultValueFormula,omitempty"`

	DefaultedOnCreate bool `xml:"defaultedOnCreate,omitempty"`

	DependentPicklist bool `xml:"dependentPicklist,omitempty"`

	DeprecatedAndHidden bool `xml:"deprecatedAndHidden,omitempty"`

	Digits int32 `xml:"digits,omitempty"`

	DisplayLocationInDecimal bool `xml:"displayLocationInDecimal,omitempty"`

	Encrypted bool `xml:"encrypted,omitempty"`

	ExternalId bool `xml:"externalId,omitempty"`

	ExtraTypeInfo string `xml:"extraTypeInfo,omitempty"`

	Filterable bool `xml:"filterable,omitempty"`

	FilteredLookupInfo *FilteredLookupInfo `xml:"filteredLookupInfo,omitempty"`

	Groupable bool `xml:"groupable,omitempty"`

	HighScaleNumber bool `xml:"highScaleNumber,omitempty"`

	HtmlFormatted bool `xml:"htmlFormatted,omitempty"`

	IdLookup bool `xml:"idLookup,omitempty"`

	InlineHelpText string `xml:"inlineHelpText,omitempty"`

	Label string `xml:"label,omitempty"`

	Length int32 `xml:"length,omitempty"`

	Mask string `xml:"mask,omitempty"`

	MaskType string `xml:"maskType,omitempty"`

	Name string `xml:"name,omitempty"`

	NameField bool `xml:"nameField,omitempty"`

	NamePointing bool `xml:"namePointing,omitempty"`

	Nillable bool `xml:"nillable,omitempty"`

	Permissionable bool `xml:"permissionable,omitempty"`

	PicklistValues []*PicklistEntry `xml:"picklistValues,omitempty"`

	Precision int32 `xml:"precision,omitempty"`

	QueryByDistance bool `xml:"queryByDistance,omitempty"`

	ReferenceTargetField string `xml:"referenceTargetField,omitempty"`

	ReferenceTo []string `xml:"referenceTo,omitempty"`

	RelationshipName string `xml:"relationshipName,omitempty"`

	RelationshipOrder int32 `xml:"relationshipOrder,omitempty"`

	RestrictedDelete bool `xml:"restrictedDelete,omitempty"`

	RestrictedPicklist bool `xml:"restrictedPicklist,omitempty"`

	Scale int32 `xml:"scale,omitempty"`

	SoapType *SoapType `xml:"soapType,omitempty"`

	Sortable bool `xml:"sortable,omitempty"`

	Type_ *FieldType `xml:"type,omitempty"`

	Unique bool `xml:"unique,omitempty"`

	Updateable bool `xml:"updateable,omitempty"`

	WriteRequiresMasterRead bool `xml:"writeRequiresMasterRead,omitempty"`
}

type PicklistEntry struct {
	Active bool `xml:"active,omitempty"`

	DefaultValue bool `xml:"defaultValue,omitempty"`

	Label string `xml:"label,omitempty"`

	ValidFor []byte `xml:"validFor,omitempty"`

	Value string `xml:"value,omitempty"`
}

type DescribeDataCategoryGroupResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeDataCategoryGroupResult"`

	CategoryCount int32 `xml:"categoryCount,omitempty"`

	Description string `xml:"description,omitempty"`

	Label string `xml:"label,omitempty"`

	Name string `xml:"name,omitempty"`

	Sobject string `xml:"sobject,omitempty"`
}

type DescribeDataCategoryGroupStructureResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeDataCategoryGroupStructureResult"`

	Description string `xml:"description,omitempty"`

	Label string `xml:"label,omitempty"`

	Name string `xml:"name,omitempty"`

	Sobject string `xml:"sobject,omitempty"`

	TopCategories []*DataCategory `xml:"topCategories,omitempty"`
}

type DataCategoryGroupSobjectTypePair struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DataCategoryGroupSobjectTypePair"`

	DataCategoryGroupName string `xml:"dataCategoryGroupName,omitempty"`

	Sobject string `xml:"sobject,omitempty"`
}

type DataCategory struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DataCategory"`

	ChildCategories []*DataCategory `xml:"childCategories,omitempty"`

	Label string `xml:"label,omitempty"`

	Name string `xml:"name,omitempty"`
}

type KnowledgeSettings struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com KnowledgeSettings"`

	DefaultLanguage string `xml:"defaultLanguage,omitempty"`

	KnowledgeEnabled bool `xml:"knowledgeEnabled,omitempty"`

	Languages []*KnowledgeLanguageItem `xml:"languages,omitempty"`
}

type KnowledgeLanguageItem struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com KnowledgeLanguageItem"`

	Active bool `xml:"active,omitempty"`

	Name string `xml:"name,omitempty"`
}

type FieldDiff struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com FieldDiff"`

	Difference *DifferenceType `xml:"difference,omitempty"`

	Name string `xml:"name,omitempty"`
}

type AdditionalInformationMap struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com AdditionalInformationMap"`

	Name string `xml:"name,omitempty"`

	Value string `xml:"value,omitempty"`
}

type MatchRecord struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com MatchRecord"`

	AdditionalInformation []*AdditionalInformationMap `xml:"additionalInformation,omitempty"`

	FieldDiffs []*FieldDiff `xml:"fieldDiffs,omitempty"`

	MatchConfidence float64 `xml:"matchConfidence,omitempty"`

	Record *SObject `xml:"record,omitempty"`
}

type MatchResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com MatchResult"`

	EntityType string `xml:"entityType,omitempty"`

	Errors []*Error `xml:"errors,omitempty"`

	MatchEngine string `xml:"matchEngine,omitempty"`

	MatchRecords []*MatchRecord `xml:"matchRecords,omitempty"`

	Rule string `xml:"rule,omitempty"`

	Size int32 `xml:"size,omitempty"`

	Success bool `xml:"success,omitempty"`
}

type DuplicateResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DuplicateResult"`

	AllowSave bool `xml:"allowSave,omitempty"`

	DuplicateRule string `xml:"duplicateRule,omitempty"`

	DuplicateRuleEntityType string `xml:"duplicateRuleEntityType,omitempty"`

	ErrorMessage string `xml:"errorMessage,omitempty"`

	MatchResults []*MatchResult `xml:"matchResults,omitempty"`
}

type DuplicateError struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DuplicateError"`

	*Error

	DuplicateResult *DuplicateResult `xml:"duplicateResult,omitempty"`
}

type DescribeNounResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeNounResult"`

	CaseValues []*NameCaseValue `xml:"caseValues,omitempty"`

	DeveloperName string `xml:"developerName,omitempty"`

	Gender *Gender `xml:"gender,omitempty"`

	Name string `xml:"name,omitempty"`

	PluralAlias string `xml:"pluralAlias,omitempty"`

	StartsWith *StartsWith `xml:"startsWith,omitempty"`
}

type NameCaseValue struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com NameCaseValue"`

	Article *Article `xml:"article,omitempty"`

	CaseType *CaseType `xml:"caseType,omitempty"`

	Number *GrammaticalNumber `xml:"number,omitempty"`

	Possessive *Possessive `xml:"possessive,omitempty"`

	Value string `xml:"value,omitempty"`
}

type FindDuplicatesResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com FindDuplicatesResult"`

	DuplicateResults []*DuplicateResult `xml:"duplicateResults,omitempty"`

	Errors []*Error `xml:"errors,omitempty"`

	Success bool `xml:"success,omitempty"`
}

type DescribeFlexiPageResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeFlexiPageResult"`

	Id string `xml:"id,omitempty"`

	Label string `xml:"label,omitempty"`

	Name string `xml:"name,omitempty"`

	QuickActionList *DescribeQuickActionListResult `xml:"quickActionList,omitempty"`

	Regions []*DescribeFlexiPageRegion `xml:"regions,omitempty"`

	SobjectType string `xml:"sobjectType,omitempty"`

	Template string `xml:"template,omitempty"`

	Type_ string `xml:"type,omitempty"`
}

type DescribeFlexiPageRegion struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeFlexiPageRegion"`

	Components []*DescribeComponentInstance `xml:"components,omitempty"`

	Name string `xml:"name,omitempty"`
}

type DescribeComponentInstance struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeComponentInstance"`

	Properties []*DescribeComponentInstanceProperty `xml:"properties,omitempty"`

	TypeName string `xml:"typeName,omitempty"`

	TypeNamespace string `xml:"typeNamespace,omitempty"`
}

type DescribeComponentInstanceProperty struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeComponentInstanceProperty"`

	Name string `xml:"name,omitempty"`

	Region *DescribeFlexiPageRegion `xml:"region,omitempty"`

	Type_ *ComponentInstancePropertyTypeEnum `xml:"type,omitempty"`

	Value string `xml:"value,omitempty"`
}

type FlexipageContext struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com FlexipageContext"`

	Type_ *FlexipageContextTypeEnum `xml:"type,omitempty"`

	Value string `xml:"value,omitempty"`
}

type DescribeAppMenuResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeAppMenuResult"`

	AppMenuItems []*DescribeAppMenuItem `xml:"appMenuItems,omitempty"`
}

type DescribeAppMenuItem struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeAppMenuItem"`

	Colors []*DescribeColor `xml:"colors,omitempty"`

	Content string `xml:"content,omitempty"`

	Icons []*DescribeIcon `xml:"icons,omitempty"`

	Label string `xml:"label,omitempty"`

	Name string `xml:"name,omitempty"`

	Type_ string `xml:"type,omitempty"`

	Url string `xml:"url,omitempty"`
}

type DescribeThemeResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeThemeResult"`

	ThemeItems []*DescribeThemeItem `xml:"themeItems,omitempty"`
}

type DescribeThemeItem struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeThemeItem"`

	Colors []*DescribeColor `xml:"colors,omitempty"`

	Icons []*DescribeIcon `xml:"icons,omitempty"`

	Name string `xml:"name,omitempty"`
}

type DescribeSoftphoneLayoutResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeSoftphoneLayoutResult"`

	CallTypes []*DescribeSoftphoneLayoutCallType `xml:"callTypes,omitempty"`

	Id string `xml:"id,omitempty"`

	Name string `xml:"name,omitempty"`
}

type DescribeSoftphoneLayoutCallType struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeSoftphoneLayoutCallType"`

	InfoFields []*DescribeSoftphoneLayoutInfoField `xml:"infoFields,omitempty"`

	Name string `xml:"name,omitempty"`

	ScreenPopOptions []*DescribeSoftphoneScreenPopOption `xml:"screenPopOptions,omitempty"`

	ScreenPopsOpenWithin string `xml:"screenPopsOpenWithin,omitempty"`

	Sections []*DescribeSoftphoneLayoutSection `xml:"sections,omitempty"`
}

type DescribeSoftphoneScreenPopOption struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeSoftphoneScreenPopOption"`

	MatchType string `xml:"matchType,omitempty"`

	ScreenPopData string `xml:"screenPopData,omitempty"`

	ScreenPopType string `xml:"screenPopType,omitempty"`
}

type DescribeSoftphoneLayoutInfoField struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeSoftphoneLayoutInfoField"`

	Name string `xml:"name,omitempty"`
}

type DescribeSoftphoneLayoutSection struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeSoftphoneLayoutSection"`

	EntityApiName string `xml:"entityApiName,omitempty"`

	Items []*DescribeSoftphoneLayoutItem `xml:"items,omitempty"`
}

type DescribeSoftphoneLayoutItem struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeSoftphoneLayoutItem"`

	ItemApiName string `xml:"itemApiName,omitempty"`
}

type DescribeCompactLayoutsResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeCompactLayoutsResult"`

	CompactLayouts []*DescribeCompactLayout `xml:"compactLayouts,omitempty"`

	DefaultCompactLayoutId string `xml:"defaultCompactLayoutId,omitempty"`

	RecordTypeCompactLayoutMappings []*RecordTypeCompactLayoutMapping `xml:"recordTypeCompactLayoutMappings,omitempty"`
}

type DescribeCompactLayout struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribeCompactLayout"`

	Actions []*DescribeLayoutButton `xml:"actions,omitempty"`

	FieldItems []*DescribeLayoutItem `xml:"fieldItems,omitempty"`

	Id string `xml:"id,omitempty"`

	ImageItems []*DescribeLayoutItem `xml:"imageItems,omitempty"`

	Label string `xml:"label,omitempty"`

	Name string `xml:"name,omitempty"`

	ObjectType string `xml:"objectType,omitempty"`
}

type RecordTypeCompactLayoutMapping struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com RecordTypeCompactLayoutMapping"`

	Available bool `xml:"available,omitempty"`

	CompactLayoutId string `xml:"compactLayoutId,omitempty"`

	CompactLayoutName string `xml:"compactLayoutName,omitempty"`

	RecordTypeId string `xml:"recordTypeId,omitempty"`

	RecordTypeName string `xml:"recordTypeName,omitempty"`
}

type DescribePathAssistantsResult struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribePathAssistantsResult"`

	PathAssistants []*DescribePathAssistant `xml:"pathAssistants,omitempty"`
}

type DescribePathAssistant struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribePathAssistant"`

	Active bool `xml:"active,omitempty"`

	ApiName string `xml:"apiName,omitempty"`

	Label string `xml:"label,omitempty"`

	PathPicklistField string `xml:"pathPicklistField,omitempty"`

	PicklistsForRecordType []*PicklistForRecordType `xml:"picklistsForRecordType,omitempty"`

	RecordTypeId string `xml:"recordTypeId,omitempty"`

	Steps []*DescribePathAssistantStep `xml:"steps,omitempty"`
}

type DescribePathAssistantStep struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribePathAssistantStep"`

	Closed bool `xml:"closed,omitempty"`

	Converted bool `xml:"converted,omitempty"`

	Fields []*DescribePathAssistantField `xml:"fields,omitempty"`

	Info string `xml:"info,omitempty"`

	LayoutSection *DescribeLayoutSection `xml:"layoutSection,omitempty"`

	PicklistLabel string `xml:"picklistLabel,omitempty"`

	PicklistValue string `xml:"picklistValue,omitempty"`

	Won bool `xml:"won,omitempty"`
}

type DescribePathAssistantField struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com DescribePathAssistantField"`

	ApiName string `xml:"apiName,omitempty"`

	Label string `xml:"label,omitempty"`

	ReadOnly bool `xml:"readOnly,omitempty"`

	Required bool `xml:"required,omitempty"`
}

type DescribeApprovalLayoutResult struct {
	ApprovalLayouts []*DescribeApprovalLayout `xml:"approvalLayouts,omitempty"`
}

type DescribeApprovalLayout struct {
	Id string `xml:"id,omitempty"`

	Label string `xml:"label,omitempty"`

	LayoutItems []*DescribeLayoutItem `xml:"layoutItems,omitempty"`

	Name string `xml:"name,omitempty"`
}

type DescribeLayoutResultResult struct {
	Layouts []*DescribeLayoutResult `xml:"layouts,omitempty"`

	RecordTypeMappings []*RecordTypeMapping `xml:"recordTypeMappings,omitempty"`

	RecordTypeSelectorRequired bool `xml:"recordTypeSelectorRequired,omitempty"`
}

type DescribeLayoutResult struct {
	ButtonLayoutSection *DescribeLayoutButtonSection `xml:"buttonLayoutSection,omitempty"`

	DetailLayoutSections []*DescribeLayoutSection `xml:"detailLayoutSections,omitempty"`

	EditLayoutSections []*DescribeLayoutSection `xml:"editLayoutSections,omitempty"`

	FeedView *DescribeLayoutFeedView `xml:"feedView,omitempty"`

	HighlightsPanelLayoutSection *DescribeLayoutSection `xml:"highlightsPanelLayoutSection,omitempty"`

	Id string `xml:"id,omitempty"`

	QuickActionList *DescribeQuickActionListResult `xml:"quickActionList,omitempty"`

	RelatedContent *RelatedContent `xml:"relatedContent,omitempty"`

	RelatedLists []*RelatedList `xml:"relatedLists,omitempty"`
}

type DescribeQuickActionListResult struct {
	QuickActionListItems []*DescribeQuickActionListItemResult `xml:"quickActionListItems,omitempty"`
}

type DescribeQuickActionListItemResult struct {
	AccessLevelRequired *ShareAccessLevel `xml:"accessLevelRequired,omitempty"`

	Colors []*DescribeColor `xml:"colors,omitempty"`

	IconUrl string `xml:"iconUrl,omitempty"`

	Icons []*DescribeIcon `xml:"icons,omitempty"`

	Label string `xml:"label,omitempty"`

	MiniIconUrl string `xml:"miniIconUrl,omitempty"`

	QuickActionName string `xml:"quickActionName,omitempty"`

	TargetSobjectType string `xml:"targetSobjectType,omitempty"`

	Type_ string `xml:"type,omitempty"`
}

type DescribeLayoutFeedView struct {
	FeedFilters []*DescribeLayoutFeedFilter `xml:"feedFilters,omitempty"`
}

type DescribeLayoutFeedFilter struct {
	Label string `xml:"label,omitempty"`

	Name string `xml:"name,omitempty"`

	Type_ *FeedLayoutFilterType `xml:"type,omitempty"`
}

type DescribeLayoutSection struct {
	Collapsed bool `xml:"collapsed,omitempty"`

	Columns int32 `xml:"columns,omitempty"`

	Heading string `xml:"heading,omitempty"`

	LayoutRows []*DescribeLayoutRow `xml:"layoutRows,omitempty"`

	LayoutSectionId string `xml:"layoutSectionId,omitempty"`

	ParentLayoutId string `xml:"parentLayoutId,omitempty"`

	Rows int32 `xml:"rows,omitempty"`

	TabOrder *TabOrderType `xml:"tabOrder,omitempty"`

	UseCollapsibleSection bool `xml:"useCollapsibleSection,omitempty"`

	UseHeading bool `xml:"useHeading,omitempty"`
}

type DescribeLayoutButtonSection struct {
	DetailButtons []*DescribeLayoutButton `xml:"detailButtons,omitempty"`
}

type DescribeLayoutRow struct {
	LayoutItems []*DescribeLayoutItem `xml:"layoutItems,omitempty"`

	NumItems int32 `xml:"numItems,omitempty"`
}

type DescribeLayoutItem struct {
	EditableForNew bool `xml:"editableForNew,omitempty"`

	EditableForUpdate bool `xml:"editableForUpdate,omitempty"`

	Label string `xml:"label,omitempty"`

	LayoutComponents []*DescribeLayoutComponent `xml:"layoutComponents,omitempty"`

	Placeholder bool `xml:"placeholder,omitempty"`

	Required bool `xml:"required,omitempty"`
}

type DescribeLayoutButton struct {
	Behavior *WebLinkWindowType `xml:"behavior,omitempty"`

	Colors []*DescribeColor `xml:"colors,omitempty"`

	Content string `xml:"content,omitempty"`

	ContentSource *WebLinkType `xml:"contentSource,omitempty"`

	Custom bool `xml:"custom,omitempty"`

	Encoding string `xml:"encoding,omitempty"`

	Height int32 `xml:"height,omitempty"`

	Icons []*DescribeIcon `xml:"icons,omitempty"`

	Label string `xml:"label,omitempty"`

	Menubar bool `xml:"menubar,omitempty"`

	Name string `xml:"name,omitempty"`

	Overridden bool `xml:"overridden,omitempty"`

	Resizeable bool `xml:"resizeable,omitempty"`

	Scrollbars bool `xml:"scrollbars,omitempty"`

	ShowsLocation bool `xml:"showsLocation,omitempty"`

	ShowsStatus bool `xml:"showsStatus,omitempty"`

	Toolbar bool `xml:"toolbar,omitempty"`

	Url string `xml:"url,omitempty"`

	Width int32 `xml:"width,omitempty"`

	WindowPosition *WebLinkPosition `xml:"windowPosition,omitempty"`
}

type DescribeLayoutComponent struct {
	DisplayLines int32 `xml:"displayLines,omitempty"`

	TabOrder int32 `xml:"tabOrder,omitempty"`

	Type_ *LayoutComponentType `xml:"type,omitempty"`

	Value string `xml:"value,omitempty"`
}

type FieldComponent struct {
	*DescribeLayoutComponent

	Field *Field `xml:"field,omitempty"`
}

type FieldLayoutComponent struct {
	*DescribeLayoutComponent

	Components []*DescribeLayoutComponent `xml:"components,omitempty"`

	FieldType *FieldType `xml:"fieldType,omitempty"`
}

type VisualforcePage struct {
	*DescribeLayoutComponent

	ShowLabel bool `xml:"showLabel,omitempty"`

	ShowScrollbars bool `xml:"showScrollbars,omitempty"`

	SuggestedHeight string `xml:"suggestedHeight,omitempty"`

	SuggestedWidth string `xml:"suggestedWidth,omitempty"`

	Url string `xml:"url,omitempty"`
}

type Canvas struct {
	*DescribeLayoutComponent

	DisplayLocation string `xml:"displayLocation,omitempty"`

	ReferenceId string `xml:"referenceId,omitempty"`

	ShowLabel bool `xml:"showLabel,omitempty"`

	ShowScrollbars bool `xml:"showScrollbars,omitempty"`

	SuggestedHeight string `xml:"suggestedHeight,omitempty"`

	SuggestedWidth string `xml:"suggestedWidth,omitempty"`
}

type ReportChartComponent struct {
	*DescribeLayoutComponent

	CacheData bool `xml:"cacheData,omitempty"`

	ContextFilterableField string `xml:"contextFilterableField,omitempty"`

	Error string `xml:"error,omitempty"`

	HideOnError bool `xml:"hideOnError,omitempty"`

	IncludeContext bool `xml:"includeContext,omitempty"`

	ShowTitle bool `xml:"showTitle,omitempty"`

	Size *ReportChartSize `xml:"size,omitempty"`
}

type AnalyticsCloudComponent struct {
	*DescribeLayoutComponent

	Error string `xml:"error,omitempty"`

	Filter string `xml:"filter,omitempty"`

	Height string `xml:"height,omitempty"`

	HideOnError bool `xml:"hideOnError,omitempty"`

	ShowSharing bool `xml:"showSharing,omitempty"`

	ShowTitle bool `xml:"showTitle,omitempty"`

	Width string `xml:"width,omitempty"`
}

type CustomLinkComponent struct {
	*DescribeLayoutComponent

	CustomLink *DescribeLayoutButton `xml:"customLink,omitempty"`
}

type NamedLayoutInfo struct {
	Name string `xml:"name,omitempty"`
}

type RecordTypeInfo struct {
	Available bool `xml:"available,omitempty"`

	DefaultRecordTypeMapping bool `xml:"defaultRecordTypeMapping,omitempty"`

	Master bool `xml:"master,omitempty"`

	Name string `xml:"name,omitempty"`

	RecordTypeId string `xml:"recordTypeId,omitempty"`
}

type RecordTypeMapping struct {
	Available bool `xml:"available,omitempty"`

	DefaultRecordTypeMapping bool `xml:"defaultRecordTypeMapping,omitempty"`

	LayoutId string `xml:"layoutId,omitempty"`

	Master bool `xml:"master,omitempty"`

	Name string `xml:"name,omitempty"`

	PicklistsForRecordType []*PicklistForRecordType `xml:"picklistsForRecordType,omitempty"`

	RecordTypeId string `xml:"recordTypeId,omitempty"`
}

type PicklistForRecordType struct {
	PicklistName string `xml:"picklistName,omitempty"`

	PicklistValues []*PicklistEntry `xml:"picklistValues,omitempty"`
}

type RelatedContent struct {
	RelatedContentItems []*DescribeRelatedContentItem `xml:"relatedContentItems,omitempty"`
}

type DescribeRelatedContentItem struct {
	DescribeLayoutItem *DescribeLayoutItem `xml:"describeLayoutItem,omitempty"`
}

type RelatedList struct {
	AccessLevelRequiredForCreate *ShareAccessLevel `xml:"accessLevelRequiredForCreate,omitempty"`

	Buttons []*DescribeLayoutButton `xml:"buttons,omitempty"`

	Columns []*RelatedListColumn `xml:"columns,omitempty"`

	Custom bool `xml:"custom,omitempty"`

	Field string `xml:"field,omitempty"`

	Label string `xml:"label,omitempty"`

	LimitRows int32 `xml:"limitRows,omitempty"`

	Name string `xml:"name,omitempty"`

	Sobject string `xml:"sobject,omitempty"`

	Sort []*RelatedListSort `xml:"sort,omitempty"`
}

type RelatedListColumn struct {
	Field string `xml:"field,omitempty"`

	FieldApiName string `xml:"fieldApiName,omitempty"`

	Format string `xml:"format,omitempty"`

	Label string `xml:"label,omitempty"`

	LookupId string `xml:"lookupId,omitempty"`

	Name string `xml:"name,omitempty"`
}

type RelatedListSort struct {
	Ascending bool `xml:"ascending,omitempty"`

	Column string `xml:"column,omitempty"`
}

type EmailFileAttachment struct {
	Body []byte `xml:"body,omitempty"`

	ContentType string `xml:"contentType,omitempty"`

	FileName string `xml:"fileName,omitempty"`

	Inline bool `xml:"inline,omitempty"`
}

type Email struct {
	BccSender bool `xml:"bccSender,omitempty"`

	EmailPriority *EmailPriority `xml:"emailPriority,omitempty"`

	ReplyTo string `xml:"replyTo,omitempty"`

	SaveAsActivity bool `xml:"saveAsActivity,omitempty"`

	SenderDisplayName string `xml:"senderDisplayName,omitempty"`

	Subject string `xml:"subject,omitempty"`

	UseSignature bool `xml:"useSignature,omitempty"`
}

type MassEmailMessage struct {
	*Email

	Description string `xml:"description,omitempty"`

	TargetObjectIds string `xml:"targetObjectIds,omitempty"`

	TemplateId string `xml:"templateId,omitempty"`

	WhatIds string `xml:"whatIds,omitempty"`
}

type SingleEmailMessage struct {
	*Email

	BccAddresses string `xml:"bccAddresses,omitempty"`

	CcAddresses string `xml:"ccAddresses,omitempty"`

	Charset string `xml:"charset,omitempty"`

	DocumentAttachments []string `xml:"documentAttachments,omitempty"`

	EntityAttachments []string `xml:"entityAttachments,omitempty"`

	FileAttachments []*EmailFileAttachment `xml:"fileAttachments,omitempty"`

	HtmlBody string `xml:"htmlBody,omitempty"`

	InReplyTo string `xml:"inReplyTo,omitempty"`

	OptOutPolicy *SendEmailOptOutPolicy `xml:"optOutPolicy,omitempty"`

	OrgWideEmailAddressId string `xml:"orgWideEmailAddressId,omitempty"`

	PlainTextBody string `xml:"plainTextBody,omitempty"`

	References string `xml:"references,omitempty"`

	TargetObjectId string `xml:"targetObjectId,omitempty"`

	TemplateId string `xml:"templateId,omitempty"`

	ToAddresses string `xml:"toAddresses,omitempty"`

	TreatBodiesAsTemplate bool `xml:"treatBodiesAsTemplate,omitempty"`

	TreatTargetObjectAsRecipient bool `xml:"treatTargetObjectAsRecipient,omitempty"`

	WhatId string `xml:"whatId,omitempty"`
}

type SendEmailResult struct {
	Errors []*SendEmailError `xml:"errors,omitempty"`

	Success bool `xml:"success,omitempty"`
}

type ListViewColumn struct {
	AscendingLabel string `xml:"ascendingLabel,omitempty"`

	DescendingLabel string `xml:"descendingLabel,omitempty"`

	FieldNameOrPath string `xml:"fieldNameOrPath,omitempty"`

	Hidden bool `xml:"hidden,omitempty"`

	Label string `xml:"label,omitempty"`

	SelectListItem string `xml:"selectListItem,omitempty"`

	SortDirection *OrderByDirection `xml:"sortDirection,omitempty"`

	SortIndex int32 `xml:"sortIndex,omitempty"`

	Sortable bool `xml:"sortable,omitempty"`

	Type_ *FieldType `xml:"type,omitempty"`
}

type ListViewOrderBy struct {
	FieldNameOrPath string `xml:"fieldNameOrPath,omitempty"`

	NullsPosition *OrderByNullsPosition `xml:"nullsPosition,omitempty"`

	SortDirection *OrderByDirection `xml:"sortDirection,omitempty"`
}

type DescribeSoqlListView struct {
	Columns []*ListViewColumn `xml:"columns,omitempty"`

	Id string `xml:"id,omitempty"`

	OrderBy []*ListViewOrderBy `xml:"orderBy,omitempty"`

	Query string `xml:"query,omitempty"`

	Scope string `xml:"scope,omitempty"`

	SobjectType string `xml:"sobjectType,omitempty"`

	WhereCondition *SoqlWhereCondition `xml:"whereCondition,omitempty"`
}

type DescribeSoqlListViewsRequest struct {
	ListViewParams []*DescribeSoqlListViewParams `xml:"listViewParams,omitempty"`
}

type DescribeSoqlListViewParams struct {
	DeveloperNameOrId string `xml:"developerNameOrId,omitempty"`

	SobjectType string `xml:"sobjectType,omitempty"`
}

type DescribeSoqlListViewResult struct {
	DescribeSoqlListViews []*DescribeSoqlListView `xml:"describeSoqlListViews,omitempty"`
}

type ExecuteListViewRequest struct {
	DeveloperNameOrId string `xml:"developerNameOrId,omitempty"`

	Limit int32 `xml:"limit,omitempty"`

	Offset int32 `xml:"offset,omitempty"`

	OrderBy []*ListViewOrderBy `xml:"orderBy,omitempty"`

	SobjectType string `xml:"sobjectType,omitempty"`
}

type ExecuteListViewResult struct {
	Columns []*ListViewColumn `xml:"columns,omitempty"`

	DeveloperName string `xml:"developerName,omitempty"`

	Done bool `xml:"done,omitempty"`

	Id string `xml:"id,omitempty"`

	Label string `xml:"label,omitempty"`

	Records []*ListViewRecord `xml:"records,omitempty"`

	Size int32 `xml:"size,omitempty"`
}

type ListViewRecord struct {
	Columns []*ListViewRecordColumn `xml:"columns,omitempty"`
}

type ListViewRecordColumn struct {
	FieldNameOrPath string `xml:"fieldNameOrPath,omitempty"`

	Value string `xml:"value,omitempty"`
}

type SoqlWhereCondition struct {
	XMLName xml.Name `xml:"urn:partner.soap.sforce.com SoqlWhereCondition"`
}

type SoqlCondition struct {
	*SoqlWhereCondition

	Field string `xml:"field,omitempty"`

	Operator *SoqlOperator `xml:"operator,omitempty"`

	Values []string `xml:"values,omitempty"`
}

type SoqlNotCondition struct {
	*SoqlWhereCondition

	Condition *SoqlWhereCondition `xml:"condition,omitempty"`
}

type SoqlConditionGroup struct {
	*SoqlWhereCondition

	Conditions []*SoqlWhereCondition `xml:"conditions,omitempty"`

	Conjunction *SoqlConjunction `xml:"conjunction,omitempty"`
}

type SoqlSubQueryCondition struct {
	*SoqlWhereCondition

	Field string `xml:"field,omitempty"`

	Operator *SoqlOperator `xml:"operator,omitempty"`

	SubQuery string `xml:"subQuery,omitempty"`
}

type DescribeSearchLayoutResult struct {
	ErrorMsg string `xml:"errorMsg,omitempty"`

	Label string `xml:"label,omitempty"`

	LimitRows int32 `xml:"limitRows,omitempty"`

	ObjectType string `xml:"objectType,omitempty"`

	SearchColumns []*DescribeColumn `xml:"searchColumns,omitempty"`
}

type DescribeColumn struct {
	Field string `xml:"field,omitempty"`

	Format string `xml:"format,omitempty"`

	Label string `xml:"label,omitempty"`

	Name string `xml:"name,omitempty"`
}

type DescribeSearchScopeOrderResult struct {
	KeyPrefix string `xml:"keyPrefix,omitempty"`

	Name string `xml:"name,omitempty"`
}

type DescribeSearchableEntityResult struct {
	Label string `xml:"label,omitempty"`

	Name string `xml:"name,omitempty"`

	PluralLabel string `xml:"pluralLabel,omitempty"`
}

type DescribeTabSetResult struct {
	Description string `xml:"description,omitempty"`

	Label string `xml:"label,omitempty"`

	LogoUrl string `xml:"logoUrl,omitempty"`

	Namespace string `xml:"namespace,omitempty"`

	Selected bool `xml:"selected,omitempty"`

	TabSetId string `xml:"tabSetId,omitempty"`

	Tabs []*DescribeTab `xml:"tabs,omitempty"`
}

type DescribeTab struct {
	Colors []*DescribeColor `xml:"colors,omitempty"`

	Custom bool `xml:"custom,omitempty"`

	IconUrl string `xml:"iconUrl,omitempty"`

	Icons []*DescribeIcon `xml:"icons,omitempty"`

	Label string `xml:"label,omitempty"`

	MiniIconUrl string `xml:"miniIconUrl,omitempty"`

	Name string `xml:"name,omitempty"`

	SobjectName string `xml:"sobjectName,omitempty"`

	Url string `xml:"url,omitempty"`
}

type DescribeColor struct {
	Color string `xml:"color,omitempty"`

	Context string `xml:"context,omitempty"`

	Theme string `xml:"theme,omitempty"`
}

type DescribeIcon struct {
	ContentType string `xml:"contentType,omitempty"`

	Height int32 `xml:"height,omitempty"`

	Theme string `xml:"theme,omitempty"`

	Url string `xml:"url,omitempty"`

	Width int32 `xml:"width,omitempty"`
}

type ActionOverride struct {
	FormFactor string `xml:"formFactor,omitempty"`

	IsAvailableInTouch bool `xml:"isAvailableInTouch,omitempty"`

	Name string `xml:"name,omitempty"`

	PageId string `xml:"pageId,omitempty"`

	Url string `xml:"url,omitempty"`
}

type RenderEmailTemplateRequest struct {
	TemplateBodies string `xml:"templateBodies,omitempty"`

	WhatId string `xml:"whatId,omitempty"`

	WhoId string `xml:"whoId,omitempty"`
}

type RenderEmailTemplateBodyResult struct {
	Errors []*RenderEmailTemplateError `xml:"errors,omitempty"`

	MergedBody string `xml:"mergedBody,omitempty"`

	Success bool `xml:"success,omitempty"`
}

type RenderEmailTemplateResult struct {
	BodyResults *RenderEmailTemplateBodyResult `xml:"bodyResults,omitempty"`

	Errors []*Error `xml:"errors,omitempty"`

	Success bool `xml:"success,omitempty"`
}

type LogInfo struct {
	Category string `xml:"category,omitempty"`

	Level string `xml:"level,omitempty"`
}

type PackageVersion struct {
	MajorNumber int32 `xml:"majorNumber,omitempty"`

	MinorNumber int32 `xml:"minorNumber,omitempty"`

	Namespace string `xml:"namespace,omitempty"`
}

type LimitInfo struct {
	Current int32 `xml:"current,omitempty"`

	Limit int32 `xml:"limit,omitempty"`

	Type_ string `xml:"type,omitempty"`
}

type OwnerChangeOption struct {
	Type_ *OwnerChangeOptionType `xml:"type,omitempty"`

	Execute bool `xml:"execute,omitempty"`
}

type ExceptionCode string

const (
	ExceptionCodeAPEX_TRIGGER_COUPLING_LIMIT ExceptionCode = "APEX_TRIGGER_COUPLING_LIMIT"

	ExceptionCodeAPI_CURRENTLY_DISABLED ExceptionCode = "API_CURRENTLY_DISABLED"

	ExceptionCodeAPI_DISABLED_FOR_ORG ExceptionCode = "API_DISABLED_FOR_ORG"

	ExceptionCodeARGUMENT_OBJECT_PARSE_ERROR ExceptionCode = "ARGUMENT_OBJECT_PARSE_ERROR"

	ExceptionCodeASYNC_OPERATION_LOCATOR ExceptionCode = "ASYNC_OPERATION_LOCATOR"

	ExceptionCodeASYNC_QUERY_UNSUPPORTED_QUERY ExceptionCode = "ASYNC_QUERY_UNSUPPORTED_QUERY"

	ExceptionCodeBATCH_PROCESSING_HALTED ExceptionCode = "BATCH_PROCESSING_HALTED"

	ExceptionCodeBIG_OBJECT_UNSUPPORTED_OPERATION ExceptionCode = "BIG_OBJECT_UNSUPPORTED_OPERATION"

	ExceptionCodeCANNOT_DELETE_ENTITY ExceptionCode = "CANNOT_DELETE_ENTITY"

	ExceptionCodeCANNOT_DELETE_OWNER ExceptionCode = "CANNOT_DELETE_OWNER"

	ExceptionCodeCANT_ADD_STANDADRD_PORTAL_USER_TO_TERRITORY ExceptionCode = "CANT_ADD_STANDADRD_PORTAL_USER_TO_TERRITORY"

	ExceptionCodeCANT_ADD_STANDARD_PORTAL_USER_TO_TERRITORY ExceptionCode = "CANT_ADD_STANDARD_PORTAL_USER_TO_TERRITORY"

	ExceptionCodeCIRCULAR_OBJECT_GRAPH ExceptionCode = "CIRCULAR_OBJECT_GRAPH"

	ExceptionCodeCLIENT_NOT_ACCESSIBLE_FOR_USER ExceptionCode = "CLIENT_NOT_ACCESSIBLE_FOR_USER"

	ExceptionCodeCLIENT_REQUIRE_UPDATE_FOR_USER ExceptionCode = "CLIENT_REQUIRE_UPDATE_FOR_USER"

	ExceptionCodeCONTENT_CUSTOM_DOWNLOAD_EXCEPTION ExceptionCode = "CONTENT_CUSTOM_DOWNLOAD_EXCEPTION"

	ExceptionCodeCONTENT_HUB_AUTHENTICATION_EXCEPTION ExceptionCode = "CONTENT_HUB_AUTHENTICATION_EXCEPTION"

	ExceptionCodeCONTENT_HUB_FILE_DOWNLOAD_EXCEPTION ExceptionCode = "CONTENT_HUB_FILE_DOWNLOAD_EXCEPTION"

	ExceptionCodeCONTENT_HUB_FILE_NOT_FOUND_EXCEPTION ExceptionCode = "CONTENT_HUB_FILE_NOT_FOUND_EXCEPTION"

	ExceptionCodeCONTENT_HUB_INVALID_OBJECT_TYPE_EXCEPTION ExceptionCode = "CONTENT_HUB_INVALID_OBJECT_TYPE_EXCEPTION"

	ExceptionCodeCONTENT_HUB_INVALID_PAGE_NUMBER_EXCEPTION ExceptionCode = "CONTENT_HUB_INVALID_PAGE_NUMBER_EXCEPTION"

	ExceptionCodeCONTENT_HUB_INVALID_PAYLOAD ExceptionCode = "CONTENT_HUB_INVALID_PAYLOAD"

	ExceptionCodeCONTENT_HUB_INVALID_RENDITION_PAGE_NUMBER_EXCEPTION ExceptionCode = "CONTENT_HUB_INVALID_RENDITION_PAGE_NUMBER_EXCEPTION"

	ExceptionCodeCONTENT_HUB_ITEM_TYPE_NOT_FOUND_EXCEPTION ExceptionCode = "CONTENT_HUB_ITEM_TYPE_NOT_FOUND_EXCEPTION"

	ExceptionCodeCONTENT_HUB_OBJECT_NOT_FOUND_EXCEPTION ExceptionCode = "CONTENT_HUB_OBJECT_NOT_FOUND_EXCEPTION"

	ExceptionCodeCONTENT_HUB_OPERATION_NOT_SUPPORTED_EXCEPTION ExceptionCode = "CONTENT_HUB_OPERATION_NOT_SUPPORTED_EXCEPTION"

	ExceptionCodeCONTENT_HUB_SECURITY_EXCEPTION ExceptionCode = "CONTENT_HUB_SECURITY_EXCEPTION"

	ExceptionCodeCONTENT_HUB_TIMEOUT_EXCEPTION ExceptionCode = "CONTENT_HUB_TIMEOUT_EXCEPTION"

	ExceptionCodeCONTENT_HUB_UNEXPECTED_EXCEPTION ExceptionCode = "CONTENT_HUB_UNEXPECTED_EXCEPTION"

	ExceptionCodeCUSTOM_METADATA_LIMIT_EXCEEDED ExceptionCode = "CUSTOM_METADATA_LIMIT_EXCEEDED"

	ExceptionCodeCUSTOM_SETTINGS_LIMIT_EXCEEDED ExceptionCode = "CUSTOM_SETTINGS_LIMIT_EXCEEDED"

	ExceptionCodeDATACLOUD_API_CLIENT_EXCEPTION ExceptionCode = "DATACLOUD_API_CLIENT_EXCEPTION"

	ExceptionCodeDATACLOUD_API_DISABLED_EXCEPTION ExceptionCode = "DATACLOUD_API_DISABLED_EXCEPTION"

	ExceptionCodeDATACLOUD_API_INVALID_QUERY_EXCEPTION ExceptionCode = "DATACLOUD_API_INVALID_QUERY_EXCEPTION"

	ExceptionCodeDATACLOUD_API_SERVER_BUSY_EXCEPTION ExceptionCode = "DATACLOUD_API_SERVER_BUSY_EXCEPTION"

	ExceptionCodeDATACLOUD_API_SERVER_EXCEPTION ExceptionCode = "DATACLOUD_API_SERVER_EXCEPTION"

	ExceptionCodeDATACLOUD_API_TIMEOUT_EXCEPTION ExceptionCode = "DATACLOUD_API_TIMEOUT_EXCEPTION"

	ExceptionCodeDATACLOUD_API_UNAVAILABLE ExceptionCode = "DATACLOUD_API_UNAVAILABLE"

	ExceptionCodeDUPLICATE_ARGUMENT_VALUE ExceptionCode = "DUPLICATE_ARGUMENT_VALUE"

	ExceptionCodeDUPLICATE_VALUE ExceptionCode = "DUPLICATE_VALUE"

	ExceptionCodeEMAIL_BATCH_SIZE_LIMIT_EXCEEDED ExceptionCode = "EMAIL_BATCH_SIZE_LIMIT_EXCEEDED"

	ExceptionCodeEMAIL_TO_CASE_INVALID_ROUTING ExceptionCode = "EMAIL_TO_CASE_INVALID_ROUTING"

	ExceptionCodeEMAIL_TO_CASE_LIMIT_EXCEEDED ExceptionCode = "EMAIL_TO_CASE_LIMIT_EXCEEDED"

	ExceptionCodeEMAIL_TO_CASE_NOT_ENABLED ExceptionCode = "EMAIL_TO_CASE_NOT_ENABLED"

	ExceptionCodeENTITY_NOT_QUERYABLE ExceptionCode = "ENTITY_NOT_QUERYABLE"

	ExceptionCodeENVIRONMENT_HUB_MEMBERSHIP_CONFLICT ExceptionCode = "ENVIRONMENT_HUB_MEMBERSHIP_CONFLICT"

	ExceptionCodeEXCEEDED_ID_LIMIT ExceptionCode = "EXCEEDED_ID_LIMIT"

	ExceptionCodeEXCEEDED_LEAD_CONVERT_LIMIT ExceptionCode = "EXCEEDED_LEAD_CONVERT_LIMIT"

	ExceptionCodeEXCEEDED_MAX_SIZE_REQUEST ExceptionCode = "EXCEEDED_MAX_SIZE_REQUEST"

	ExceptionCodeEXCEEDED_MAX_SOBJECTS ExceptionCode = "EXCEEDED_MAX_SOBJECTS"

	ExceptionCodeEXCEEDED_MAX_TYPES_LIMIT ExceptionCode = "EXCEEDED_MAX_TYPES_LIMIT"

	ExceptionCodeEXCEEDED_QUOTA ExceptionCode = "EXCEEDED_QUOTA"

	ExceptionCodeEXTERNAL_OBJECT_AUTHENTICATION_EXCEPTION ExceptionCode = "EXTERNAL_OBJECT_AUTHENTICATION_EXCEPTION"

	ExceptionCodeEXTERNAL_OBJECT_CONNECTION_EXCEPTION ExceptionCode = "EXTERNAL_OBJECT_CONNECTION_EXCEPTION"

	ExceptionCodeEXTERNAL_OBJECT_EXCEPTION ExceptionCode = "EXTERNAL_OBJECT_EXCEPTION"

	ExceptionCodeEXTERNAL_OBJECT_UNSUPPORTED_EXCEPTION ExceptionCode = "EXTERNAL_OBJECT_UNSUPPORTED_EXCEPTION"

	ExceptionCodeFEDERATED_SEARCH_ERROR ExceptionCode = "FEDERATED_SEARCH_ERROR"

	ExceptionCodeFEED_NOT_ENABLED_FOR_OBJECT ExceptionCode = "FEED_NOT_ENABLED_FOR_OBJECT"

	ExceptionCodeFUNCTIONALITY_NOT_ENABLED ExceptionCode = "FUNCTIONALITY_NOT_ENABLED"

	ExceptionCodeFUNCTIONALITY_TEMPORARILY_UNAVAILABLE ExceptionCode = "FUNCTIONALITY_TEMPORARILY_UNAVAILABLE"

	ExceptionCodeILLEGAL_QUERY_PARAMETER_VALUE ExceptionCode = "ILLEGAL_QUERY_PARAMETER_VALUE"

	ExceptionCodeINACTIVE_OWNER_OR_USER ExceptionCode = "INACTIVE_OWNER_OR_USER"

	ExceptionCodeINACTIVE_PORTAL ExceptionCode = "INACTIVE_PORTAL"

	ExceptionCodeINSERT_UPDATE_DELETE_NOT_ALLOWED_DURING_MAINTENANCE ExceptionCode = "INSERT_UPDATE_DELETE_NOT_ALLOWED_DURING_MAINTENANCE"

	ExceptionCodeINSUFFICIENT_ACCESS ExceptionCode = "INSUFFICIENT_ACCESS"

	ExceptionCodeINTERNAL_CANVAS_ERROR ExceptionCode = "INTERNAL_CANVAS_ERROR"

	ExceptionCodeINVALID_ASSIGNMENT_RULE ExceptionCode = "INVALID_ASSIGNMENT_RULE"

	ExceptionCodeINVALID_BATCH_REQUEST ExceptionCode = "INVALID_BATCH_REQUEST"

	ExceptionCodeINVALID_BATCH_SIZE ExceptionCode = "INVALID_BATCH_SIZE"

	ExceptionCodeINVALID_CLIENT ExceptionCode = "INVALID_CLIENT"

	ExceptionCodeINVALID_CROSS_REFERENCE_KEY ExceptionCode = "INVALID_CROSS_REFERENCE_KEY"

	ExceptionCodeINVALID_DATE_FORMAT ExceptionCode = "INVALID_DATE_FORMAT"

	ExceptionCodeINVALID_FIELD ExceptionCode = "INVALID_FIELD"

	ExceptionCodeINVALID_FILTER_LANGUAGE ExceptionCode = "INVALID_FILTER_LANGUAGE"

	ExceptionCodeINVALID_FILTER_VALUE ExceptionCode = "INVALID_FILTER_VALUE"

	ExceptionCodeINVALID_ID_FIELD ExceptionCode = "INVALID_ID_FIELD"

	ExceptionCodeINVALID_INPUT_COMBINATION ExceptionCode = "INVALID_INPUT_COMBINATION"

	ExceptionCodeINVALID_LOCALE_LANGUAGE ExceptionCode = "INVALID_LOCALE_LANGUAGE"

	ExceptionCodeINVALID_LOCATOR ExceptionCode = "INVALID_LOCATOR"

	ExceptionCodeINVALID_LOGIN ExceptionCode = "INVALID_LOGIN"

	ExceptionCodeINVALID_MULTIPART_REQUEST ExceptionCode = "INVALID_MULTIPART_REQUEST"

	ExceptionCodeINVALID_NEW_PASSWORD ExceptionCode = "INVALID_NEW_PASSWORD"

	ExceptionCodeINVALID_OPERATION ExceptionCode = "INVALID_OPERATION"

	ExceptionCodeINVALID_OPERATION_WITH_EXPIRED_PASSWORD ExceptionCode = "INVALID_OPERATION_WITH_EXPIRED_PASSWORD"

	ExceptionCodeINVALID_PACKAGE_VERSION ExceptionCode = "INVALID_PACKAGE_VERSION"

	ExceptionCodeINVALID_PAGING_OPTION ExceptionCode = "INVALID_PAGING_OPTION"

	ExceptionCodeINVALID_QUERY_FILTER_OPERATOR ExceptionCode = "INVALID_QUERY_FILTER_OPERATOR"

	ExceptionCodeINVALID_QUERY_LOCATOR ExceptionCode = "INVALID_QUERY_LOCATOR"

	ExceptionCodeINVALID_QUERY_SCOPE ExceptionCode = "INVALID_QUERY_SCOPE"

	ExceptionCodeINVALID_REPLICATION_DATE ExceptionCode = "INVALID_REPLICATION_DATE"

	ExceptionCodeINVALID_SEARCH ExceptionCode = "INVALID_SEARCH"

	ExceptionCodeINVALID_SEARCH_SCOPE ExceptionCode = "INVALID_SEARCH_SCOPE"

	ExceptionCodeINVALID_SESSION_ID ExceptionCode = "INVALID_SESSION_ID"

	ExceptionCodeINVALID_SOAP_HEADER ExceptionCode = "INVALID_SOAP_HEADER"

	ExceptionCodeINVALID_SORT_OPTION ExceptionCode = "INVALID_SORT_OPTION"

	ExceptionCodeINVALID_SSO_GATEWAY_URL ExceptionCode = "INVALID_SSO_GATEWAY_URL"

	ExceptionCodeINVALID_TYPE ExceptionCode = "INVALID_TYPE"

	ExceptionCodeINVALID_TYPE_FOR_OPERATION ExceptionCode = "INVALID_TYPE_FOR_OPERATION"

	ExceptionCodeJIGSAW_ACTION_DISABLED ExceptionCode = "JIGSAW_ACTION_DISABLED"

	ExceptionCodeJIGSAW_IMPORT_LIMIT_EXCEEDED ExceptionCode = "JIGSAW_IMPORT_LIMIT_EXCEEDED"

	ExceptionCodeJIGSAW_REQUEST_NOT_SUPPORTED ExceptionCode = "JIGSAW_REQUEST_NOT_SUPPORTED"

	ExceptionCodeJSON_PARSER_ERROR ExceptionCode = "JSON_PARSER_ERROR"

	ExceptionCodeKEY_HAS_BEEN_DESTROYED ExceptionCode = "KEY_HAS_BEEN_DESTROYED"

	ExceptionCodeLICENSING_DATA_ERROR ExceptionCode = "LICENSING_DATA_ERROR"

	ExceptionCodeLICENSING_UNKNOWN_ERROR ExceptionCode = "LICENSING_UNKNOWN_ERROR"

	ExceptionCodeLIMIT_EXCEEDED ExceptionCode = "LIMIT_EXCEEDED"

	ExceptionCodeLOGIN_CHALLENGE_ISSUED ExceptionCode = "LOGIN_CHALLENGE_ISSUED"

	ExceptionCodeLOGIN_CHALLENGE_PENDING ExceptionCode = "LOGIN_CHALLENGE_PENDING"

	ExceptionCodeLOGIN_DURING_RESTRICTED_DOMAIN ExceptionCode = "LOGIN_DURING_RESTRICTED_DOMAIN"

	ExceptionCodeLOGIN_DURING_RESTRICTED_TIME ExceptionCode = "LOGIN_DURING_RESTRICTED_TIME"

	ExceptionCodeLOGIN_MUST_USE_SECURITY_TOKEN ExceptionCode = "LOGIN_MUST_USE_SECURITY_TOKEN"

	ExceptionCodeMALFORMED_ID ExceptionCode = "MALFORMED_ID"

	ExceptionCodeMALFORMED_QUERY ExceptionCode = "MALFORMED_QUERY"

	ExceptionCodeMALFORMED_SEARCH ExceptionCode = "MALFORMED_SEARCH"

	ExceptionCodeMISSING_ARGUMENT ExceptionCode = "MISSING_ARGUMENT"

	ExceptionCodeMISSING_RECORD ExceptionCode = "MISSING_RECORD"

	ExceptionCodeMODIFIED ExceptionCode = "MODIFIED"

	ExceptionCodeMUTUAL_AUTHENTICATION_FAILED ExceptionCode = "MUTUAL_AUTHENTICATION_FAILED"

	ExceptionCodeNOT_ACCEPTABLE ExceptionCode = "NOT_ACCEPTABLE"

	ExceptionCodeNOT_MODIFIED ExceptionCode = "NOT_MODIFIED"

	ExceptionCodeNO_ACTIVE_DUPLICATE_RULE ExceptionCode = "NO_ACTIVE_DUPLICATE_RULE"

	ExceptionCodeNO_SOFTPHONE_LAYOUT ExceptionCode = "NO_SOFTPHONE_LAYOUT"

	ExceptionCodeNUMBER_OUTSIDE_VALID_RANGE ExceptionCode = "NUMBER_OUTSIDE_VALID_RANGE"

	ExceptionCodeOPERATION_TOO_LARGE ExceptionCode = "OPERATION_TOO_LARGE"

	ExceptionCodeORG_IN_MAINTENANCE ExceptionCode = "ORG_IN_MAINTENANCE"

	ExceptionCodeORG_IS_DOT_ORG ExceptionCode = "ORG_IS_DOT_ORG"

	ExceptionCodeORG_IS_SIGNING_UP ExceptionCode = "ORG_IS_SIGNING_UP"

	ExceptionCodeORG_LOCKED ExceptionCode = "ORG_LOCKED"

	ExceptionCodeORG_NOT_OWNED_BY_INSTANCE ExceptionCode = "ORG_NOT_OWNED_BY_INSTANCE"

	ExceptionCodePASSWORD_LOCKOUT ExceptionCode = "PASSWORD_LOCKOUT"

	ExceptionCodePORTAL_NO_ACCESS ExceptionCode = "PORTAL_NO_ACCESS"

	ExceptionCodePOST_BODY_PARSE_ERROR ExceptionCode = "POST_BODY_PARSE_ERROR"

	ExceptionCodeQUERY_TIMEOUT ExceptionCode = "QUERY_TIMEOUT"

	ExceptionCodeQUERY_TOO_COMPLICATED ExceptionCode = "QUERY_TOO_COMPLICATED"

	ExceptionCodeREQUEST_LIMIT_EXCEEDED ExceptionCode = "REQUEST_LIMIT_EXCEEDED"

	ExceptionCodeREQUEST_RUNNING_TOO_LONG ExceptionCode = "REQUEST_RUNNING_TOO_LONG"

	ExceptionCodeSERVER_UNAVAILABLE ExceptionCode = "SERVER_UNAVAILABLE"

	ExceptionCodeSERVICE_DESK_NOT_ENABLED ExceptionCode = "SERVICE_DESK_NOT_ENABLED"

	ExceptionCodeSOCIALCRM_FEEDSERVICE_API_CLIENT_EXCEPTION ExceptionCode = "SOCIALCRM_FEEDSERVICE_API_CLIENT_EXCEPTION"

	ExceptionCodeSOCIALCRM_FEEDSERVICE_API_SERVER_EXCEPTION ExceptionCode = "SOCIALCRM_FEEDSERVICE_API_SERVER_EXCEPTION"

	ExceptionCodeSOCIALCRM_FEEDSERVICE_API_UNAVAILABLE ExceptionCode = "SOCIALCRM_FEEDSERVICE_API_UNAVAILABLE"

	ExceptionCodeSSO_SERVICE_DOWN ExceptionCode = "SSO_SERVICE_DOWN"

	ExceptionCodeSST_ADMIN_FILE_DOWNLOAD_EXCEPTION ExceptionCode = "SST_ADMIN_FILE_DOWNLOAD_EXCEPTION"

	ExceptionCodeTOO_MANY_APEX_REQUESTS ExceptionCode = "TOO_MANY_APEX_REQUESTS"

	ExceptionCodeTOO_MANY_RECIPIENTS ExceptionCode = "TOO_MANY_RECIPIENTS"

	ExceptionCodeTOO_MANY_RECORDS ExceptionCode = "TOO_MANY_RECORDS"

	ExceptionCodeTRIAL_EXPIRED ExceptionCode = "TRIAL_EXPIRED"

	ExceptionCodeTXN_SECURITY_END_A_SESSION ExceptionCode = "TXN_SECURITY_END_A_SESSION"

	ExceptionCodeTXN_SECURITY_NO_ACCESS ExceptionCode = "TXN_SECURITY_NO_ACCESS"

	ExceptionCodeTXN_SECURITY_TWO_FA_REQUIRED ExceptionCode = "TXN_SECURITY_TWO_FA_REQUIRED"

	ExceptionCodeUNABLE_TO_LOCK_ROW ExceptionCode = "UNABLE_TO_LOCK_ROW"

	ExceptionCodeUNKNOWN_ATTACHMENT_EXCEPTION ExceptionCode = "UNKNOWN_ATTACHMENT_EXCEPTION"

	ExceptionCodeUNKNOWN_EXCEPTION ExceptionCode = "UNKNOWN_EXCEPTION"

	ExceptionCodeUNKNOWN_ORG_SETTING ExceptionCode = "UNKNOWN_ORG_SETTING"

	ExceptionCodeUNSUPPORTED_API_VERSION ExceptionCode = "UNSUPPORTED_API_VERSION"

	ExceptionCodeUNSUPPORTED_ATTACHMENT_ENCODING ExceptionCode = "UNSUPPORTED_ATTACHMENT_ENCODING"

	ExceptionCodeUNSUPPORTED_CLIENT ExceptionCode = "UNSUPPORTED_CLIENT"

	ExceptionCodeUNSUPPORTED_MEDIA_TYPE ExceptionCode = "UNSUPPORTED_MEDIA_TYPE"

	ExceptionCodeXML_PARSER_ERROR ExceptionCode = "XML_PARSER_ERROR"
)

type FaultCode string

const (
	FaultCodeFnsAPEX_TRIGGER_COUPLING_LIMIT FaultCode = "fns:APEX_TRIGGER_COUPLING_LIMIT"

	FaultCodeFnsAPI_CURRENTLY_DISABLED FaultCode = "fns:API_CURRENTLY_DISABLED"

	FaultCodeFnsAPI_DISABLED_FOR_ORG FaultCode = "fns:API_DISABLED_FOR_ORG"

	FaultCodeFnsARGUMENT_OBJECT_PARSE_ERROR FaultCode = "fns:ARGUMENT_OBJECT_PARSE_ERROR"

	FaultCodeFnsASYNC_OPERATION_LOCATOR FaultCode = "fns:ASYNC_OPERATION_LOCATOR"

	FaultCodeFnsASYNC_QUERY_UNSUPPORTED_QUERY FaultCode = "fns:ASYNC_QUERY_UNSUPPORTED_QUERY"

	FaultCodeFnsBATCH_PROCESSING_HALTED FaultCode = "fns:BATCH_PROCESSING_HALTED"

	FaultCodeFnsBIG_OBJECT_UNSUPPORTED_OPERATION FaultCode = "fns:BIG_OBJECT_UNSUPPORTED_OPERATION"

	FaultCodeFnsCANNOT_DELETE_ENTITY FaultCode = "fns:CANNOT_DELETE_ENTITY"

	FaultCodeFnsCANNOT_DELETE_OWNER FaultCode = "fns:CANNOT_DELETE_OWNER"

	FaultCodeFnsCANT_ADD_STANDADRD_PORTAL_USER_TO_TERRITORY FaultCode = "fns:CANT_ADD_STANDADRD_PORTAL_USER_TO_TERRITORY"

	FaultCodeFnsCANT_ADD_STANDARD_PORTAL_USER_TO_TERRITORY FaultCode = "fns:CANT_ADD_STANDARD_PORTAL_USER_TO_TERRITORY"

	FaultCodeFnsCIRCULAR_OBJECT_GRAPH FaultCode = "fns:CIRCULAR_OBJECT_GRAPH"

	FaultCodeFnsCLIENT_NOT_ACCESSIBLE_FOR_USER FaultCode = "fns:CLIENT_NOT_ACCESSIBLE_FOR_USER"

	FaultCodeFnsCLIENT_REQUIRE_UPDATE_FOR_USER FaultCode = "fns:CLIENT_REQUIRE_UPDATE_FOR_USER"

	FaultCodeFnsCONTENT_CUSTOM_DOWNLOAD_EXCEPTION FaultCode = "fns:CONTENT_CUSTOM_DOWNLOAD_EXCEPTION"

	FaultCodeFnsCONTENT_HUB_AUTHENTICATION_EXCEPTION FaultCode = "fns:CONTENT_HUB_AUTHENTICATION_EXCEPTION"

	FaultCodeFnsCONTENT_HUB_FILE_DOWNLOAD_EXCEPTION FaultCode = "fns:CONTENT_HUB_FILE_DOWNLOAD_EXCEPTION"

	FaultCodeFnsCONTENT_HUB_FILE_NOT_FOUND_EXCEPTION FaultCode = "fns:CONTENT_HUB_FILE_NOT_FOUND_EXCEPTION"

	FaultCodeFnsCONTENT_HUB_INVALID_OBJECT_TYPE_EXCEPTION FaultCode = "fns:CONTENT_HUB_INVALID_OBJECT_TYPE_EXCEPTION"

	FaultCodeFnsCONTENT_HUB_INVALID_PAGE_NUMBER_EXCEPTION FaultCode = "fns:CONTENT_HUB_INVALID_PAGE_NUMBER_EXCEPTION"

	FaultCodeFnsCONTENT_HUB_INVALID_PAYLOAD FaultCode = "fns:CONTENT_HUB_INVALID_PAYLOAD"

	FaultCodeFnsCONTENT_HUB_INVALID_RENDITION_PAGE_NUMBER_EXCEPTION FaultCode = "fns:CONTENT_HUB_INVALID_RENDITION_PAGE_NUMBER_EXCEPTION"

	FaultCodeFnsCONTENT_HUB_ITEM_TYPE_NOT_FOUND_EXCEPTION FaultCode = "fns:CONTENT_HUB_ITEM_TYPE_NOT_FOUND_EXCEPTION"

	FaultCodeFnsCONTENT_HUB_OBJECT_NOT_FOUND_EXCEPTION FaultCode = "fns:CONTENT_HUB_OBJECT_NOT_FOUND_EXCEPTION"

	FaultCodeFnsCONTENT_HUB_OPERATION_NOT_SUPPORTED_EXCEPTION FaultCode = "fns:CONTENT_HUB_OPERATION_NOT_SUPPORTED_EXCEPTION"

	FaultCodeFnsCONTENT_HUB_SECURITY_EXCEPTION FaultCode = "fns:CONTENT_HUB_SECURITY_EXCEPTION"

	FaultCodeFnsCONTENT_HUB_TIMEOUT_EXCEPTION FaultCode = "fns:CONTENT_HUB_TIMEOUT_EXCEPTION"

	FaultCodeFnsCONTENT_HUB_UNEXPECTED_EXCEPTION FaultCode = "fns:CONTENT_HUB_UNEXPECTED_EXCEPTION"

	FaultCodeFnsCUSTOM_METADATA_LIMIT_EXCEEDED FaultCode = "fns:CUSTOM_METADATA_LIMIT_EXCEEDED"

	FaultCodeFnsCUSTOM_SETTINGS_LIMIT_EXCEEDED FaultCode = "fns:CUSTOM_SETTINGS_LIMIT_EXCEEDED"

	FaultCodeFnsDATACLOUD_API_CLIENT_EXCEPTION FaultCode = "fns:DATACLOUD_API_CLIENT_EXCEPTION"

	FaultCodeFnsDATACLOUD_API_DISABLED_EXCEPTION FaultCode = "fns:DATACLOUD_API_DISABLED_EXCEPTION"

	FaultCodeFnsDATACLOUD_API_INVALID_QUERY_EXCEPTION FaultCode = "fns:DATACLOUD_API_INVALID_QUERY_EXCEPTION"

	FaultCodeFnsDATACLOUD_API_SERVER_BUSY_EXCEPTION FaultCode = "fns:DATACLOUD_API_SERVER_BUSY_EXCEPTION"

	FaultCodeFnsDATACLOUD_API_SERVER_EXCEPTION FaultCode = "fns:DATACLOUD_API_SERVER_EXCEPTION"

	FaultCodeFnsDATACLOUD_API_TIMEOUT_EXCEPTION FaultCode = "fns:DATACLOUD_API_TIMEOUT_EXCEPTION"

	FaultCodeFnsDATACLOUD_API_UNAVAILABLE FaultCode = "fns:DATACLOUD_API_UNAVAILABLE"

	FaultCodeFnsDUPLICATE_ARGUMENT_VALUE FaultCode = "fns:DUPLICATE_ARGUMENT_VALUE"

	FaultCodeFnsDUPLICATE_VALUE FaultCode = "fns:DUPLICATE_VALUE"

	FaultCodeFnsEMAIL_BATCH_SIZE_LIMIT_EXCEEDED FaultCode = "fns:EMAIL_BATCH_SIZE_LIMIT_EXCEEDED"

	FaultCodeFnsEMAIL_TO_CASE_INVALID_ROUTING FaultCode = "fns:EMAIL_TO_CASE_INVALID_ROUTING"

	FaultCodeFnsEMAIL_TO_CASE_LIMIT_EXCEEDED FaultCode = "fns:EMAIL_TO_CASE_LIMIT_EXCEEDED"

	FaultCodeFnsEMAIL_TO_CASE_NOT_ENABLED FaultCode = "fns:EMAIL_TO_CASE_NOT_ENABLED"

	FaultCodeFnsENTITY_NOT_QUERYABLE FaultCode = "fns:ENTITY_NOT_QUERYABLE"

	FaultCodeFnsENVIRONMENT_HUB_MEMBERSHIP_CONFLICT FaultCode = "fns:ENVIRONMENT_HUB_MEMBERSHIP_CONFLICT"

	FaultCodeFnsEXCEEDED_ID_LIMIT FaultCode = "fns:EXCEEDED_ID_LIMIT"

	FaultCodeFnsEXCEEDED_LEAD_CONVERT_LIMIT FaultCode = "fns:EXCEEDED_LEAD_CONVERT_LIMIT"

	FaultCodeFnsEXCEEDED_MAX_SIZE_REQUEST FaultCode = "fns:EXCEEDED_MAX_SIZE_REQUEST"

	FaultCodeFnsEXCEEDED_MAX_SOBJECTS FaultCode = "fns:EXCEEDED_MAX_SOBJECTS"

	FaultCodeFnsEXCEEDED_MAX_TYPES_LIMIT FaultCode = "fns:EXCEEDED_MAX_TYPES_LIMIT"

	FaultCodeFnsEXCEEDED_QUOTA FaultCode = "fns:EXCEEDED_QUOTA"

	FaultCodeFnsEXTERNAL_OBJECT_AUTHENTICATION_EXCEPTION FaultCode = "fns:EXTERNAL_OBJECT_AUTHENTICATION_EXCEPTION"

	FaultCodeFnsEXTERNAL_OBJECT_CONNECTION_EXCEPTION FaultCode = "fns:EXTERNAL_OBJECT_CONNECTION_EXCEPTION"

	FaultCodeFnsEXTERNAL_OBJECT_EXCEPTION FaultCode = "fns:EXTERNAL_OBJECT_EXCEPTION"

	FaultCodeFnsEXTERNAL_OBJECT_UNSUPPORTED_EXCEPTION FaultCode = "fns:EXTERNAL_OBJECT_UNSUPPORTED_EXCEPTION"

	FaultCodeFnsFEDERATED_SEARCH_ERROR FaultCode = "fns:FEDERATED_SEARCH_ERROR"

	FaultCodeFnsFEED_NOT_ENABLED_FOR_OBJECT FaultCode = "fns:FEED_NOT_ENABLED_FOR_OBJECT"

	FaultCodeFnsFUNCTIONALITY_NOT_ENABLED FaultCode = "fns:FUNCTIONALITY_NOT_ENABLED"

	FaultCodeFnsFUNCTIONALITY_TEMPORARILY_UNAVAILABLE FaultCode = "fns:FUNCTIONALITY_TEMPORARILY_UNAVAILABLE"

	FaultCodeFnsILLEGAL_QUERY_PARAMETER_VALUE FaultCode = "fns:ILLEGAL_QUERY_PARAMETER_VALUE"

	FaultCodeFnsINACTIVE_OWNER_OR_USER FaultCode = "fns:INACTIVE_OWNER_OR_USER"

	FaultCodeFnsINACTIVE_PORTAL FaultCode = "fns:INACTIVE_PORTAL"

	FaultCodeFnsINSERT_UPDATE_DELETE_NOT_ALLOWED_DURING_MAINTENANCE FaultCode = "fns:INSERT_UPDATE_DELETE_NOT_ALLOWED_DURING_MAINTENANCE"

	FaultCodeFnsINSUFFICIENT_ACCESS FaultCode = "fns:INSUFFICIENT_ACCESS"

	FaultCodeFnsINTERNAL_CANVAS_ERROR FaultCode = "fns:INTERNAL_CANVAS_ERROR"

	FaultCodeFnsINVALID_ASSIGNMENT_RULE FaultCode = "fns:INVALID_ASSIGNMENT_RULE"

	FaultCodeFnsINVALID_BATCH_REQUEST FaultCode = "fns:INVALID_BATCH_REQUEST"

	FaultCodeFnsINVALID_BATCH_SIZE FaultCode = "fns:INVALID_BATCH_SIZE"

	FaultCodeFnsINVALID_CLIENT FaultCode = "fns:INVALID_CLIENT"

	FaultCodeFnsINVALID_CROSS_REFERENCE_KEY FaultCode = "fns:INVALID_CROSS_REFERENCE_KEY"

	FaultCodeFnsINVALID_DATE_FORMAT FaultCode = "fns:INVALID_DATE_FORMAT"

	FaultCodeFnsINVALID_FIELD FaultCode = "fns:INVALID_FIELD"

	FaultCodeFnsINVALID_FILTER_LANGUAGE FaultCode = "fns:INVALID_FILTER_LANGUAGE"

	FaultCodeFnsINVALID_FILTER_VALUE FaultCode = "fns:INVALID_FILTER_VALUE"

	FaultCodeFnsINVALID_ID_FIELD FaultCode = "fns:INVALID_ID_FIELD"

	FaultCodeFnsINVALID_INPUT_COMBINATION FaultCode = "fns:INVALID_INPUT_COMBINATION"

	FaultCodeFnsINVALID_LOCALE_LANGUAGE FaultCode = "fns:INVALID_LOCALE_LANGUAGE"

	FaultCodeFnsINVALID_LOCATOR FaultCode = "fns:INVALID_LOCATOR"

	FaultCodeFnsINVALID_LOGIN FaultCode = "fns:INVALID_LOGIN"

	FaultCodeFnsINVALID_MULTIPART_REQUEST FaultCode = "fns:INVALID_MULTIPART_REQUEST"

	FaultCodeFnsINVALID_NEW_PASSWORD FaultCode = "fns:INVALID_NEW_PASSWORD"

	FaultCodeFnsINVALID_OPERATION FaultCode = "fns:INVALID_OPERATION"

	FaultCodeFnsINVALID_OPERATION_WITH_EXPIRED_PASSWORD FaultCode = "fns:INVALID_OPERATION_WITH_EXPIRED_PASSWORD"

	FaultCodeFnsINVALID_PACKAGE_VERSION FaultCode = "fns:INVALID_PACKAGE_VERSION"

	FaultCodeFnsINVALID_PAGING_OPTION FaultCode = "fns:INVALID_PAGING_OPTION"

	FaultCodeFnsINVALID_QUERY_FILTER_OPERATOR FaultCode = "fns:INVALID_QUERY_FILTER_OPERATOR"

	FaultCodeFnsINVALID_QUERY_LOCATOR FaultCode = "fns:INVALID_QUERY_LOCATOR"

	FaultCodeFnsINVALID_QUERY_SCOPE FaultCode = "fns:INVALID_QUERY_SCOPE"

	FaultCodeFnsINVALID_REPLICATION_DATE FaultCode = "fns:INVALID_REPLICATION_DATE"

	FaultCodeFnsINVALID_SEARCH FaultCode = "fns:INVALID_SEARCH"

	FaultCodeFnsINVALID_SEARCH_SCOPE FaultCode = "fns:INVALID_SEARCH_SCOPE"

	FaultCodeFnsINVALID_SESSION_ID FaultCode = "fns:INVALID_SESSION_ID"

	FaultCodeFnsINVALID_SOAP_HEADER FaultCode = "fns:INVALID_SOAP_HEADER"

	FaultCodeFnsINVALID_SORT_OPTION FaultCode = "fns:INVALID_SORT_OPTION"

	FaultCodeFnsINVALID_SSO_GATEWAY_URL FaultCode = "fns:INVALID_SSO_GATEWAY_URL"

	FaultCodeFnsINVALID_TYPE FaultCode = "fns:INVALID_TYPE"

	FaultCodeFnsINVALID_TYPE_FOR_OPERATION FaultCode = "fns:INVALID_TYPE_FOR_OPERATION"

	FaultCodeFnsJIGSAW_ACTION_DISABLED FaultCode = "fns:JIGSAW_ACTION_DISABLED"

	FaultCodeFnsJIGSAW_IMPORT_LIMIT_EXCEEDED FaultCode = "fns:JIGSAW_IMPORT_LIMIT_EXCEEDED"

	FaultCodeFnsJIGSAW_REQUEST_NOT_SUPPORTED FaultCode = "fns:JIGSAW_REQUEST_NOT_SUPPORTED"

	FaultCodeFnsJSON_PARSER_ERROR FaultCode = "fns:JSON_PARSER_ERROR"

	FaultCodeFnsKEY_HAS_BEEN_DESTROYED FaultCode = "fns:KEY_HAS_BEEN_DESTROYED"

	FaultCodeFnsLICENSING_DATA_ERROR FaultCode = "fns:LICENSING_DATA_ERROR"

	FaultCodeFnsLICENSING_UNKNOWN_ERROR FaultCode = "fns:LICENSING_UNKNOWN_ERROR"

	FaultCodeFnsLIMIT_EXCEEDED FaultCode = "fns:LIMIT_EXCEEDED"

	FaultCodeFnsLOGIN_CHALLENGE_ISSUED FaultCode = "fns:LOGIN_CHALLENGE_ISSUED"

	FaultCodeFnsLOGIN_CHALLENGE_PENDING FaultCode = "fns:LOGIN_CHALLENGE_PENDING"

	FaultCodeFnsLOGIN_DURING_RESTRICTED_DOMAIN FaultCode = "fns:LOGIN_DURING_RESTRICTED_DOMAIN"

	FaultCodeFnsLOGIN_DURING_RESTRICTED_TIME FaultCode = "fns:LOGIN_DURING_RESTRICTED_TIME"

	FaultCodeFnsLOGIN_MUST_USE_SECURITY_TOKEN FaultCode = "fns:LOGIN_MUST_USE_SECURITY_TOKEN"

	FaultCodeFnsMALFORMED_ID FaultCode = "fns:MALFORMED_ID"

	FaultCodeFnsMALFORMED_QUERY FaultCode = "fns:MALFORMED_QUERY"

	FaultCodeFnsMALFORMED_SEARCH FaultCode = "fns:MALFORMED_SEARCH"

	FaultCodeFnsMISSING_ARGUMENT FaultCode = "fns:MISSING_ARGUMENT"

	FaultCodeFnsMISSING_RECORD FaultCode = "fns:MISSING_RECORD"

	FaultCodeFnsMODIFIED FaultCode = "fns:MODIFIED"

	FaultCodeFnsMUTUAL_AUTHENTICATION_FAILED FaultCode = "fns:MUTUAL_AUTHENTICATION_FAILED"

	FaultCodeFnsNOT_ACCEPTABLE FaultCode = "fns:NOT_ACCEPTABLE"

	FaultCodeFnsNOT_MODIFIED FaultCode = "fns:NOT_MODIFIED"

	FaultCodeFnsNO_ACTIVE_DUPLICATE_RULE FaultCode = "fns:NO_ACTIVE_DUPLICATE_RULE"

	FaultCodeFnsNO_SOFTPHONE_LAYOUT FaultCode = "fns:NO_SOFTPHONE_LAYOUT"

	FaultCodeFnsNUMBER_OUTSIDE_VALID_RANGE FaultCode = "fns:NUMBER_OUTSIDE_VALID_RANGE"

	FaultCodeFnsOPERATION_TOO_LARGE FaultCode = "fns:OPERATION_TOO_LARGE"

	FaultCodeFnsORG_IN_MAINTENANCE FaultCode = "fns:ORG_IN_MAINTENANCE"

	FaultCodeFnsORG_IS_DOT_ORG FaultCode = "fns:ORG_IS_DOT_ORG"

	FaultCodeFnsORG_IS_SIGNING_UP FaultCode = "fns:ORG_IS_SIGNING_UP"

	FaultCodeFnsORG_LOCKED FaultCode = "fns:ORG_LOCKED"

	FaultCodeFnsORG_NOT_OWNED_BY_INSTANCE FaultCode = "fns:ORG_NOT_OWNED_BY_INSTANCE"

	FaultCodeFnsPASSWORD_LOCKOUT FaultCode = "fns:PASSWORD_LOCKOUT"

	FaultCodeFnsPORTAL_NO_ACCESS FaultCode = "fns:PORTAL_NO_ACCESS"

	FaultCodeFnsPOST_BODY_PARSE_ERROR FaultCode = "fns:POST_BODY_PARSE_ERROR"

	FaultCodeFnsQUERY_TIMEOUT FaultCode = "fns:QUERY_TIMEOUT"

	FaultCodeFnsQUERY_TOO_COMPLICATED FaultCode = "fns:QUERY_TOO_COMPLICATED"

	FaultCodeFnsREQUEST_LIMIT_EXCEEDED FaultCode = "fns:REQUEST_LIMIT_EXCEEDED"

	FaultCodeFnsREQUEST_RUNNING_TOO_LONG FaultCode = "fns:REQUEST_RUNNING_TOO_LONG"

	FaultCodeFnsSERVER_UNAVAILABLE FaultCode = "fns:SERVER_UNAVAILABLE"

	FaultCodeFnsSERVICE_DESK_NOT_ENABLED FaultCode = "fns:SERVICE_DESK_NOT_ENABLED"

	FaultCodeFnsSOCIALCRM_FEEDSERVICE_API_CLIENT_EXCEPTION FaultCode = "fns:SOCIALCRM_FEEDSERVICE_API_CLIENT_EXCEPTION"

	FaultCodeFnsSOCIALCRM_FEEDSERVICE_API_SERVER_EXCEPTION FaultCode = "fns:SOCIALCRM_FEEDSERVICE_API_SERVER_EXCEPTION"

	FaultCodeFnsSOCIALCRM_FEEDSERVICE_API_UNAVAILABLE FaultCode = "fns:SOCIALCRM_FEEDSERVICE_API_UNAVAILABLE"

	FaultCodeFnsSSO_SERVICE_DOWN FaultCode = "fns:SSO_SERVICE_DOWN"

	FaultCodeFnsSST_ADMIN_FILE_DOWNLOAD_EXCEPTION FaultCode = "fns:SST_ADMIN_FILE_DOWNLOAD_EXCEPTION"

	FaultCodeFnsTOO_MANY_APEX_REQUESTS FaultCode = "fns:TOO_MANY_APEX_REQUESTS"

	FaultCodeFnsTOO_MANY_RECIPIENTS FaultCode = "fns:TOO_MANY_RECIPIENTS"

	FaultCodeFnsTOO_MANY_RECORDS FaultCode = "fns:TOO_MANY_RECORDS"

	FaultCodeFnsTRIAL_EXPIRED FaultCode = "fns:TRIAL_EXPIRED"

	FaultCodeFnsTXN_SECURITY_END_A_SESSION FaultCode = "fns:TXN_SECURITY_END_A_SESSION"

	FaultCodeFnsTXN_SECURITY_NO_ACCESS FaultCode = "fns:TXN_SECURITY_NO_ACCESS"

	FaultCodeFnsTXN_SECURITY_TWO_FA_REQUIRED FaultCode = "fns:TXN_SECURITY_TWO_FA_REQUIRED"

	FaultCodeFnsUNABLE_TO_LOCK_ROW FaultCode = "fns:UNABLE_TO_LOCK_ROW"

	FaultCodeFnsUNKNOWN_ATTACHMENT_EXCEPTION FaultCode = "fns:UNKNOWN_ATTACHMENT_EXCEPTION"

	FaultCodeFnsUNKNOWN_EXCEPTION FaultCode = "fns:UNKNOWN_EXCEPTION"

	FaultCodeFnsUNKNOWN_ORG_SETTING FaultCode = "fns:UNKNOWN_ORG_SETTING"

	FaultCodeFnsUNSUPPORTED_API_VERSION FaultCode = "fns:UNSUPPORTED_API_VERSION"

	FaultCodeFnsUNSUPPORTED_ATTACHMENT_ENCODING FaultCode = "fns:UNSUPPORTED_ATTACHMENT_ENCODING"

	FaultCodeFnsUNSUPPORTED_CLIENT FaultCode = "fns:UNSUPPORTED_CLIENT"

	FaultCodeFnsUNSUPPORTED_MEDIA_TYPE FaultCode = "fns:UNSUPPORTED_MEDIA_TYPE"

	FaultCodeFnsXML_PARSER_ERROR FaultCode = "fns:XML_PARSER_ERROR"
)

type ApiFault struct {
	XMLName xml.Name `xml:"urn:fault.partner.soap.sforce.com ApiFault"`

	ExceptionCode *ExceptionCode `xml:"exceptionCode,omitempty"`

	ExceptionMessage string `xml:"exceptionMessage,omitempty"`

	ExtendedErrorDetails []*ExtendedErrorDetails `xml:"extendedErrorDetails,omitempty"`
}

type ApiQueryFault struct {
	XMLName xml.Name `xml:"urn:fault.partner.soap.sforce.com ApiQueryFault"`

	*ApiFault

	Row int32 `xml:"row,omitempty"`

	Column int32 `xml:"column,omitempty"`
}

type LoginFault struct {
	XMLName xml.Name `xml:"urn:fault.partner.soap.sforce.com LoginFault"`

	*ApiFault
}

type InvalidQueryLocatorFault struct {
	XMLName xml.Name `xml:"urn:fault.partner.soap.sforce.com InvalidQueryLocatorFault"`

	*ApiFault
}

type InvalidNewPasswordFault struct {
	XMLName xml.Name `xml:"urn:fault.partner.soap.sforce.com InvalidNewPasswordFault"`

	*ApiFault
}

type InvalidIdFault struct {
	XMLName xml.Name `xml:"urn:fault.partner.soap.sforce.com InvalidIdFault"`

	*ApiFault
}

type UnexpectedErrorFault struct {
	XMLName xml.Name `xml:"urn:fault.partner.soap.sforce.com UnexpectedErrorFault"`

	*ApiFault
}

type InvalidFieldFault struct {
	XMLName xml.Name `xml:"urn:fault.partner.soap.sforce.com InvalidFieldFault"`

	*ApiQueryFault
}

type InvalidSObjectFault struct {
	XMLName xml.Name `xml:"urn:fault.partner.soap.sforce.com InvalidSObjectFault"`

	*ApiQueryFault
}

type MalformedQueryFault struct {
	XMLName xml.Name `xml:"urn:fault.partner.soap.sforce.com MalformedQueryFault"`

	*ApiQueryFault
}

type MalformedSearchFault struct {
	XMLName xml.Name `xml:"urn:fault.partner.soap.sforce.com MalformedSearchFault"`

	*ApiQueryFault
}

type Soap struct {
	client *SOAPClient
	responseHeader *ResponseSOAPHeader
}

func NewSoap(url string, tls bool, auth *BasicAuth) *Soap {
	if url == "" {
		url = "https://login.salesforce.com/services/Soap/u/38.0"
	}
	client := NewSOAPClient(url, tls, auth)

	return &Soap{
		client: client,
		responseHeader: &ResponseSOAPHeader {
			info: &LimitInfoHeader{},
		},
	}
}

func NewSoapWithTLSConfig(url string, tlsCfg *tls.Config, auth *BasicAuth) *Soap {
	if url == "" {
		url = "https://login.salesforce.com/services/Soap/u/38.0"
	}
	client := NewSOAPClientWithTLSConfig(url, tlsCfg, auth)

	return &Soap{
		client: client,
	}
}

func (service *Soap) AddHeader(header interface{}) {
	service.client.AddHeader(header)
}

func (service *Soap) SetHeader(headers []interface{}) {
	service.ClearHeader()
	for _, header := range headers {
		service.AddHeader(header)
	}
}

func (service *Soap) ClearHeader() {
	service.client.ClearHeader()
}

func (service *Soap) SetDebug(debug bool) {
	service.client.SetDebug(debug)
}

func (service *Soap) SetLogger(logger io.Writer) {
	service.client.SetLogger(logger)
}

func (service *Soap) SetGzip(gz bool) {
	service.client.SetGzip(gz)
}

func (service *Soap) GetInfo() *LimitInfoHeader {
	return service.responseHeader.info
}

// Error can be either of the following types:
//
//   - LoginFault
//   - UnexpectedErrorFault
//   - InvalidIdFault
/* Login to the Salesforce.com SOAP Api */
func (service *Soap) Login(request *Login) (*LoginResponse, error) {
	response := new(LoginResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidSObjectFault
//   - UnexpectedErrorFault
/* Describe an sObject */
func (service *Soap) DescribeSObject(request *DescribeSObject) (*DescribeSObjectResponse, error) {
	response := new(DescribeSObjectResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidSObjectFault
//   - UnexpectedErrorFault
/* Describe multiple sObjects (upto 100) */
func (service *Soap) DescribeSObjects(request *DescribeSObjects) (*DescribeSObjectsResponse, error) {
	response := new(DescribeSObjectsResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
/* Describe the Global state */
func (service *Soap) DescribeGlobal(request *DescribeGlobal) (*DescribeGlobalResponse, error) {
	response := new(DescribeGlobalResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidSObjectFault
//   - UnexpectedErrorFault
/* Describe all the data category groups available for a given set of types */
func (service *Soap) DescribeDataCategoryGroups(request *DescribeDataCategoryGroups) (*DescribeDataCategoryGroupsResponse, error) {
	response := new(DescribeDataCategoryGroupsResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidSObjectFault
//   - UnexpectedErrorFault
/* Describe the data category group structures for a given set of pair of types and data category group name */
func (service *Soap) DescribeDataCategoryGroupStructures(request *DescribeDataCategoryGroupStructures) (*DescribeDataCategoryGroupStructuresResponse, error) {
	response := new(DescribeDataCategoryGroupStructuresResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
/* Describes your Knowledge settings, such as if knowledgeEnabled is on or off, its default language and supported languages */
func (service *Soap) DescribeKnowledgeSettings(request *DescribeKnowledgeSettings) (*DescribeKnowledgeSettingsResponse, error) {
	response := new(DescribeKnowledgeSettingsResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
//   - InvalidIdFault
/* Describe a list of FlexiPage and their contents */
func (service *Soap) DescribeFlexiPages(request *DescribeFlexiPages) (*DescribeFlexiPagesResponse, error) {
	response := new(DescribeFlexiPagesResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
/* Describe the items in an AppMenu */
func (service *Soap) DescribeAppMenu(request *DescribeAppMenu) (*DescribeAppMenuResponse, error) {
	response := new(DescribeAppMenuResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
/* Describe Gloal and Themes */
func (service *Soap) DescribeGlobalTheme(request *DescribeGlobalTheme) (*DescribeGlobalThemeResponse, error) {
	response := new(DescribeGlobalThemeResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
/* Describe Themes */
func (service *Soap) DescribeTheme(request *DescribeTheme) (*DescribeThemeResponse, error) {
	response := new(DescribeThemeResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidSObjectFault
//   - UnexpectedErrorFault
//   - InvalidIdFault
/* Describe the layout of the given sObject or the given actionable global page. */
func (service *Soap) DescribeLayout(request *DescribeLayout) (*DescribeLayoutResponse, error) {
	response := new(DescribeLayoutResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
/* Describe the layout of the SoftPhone */
func (service *Soap) DescribeSoftphoneLayout(request *DescribeSoftphoneLayout) (*DescribeSoftphoneLayoutResponse, error) {
	response := new(DescribeSoftphoneLayoutResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidSObjectFault
//   - UnexpectedErrorFault
/* Describe the search view of an sObject */
func (service *Soap) DescribeSearchLayouts(request *DescribeSearchLayouts) (*DescribeSearchLayoutsResponse, error) {
	response := new(DescribeSearchLayoutsResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/* Describe a list of entity names that reflects the current user's searchable entities */
func (service *Soap) DescribeSearchableEntities(request *DescribeSearchableEntities) (*DescribeSearchableEntitiesResponse, error) {
	response := new(DescribeSearchableEntitiesResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/* Describe a list of objects representing the order and scope of objects on a users search result page */
func (service *Soap) DescribeSearchScopeOrder(request *DescribeSearchScopeOrder) (*DescribeSearchScopeOrderResponse, error) {
	response := new(DescribeSearchScopeOrderResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/* Describe the compact layouts of the given sObject */
func (service *Soap) DescribeCompactLayouts(request *DescribeCompactLayouts) (*DescribeCompactLayoutsResponse, error) {
	response := new(DescribeCompactLayoutsResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/* Describe the Path Assistants for the given sObject and optionally RecordTypes */
func (service *Soap) DescribePathAssistants(request *DescribePathAssistants) (*DescribePathAssistantsResponse, error) {
	response := new(DescribePathAssistantsResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/* Describe the approval layouts of the given sObject */
func (service *Soap) DescribeApprovalLayout(request *DescribeApprovalLayout) (*DescribeApprovalLayoutResponse, error) {
	response := new(DescribeApprovalLayoutResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidSObjectFault
//   - UnexpectedErrorFault
/* Describe the ListViews as SOQL metadata for the generation of SOQL. */
func (service *Soap) DescribeSoqlListViews(request *DescribeSoqlListViews) (*DescribeSoqlListViewsResponse, error) {
	response := new(DescribeSoqlListViewsResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/* Execute the specified list view and return the presentation-ready results. */
func (service *Soap) ExecuteListView(request *ExecuteListView) (*ExecuteListViewResponse, error) {
	response := new(ExecuteListViewResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidSObjectFault
//   - UnexpectedErrorFault
/* Describe the ListViews of a SObject as SOQL metadata for the generation of SOQL. */
func (service *Soap) DescribeSObjectListViews(request *DescribeSObjectListViews) (*DescribeSObjectListViewsResponse, error) {
	response := new(DescribeSObjectListViewsResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
/* Describe the tabs that appear on a users page */
func (service *Soap) DescribeTabs(request *DescribeTabs) (*DescribeTabsResponse, error) {
	response := new(DescribeTabsResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
/* Describe all tabs available to a user */
func (service *Soap) DescribeAllTabs(request *DescribeAllTabs) (*DescribeAllTabsResponse, error) {
	response := new(DescribeAllTabsResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/* Describe the primary compact layouts for the sObjects requested */
func (service *Soap) DescribePrimaryCompactLayouts(request *DescribePrimaryCompactLayouts) (*DescribePrimaryCompactLayoutsResponse, error) {
	response := new(DescribePrimaryCompactLayoutsResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidSObjectFault
//   - UnexpectedErrorFault
//   - InvalidIdFault
//   - InvalidFieldFault
/* Create a set of new sObjects */
func (service *Soap) Create(request *Create) (*CreateResponse, error) {
	response := new(CreateResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidSObjectFault
//   - UnexpectedErrorFault
//   - InvalidIdFault
//   - InvalidFieldFault
/* Update a set of sObjects */
func (service *Soap) Update(request *Update) (*UpdateResponse, error) {
	response := new(UpdateResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidSObjectFault
//   - UnexpectedErrorFault
//   - InvalidIdFault
//   - InvalidFieldFault
/* Update or insert a set of sObjects based on object id */
func (service *Soap) Upsert(request *Upsert) (*UpsertResponse, error) {
	response := new(UpsertResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidSObjectFault
//   - UnexpectedErrorFault
//   - InvalidIdFault
//   - InvalidFieldFault
/* Merge and update a set of sObjects based on object id */
func (service *Soap) Merge(request *Merge) (*MergeResponse, error) {
	response := new(MergeResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
/* Delete a set of sObjects */
func (service *Soap) Delete(request *Delete) (*DeleteResponse, error) {
	response := new(DeleteResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
/* Undelete a set of sObjects */
func (service *Soap) Undelete(request *Undelete) (*UndeleteResponse, error) {
	response := new(UndeleteResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
/* Empty a set of sObjects from the recycle bin */
func (service *Soap) EmptyRecycleBin(request *EmptyRecycleBin) (*EmptyRecycleBinResponse, error) {
	response := new(EmptyRecycleBinResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidSObjectFault
//   - InvalidFieldFault
//   - MalformedQueryFault
//   - UnexpectedErrorFault
//   - InvalidIdFault
/* Get a set of sObjects */
func (service *Soap) Retrieve(request *Retrieve) (*RetrieveResponse, error) {
	response := new(RetrieveResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
//   - InvalidIdFault
/* Submit an entity to a workflow process or process a workitem */
func (service *Soap) Process(request *Process) (*ProcessResponse, error) {
	response := new(ProcessResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
/* convert a set of leads */
func (service *Soap) ConvertLead(request *ConvertLead) (*ConvertLeadResponse, error) {
	response := new(ConvertLeadResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
/* Logout the current user, invalidating the current session. */
func (service *Soap) Logout(request *Logout) (*LogoutResponse, error) {
	response := new(LogoutResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
/* Logs out and invalidates session ids */
func (service *Soap) InvalidateSessions(request *InvalidateSessions) (*InvalidateSessionsResponse, error) {
	response := new(InvalidateSessionsResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidSObjectFault
//   - UnexpectedErrorFault
/* Get the IDs for deleted sObjects */
func (service *Soap) GetDeleted(request *GetDeleted) (*GetDeletedResponse, error) {
	response := new(GetDeletedResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidSObjectFault
//   - UnexpectedErrorFault
/* Get the IDs for updated sObjects */
func (service *Soap) GetUpdated(request *GetUpdated) (*GetUpdatedResponse, error) {
	response := new(GetUpdatedResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidSObjectFault
//   - InvalidFieldFault
//   - MalformedQueryFault
//   - InvalidIdFault
//   - UnexpectedErrorFault
//   - InvalidQueryLocatorFault
/* Create a Query Cursor */
func (service *Soap) Query(request *Query) (*QueryResponse, error) {
	response := new(QueryResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidSObjectFault
//   - InvalidFieldFault
//   - MalformedQueryFault
//   - InvalidIdFault
//   - UnexpectedErrorFault
//   - InvalidQueryLocatorFault
/* Create a Query Cursor, including deleted sObjects */
func (service *Soap) QueryAll(request *QueryAll) (*QueryAllResponse, error) {
	response := new(QueryAllResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidQueryLocatorFault
//   - UnexpectedErrorFault
//   - InvalidFieldFault
//   - MalformedQueryFault
/* Gets the next batch of sObjects from a query */
func (service *Soap) QueryMore(request *QueryMore) (*QueryMoreResponse, error) {
	response := new(QueryMoreResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidSObjectFault
//   - InvalidFieldFault
//   - MalformedSearchFault
//   - UnexpectedErrorFault
/* Search for sObjects */
func (service *Soap) Search(request *Search) (*SearchResponse, error) {
	response := new(SearchResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
/* Gets server timestamp */
func (service *Soap) GetServerTimestamp(request *GetServerTimestamp) (*GetServerTimestampResponse, error) {
	response := new(GetServerTimestampResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidIdFault
//   - InvalidNewPasswordFault
//   - UnexpectedErrorFault
/* Set a user's password */
func (service *Soap) SetPassword(request *SetPassword) (*SetPasswordResponse, error) {
	response := new(SetPasswordResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidIdFault
//   - UnexpectedErrorFault
/* Reset a user's password */
func (service *Soap) ResetPassword(request *ResetPassword) (*ResetPasswordResponse, error) {
	response := new(ResetPasswordResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
/* Returns standard information relevant to the current user */
func (service *Soap) GetUserInfo(request *GetUserInfo) (*GetUserInfoResponse, error) {
	response := new(GetUserInfoResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
/* Send existing draft EmailMessage */
func (service *Soap) SendEmailMessage(request *SendEmailMessage) (*SendEmailMessageResponse, error) {
	response := new(SendEmailMessageResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
/* Send outbound email */
func (service *Soap) SendEmail(request *SendEmail) (*SendEmailResponse, error) {
	response := new(SendEmailResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - UnexpectedErrorFault
/* Perform a template merge on one or more blocks of text. */
func (service *Soap) RenderEmailTemplate(request *RenderEmailTemplate) (*RenderEmailTemplateResponse, error) {
	response := new(RenderEmailTemplateResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/* Perform a series of predefined actions such as quick create or log a task */
func (service *Soap) PerformQuickActions(request *PerformQuickActions) (*PerformQuickActionsResponse, error) {
	response := new(PerformQuickActionsResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/* Describe the details of a series of quick actions */
func (service *Soap) DescribeQuickActions(request *DescribeQuickActions) (*DescribeQuickActionsResponse, error) {
	response := new(DescribeQuickActionsResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/* Describe the details of a series of quick actions available for the given contextType */
func (service *Soap) DescribeAvailableQuickActions(request *DescribeAvailableQuickActions) (*DescribeAvailableQuickActionsResponse, error) {
	response := new(DescribeAvailableQuickActionsResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/* Retreive the template sobjects, if appropriate, for the given quick action names in a given context */
func (service *Soap) RetrieveQuickActionTemplates(request *RetrieveQuickActionTemplates) (*RetrieveQuickActionTemplatesResponse, error) {
	response := new(RetrieveQuickActionTemplatesResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/* Describe visualforce for an org */
func (service *Soap) DescribeVisualForce(request *DescribeVisualForce) (*DescribeVisualForceResponse, error) {
	response := new(DescribeVisualForceResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidSObjectFault
//   - UnexpectedErrorFault
//   - InvalidFieldFault
/* Find duplicates for a set of sObjects */
func (service *Soap) FindDuplicates(request *FindDuplicates) (*FindDuplicatesResponse, error) {
	response := new(FindDuplicatesResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/* Return the renameable nouns from the server for use in presentation using the salesforce grammar engine */
func (service *Soap) DescribeNouns(request *DescribeNouns) (*DescribeNounsResponse, error) {
	response := new(DescribeNounsResponse)
	err := service.client.Call(request, response, service.responseHeader)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *Soap) SetServerUrl(url string) {
	service.client.SetServerUrl(url)
}

func (service *Soap) GetServerUrl() string {
	return service.client.GetServerUrl()
}

var timeout = time.Duration(30 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Header  *SOAPHeader
	Body    SOAPBody
}

type SOAPHeader struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`

	Items []interface{} `xml:",omitempty"`
	response *ResponseSOAPHeader
}

type ResponseSOAPHeader struct {
	info *LimitInfoHeader
}

type SOAPBody struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`

	Fault   *SOAPFault  `xml:",omitempty"`
	Content interface{} `xml:",omitempty"`
}

type SOAPFault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`

	Code   string `xml:"faultcode,omitempty"`
	String string `xml:"faultstring,omitempty"`
	Actor  string `xml:"faultactor,omitempty"`
	Detail string `xml:"detail,omitempty"`
}

const (
	// Predefined WSS namespaces to be used in
	WssNsWSSE string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd"
	WssNsWSU  string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd"
	WssNsType string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText"
)

type WSSSecurityHeader struct {
	XMLName   xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ wsse:Security"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`

	MustUnderstand string `xml:"mustUnderstand,attr,omitempty"`

	Token *WSSUsernameToken `xml:",omitempty"`
}

type WSSUsernameToken struct {
	XMLName   xml.Name `xml:"wsse:UsernameToken"`
	XmlNSWsu  string   `xml:"xmlns:wsu,attr"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`

	Id string `xml:"wsu:Id,attr,omitempty"`

	Username *WSSUsername `xml:",omitempty"`
	Password *WSSPassword `xml:",omitempty"`
}

type WSSUsername struct {
	XMLName   xml.Name `xml:"wsse:Username"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`

	Data string `xml:",chardata"`
}

type WSSPassword struct {
	XMLName   xml.Name `xml:"wsse:Password"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`
	XmlNSType string   `xml:"Type,attr"`

	Data string `xml:",chardata"`
}

type BasicAuth struct {
	Login    string
	Password string
}

type SOAPClient struct {
	url     string
	tlsCfg  *tls.Config
	auth    *BasicAuth
	headers []interface{}
	logger  io.Writer
	debug   bool
	gzip    bool
}

// **********
// Accepted solution from http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
// Author: Icza - http://stackoverflow.com/users/1705598/icza

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStringBytesMaskImprSrc(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// **********

func NewWSSSecurityHeader(user, pass, mustUnderstand string) *WSSSecurityHeader {
	hdr := &WSSSecurityHeader{XmlNSWsse: WssNsWSSE, MustUnderstand: mustUnderstand}
	hdr.Token = &WSSUsernameToken{XmlNSWsu: WssNsWSU, XmlNSWsse: WssNsWSSE, Id: "UsernameToken-" + randStringBytesMaskImprSrc(9)}
	hdr.Token.Username = &WSSUsername{XmlNSWsse: WssNsWSSE, Data: user}
	hdr.Token.Password = &WSSPassword{XmlNSWsse: WssNsWSSE, XmlNSType: WssNsType, Data: pass}
	return hdr
}

func (b *SOAPHeader) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var (
		token    xml.Token
		err      error
		consumed bool
	)

Loop:
	for {
		if token, err = d.Token(); err != nil {
			return err
		}

		if token == nil {
			break
		}

		switch se := token.(type) {
		case xml.StartElement:
			if consumed {
				return xml.UnmarshalError("Found multiple elements inside SOAP body; not wrapped-document/literal WS-I compliant")
			} else {
				if se.Name.Local == "LimitInfoHeader" && b.response.info != nil {
					if err = d.DecodeElement(b.response.info, &se); err != nil {
						return err
					}
				}

				consumed = true
			}
		case xml.EndElement:
			break Loop
		}
	}
	return nil
}

func (b *SOAPBody) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if b.Content == nil {
		return xml.UnmarshalError("Content must be a pointer to a struct")
	}

	var (
		token    xml.Token
		err      error
		consumed bool
	)

Loop:
	for {
		if token, err = d.Token(); err != nil {
			return err
		}

		if token == nil {
			break
		}

		switch se := token.(type) {
		case xml.StartElement:
			if consumed {
				return xml.UnmarshalError("Found multiple elements inside SOAP body; not wrapped-document/literal WS-I compliant")
			} else if se.Name.Space == "http://schemas.xmlsoap.org/soap/envelope/" && se.Name.Local == "Fault" {
				b.Fault = &SOAPFault{}
				b.Content = nil

				err = d.DecodeElement(b.Fault, &se)
				if err != nil {
					return err
				}

				consumed = true
			} else {
				if err = d.DecodeElement(b.Content, &se); err != nil {
					return err
				}

				consumed = true
			}
		case xml.EndElement:
			break Loop
		}
	}

	return nil
}

func (f *SOAPFault) Error() string {
	return f.String
}

func NewSOAPClient(url string, insecureSkipVerify bool, auth *BasicAuth) *SOAPClient {
	tlsCfg := &tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
	}
	return NewSOAPClientWithTLSConfig(url, tlsCfg, auth)
}

func NewSOAPClientWithTLSConfig(url string, tlsCfg *tls.Config, auth *BasicAuth) *SOAPClient {
	return &SOAPClient{
		url:    url,
		tlsCfg: tlsCfg,
		auth:   auth,
		logger: os.Stdout,
		debug:  false,
		gzip:   true,
	}
}

func (s *SOAPClient) SetDebug(debug bool) {
	s.debug = debug
}

func (s *SOAPClient) SetLogger(logger io.Writer) {
	s.logger = logger
}

func (s *SOAPClient) SetGzip(gz bool) {
	s.gzip = gz
}

func (s *SOAPClient) AddHeader(header interface{}) {
	s.headers = append(s.headers, header)
}

func (s *SOAPClient) ClearHeader() {
	s.headers = nil
}

func (s *SOAPClient) Call(request, response interface{}, responseHeader *ResponseSOAPHeader) error {
	envelope := SOAPEnvelope{}

	if s.headers != nil && len(s.headers) > 0 {
		soapHeader := &SOAPHeader{Items: make([]interface{}, len(s.headers))}
		copy(soapHeader.Items, s.headers)
		envelope.Header = soapHeader
	}

	envelope.Body.Content = request
	buffer := new(bytes.Buffer)

	encoder := xml.NewEncoder(buffer)
	//encoder.Indent("  ", "    ")

	if err := encoder.Encode(envelope); err != nil {
		return err
	}

	if err := encoder.Flush(); err != nil {
		return err
	}

	if s.debug {
		s.logger.Write(buffer.Bytes())
		s.logger.Write([]byte("\n"))
	}

	req, err := s.createRequest(buffer)
	if err != nil {
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: s.tlsCfg,
		Dial:            dialTimeout,
	}

	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	rawbody, err := getRawBody(res)
	if err != nil {
		return err
	}

	if s.debug {
		s.logger.Write(rawbody)
		s.logger.Write([]byte("\n"))
	}

	respEnvelope := new(SOAPEnvelope)
	responseHeader.info = &LimitInfoHeader{}
	header := SOAPHeader{response: responseHeader}
	respEnvelope.Header = &header
	respEnvelope.Body = SOAPBody{Content: response}
	err = xml.Unmarshal(rawbody, respEnvelope)
	if err != nil {
		return err
	}

	fault := respEnvelope.Body.Fault
	if fault != nil {
		return fault
	}

	return nil
}

func (s *SOAPClient) SetServerUrl(url string) {
	s.url = url
}

func (s *SOAPClient) GetServerUrl() string {
	return s.url
}

func (s *SOAPClient) createRequest(buffer *bytes.Buffer) (*http.Request, error) {
	var req *http.Request
	var err error
	if s.gzip {
		gzipBuffer := new(bytes.Buffer)
		gw := gzip.NewWriter(gzipBuffer)
		_, err = gw.Write(buffer.Bytes())
		if err != nil {
			return nil, err
		}
		err := gw.Close()
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest("POST", s.url, gzipBuffer)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Content-Encoding", "gzip")
		req.Header.Add("Accept-Encoding", "gzip")
	} else {
		req, err = http.NewRequest("POST", s.url, buffer)
		if err != nil {
			return nil, err
		}
	}

	if s.auth != nil {
		req.SetBasicAuth(s.auth.Login, s.auth.Password)
	}

	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	req.Header.Add("SOAPAction", "''")

	req.Header.Set("User-Agent", "gowsdl/0.1")
	req.Close = true
	return req, nil
}

func getRawBody(res *http.Response) ([]byte, error) {
	if res.Header.Get("Content-Encoding") == "gzip" {
		buf := new(bytes.Buffer)
		gr, err := gzip.NewReader(res.Body)
		if err != nil {
			return nil, err
		}
		_, err = buf.ReadFrom(gr)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	return ioutil.ReadAll(res.Body)
}
