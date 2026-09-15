
# ApiV2AuthLoginPostRequest


## Properties

Name | Type
------------ | -------------
`username` | string
`password` | string
`totp` | string
`backupCode` | string

## Example

```typescript
import type { ApiV2AuthLoginPostRequest } from ''

// TODO: Update the object below with actual values
const example = {
  "username": admin,
  "password": password123,
  "totp": 123456,
  "backupCode": abc-xyz,
} satisfies ApiV2AuthLoginPostRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV2AuthLoginPostRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


