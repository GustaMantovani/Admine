
# User


## Properties

Name | Type
------------ | -------------
`userId` | number
`created` | string
`username` | string
`enabled` | boolean
`superuser` | boolean
`lang` | string

## Example

```typescript
import type { User } from ''

// TODO: Update the object below with actual values
const example = {
  "userId": 1,
  "created": 05/18/2022, 01:43:56,
  "username": admin,
  "enabled": true,
  "superuser": true,
  "lang": en_EN,
} satisfies User

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as User
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


