# RolesApi

All URIs are relative to *https://api.example.com*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**apiV2RolesGet**](RolesApi.md#apiv2rolesget) | **GET** /api/v2/roles | Get all Roles |
| [**apiV2RolesPost**](RolesApi.md#apiv2rolespostoperation) | **POST** /api/v2/roles | Create a role |
| [**apiV2RolesRoleIDDelete**](RolesApi.md#apiv2rolesroleiddelete) | **DELETE** /api/v2/roles/{roleID} | Delete a role. |
| [**apiV2RolesRoleIDGet**](RolesApi.md#apiv2rolesroleidget) | **GET** /api/v2/roles/{roleID} | Get a role |
| [**apiV2RolesRoleIDPatch**](RolesApi.md#apiv2rolesroleidpatchoperation) | **PATCH** /api/v2/roles/{roleID} | Modify a role. |
| [**apiV2RolesRoleIDServersGet**](RolesApi.md#apiv2rolesroleidserversget) | **GET** /api/v2/roles/{roleID}/servers | Get a Role\&#39;s servers |
| [**apiV2RolesRoleIDUsersGet**](RolesApi.md#apiv2rolesroleidusersget) | **GET** /api/v2/roles/{roleID}/users | Get a role\&#39;s user |



## apiV2RolesGet

> ApiV2RolesGet200Response apiV2RolesGet(ids)

Get all Roles

Get all roles. Requires superuser permission.

### Example

```ts
import {
  Configuration,
  RolesApi,
} from '';
import type { ApiV2RolesGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new RolesApi(config);

  const body = {
    // boolean | Returns only role IDs, not additional role information (optional)
    ids: true,
  } satisfies ApiV2RolesGetRequest;

  try {
    const data = await api.apiV2RolesGet(body);
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
| **ids** | `boolean` | Returns only role IDs, not additional role information | [Optional] [Defaults to `undefined`] |

### Return type

[**ApiV2RolesGet200Response**](ApiV2RolesGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Example response body for /api/v2/roles |  -  |
| **403** | Access denied response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2RolesPost

> ApiV2RolesPost200Response apiV2RolesPost(apiV2RolesPostRequest)

Create a role

Requires superuser permission

### Example

```ts
import {
  Configuration,
  RolesApi,
} from '';
import type { ApiV2RolesPostOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new RolesApi(config);

  const body = {
    // ApiV2RolesPostRequest
    apiV2RolesPostRequest: ...,
  } satisfies ApiV2RolesPostOperationRequest;

  try {
    const data = await api.apiV2RolesPost(body);
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
| **apiV2RolesPostRequest** | [ApiV2RolesPostRequest](ApiV2RolesPostRequest.md) |  | |

### Return type

[**ApiV2RolesPost200Response**](ApiV2RolesPost200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully created the role |  -  |
| **403** | Access denied response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2RolesRoleIDDelete

> ApiV2RolesRoleIDDelete200Response apiV2RolesRoleIDDelete(roleID)

Delete a role.

Deletes the specified role. Requires superuser permission.

### Example

```ts
import {
  Configuration,
  RolesApi,
} from '';
import type { ApiV2RolesRoleIDDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new RolesApi(config);

  const body = {
    // number | ID of role to delete.
    roleID: 1,
  } satisfies ApiV2RolesRoleIDDeleteRequest;

  try {
    const data = await api.apiV2RolesRoleIDDelete(body);
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
| **roleID** | `number` | ID of role to delete. | [Defaults to `undefined`] |

### Return type

[**ApiV2RolesRoleIDDelete200Response**](ApiV2RolesRoleIDDelete200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully deleted role. |  -  |
| **403** | Invalid user or permissions for action. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2RolesRoleIDGet

> ApiV2RolesRoleIDGet200Response apiV2RolesRoleIDGet(roleID)

Get a role

Get information about a role. Requires superuser permission.

### Example

```ts
import {
  Configuration,
  RolesApi,
} from '';
import type { ApiV2RolesRoleIDGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new RolesApi(config);

  const body = {
    // number | Role ID of the role to query
    roleID: 56,
  } satisfies ApiV2RolesRoleIDGetRequest;

  try {
    const data = await api.apiV2RolesRoleIDGet(body);
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
| **roleID** | `number` | Role ID of the role to query | [Defaults to `undefined`] |

### Return type

[**ApiV2RolesRoleIDGet200Response**](ApiV2RolesRoleIDGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully got a role. |  -  |
| **403** | Invalid user or permissions for action. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2RolesRoleIDPatch

> StatusOK apiV2RolesRoleIDPatch(roleID, apiV2RolesRoleIDPatchRequest)

Modify a role.

Modify the specified role. Requires superuser permission.

### Example

```ts
import {
  Configuration,
  RolesApi,
} from '';
import type { ApiV2RolesRoleIDPatchOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new RolesApi(config);

  const body = {
    // number | ID of role to update.
    roleID: 1,
    // ApiV2RolesRoleIDPatchRequest
    apiV2RolesRoleIDPatchRequest: ...,
  } satisfies ApiV2RolesRoleIDPatchOperationRequest;

  try {
    const data = await api.apiV2RolesRoleIDPatch(body);
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
| **roleID** | `number` | ID of role to update. | [Defaults to `undefined`] |
| **apiV2RolesRoleIDPatchRequest** | [ApiV2RolesRoleIDPatchRequest](ApiV2RolesRoleIDPatchRequest.md) |  | |

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
| **200** | Successfully updated role. |  -  |
| **403** | Invalid user or permissions for action. |  -  |
| **400** | Invalid request. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2RolesRoleIDServersGet

> ApiV2RolesRoleIDServersGet200Response apiV2RolesRoleIDServersGet(roleID, ids)

Get a Role\&#39;s servers

Get all the servers that a role has access to. Requires superuser permissions

### Example

```ts
import {
  Configuration,
  RolesApi,
} from '';
import type { ApiV2RolesRoleIDServersGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new RolesApi(config);

  const body = {
    // number | Role ID of the role to query
    roleID: 56,
    // boolean (optional)
    ids: true,
  } satisfies ApiV2RolesRoleIDServersGetRequest;

  try {
    const data = await api.apiV2RolesRoleIDServersGet(body);
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
| **roleID** | `number` | Role ID of the role to query | [Defaults to `undefined`] |
| **ids** | `boolean` |  | [Optional] [Defaults to `undefined`] |

### Return type

[**ApiV2RolesRoleIDServersGet200Response**](ApiV2RolesRoleIDServersGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully got a role\&#39;s servers |  -  |
| **403** | Invalid user or permissions for action. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2RolesRoleIDUsersGet

> ApiV2CraftyConfigGet200Response apiV2RolesRoleIDUsersGet(roleID)

Get a role\&#39;s user

Get all the users of a specified role. Requires superuser.

### Example

```ts
import {
  Configuration,
  RolesApi,
} from '';
import type { ApiV2RolesRoleIDUsersGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new RolesApi(config);

  const body = {
    // number
    roleID: 1,
  } satisfies ApiV2RolesRoleIDUsersGetRequest;

  try {
    const data = await api.apiV2RolesRoleIDUsersGet(body);
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
| **roleID** | `number` |  | [Defaults to `undefined`] |

### Return type

[**ApiV2CraftyConfigGet200Response**](ApiV2CraftyConfigGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully got a role\&#39;s users |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

