# ServersApi

All URIs are relative to *https://api.example.com*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**apiV2ServerServerIDWebhookWebhookIDDelete**](ServersApi.md#apiv2serverserveridwebhookwebhookiddelete) | **DELETE** /api/v2/server/{serverID}/webhook/{webhookID} | Delete a webhook |
| [**apiV2ServerServerIDWebhookWebhookIDGet**](ServersApi.md#apiv2serverserveridwebhookwebhookidget) | **GET** /api/v2/server/{serverID}/webhook/{webhookID} | Get a webhook |
| [**apiV2ServerServerIDWebhookWebhookIDPatch**](ServersApi.md#apiv2serverserveridwebhookwebhookidpatchoperation) | **PATCH** /api/v2/server/{serverID}/webhook/{webhookID} | Modify a Server Webhook |
| [**apiV2ServersGet**](ServersApi.md#apiv2serversget) | **GET** /api/v2/servers | Get all servers you have access to. |
| [**apiV2ServersPost**](ServersApi.md#apiv2serverspostoperation) | **POST** /api/v2/servers | Create a server |
| [**apiV2ServersServerIDActionActionPost**](ServersApi.md#apiv2serversserveridactionactionpost) | **POST** /api/v2/servers/{serverID}/action/{action} | Send an action to a server |
| [**apiV2ServersServerIDDelete**](ServersApi.md#apiv2serversserveriddelete) | **DELETE** /api/v2/servers/{serverID} | Delete a server |
| [**apiV2ServersServerIDGet**](ServersApi.md#apiv2serversserveridget) | **GET** /api/v2/servers/{serverID} | Get a server\&#39;s data |
| [**apiV2ServersServerIDLogsGet**](ServersApi.md#apiv2serversserveridlogsget) | **GET** /api/v2/servers/{serverID}/logs | Get a server\&#39;s logs |
| [**apiV2ServersServerIDPatch**](ServersApi.md#apiv2serversserveridpatchoperation) | **PATCH** /api/v2/servers/{serverID} | Modify a server |
| [**apiV2ServersServerIDPublicGet**](ServersApi.md#apiv2serversserveridpublicget) | **GET** /api/v2/servers/{serverID}/public | Get a server\&#39;s public data |
| [**apiV2ServersServerIDStatsGet**](ServersApi.md#apiv2serversserveridstatsget) | **GET** /api/v2/servers/{serverID}/stats | Get a server\&#39;s statistics |
| [**apiV2ServersServerIDStdinPost**](ServersApi.md#apiv2serversserveridstdinpost) | **POST** /api/v2/servers/{serverID}/stdin | Send a STDin command to the server. |
| [**apiV2ServersServerIDTasksPost**](ServersApi.md#apiv2serversserveridtaskspostoperation) | **POST** /api/v2/servers/{serverID}/tasks | Create a schedule for a server |
| [**apiV2ServersServerIDTasksTaskIDDelete**](ServersApi.md#apiv2serversserveridtaskstaskiddelete) | **DELETE** /api/v2/servers/{serverID}/tasks/{TaskID} | Remove a schedule |
| [**apiV2ServersServerIDTasksTaskIDPatch**](ServersApi.md#apiv2serversserveridtaskstaskidpatchoperation) | **PATCH** /api/v2/servers/{serverID}/tasks/{TaskID} | Modify a Schedule for a Server |
| [**apiV2ServersServerIDUsersGet**](ServersApi.md#apiv2serversserveridusersget) | **GET** /api/v2/servers/{serverID}/users | Get all users with access to a server |
| [**apiV2ServersServerIDWebhookGet**](ServersApi.md#apiv2serversserveridwebhookget) | **GET** /api/v2/servers/{serverID}/webhook | Get all server webhooks |
| [**apiV2ServersServerIDWebhookPost**](ServersApi.md#apiv2serversserveridwebhookpostoperation) | **POST** /api/v2/servers/{serverID}/webhook | Create a server webhook |



## apiV2ServerServerIDWebhookWebhookIDDelete

> StatusOK apiV2ServerServerIDWebhookWebhookIDDelete(serverID, webhookID)

Delete a webhook

Delete a webhook by ID. Requires config permission to the server.

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '';
import type { ApiV2ServerServerIDWebhookWebhookIDDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ServersApi(config);

  const body = {
    // string | Target server\'s ID.
    serverID: 5adcec83-684c-4555-a7e9-9d913203d07e,
    // number | Target webhook ID
    webhookID: 1,
  } satisfies ApiV2ServerServerIDWebhookWebhookIDDeleteRequest;

  try {
    const data = await api.apiV2ServerServerIDWebhookWebhookIDDelete(body);
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
| **serverID** | `string` | Target server\&#39;s ID. | [Defaults to `undefined`] |
| **webhookID** | `number` | Target webhook ID | [Defaults to `undefined`] |

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
| **200** | Successfully deleted webhook. |  -  |
| **403** | Failed authentication or invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2ServerServerIDWebhookWebhookIDGet

> ApiV2ServerServerIDWebhookWebhookIDGet200Response apiV2ServerServerIDWebhookWebhookIDGet(serverID, webhookID)

Get a webhook

Get a webhook by ID. Requires config permission.

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '';
import type { ApiV2ServerServerIDWebhookWebhookIDGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ServersApi(config);

  const body = {
    // string | Target server\'s ID.
    serverID: 5adcec83-684c-4555-a7e9-9d913203d07e,
    // number | Target webhook ID
    webhookID: 1,
  } satisfies ApiV2ServerServerIDWebhookWebhookIDGetRequest;

  try {
    const data = await api.apiV2ServerServerIDWebhookWebhookIDGet(body);
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
| **serverID** | `string` | Target server\&#39;s ID. | [Defaults to `undefined`] |
| **webhookID** | `number` | Target webhook ID | [Defaults to `undefined`] |

### Return type

[**ApiV2ServerServerIDWebhookWebhookIDGet200Response**](ApiV2ServerServerIDWebhookWebhookIDGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully got webhook details. |  -  |
| **403** | Failed authentication or Invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2ServerServerIDWebhookWebhookIDPatch

> StatusError apiV2ServerServerIDWebhookWebhookIDPatch(serverID, webhookID, apiV2ServerServerIDWebhookWebhookIDPatchRequest)

Modify a Server Webhook

Modify a server webhook. Requires config permission.

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '';
import type { ApiV2ServerServerIDWebhookWebhookIDPatchOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ServersApi(config);

  const body = {
    // string | Target server\'s ID.
    serverID: 5adcec83-684c-4555-a7e9-9d913203d07e,
    // number | Target webhook ID
    webhookID: 1,
    // ApiV2ServerServerIDWebhookWebhookIDPatchRequest (optional)
    apiV2ServerServerIDWebhookWebhookIDPatchRequest: ...,
  } satisfies ApiV2ServerServerIDWebhookWebhookIDPatchOperationRequest;

  try {
    const data = await api.apiV2ServerServerIDWebhookWebhookIDPatch(body);
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
| **serverID** | `string` | Target server\&#39;s ID. | [Defaults to `undefined`] |
| **webhookID** | `number` | Target webhook ID | [Defaults to `undefined`] |
| **apiV2ServerServerIDWebhookWebhookIDPatchRequest** | [ApiV2ServerServerIDWebhookWebhookIDPatchRequest](ApiV2ServerServerIDWebhookWebhookIDPatchRequest.md) |  | [Optional] |

### Return type

[**StatusError**](StatusError.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully modified webhook. |  -  |
| **403** | Failed authentication or Invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2ServersGet

> ApiV2ServersGet200Response apiV2ServersGet()

Get all servers you have access to.

Get all servers the user has access to. Requires any valid user token.

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '';
import type { ApiV2ServersGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ServersApi(config);

  try {
    const data = await api.apiV2ServersGet();
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

[**ApiV2ServersGet200Response**](ApiV2ServersGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Returns all servers user has access to. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2ServersPost

> ApiV2ServersPost200Response apiV2ServersPost(apiV2ServersPostRequest)

Create a server

Create a server. Requires server creation permission.

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '';
import type { ApiV2ServersPostOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ServersApi(config);

  const body = {
    // ApiV2ServersPostRequest
    apiV2ServersPostRequest: ...,
  } satisfies ApiV2ServersPostOperationRequest;

  try {
    const data = await api.apiV2ServersPost(body);
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
| **apiV2ServersPostRequest** | [ApiV2ServersPostRequest](ApiV2ServersPostRequest.md) |  | |

### Return type

[**ApiV2ServersPost200Response**](ApiV2ServersPost200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully created new server |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2ServersServerIDActionActionPost

> ApiV2ServersServerIDActionActionPost200Response apiV2ServersServerIDActionActionPost(serverID, action)

Send an action to a server

Send an action to a server. Requires command permission on server.

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '';
import type { ApiV2ServersServerIDActionActionPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ServersApi(config);

  const body = {
    // string | Target server\'s ID.
    serverID: 5adcec83-684c-4555-a7e9-9d913203d07e,
    // 'clone_server' | 'start_server' | 'stop_server' | 'restart_server' | 'kill_server' | 'backup_server' | 'update_executable' | Action to send to server.
    action: action_example,
  } satisfies ApiV2ServersServerIDActionActionPostRequest;

  try {
    const data = await api.apiV2ServersServerIDActionActionPost(body);
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
| **serverID** | `string` | Target server\&#39;s ID. | [Defaults to `undefined`] |
| **action** | `clone_server`, `start_server`, `stop_server`, `restart_server`, `kill_server`, `backup_server`, `update_executable` | Action to send to server. | [Defaults to `undefined`] [Enum: clone_server, start_server, stop_server, restart_server, kill_server, backup_server, update_executable] |

### Return type

[**ApiV2ServersServerIDActionActionPost200Response**](ApiV2ServersServerIDActionActionPost200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully send command to target server. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2ServersServerIDDelete

> StatusOK apiV2ServersServerIDDelete(serverID)

Delete a server

Delete a server. Requires config permission to target server.

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '';
import type { ApiV2ServersServerIDDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ServersApi(config);

  const body = {
    // string | Target server\'s ID.
    serverID: 5adcec83-684c-4555-a7e9-9d913203d07e,
  } satisfies ApiV2ServersServerIDDeleteRequest;

  try {
    const data = await api.apiV2ServersServerIDDelete(body);
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
| **serverID** | `string` | Target server\&#39;s ID. | [Defaults to `undefined`] |

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
| **200** | Successfully deleted target server. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2ServersServerIDGet

> ApiV2ServersServerIDGet200Response apiV2ServersServerIDGet(serverID)

Get a server\&#39;s data

Get all information from a server. Requires access to the server.

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '';
import type { ApiV2ServersServerIDGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ServersApi(config);

  const body = {
    // string | Target server\'s ID.
    serverID: 5adcec83-684c-4555-a7e9-9d913203d07e,
  } satisfies ApiV2ServersServerIDGetRequest;

  try {
    const data = await api.apiV2ServersServerIDGet(body);
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
| **serverID** | `string` | Target server\&#39;s ID. | [Defaults to `undefined`] |

### Return type

[**ApiV2ServersServerIDGet200Response**](ApiV2ServersServerIDGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully got server information. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2ServersServerIDLogsGet

> ApiV2ServersServerIDLogsGet200Response apiV2ServersServerIDLogsGet(serverID, file, colors, raw, html)

Get a server\&#39;s logs

Get a server\&#39;s logs. Requires logs permission to the server.

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '';
import type { ApiV2ServersServerIDLogsGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ServersApi(config);

  const body = {
    // string | Target server\'s ID.
    serverID: 5adcec83-684c-4555-a7e9-9d913203d07e,
    // boolean | Whether to read the log file or stdout. Defaults to false (standard out). (optional)
    file: false,
    // boolean | Whether to add HTML coloring or not. Defaults to false (no color). (optional)
    colors: false,
    // boolean | Whether to disable ANSI stripping or not. Defaults to false (no ANSI stripping). (optional)
    raw: false,
    // boolean | Whether to output in HTML with &lt;br/&gt;\'s or with the standard API JSON format with an array of lines. Defaults to false (optional)
    html: false,
  } satisfies ApiV2ServersServerIDLogsGetRequest;

  try {
    const data = await api.apiV2ServersServerIDLogsGet(body);
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
| **serverID** | `string` | Target server\&#39;s ID. | [Defaults to `undefined`] |
| **file** | `boolean` | Whether to read the log file or stdout. Defaults to false (standard out). | [Optional] [Defaults to `undefined`] |
| **colors** | `boolean` | Whether to add HTML coloring or not. Defaults to false (no color). | [Optional] [Defaults to `undefined`] |
| **raw** | `boolean` | Whether to disable ANSI stripping or not. Defaults to false (no ANSI stripping). | [Optional] [Defaults to `undefined`] |
| **html** | `boolean` | Whether to output in HTML with &amp;lt;br/&amp;gt;\&#39;s or with the standard API JSON format with an array of lines. Defaults to false | [Optional] [Defaults to `undefined`] |

### Return type

[**ApiV2ServersServerIDLogsGet200Response**](ApiV2ServersServerIDLogsGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully got server logs. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2ServersServerIDPatch

> StatusOK apiV2ServersServerIDPatch(serverID, apiV2ServersServerIDPatchRequest)

Modify a server

Modify a server. Requires config permission to target server.

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '';
import type { ApiV2ServersServerIDPatchOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ServersApi(config);

  const body = {
    // string | Server ID of server to modify
    serverID: f53b7b9d-48cc-40f1-89ce-331bda9458ba,
    // ApiV2ServersServerIDPatchRequest (optional)
    apiV2ServersServerIDPatchRequest: ...,
  } satisfies ApiV2ServersServerIDPatchOperationRequest;

  try {
    const data = await api.apiV2ServersServerIDPatch(body);
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
| **serverID** | `string` | Server ID of server to modify | [Defaults to `undefined`] |
| **apiV2ServersServerIDPatchRequest** | [ApiV2ServersServerIDPatchRequest](ApiV2ServersServerIDPatchRequest.md) |  | [Optional] |

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
| **200** | Successfully modified server. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2ServersServerIDPublicGet

> ApiV2ServersServerIDPublicGet200Response apiV2ServersServerIDPublicGet(serverID)

Get a server\&#39;s public data

Get a server\&#39;s public data. Requires logged in permission.

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '';
import type { ApiV2ServersServerIDPublicGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ServersApi(config);

  const body = {
    // string | Target server\'s ID.
    serverID: 5adcec83-684c-4555-a7e9-9d913203d07e,
  } satisfies ApiV2ServersServerIDPublicGetRequest;

  try {
    const data = await api.apiV2ServersServerIDPublicGet(body);
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
| **serverID** | `string` | Target server\&#39;s ID. | [Defaults to `undefined`] |

### Return type

[**ApiV2ServersServerIDPublicGet200Response**](ApiV2ServersServerIDPublicGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully got server\&#39;s public data. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2ServersServerIDStatsGet

> ApiV2ServersServerIDStatsGet200Response apiV2ServersServerIDStatsGet(serverID)

Get a server\&#39;s statistics

Get a server\&#39;s statistics. Requires access to a server.

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '';
import type { ApiV2ServersServerIDStatsGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ServersApi(config);

  const body = {
    // string | Target server\'s ID.
    serverID: 5adcec83-684c-4555-a7e9-9d913203d07e,
  } satisfies ApiV2ServersServerIDStatsGetRequest;

  try {
    const data = await api.apiV2ServersServerIDStatsGet(body);
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
| **serverID** | `string` | Target server\&#39;s ID. | [Defaults to `undefined`] |

### Return type

[**ApiV2ServersServerIDStatsGet200Response**](ApiV2ServersServerIDStatsGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully got server stats |  -  |
| **403** | Authentication incorrect or failed. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2ServersServerIDStdinPost

> StatusOK apiV2ServersServerIDStdinPost(serverID, body)

Send a STDin command to the server.

Sends a STDin command to the server. Requires command permission for the server. This is the equivalent of server console commands. Don\&#39;t sent the slash!

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '';
import type { ApiV2ServersServerIDStdinPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ServersApi(config);

  const body = {
    // string | Target server\'s ID.
    serverID: 5adcec83-684c-4555-a7e9-9d913203d07e,
    // string
    body: body_example,
  } satisfies ApiV2ServersServerIDStdinPostRequest;

  try {
    const data = await api.apiV2ServersServerIDStdinPost(body);
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
| **serverID** | `string` | Target server\&#39;s ID. | [Defaults to `undefined`] |
| **body** | `string` |  | |

### Return type

[**StatusOK**](StatusOK.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: `text/plain`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully send STDin to the server. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2ServersServerIDTasksPost

> ApiV2ServersServerIDTasksPost200Response apiV2ServersServerIDTasksPost(serverID, apiV2ServersServerIDTasksPostRequest)

Create a schedule for a server

Create a schedule for a server. Requires schedule permission to the server.

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '';
import type { ApiV2ServersServerIDTasksPostOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ServersApi(config);

  const body = {
    // string | Target server\'s ID.
    serverID: 5adcec83-684c-4555-a7e9-9d913203d07e,
    // ApiV2ServersServerIDTasksPostRequest
    apiV2ServersServerIDTasksPostRequest: ...,
  } satisfies ApiV2ServersServerIDTasksPostOperationRequest;

  try {
    const data = await api.apiV2ServersServerIDTasksPost(body);
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
| **serverID** | `string` | Target server\&#39;s ID. | [Defaults to `undefined`] |
| **apiV2ServersServerIDTasksPostRequest** | [ApiV2ServersServerIDTasksPostRequest](ApiV2ServersServerIDTasksPostRequest.md) |  | |

### Return type

[**ApiV2ServersServerIDTasksPost200Response**](ApiV2ServersServerIDTasksPost200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully created new schedule |  -  |
| **403** | Failed authentication or Invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2ServersServerIDTasksTaskIDDelete

> StatusOK apiV2ServersServerIDTasksTaskIDDelete(serverID, taskID)

Remove a schedule

Remove a schedule. Requires server schedule permission.

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '';
import type { ApiV2ServersServerIDTasksTaskIDDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ServersApi(config);

  const body = {
    // string | Target server\'s ID.
    serverID: 5adcec83-684c-4555-a7e9-9d913203d07e,
    // number | Target Task ID
    taskID: 1,
  } satisfies ApiV2ServersServerIDTasksTaskIDDeleteRequest;

  try {
    const data = await api.apiV2ServersServerIDTasksTaskIDDelete(body);
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
| **serverID** | `string` | Target server\&#39;s ID. | [Defaults to `undefined`] |
| **taskID** | `number` | Target Task ID | [Defaults to `undefined`] |

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
| **200** | Successfully deleted schedule. |  -  |
| **403** | Failed authentication or Invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2ServersServerIDTasksTaskIDPatch

> StatusOK apiV2ServersServerIDTasksTaskIDPatch(serverID, taskID, apiV2ServersServerIDTasksTaskIDPatchRequest)

Modify a Schedule for a Server

Modify a schedule for a server. Requires Schedule permission for the target server.

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '';
import type { ApiV2ServersServerIDTasksTaskIDPatchOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ServersApi(config);

  const body = {
    // string | Target server\'s ID.
    serverID: 5adcec83-684c-4555-a7e9-9d913203d07e,
    // number | Target Task ID
    taskID: 1,
    // ApiV2ServersServerIDTasksTaskIDPatchRequest (optional)
    apiV2ServersServerIDTasksTaskIDPatchRequest: ...,
  } satisfies ApiV2ServersServerIDTasksTaskIDPatchOperationRequest;

  try {
    const data = await api.apiV2ServersServerIDTasksTaskIDPatch(body);
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
| **serverID** | `string` | Target server\&#39;s ID. | [Defaults to `undefined`] |
| **taskID** | `number` | Target Task ID | [Defaults to `undefined`] |
| **apiV2ServersServerIDTasksTaskIDPatchRequest** | [ApiV2ServersServerIDTasksTaskIDPatchRequest](ApiV2ServersServerIDTasksTaskIDPatchRequest.md) |  | [Optional] |

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
| **200** | Successfully modified schedule. |  -  |
| **403** | Failed authentication or Invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2ServersServerIDUsersGet

> ApiV2ServersServerIDUsersGet200Response apiV2ServersServerIDUsersGet(serverID)

Get all users with access to a server

Get all users with access to a server. Requires access to a server, user config permission, and roles config permission.

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '';
import type { ApiV2ServersServerIDUsersGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ServersApi(config);

  const body = {
    // string | Target server\'s ID.
    serverID: 5adcec83-684c-4555-a7e9-9d913203d07e,
  } satisfies ApiV2ServersServerIDUsersGetRequest;

  try {
    const data = await api.apiV2ServersServerIDUsersGet(body);
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
| **serverID** | `string` | Target server\&#39;s ID. | [Defaults to `undefined`] |

### Return type

[**ApiV2ServersServerIDUsersGet200Response**](ApiV2ServersServerIDUsersGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully got users with access to a server. |  -  |
| **403** | Failed authentication or Invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2ServersServerIDWebhookGet

> ApiV2ServersServerIDWebhookGet200Response apiV2ServersServerIDWebhookGet(serverID)

Get all server webhooks

Get all server webhooks. Requires config permission to the server.

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '';
import type { ApiV2ServersServerIDWebhookGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ServersApi(config);

  const body = {
    // string | Target server\'s ID.
    serverID: 5adcec83-684c-4555-a7e9-9d913203d07e,
  } satisfies ApiV2ServersServerIDWebhookGetRequest;

  try {
    const data = await api.apiV2ServersServerIDWebhookGet(body);
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
| **serverID** | `string` | Target server\&#39;s ID. | [Defaults to `undefined`] |

### Return type

[**ApiV2ServersServerIDWebhookGet200Response**](ApiV2ServersServerIDWebhookGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully got all server webhooks. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2ServersServerIDWebhookPost

> ApiV2ServersServerIDWebhookPost200Response apiV2ServersServerIDWebhookPost(serverID, apiV2ServersServerIDWebhookPostRequest)

Create a server webhook

Create a server webhook. Requires server config permission.

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '';
import type { ApiV2ServersServerIDWebhookPostOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ServersApi(config);

  const body = {
    // string | Target server\'s ID.
    serverID: 5adcec83-684c-4555-a7e9-9d913203d07e,
    // ApiV2ServersServerIDWebhookPostRequest (optional)
    apiV2ServersServerIDWebhookPostRequest: ...,
  } satisfies ApiV2ServersServerIDWebhookPostOperationRequest;

  try {
    const data = await api.apiV2ServersServerIDWebhookPost(body);
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
| **serverID** | `string` | Target server\&#39;s ID. | [Defaults to `undefined`] |
| **apiV2ServersServerIDWebhookPostRequest** | [ApiV2ServersServerIDWebhookPostRequest](ApiV2ServersServerIDWebhookPostRequest.md) |  | [Optional] |

### Return type

[**ApiV2ServersServerIDWebhookPost200Response**](ApiV2ServersServerIDWebhookPost200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully created webhook. |  -  |
| **403** | Failed authentication or Invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

