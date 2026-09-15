
# CreateJavaServerMinecraftJavaCreateData


## Properties

Name | Type
------------ | -------------
`createType` | string
`downloadJarCreateData` | [CreateJavaServerMinecraftJavaCreateDataDownloadJarCreateData](CreateJavaServerMinecraftJavaCreateDataDownloadJarCreateData.md)

## Example

```typescript
import type { CreateJavaServerMinecraftJavaCreateData } from ''

// TODO: Update the object below with actual values
const example = {
  "createType": download_jar,
  "downloadJarCreateData": null,
} satisfies CreateJavaServerMinecraftJavaCreateData

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as CreateJavaServerMinecraftJavaCreateData
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


