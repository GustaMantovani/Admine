
# Webhook


## Properties

Name | Type
------------ | -------------
`webhookType` | string
`name` | string
`url` | string
`botName` | string
`trigger` | string
`body` | string
`enabled` | boolean

## Example

```typescript
import type { Webhook } from ''

// TODO: Update the object below with actual values
const example = {
  "webhookType": Discord,
  "name": Example webhook,
  "url": https://webhook.example.com/,
  "botName": Example bot,
  "trigger": server_start,
  "body": Webhook body text.,
  "enabled": true,
} satisfies Webhook

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Webhook
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


