# CraftyApi

All URIs are relative to *https://api.example.com*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**apiV2CraftyAnnouncementsGet**](CraftyApi.md#apiv2craftyannouncementsget) | **GET** /api/v2/crafty/announcements | Get announcements |
| [**apiV2CraftyAnnouncementsPost**](CraftyApi.md#apiv2craftyannouncementspostoperation) | **POST** /api/v2/crafty/announcements | Clear a notification. |
| [**apiV2CraftyCheckGet**](CraftyApi.md#apiv2craftycheckget) | **GET** /api/v2/crafty/check | Status check |
| [**apiV2CraftyConfigCustomizeGet**](CraftyApi.md#apiv2craftyconfigcustomizeget) | **GET** /api/v2/crafty/config/customize | Currently unused |
| [**apiV2CraftyConfigCustomizePatch**](CraftyApi.md#apiv2craftyconfigcustomizepatchoperation) | **PATCH** /api/v2/crafty/config/customize | Update Crafty customizations |
| [**apiV2CraftyConfigGet**](CraftyApi.md#apiv2craftyconfigget) | **GET** /api/v2/crafty/config | Get Crafty\&#39;s config |
| [**apiV2CraftyConfigPatch**](CraftyApi.md#apiv2craftyconfigpatchoperation) | **PATCH** /api/v2/crafty/config | Modify config.jso |
| [**apiV2CraftyResetPassGet**](CraftyApi.md#apiv2craftyresetpassget) | **GET** /api/v2/crafty/resetPass | Initiate password recovery system. |



## apiV2CraftyAnnouncementsGet

> ApiV2CraftyAnnouncementsGet200Response apiV2CraftyAnnouncementsGet()

Get announcements

Get announcements that are shown to the user

### Example

```ts
import {
  Configuration,
  CraftyApi,
} from '';
import type { ApiV2CraftyAnnouncementsGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new CraftyApi(config);

  try {
    const data = await api.apiV2CraftyAnnouncementsGet();
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

[**ApiV2CraftyAnnouncementsGet200Response**](ApiV2CraftyAnnouncementsGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully got announcements for user. |  -  |
| **403** | Failed authentication or invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2CraftyAnnouncementsPost

> ApiV2CraftyAnnouncementsPost200Response apiV2CraftyAnnouncementsPost(apiV2CraftyAnnouncementsPostRequest)

Clear a notification.

Clear a notification. It will no longer show after clearing.

### Example

```ts
import {
  Configuration,
  CraftyApi,
} from '';
import type { ApiV2CraftyAnnouncementsPostOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new CraftyApi(config);

  const body = {
    // ApiV2CraftyAnnouncementsPostRequest (optional)
    apiV2CraftyAnnouncementsPostRequest: ...,
  } satisfies ApiV2CraftyAnnouncementsPostOperationRequest;

  try {
    const data = await api.apiV2CraftyAnnouncementsPost(body);
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
| **apiV2CraftyAnnouncementsPostRequest** | [ApiV2CraftyAnnouncementsPostRequest](ApiV2CraftyAnnouncementsPostRequest.md) |  | [Optional] |

### Return type

[**ApiV2CraftyAnnouncementsPost200Response**](ApiV2CraftyAnnouncementsPost200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully cleared notification. |  -  |
| **403** | Failed authentication or invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2CraftyCheckGet

> StatusOK apiV2CraftyCheckGet()

Status check

Crafty status/alive check.

### Example

```ts
import {
  Configuration,
  CraftyApi,
} from '';
import type { ApiV2CraftyCheckGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new CraftyApi(config);

  try {
    const data = await api.apiV2CraftyCheckGet();
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
| **200** | Successful status check |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2CraftyConfigCustomizeGet

> apiV2CraftyConfigCustomizeGet(ids)

Currently unused

Get Crafty customizations. Currently unused. Requires superuser.

### Example

```ts
import {
  Configuration,
  CraftyApi,
} from '';
import type { ApiV2CraftyConfigCustomizeGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new CraftyApi(config);

  const body = {
    // boolean (optional)
    ids: true,
  } satisfies ApiV2CraftyConfigCustomizeGetRequest;

  try {
    const data = await api.apiV2CraftyConfigCustomizeGet(body);
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
| **ids** | `boolean` |  | [Optional] [Defaults to `undefined`] |

### Return type

`void` (Empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **403** | Failed authentication or invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2CraftyConfigCustomizePatch

> ApiV2CraftyConfigCustomizePatch200Response apiV2CraftyConfigCustomizePatch(apiV2CraftyConfigCustomizePatchRequest)

Update Crafty customizations

Update Crafty customizations. Requires superuser permissions.

### Example

```ts
import {
  Configuration,
  CraftyApi,
} from '';
import type { ApiV2CraftyConfigCustomizePatchOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new CraftyApi(config);

  const body = {
    // ApiV2CraftyConfigCustomizePatchRequest (optional)
    apiV2CraftyConfigCustomizePatchRequest: ...,
  } satisfies ApiV2CraftyConfigCustomizePatchOperationRequest;

  try {
    const data = await api.apiV2CraftyConfigCustomizePatch(body);
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
| **apiV2CraftyConfigCustomizePatchRequest** | [ApiV2CraftyConfigCustomizePatchRequest](ApiV2CraftyConfigCustomizePatchRequest.md) |  | [Optional] |

### Return type

[**ApiV2CraftyConfigCustomizePatch200Response**](ApiV2CraftyConfigCustomizePatch200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully updated login photo. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2CraftyConfigGet

> ApiV2CraftyConfigGet200Response apiV2CraftyConfigGet(ids)

Get Crafty\&#39;s config

Get Crafty\&#39;s config. Requires superuser permission. Currently unused.

### Example

```ts
import {
  Configuration,
  CraftyApi,
} from '';
import type { ApiV2CraftyConfigGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new CraftyApi(config);

  const body = {
    // boolean | Get only ids (optional)
    ids: true,
  } satisfies ApiV2CraftyConfigGetRequest;

  try {
    const data = await api.apiV2CraftyConfigGet(body);
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
| **ids** | `boolean` | Get only ids | [Optional] [Defaults to `undefined`] |

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
| **200** | Successfully got user\&#39;s roles. |  -  |
| **403** | Failed authentication or invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2CraftyConfigPatch

> StatusOK apiV2CraftyConfigPatch(apiV2CraftyConfigPatchRequest)

Modify config.jso

Modify config.json. Requires superuser permission.

### Example

```ts
import {
  Configuration,
  CraftyApi,
} from '';
import type { ApiV2CraftyConfigPatchOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new CraftyApi(config);

  const body = {
    // ApiV2CraftyConfigPatchRequest (optional)
    apiV2CraftyConfigPatchRequest: ...,
  } satisfies ApiV2CraftyConfigPatchOperationRequest;

  try {
    const data = await api.apiV2CraftyConfigPatch(body);
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
| **apiV2CraftyConfigPatchRequest** | [ApiV2CraftyConfigPatchRequest](ApiV2CraftyConfigPatchRequest.md) |  | [Optional] |

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
| **200** | Successfully updated config.json |  -  |
| **403** | Failed authentication or invalid authentication. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV2CraftyResetPassGet

> ApiV2CraftyResetPassGet200Response apiV2CraftyResetPassGet()

Initiate password recovery system.

Initial password recovery. No permissions required.

### Example

```ts
import {
  Configuration,
  CraftyApi,
} from '';
import type { ApiV2CraftyResetPassGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new CraftyApi(config);

  try {
    const data = await api.apiV2CraftyResetPassGet();
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

[**ApiV2CraftyResetPassGet200Response**](ApiV2CraftyResetPassGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully started password reset procedure. |  -  |
| **425** | Request sent before end of previous anti-lockout. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

