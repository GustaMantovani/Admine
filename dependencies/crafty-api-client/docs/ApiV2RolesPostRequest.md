
# ApiV2RolesPostRequest


## Properties

Name | Type
------------ | -------------
`name` | string
`server` | [Array&lt;ApiV2RolesPostRequestServerInner&gt;](ApiV2RolesPostRequestServerInner.md)
`manager` | number
`mfaRequired` | boolean

## Example

```typescript
import type { ApiV2RolesPostRequest } from ''

// TODO: Update the object below with actual values
const example = {
  "name": Example role,
  "server": null,
  "manager": 1,
  "mfaRequired": false,
} satisfies ApiV2RolesPostRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV2RolesPostRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


