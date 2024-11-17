# \OrganizationsApi

All URIs are relative to *https://api.zerotier.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**accept_invitation**](OrganizationsApi.md#accept_invitation) | **POST** /org-invitation/{inviteID} | Accept organization invitation
[**decline_invitation**](OrganizationsApi.md#decline_invitation) | **DELETE** /org-invitation/{inviteID} | Decline organization invitation
[**get_invitation_by_id**](OrganizationsApi.md#get_invitation_by_id) | **GET** /org-invitation/{inviteID} | Get organization invitation
[**get_organization**](OrganizationsApi.md#get_organization) | **GET** /org | Get the current user's organization
[**get_organization_by_id**](OrganizationsApi.md#get_organization_by_id) | **GET** /org/{orgID} | Get organization by ID
[**get_organization_invitation_list**](OrganizationsApi.md#get_organization_invitation_list) | **GET** /org-invitation | Get list of organization invitations
[**get_organization_members**](OrganizationsApi.md#get_organization_members) | **GET** /org/{orgID}/user | Get list of organization members
[**invite_user_by_email**](OrganizationsApi.md#invite_user_by_email) | **POST** /org-invitation | Invite a user to your organization by email



## accept_invitation

> models::OrganizationInvitation accept_invitation(invite_id)
Accept organization invitation

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**invite_id** | **String** | Invitation ID | [required] |

### Return type

[**models::OrganizationInvitation**](OrganizationInvitation.md)

### Authorization

[tokenAuth](../README.md#tokenAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## decline_invitation

> decline_invitation(invite_id)
Decline organization invitation

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**invite_id** | **String** | Invitation ID | [required] |

### Return type

 (empty response body)

### Authorization

[tokenAuth](../README.md#tokenAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## get_invitation_by_id

> models::OrganizationInvitation get_invitation_by_id(invite_id)
Get organization invitation

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**invite_id** | **String** | Invitation ID | [required] |

### Return type

[**models::OrganizationInvitation**](OrganizationInvitation.md)

### Authorization

[tokenAuth](../README.md#tokenAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## get_organization

> models::Organization get_organization()
Get the current user's organization

### Parameters

This endpoint does not need any parameter.

### Return type

[**models::Organization**](Organization.md)

### Authorization

[tokenAuth](../README.md#tokenAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## get_organization_by_id

> models::Organization get_organization_by_id(org_id)
Get organization by ID

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**org_id** | **String** | Organization ID | [required] |

### Return type

[**models::Organization**](Organization.md)

### Authorization

[tokenAuth](../README.md#tokenAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## get_organization_invitation_list

> Vec<models::OrganizationInvitation> get_organization_invitation_list()
Get list of organization invitations

### Parameters

This endpoint does not need any parameter.

### Return type

[**Vec<models::OrganizationInvitation>**](OrganizationInvitation.md)

### Authorization

[tokenAuth](../README.md#tokenAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## get_organization_members

> models::OrganizationMember get_organization_members(org_id)
Get list of organization members

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**org_id** | **String** | Organization ID | [required] |

### Return type

[**models::OrganizationMember**](OrganizationMember.md)

### Authorization

[tokenAuth](../README.md#tokenAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## invite_user_by_email

> models::OrganizationInvitation invite_user_by_email(organization_invitation)
Invite a user to your organization by email

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**organization_invitation** | [**OrganizationInvitation**](OrganizationInvitation.md) | Organization Invitation JSON object | [required] |

### Return type

[**models::OrganizationInvitation**](OrganizationInvitation.md)

### Authorization

[tokenAuth](../README.md#tokenAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

