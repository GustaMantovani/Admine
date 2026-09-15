
# Server


## Properties

Name | Type
------------ | -------------
`serverId` | string
`created` | string
`serverUuid` | string
`serverName` | string
`path` | string
`backupPath` | string
`executable` | string
`logPath` | string
`executionCommand` | string
`autoStart` | boolean
`autoStartDelay` | number
`crashDetection` | boolean
`stopCommand` | string
`executableUpdateUrl` | string
`serverIp` | string
`serverPort` | number
`logsDeleteAfter` | number
`type` | string

## Example

```typescript
import type { Server } from ''

// TODO: Update the object below with actual values
const example = {
  "serverId": 5adcec83-684c-4555-a7e9-9d913203d07e,
  "created": 2022-05-18T22:36:04.751211,
  "serverUuid": 5adcec83-684c-4555-a7e9-9d913203d07e,
  "serverName": Example Server,
  "path": /servers/e7e107ed-b4fd-4f87-b2dd-b66b4930706e,
  "backupPath": /backups/e7e107ed-b4fd-4f87-b2dd-b66b4930706e,
  "executable": paper-1.18.2.jar,
  "logPath": /servers/e7e107ed-b4fd-4f87-b2dd-b66b4930706e/logs/latest.log,
  "executionCommand": java -Xms1000M -Xmx2000M -jar /servers/e7e107ed-b4fd-4f87-b2dd-b66b4930706e/paper-1.18.2.jar nogui,
  "autoStart": false,
  "autoStartDelay": 10,
  "crashDetection": false,
  "stopCommand": stop,
  "executableUpdateUrl": null,
  "serverIp": 127.0.0.1,
  "serverPort": 25565,
  "logsDeleteAfter": 0,
  "type": minecraft-java,
} satisfies Server

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Server
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


