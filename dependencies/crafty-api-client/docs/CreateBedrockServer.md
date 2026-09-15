
# CreateBedrockServer


## Properties

Name | Type
------------ | -------------
`name` | string
`monitoringType` | string
`minecraftBedrockMonitoringData` | [CreateBedrockServerMinecraftBedrockMonitoringData](CreateBedrockServerMinecraftBedrockMonitoringData.md)
`createType` | string
`minecraftBedrockCreateData` | [CreateBedrockServerMinecraftBedrockCreateData](CreateBedrockServerMinecraftBedrockCreateData.md)

## Example

```typescript
import type { CreateBedrockServer } from ''

// TODO: Update the object below with actual values
const example = {
  "name": Example Bedrock server.,
  "monitoringType": minecraft_bedrock,
  "minecraftBedrockMonitoringData": null,
  "createType": minecraft_bedrock,
  "minecraftBedrockCreateData": null,
} satisfies CreateBedrockServer

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as CreateBedrockServer
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


