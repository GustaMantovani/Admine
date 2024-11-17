# \UserApi

All URIs are relative to *https://api.zerotier.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**add_api_token**](UserApi.md#add_api_token) | **POST** /user/{userID}/token | Add an API token
[**delete_api_token**](UserApi.md#delete_api_token) | **DELETE** /user/{userID}/token/{tokenName} | Delete API Token
[**delete_user_by_id**](UserApi.md#delete_user_by_id) | **DELETE** /user/{userID} | Delete user
[**get_user_by_id**](UserApi.md#get_user_by_id) | **GET** /user/{userID} | Get user record
[**update_user_by_id**](UserApi.md#update_user_by_id) | **POST** /user/{userID} | Update user record (SMS number or Display Name only)



## add_api_token

> models::ApiToken add_api_token(user_id, api_token)
Add an API token

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**user_id** | **String** | User ID | [required] |
**api_token** | [**ApiToken**](ApiToken.md) | APIToken JSON object | [required] |

### Return type

[**models::ApiToken**](APIToken.md)

### Authorization

[tokenAuth](../README.md#tokenAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## delete_api_token

> delete_api_token(user_id, token_name)
Delete API Token

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**user_id** | **String** | User ID | [required] |
**token_name** | **String** | Token Name | [required] |

### Return type

 (empty response body)

### Authorization

[tokenAuth](../README.md#tokenAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## delete_user_by_id

> delete_user_by_id(user_id)
Delete user

Deletes the user and all associated networks.  This is not reversible. Delete at your own risk.

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**user_id** | **String** | User ID | [required] |

### Return type

 (empty response body)

### Authorization

[tokenAuth](../README.md#tokenAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## get_user_by_id

> models::User get_user_by_id(user_id)
Get user record

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**user_id** | **String** | User ID | [required] |

### Return type

[**models::User**](User.md)

### Authorization

[tokenAuth](../README.md#tokenAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## update_user_by_id

> models::User update_user_by_id(user_id, user)
Update user record (SMS number or Display Name only)

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**user_id** | **String** | User ID | [required] |
**user** | [**User**](User.md) | User object JSON | [required] |

### Return type

[**models::User**](User.md)

### Authorization

[tokenAuth](../README.md#tokenAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

