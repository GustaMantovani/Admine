# NetworkSsoConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**enabled** | Option<**bool**> | SSO enabled/disabled on network | [optional]
**mode** | Option<**String**> | SSO mode.  One of: `default`, `email`, `group` | [optional]
**client_id** | Option<**String**> | SSO client ID.  Client ID must be already configured in the Org | [optional]
**issuer** | Option<**String**> | URL of the OIDC issuer | [optional][readonly]
**provider** | Option<**String**> | Provider type | [optional][readonly]
**authorization_endpoint** | Option<**String**> | Authorization URL endpoint | [optional][readonly]
**allow_list** | Option<**Vec<String>**> | List of email addresses or group memberships that may SSO auth onto the network | [optional]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


