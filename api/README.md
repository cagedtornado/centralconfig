# API server information

API path              | Description
----------            | -----------
/                     | Redirects to the /UI path (and the web based user interface)
/config/get           | Gets a single configuration item
/config/set           | Sets (creates or updates) a configuration item
/config/remove        | Removes a configuration item
/config/getall        | Gets all configuration items
/config/getallforapp  | Get all configuration items for a single application (plus the default * application)
/applications/getall  | Get all applications

######Requests
Most API operations expect a configitem object in the POST body that will be used to either filter (in a get operation), update or create (in a set operation), or remove an item (in a remove operation).  

For example:

```json
{
    "application" : "WickedCool",
    "name": "TestItem42",
    "value": "Magic!"
}
```

######Responses
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


#### /config/get

This operation retrieves a single configuration item

Example request:
```json
{
  "name": "test"
}
```


