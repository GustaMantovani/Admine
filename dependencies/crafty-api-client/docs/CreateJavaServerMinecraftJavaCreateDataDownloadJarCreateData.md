
# CreateJavaServerMinecraftJavaCreateDataDownloadJarCreateData


## Properties

Name | Type
------------ | -------------
`category` | string
`type` | string
`version` | string
`memMin` | number
`memMax` | number
`serverPropertiesPort` | number

## Example

```typescript
import type { CreateJavaServerMinecraftJavaCreateDataDownloadJarCreateData } from ''

// TODO: Update the object below with actual values
const example = {
  "category": mc_java_servers,
  "type": paper,
  "version": 1.18.2,
  "memMin": 1,
  "memMax": 2,
  "serverPropertiesPort": 25565,
} satisfies CreateJavaServerMinecraftJavaCreateDataDownloadJarCreateData

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as CreateJavaServerMinecraftJavaCreateDataDownloadJarCreateData
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


