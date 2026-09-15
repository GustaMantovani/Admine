
# ApiV2UsersUserIDPublicGet200ResponseData


## Properties

Name | Type
------------ | -------------
`userId` | number
`created` | string
`username` | string
`enabled` | boolean
`superuser` | boolean
`lang` | string
`roles` | Array&lt;number&gt;

## Example

```typescript
import type { ApiV2UsersUserIDPublicGet200ResponseData } from ''

// TODO: Update the object below with actual values
const example = {
  "userId": 3,
  "created": 05/24/2022, 09:44:33,
  "username": example,
  "enabled": true,
  "superuser": true,
  "lang": fi_FI,
  "roles": null,
} satisfies ApiV2UsersUserIDPublicGet200ResponseData

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV2UsersUserIDPublicGet200ResponseData
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


