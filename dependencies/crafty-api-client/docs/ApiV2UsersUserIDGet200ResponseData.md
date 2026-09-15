
# ApiV2UsersUserIDGet200ResponseData


## Properties

Name | Type
------------ | -------------
`userId` | number
`created` | string
`lastLogin` | string
`lastUpdate` | string
`lastIp` | string
`username` | string
`enabled` | boolean
`superuser` | boolean
`lang` | string
`supportLogs` | string
`serverOrder` | string
`preparing` | boolean
`hints` | boolean
`roles` | [Array&lt;Role&gt;](Role.md)

## Example

```typescript
import type { ApiV2UsersUserIDGet200ResponseData } from ''

// TODO: Update the object below with actual values
const example = {
  "userId": 3,
  "created": 05/24/2022, 09:44:33,
  "lastLogin": 2022-05-24T09:44:33.646325,
  "lastUpdate": 05/24/2022, 09:44:33,
  "lastIp": 127.0.0.1,
  "username": example,
  "enabled": true,
  "superuser": true,
  "lang": fi_FI,
  "supportLogs": ,
  "serverOrder": ,
  "preparing": false,
  "hints": true,
  "roles": null,
} satisfies ApiV2UsersUserIDGet200ResponseData

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV2UsersUserIDGet200ResponseData
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


