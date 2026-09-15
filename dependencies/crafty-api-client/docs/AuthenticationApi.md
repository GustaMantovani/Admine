# AuthenticationApi

All URIs are relative to *https://api.example.com*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**apiV2AuthInvalidateTokensPost**](AuthenticationApi.md#apiv2authinvalidatetokenspost) | **POST** /api/v2/auth/invalidate_tokens | Sign out all user\&#39;s sessions |
| [**apiV2AuthLoginPost**](AuthenticationApi.md#apiv2authloginpostoperation) | **POST** /api/v2/auth/login | Log in and get a token |



## apiV2AuthInvalidateTokensPost

> StatusOK apiV2AuthInvalidateTokensPost()

Sign out all user\&#39;s sessions

Revokes all sessions of the signed in user. Requires valid login session or bearer token.

### Example

```ts
import {
  Configuration,
  AuthenticationApi,
} from '';
import type { ApiV2AuthInvalidateTokensPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new AuthenticationApi(config);

  try {
    const data = await api.apiV2AuthInvalidateTokensPost();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**StatusOK**](StatusOK.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully invalidated all user sessions. |  -  |
| **403** | Invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2AuthLoginPost

> ApiV2AuthLoginPost200Response apiV2AuthLoginPost(apiV2AuthLoginPostRequest)

Log in and get a token

Login with username and password and get an API token.

### Example

```ts
import {
  Configuration,
  AuthenticationApi,
} from '';
import type { ApiV2AuthLoginPostOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new AuthenticationApi();

  const body = {
    // ApiV2AuthLoginPostRequest
    apiV2AuthLoginPostRequest: ...,
  } satisfies ApiV2AuthLoginPostOperationRequest;

  try {
    const data = await api.apiV2AuthLoginPost(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **apiV2AuthLoginPostRequest** | [ApiV2AuthLoginPostRequest](ApiV2AuthLoginPostRequest.md) |  | |

### Return type

[**ApiV2AuthLoginPost200Response**](ApiV2AuthLoginPost200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successful authentication |  -  |
| **400** | Invalid credentials for format. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

