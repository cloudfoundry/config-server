## API Docs (WIP)
This document describes the APIs exposed by the **Config Server**.

### 1 - Get Key Value
```
GET /v1/data/:key-path
```

| Name | Description |
| ---- | ----------- |
| key-path | Full path to key |

##### Response Body
`Content-Type: application/json`

| Name | Type | Description |
| ---- | ---- | ----------- |
| id | string | Unique Id |
| path | string | Full path to key |
| value | JSON Object | Any valid JSON object |

##### Response Codes
| Code   | Description |
| ------ | ----------- |
| 200 | Status OK |
| 400 | Bad Request |
| 401 | Not Authorized |
| 404 | Key not found |
| 500 | Server Error |

##### Sample Requests/Responses

`GET /v1/data/color`

Response:
``` JSON
{
  "id": "some_id",
  "path": "color",
  "value": "blue"
}
```

--

`GET /v1/data/server/tomcat/port`

Response:
``` JSON
{
  "id": "some_id",
  "path": "server/tomcat/port",
  "value": 8080
}
```

--

`GET /v1/data/server/tomcat/cert`

Response:
``` JSON
{
  "id": "some_id",
  "path": "server/tomcat/port",
  "value": {
    "cert": "my-cert",
    "private-key": "my private key"
  }
}
```

---

### 2 - Set Key Value
```
PUT /v1/data/:key-path
```

| Name | Description |
| ---- | ----------- |
| key-path | Full path to key |

##### Request Body
`Content-Type: application/json`

| Name | Type | Description |
| ---- | ---- | ----------- |
| value | JSON Object | Any valid JSON object |

##### Sample Request

`PUT /v1/data/full/path/to/key`

Request Body:
```
{
  "value": "happy value"
}
```

##### Response Body
`Content-Type: application/json`

| Name | Type | Description |
| ---- | ---- | ----------- |
| id | string | Unique Id |
| path | string | Full path to key |
| value | JSON Object | Any valid JSON object |

##### Response Codes
| Code | Description |
| ---- | ----------- |
| 200 | Call successful - key value was added |
| 400 | Bad Request |
| 401 | Not Authorized |
| 415 | Unsupported Media Type |
| 500 | Server Error |

---

### 3 - Generate password/certificate

```
POST /v1/data/:key-path
```

| Name | Description |
| ---- | ----------- |
| key-path | Full path to key |

##### Request Body
`Content-Type: application/json`

| Name | Type | Valid Values | Description |
| ---- | ---- | ------------ | ----------- |
| type | String | password, certificate | The type of data to generate |
| parameters | JSON Object | | See below for valid parameters |

###### Request body extra parameters values
| For type | Name | Type |
| -------- | ---- | ---- |
| certificate | common_name | String |
| certificate | alternative_names | Array of Strings |

##### Sample Requests

###### Password Generation
`POST /v1/data/mypasswd`

Request Body:
``` JSON
{
  "type": "password"
}
```

###### Certificate Generation
`POST /v1/data/mycert`

Request Body:
``` JSON
{
  "type": "certificate",
  "parameters": {
    "common_name": "bosh.io",
    "alternative_names": ["bosh.io", "blah.bosh.io"]
  }
}
```

##### Response Body
`Content-Type: application/json`

It returns an array of the following object:

| Name | Type | Description |
| ---- | ---- | ----------- |
| id | string | Unique Id |
| path | string | key  |
| value | JSON Object | value generated |

##### Response Codes
| Code | Description |
| ---- | ----------- |
| 201 | Call successful |
| 400 | Bad Request |
| 401 | Not Authorized |
| 415 | Unsupported Media Type |
| 500 | Server Error |

##### Sample Response
###### Password
```
{
  "id": "some_id",
  "path": "/mypasswd",
  "value":"49cek4ow75ev5zw4t3v3"
}
```
###### Certificate
``` 
{
  "id": "some_id",
  "path":"/mycert",
  "value": {
    "ca" : "---- Root CA Certificate ----",
    "certificate": "---- Generated Certificate. Signed by rootCA ----",
    "private_key": "---- Private key for the Generated certificate ----"
  }
}
```

### 4 - Delete Key
```
DELETE /v1/data/:key-path
```

| Name | Description |
| ---- | ----------- |
| key-path | Full path to key |


##### Sample Request

`DELETE /v1/data/full/path/to/key`

##### Response Codes
| Code | Description |
| ---- | ----------- |
| 204 | Call successful - key was deleted |
| 400 | Bad Request |
| 401 | Not Authorized |
| 404 | Not Found |
| 500 | Server Error |
