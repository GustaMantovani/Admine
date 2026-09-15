
# ApiV2CraftyConfigPatchRequest


## Properties

Name | Type
------------ | -------------
`httpsPort` | number
`language` | string
`cookieExpire` | number
`showErrors` | boolean
`historyMaxAge` | number
`statusUpdateFrequencySeconds` | number
`deleteDefaultJson` | boolean
`showContributeLink` | boolean
`virtualTerminalLines` | number
`maxLogLines` | number
`maxAuditEntries` | number
`disabledLanguageFiles` | Array&lt;string&gt;
`streamSizeGB` | number
`keywords` | Array&lt;string&gt;
`allowNsfwProfilePictures` | boolean
`enableUserSelfDelete` | boolean
`resetSecretsOnNextBoot` | boolean
`monitoredMounts` | Array&lt;string&gt;
`dirSizePollFreqMinutes` | number
`craftyLogsDeleteAfterDays` | number
`bigBucketRepo` | string
`enableOtpSkew` | boolean
`maxLoginAttempts` | number
`superMFA` | boolean

## Example

```typescript
import type { ApiV2CraftyConfigPatchRequest } from ''

// TODO: Update the object below with actual values
const example = {
  "httpsPort": 8443,
  "language": en_EN,
  "cookieExpire": 30,
  "showErrors": true,
  "historyMaxAge": 0,
  "statusUpdateFrequencySeconds": 30,
  "deleteDefaultJson": false,
  "showContributeLink": true,
  "virtualTerminalLines": 300,
  "maxLogLines": 700,
  "maxAuditEntries": 7,
  "disabledLanguageFiles": null,
  "streamSizeGB": 1,
  "keywords": null,
  "allowNsfwProfilePictures": true,
  "enableUserSelfDelete": false,
  "resetSecretsOnNextBoot": false,
  "monitoredMounts": null,
  "dirSizePollFreqMinutes": 5,
  "craftyLogsDeleteAfterDays": 5,
  "bigBucketRepo": https://jars.arcadiatech.org,
  "enableOtpSkew": false,
  "maxLoginAttempts": 3,
  "superMFA": false,
} satisfies ApiV2CraftyConfigPatchRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV2CraftyConfigPatchRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


