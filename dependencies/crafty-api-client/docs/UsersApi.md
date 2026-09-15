# UsersApi

All URIs are relative to *https://api.example.com*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**apiV2UsersGet**](UsersApi.md#apiv2usersget) | **GET** /api/v2/users | Get all the users |
| [**apiV2UsersPost**](UsersApi.md#apiv2userspost) | **POST** /api/v2/users | Create a user |
| [**apiV2UsersUserIDDelete**](UsersApi.md#apiv2usersuseriddelete) | **DELETE** /api/v2/users/{userID} | Delete a user |
| [**apiV2UsersUserIDGet**](UsersApi.md#apiv2usersuseridget) | **GET** /api/v2/users/{userID} | Get a user |
| [**apiV2UsersUserIDPatch**](UsersApi.md#apiv2usersuseridpatchoperation) | **PATCH** /api/v2/users/{userID} | Modify a user |
| [**apiV2UsersUserIDPermissionsGet**](UsersApi.md#apiv2usersuseridpermissionsget) | **GET** /api/v2/users/{userID}/permissions | Get a user\&#39;s Crafty permissions. |
| [**apiV2UsersUserIDPfpGet**](UsersApi.md#apiv2usersuseridpfpget) | **GET** /api/v2/users/{userID}/pfp | Get a user\&#39;s profile picture |
| [**apiV2UsersUserIDPublicGet**](UsersApi.md#apiv2usersuseridpublicget) | **GET** /api/v2/users/{userID}/public | Get a user\&#39;s public data |



## apiV2UsersGet

> ApiV2UsersGet200Response apiV2UsersGet(ids)

Get all the users

Get all the users. Requires user config permission.

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { ApiV2UsersGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new UsersApi(config);

  const body = {
    // boolean | If true, only returns IDs, no other information. (optional)
    ids: false,
  } satisfies ApiV2UsersGetRequest;

  try {
    const data = await api.apiV2UsersGet(body);
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
| **ids** | `boolean` | If true, only returns IDs, no other information. | [Optional] [Defaults to `undefined`] |

### Return type

[**ApiV2UsersGet200Response**](ApiV2UsersGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully got all users. |  -  |
| **403** | Failed authentication or invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2UsersPost

> ApiV2UsersPost200Response apiV2UsersPost(user)

Create a user

Create a user. Requires user config permission.

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { ApiV2UsersPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new UsersApi(config);

  const body = {
    // User (optional)
    user: ...,
  } satisfies ApiV2UsersPostRequest;

  try {
    const data = await api.apiV2UsersPost(body);
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
| **user** | [User](User.md) |  | [Optional] |

### Return type

[**ApiV2UsersPost200Response**](ApiV2UsersPost200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully created user. |  -  |
| **403** | Failed authentication or invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2UsersUserIDDelete

> StatusOK apiV2UsersUserIDDelete(userID)

Delete a user

Get a user. Requires user config permission if not viewing self.

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { ApiV2UsersUserIDDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new UsersApi(config);

  const body = {
    // number | Target User ID
    userID: 1,
  } satisfies ApiV2UsersUserIDDeleteRequest;

  try {
    const data = await api.apiV2UsersUserIDDelete(body);
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
| **userID** | `number` | Target User ID | [Defaults to `undefined`] |

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
| **200** | Successfully deleted user |  -  |
| **403** | Failed authentication or invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2UsersUserIDGet

> ApiV2UsersUserIDGet200Response apiV2UsersUserIDGet(userID)

Get a user

Get a user. Requires user config permission if not viewing self.

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { ApiV2UsersUserIDGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new UsersApi(config);

  const body = {
    // number | Target User ID
    userID: 1,
  } satisfies ApiV2UsersUserIDGetRequest;

  try {
    const data = await api.apiV2UsersUserIDGet(body);
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
| **userID** | `number` | Target User ID | [Defaults to `undefined`] |

### Return type

[**ApiV2UsersUserIDGet200Response**](ApiV2UsersUserIDGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully got user. |  -  |
| **403** | Failed authentication or invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2UsersUserIDPatch

> StatusOK apiV2UsersUserIDPatch(userID, apiV2UsersUserIDPatchRequest)

Modify a user

Get a user. Requires user config permission if not viewing self.

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { ApiV2UsersUserIDPatchOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new UsersApi(config);

  const body = {
    // number | Target User ID
    userID: 1,
    // ApiV2UsersUserIDPatchRequest (optional)
    apiV2UsersUserIDPatchRequest: ...,
  } satisfies ApiV2UsersUserIDPatchOperationRequest;

  try {
    const data = await api.apiV2UsersUserIDPatch(body);
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
| **userID** | `number` | Target User ID | [Defaults to `undefined`] |
| **apiV2UsersUserIDPatchRequest** | [ApiV2UsersUserIDPatchRequest](ApiV2UsersUserIDPatchRequest.md) |  | [Optional] |

### Return type

[**StatusOK**](StatusOK.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully updated user. |  -  |
| **403** | Failed authentication or invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2UsersUserIDPermissionsGet

> ApiV2UsersUserIDPermissionsGet200Response apiV2UsersUserIDPermissionsGet(userID)

Get a user\&#39;s Crafty permissions.

Get a user\&#39;s Crafty permissions. Requires user config permission if not viewing self.

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { ApiV2UsersUserIDPermissionsGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new UsersApi(config);

  const body = {
    // number | Target User ID
    userID: 1,
  } satisfies ApiV2UsersUserIDPermissionsGetRequest;

  try {
    const data = await api.apiV2UsersUserIDPermissionsGet(body);
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
| **userID** | `number` | Target User ID | [Defaults to `undefined`] |

### Return type

[**ApiV2UsersUserIDPermissionsGet200Response**](ApiV2UsersUserIDPermissionsGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully got a user\&#39;s permissions. |  -  |
| **403** | Failed authentication or invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2UsersUserIDPfpGet

> ApiV2UsersUserIDPfpGet200Response apiV2UsersUserIDPfpGet(userID)

Get a user\&#39;s profile picture

Get a user\&#39;s profile picture. Requires logged in.

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { ApiV2UsersUserIDPfpGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new UsersApi(config);

  const body = {
    // number | Target User ID
    userID: 1,
  } satisfies ApiV2UsersUserIDPfpGetRequest;

  try {
    const data = await api.apiV2UsersUserIDPfpGet(body);
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
| **userID** | `number` | Target User ID | [Defaults to `undefined`] |

### Return type

[**ApiV2UsersUserIDPfpGet200Response**](ApiV2UsersUserIDPfpGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully got user\&#39;s profile picture. |  -  |
| **403** | Failed authentication or invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2UsersUserIDPublicGet

> ApiV2UsersUserIDPublicGet200Response apiV2UsersUserIDPublicGet(userID)

Get a user\&#39;s public data

Get a user\&#39;s public data. Requires logged in.

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { ApiV2UsersUserIDPublicGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new UsersApi(config);

  const body = {
    // number | Target User ID
    userID: 1,
  } satisfies ApiV2UsersUserIDPublicGetRequest;

  try {
    const data = await api.apiV2UsersUserIDPublicGet(body);
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
| **userID** | `number` | Target User ID | [Defaults to `undefined`] |

### Return type

[**ApiV2UsersUserIDPublicGet200Response**](ApiV2UsersUserIDPublicGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully got user\&#39;s public data. |  -  |
| **403** | Failed authentication or invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

