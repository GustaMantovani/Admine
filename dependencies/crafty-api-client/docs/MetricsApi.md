# MetricsApi

All URIs are relative to *https://api.example.com*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**metricsGet**](MetricsApi.md#metricsget) | **GET** /metrics | Global Crafty metrics |
| [**metricsHostGet**](MetricsApi.md#metricshostget) | **GET** /metrics/host | Host system metrics |
| [**metricsServersServerIDGet**](MetricsApi.md#metricsserversserveridget) | **GET** /metrics/servers/{serverID} | Metrics for a specific managed server |



## metricsGet

> string metricsGet()

Global Crafty metrics

Returns Crafty\&#39;s global Prometheus/OpenMetrics registry. Includes Python runtime metrics and Crafty instance metadata. 

### Example

```ts
import {
  Configuration,
  MetricsApi,
} from '';
import type { MetricsGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new MetricsApi(config);

  try {
    const data = await api.metricsGet();
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

**string**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `text/plain`, `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Metrics output in OpenMetrics format. |  -  |
| **403** | Authentication failed. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## metricsHostGet

> string metricsHostGet()

Host system metrics

Returns host-level CPU and memory metrics from the machine running Crafty. Output is Prometheus/OpenMetrics format. 

### Example

```ts
import {
  Configuration,
  MetricsApi,
} from '';
import type { MetricsHostGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new MetricsApi(config);

  try {
    const data = await api.metricsHostGet();
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

**string**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `text/plain`, `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Host metrics in OpenMetrics format. |  -  |
| **400** | User lacks global access permissions. |  -  |
| **403** | Authentication failed. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## metricsServersServerIDGet

> string metricsServersServerIDGet(serverID)

Metrics for a specific managed server

Returns detailed Prometheus/OpenMetrics metrics for a Crafty-managed server, including CPU, memory, version, and player information. Requires that the authenticated user has access to the target server. 

### Example

```ts
import {
  Configuration,
  MetricsApi,
} from '';
import type { MetricsServersServerIDGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new MetricsApi(config);

  const body = {
    // string | Target server\'s ID.
    serverID: 5adcec83-684c-4555-a7e9-9d913203d07e,
  } satisfies MetricsServersServerIDGetRequest;

  try {
    const data = await api.metricsServersServerIDGet(body);
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

**string**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `text/plain`, `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Server metrics in OpenMetrics format. |  -  |
| **400** | User does not have permission for this server. |  -  |
| **403** | Authentication failed. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

