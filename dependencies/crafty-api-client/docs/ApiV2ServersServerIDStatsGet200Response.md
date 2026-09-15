
# ApiV2ServersServerIDStatsGet200Response


## Properties

Name | Type
------------ | -------------
`status` | string
`data` | [Server](Server.md)
`started` | string
`running` | boolean
`cpu` | number
`mem` | string
`memPercent` | string
`worldName` | string
`worldSize` | string
`serverPort` | number
`intPingResults` | string
`online` | number
`max` | number
`players` | string
`desc` | string
`icon` | string
`version` | string
`updating` | boolean
`waitingStart` | boolean
`firstRun` | boolean
`crashed` | boolean

## Example

```typescript
import type { ApiV2ServersServerIDStatsGet200Response } from ''

// TODO: Update the object below with actual values
const example = {
  "status": ok,
  "data": null,
  "started": 2022-05-25 15:44:05,
  "running": true,
  "cpu": null,
  "mem": 42.2MB,
  "memPercent": 5.0,
  "worldName": Example server,
  "worldSize": 10.3MB,
  "serverPort": 25565,
  "intPingResults": True,
  "online": 5,
  "max": 20,
  "players": [],
  "desc": A Minecraft Server,
  "icon": ,
  "version": 1.21.3,
  "updating": false,
  "waitingStart": false,
  "firstRun": false,
  "crashed": false,
} satisfies ApiV2ServersServerIDStatsGet200Response

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV2ServersServerIDStatsGet200Response
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


