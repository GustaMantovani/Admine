
# CreateBedrockServerMinecraftBedrockCreateData


## Properties

Name | Type
------------ | -------------
`createType` | string
`downloadExeCreateData` | [CreateBedrockServerMinecraftBedrockCreateDataDownloadExeCreateData](CreateBedrockServerMinecraftBedrockCreateDataDownloadExeCreateData.md)

## Example

```typescript
import type { CreateBedrockServerMinecraftBedrockCreateData } from ''

// TODO: Update the object below with actual values
const example = {
  "createType": download_exe,
  "downloadExeCreateData": null,
} satisfies CreateBedrockServerMinecraftBedrockCreateData

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as CreateBedrockServerMinecraftBedrockCreateData
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


