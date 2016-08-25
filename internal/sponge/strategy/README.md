# Strategy

import "github.com/coralproject/shelf/internal/sponge/strategy"

It handles the loading and distribution of configuration related with external sources. It has the translation from the external database to our coral schema.

Explaining how to write a strategy file.

## Name

The name of the strategy that we are describing.

### DateTimeFormat

If your source have date time fields, we need to know how to parse them. You should write the representation of 2006 Mon Jan 2 15:04:05 in the desired format. More more info read [pkg-constants](https://golang.org/pkg/time/#pkg-constants)

### Entities

It describe all the different entities we have at the Coral database and how to do its transformations.

#### Name of the entity

Example:
```
"asset": {
  "Foreign": "crnr_asset",
  "Local": "asset",
  "OrderBy": "assetid",
  "Fields": [
    {
      "foreign": "assetid",
      "local": "asset_id",
      "relation": "Source",
      "type": "int"
    },
    {
      "foreign": "asseturl",
      "local": "url",
      "relation": "Passthrough",
      "type": "int"
    },
    {
      "foreign": "updatedate",
      "local": "date_updated",
      "relation": "ParseTimeDate",
      "type": "dateTime"
    },
    {
      "foreign": "createdate",
      "local": "date_created",
      "relation": "ParseTimeDate",
      "type": "dateTime"
    }
  ]
}
```

##### Foreign

The name of the foreign entity.

##### Local

The collection to where we are importing this entity into.

##### OrderBy

A default order by when quering the foreign source.

##### Fields

All the fields that are being mapped.

###### Foreign

The name of the field in the foreign entity.

###### Local

The name of the field in our local database.

###### Relation

The relationship between the foreign field and the local one. We have this options:
- Passthrough: when the value is the same
- Source: when it needs to be added to our source struct for the local collection (the original identifiers have to go into source)
- ParseTimeDate: when we need to parse the foreign value as date time.
- Constant: when the local field should always be the same value. In this case we will have "foreign" blank and we will have other field called "value" with the value of the local field.
- SubDocument: when the local field has an array of documents in one of the fields.
- Status: when the field need to be translated based on the status map that is declared in that same strategy file for the entity.

###### Type

The type of the value we are converting.

- String
- Timedate
