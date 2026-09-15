
# CreateJavaServer


## Properties

Name | Type
------------ | -------------
`name` | string
`monitoringType` | string
`minecraftJavaMonitoringData` | [CreateJavaServerMinecraftJavaMonitoringData](CreateJavaServerMinecraftJavaMonitoringData.md)
`createType` | string
`minecraftJavaCreateData` | [CreateJavaServerMinecraftJavaCreateData](CreateJavaServerMinecraftJavaCreateData.md)

## Example

```typescript
import type { CreateJavaServer } from ''

// TODO: Update the object below with actual values
const example = {
  "name": Example Java server.,
  "monitoringType": minecraft_java,
  "minecraftJavaMonitoringData": null,
  "createType": minecraft_java,
  "minecraftJavaCreateData": null,
} satisfies CreateJavaServer

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as CreateJavaServer
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


