
# Role


## Properties

Name | Type
------------ | -------------
`roleId` | number
`created` | string
`lastUpdate` | string
`roleName` | string
`manager` | number
`mfaRequired` | boolean

## Example

```typescript
import type { Role } from ''

// TODO: Update the object below with actual values
const example = {
  "roleId": 1,
  "created": 05/18/2022, 16:25:43,
  "lastUpdate": 2022-05-18T16:25:45.377352,
  "roleName": test 1234,
  "manager": 1,
  "mfaRequired": false,
} satisfies Role

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Role
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


