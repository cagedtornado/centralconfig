# API server information

API path              | Description
----------            | -----------
/                     | Redirects to the /UI path (and the web based user interface)
[/config/get](https://github.com/danesparza/centralconfig/tree/master/api#configget)           | Gets a single configuration item
[/config/set](https://github.com/danesparza/centralconfig/tree/master/api#configset)           | Sets (creates or updates) a configuration item
[/config/remove](https://github.com/danesparza/centralconfig/tree/master/api#configremove)        | Removes a configuration item
[/config/getall](https://github.com/danesparza/centralconfig/tree/master/api#configgetall)        | Gets all configuration items
[/config/getallforapp](https://github.com/danesparza/centralconfig/tree/master/api#configgetallforapp)  | Get all configuration items for a single application (plus the default * application)
[/applications/getall](https://github.com/danesparza/centralconfig/tree/master/api#applicationsgetall)  | Get all applications

####Requests
Most API operations expect a configitem object in the POST body that will be used to either filter (in a get operation), update or create (in a set operation), or remove an item (in a remove operation).  

For example:

```json
{
    "application" : "WickedCool",
    "name": "TestItem42",
    "value": "Magic!"
}
```

####Responses
All operations will return an object that contain the fields status, message, and data.  

For Example:
```json
{
  "status": 200,
  "message": "Config items found",
  "data": [
    {
      "id": 1,
      "application": "*",
      "machine": "",
      "name": "Environment",
      "value": "DEV",
      "updated": "2016-04-26T09:11:17.897Z"
    },
    {
      "id": 2,
      "application": "TestApp",
      "machine": "",
      "name": "AppUser",
      "value": "TestApp_dev",
      "updated": "2016-04-26T09:11:40.34Z"
    }
  ]
}
```


### /config/get

This operation retrieves a single configuration item.  If it doesn't exist for the given application, it attemps to get it for the default application (*)

######Example request:
```json
{
  "name": "test"
}
```

######Example response:
```json
{
  "name": "test"
}
```

### /config/set

This operation sets the value of a single configuration item

######Example request:
```json
{
  "name": "test"
}
```

######Example response:
```json
{
  "name": "test"
}
```

### /config/remove

This operation removes a single configuration item

######Example request:
```json
{
  "name": "test"
}
```

######Example response:
```json
{
  "name": "test"
}
```

### /config/getall

This operation retrieves all configuration items. 

This is a `GET` request.

######Example response:
```json
{
  "status": 200,
  "message": "Config items found",
  "data": [
    {
      "id": 2,
      "application": "MyApp",
      "machine": "",
      "name": "Another",
      "value": "Value1",
      "updated": "2016-06-09T15:57:07.1052893-04:00"
    },
    {
      "id": 1,
      "application": "MyApp",
      "machine": "",
      "name": "Environment",
      "value": "DEV",
      "updated": "2016-06-09T15:56:51.4547244-04:00"
    },
    {
      "id": 4,
      "application": "SomeOtherAppEntirely",
      "machine": "",
      "name": "TheName",
      "value": "TheValue",
      "updated": "2016-06-09T15:59:53.6649436-04:00"
    }
  ]
}
```

### /config/getallforapp

This operation retrieves all configuration items for a specified application

######Example request:
```json
{
  "name": "test"
}
```

######Example response:
```json
{
  "name": "test"
}
```

### /applications/getall

This operation retrieves all applications

######Example request:
```json
{
  "name": "test"
}
```

######Example response:
```json
{
  "name": "test"
}
```
