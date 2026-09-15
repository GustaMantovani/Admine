
# ApiV2ServersPostRequest


## Properties

Name | Type
------------ | -------------
`name` | string
`monitoringType` | string
`minecraftJavaMonitoringData` | [CreateJavaServerMinecraftJavaMonitoringData](CreateJavaServerMinecraftJavaMonitoringData.md)
`createType` | string
`minecraftJavaCreateData` | [CreateJavaServerMinecraftJavaCreateData](CreateJavaServerMinecraftJavaCreateData.md)
`minecraftBedrockMonitoringData` | [CreateBedrockServerMinecraftBedrockMonitoringData](CreateBedrockServerMinecraftBedrockMonitoringData.md)
`minecraftBedrockCreateData` | [CreateBedrockServerMinecraftBedrockCreateData](CreateBedrockServerMinecraftBedrockCreateData.md)
`hytaleMonitoringData` | [CreateHytaleServerHytaleMonitoringData](CreateHytaleServerHytaleMonitoringData.md)
`hytaleCreateData` | [CreateHytaleServerHytaleCreateData](CreateHytaleServerHytaleCreateData.md)

## Example

```typescript
import type { ApiV2ServersPostRequest } from ''

// TODO: Update the object below with actual values
const example = {
  "name": Example Hytale server.,
  "monitoringType": hytale,
  "minecraftJavaMonitoringData": null,
  "createType": hytale,
  "minecraftJavaCreateData": null,
  "minecraftBedrockMonitoringData": null,
  "minecraftBedrockCreateData": null,
  "hytaleMonitoringData": null,
  "hytaleCreateData": null,
} satisfies ApiV2ServersPostRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV2ServersPostRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


