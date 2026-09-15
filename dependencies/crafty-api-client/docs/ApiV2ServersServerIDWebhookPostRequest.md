
# ApiV2ServersServerIDWebhookPostRequest


## Properties

Name | Type
------------ | -------------
`webhookType` | string
`name` | string
`url` | string
`botName` | string
`trigger` | Array&lt;string&gt;
`body` | string
`enabled` | boolean

## Example

```typescript
import type { ApiV2ServersServerIDWebhookPostRequest } from ''

// TODO: Update the object below with actual values
const example = {
  "webhookType": Discord,
  "name": My example Webhook,
  "url": https://webhook-url.example.com,
  "botName": Crafty Bot,
  "trigger": null,
  "body": Server Has Started,
  "enabled": true,
} satisfies ApiV2ServersServerIDWebhookPostRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ApiV2ServersServerIDWebhookPostRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


