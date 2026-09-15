
# ApiV2ServersServerIDTasksPostRequest


## Properties

Name | Type
------------ | -------------
`name` | string
`enabled` | boolean
`action` | string
`interval` | number
`intervalType` | string
`command` | string
`parent` | string
`oneTime` | boolean
`cronString` | string
`delay` | number

## Example

```typescript
import type { ApiV2ServersServerIDTasksPostRequest } from ''

// TODO: Update the object below with actual values
const example = {
  "name": Example Task,
  "enabled": true,
  "action": custom_command,
  "interval": 1,
  "intervalType": minutes,
  "command": say hello world,
  "parent": null,
  "oneTime": false,
  "cronString": ,
  "delay": 0,
} satisfies ApiV2ServersServerIDTasksPostRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV2ServersServerIDTasksPostRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


