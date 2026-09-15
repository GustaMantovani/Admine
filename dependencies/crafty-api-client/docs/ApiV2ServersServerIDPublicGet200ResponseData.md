
# ApiV2ServersServerIDPublicGet200ResponseData


## Properties

Name | Type
------------ | -------------
`serverId` | string
`created` | string
`serverName` | string
`type` | string

## Example

```typescript
import type { ApiV2ServersServerIDPublicGet200ResponseData } from ''

// TODO: Update the object below with actual values
const example = {
  "serverId": f53b7b9d-48cc-40f1-89ce-331bda9458ba,
  "created": 2022-05-25T01:24:22.427327,
  "serverName": Example server,
  "type": minecraft-java,
} satisfies ApiV2ServersServerIDPublicGet200ResponseData

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV2ServersServerIDPublicGet200ResponseData
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


