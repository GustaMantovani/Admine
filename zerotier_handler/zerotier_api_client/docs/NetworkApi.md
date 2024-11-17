# \NetworkApi

All URIs are relative to *https://api.zerotier.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**delete_network**](NetworkApi.md#delete_network) | **DELETE** /network/{networkID} | delete network
[**get_network_by_id**](NetworkApi.md#get_network_by_id) | **GET** /network/{networkID} | Get network by ID
[**get_network_list**](NetworkApi.md#get_network_list) | **GET** /network | Returns a list of Networks you have access to.
[**new_network**](NetworkApi.md#new_network) | **POST** /network | Create a new network.
[**update_network**](NetworkApi.md#update_network) | **POST** /network/{networkID} | update network configuration



## delete_network

> delete_network(network_id)
delete network

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**network_id** | **String** | ID of the network | [required] |

### Return type

 (empty response body)

### Authorization

[tokenAuth](../README.md#tokenAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## get_network_by_id

> models::Network get_network_by_id(network_id)
Get network by ID

Returns a single network

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**network_id** | **String** | ID of the network to return | [required] |

### Return type

[**models::Network**](Network.md)

### Authorization

[tokenAuth](../README.md#tokenAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## get_network_list

> Vec<models::Network> get_network_list()
Returns a list of Networks you have access to.

### Parameters

This endpoint does not need any parameter.

### Return type

[**Vec<models::Network>**](Network.md)

### Authorization

[tokenAuth](../README.md#tokenAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## new_network

> models::Network new_network(body)
Create a new network.

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**body** | **serde_json::Value** | empty JSON object | [required] |

### Return type

[**models::Network**](Network.md)

### Authorization

[tokenAuth](../README.md#tokenAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


## update_network

> models::Network update_network(network_id, network)
update network configuration

### Parameters


Name | Type | Description  | Required | Notes
------------- | ------------- | ------------- | ------------- | -------------
**network_id** | **String** | ID of the network to change | [required] |
**network** | [**Network**](Network.md) | Network object JSON | [required] |

### Return type

[**models::Network**](Network.md)

### Authorization

[tokenAuth](../README.md#tokenAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

