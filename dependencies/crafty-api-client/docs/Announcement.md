
# Announcement


## Properties

Name | Type
------------ | -------------
`id` | string
`date` | string
`title` | string
`desc` | string
`link` | string

## Example

```typescript
import type { Announcement } from ''

// TODO: Update the object below with actual values
const example = {
  "id": dbc0c2aa-690a-46c3-aa67-db8e3bde0cd0,
  "date": 2023/10/18,
  "title": Crafty Notification Center,
  "desc": This is the Crafty notification center. This allows Arcadia Technology and Crafty Controller send notifications directly to you!,
  "link": https://discord.com/channels/640894866934726656/654307216970612746,
} satisfies Announcement

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Announcement
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


