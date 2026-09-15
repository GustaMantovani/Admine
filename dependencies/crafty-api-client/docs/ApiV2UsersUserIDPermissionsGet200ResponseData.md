
# ApiV2UsersUserIDPermissionsGet200ResponseData


## Properties

Name | Type
------------ | -------------
`permissions` | string
`counters` | [ApiV2UsersUserIDPermissionsGet200ResponseDataCounters](ApiV2UsersUserIDPermissionsGet200ResponseDataCounters.md)
`limits` | [ApiV2UsersUserIDPermissionsGet200ResponseDataLimits](ApiV2UsersUserIDPermissionsGet200ResponseDataLimits.md)

## Example

```typescript
import type { ApiV2UsersUserIDPermissionsGet200ResponseData } from ''

// TODO: Update the object below with actual values
const example = {
  "permissions": 110,
  "counters": null,
  "limits": null,
} satisfies ApiV2UsersUserIDPermissionsGet200ResponseData

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV2UsersUserIDPermissionsGet200ResponseData
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


