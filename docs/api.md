## API Docs (WIP)
This document describes the APIs exposed by the **Config Server**.

### 1 - Get Name Value
```
GET /v1/data/:name
```

| Name | Description |
| ---- | ----------- |
| name | Full path |

##### Response Body
`Content-Type: application/json`

| Name | Type | Description |
| ---- | ---- | ----------- |
| id | string | Unique Id |
| path | string | Full path |
| value | JSON Object | Any valid JSON object |

##### Response Codes
| Code   | Description |
| ------ | ----------- |
| 200 | Status OK |
| 400 | Bad Request |
| 401 | Not Authorized |
| 404 | Name not found |
| 500 | Server Error |

##### Sample Requests/Responses

`GET /v1/data/color`

Response:
``` JSON
{
  "id": "some_id",
  "name": "color",
  "value": "blue"
}
```

--

`GET /v1/data/server/tomcat/port`

Response:
``` JSON
{
  "id": "some_id",
  "name": "server/tomcat/port",
  "value": 8080
}
```

--

`GET /v1/data/server/tomcat/cert`

Response:
``` JSON
{
  "id": "some_id",
  "name": "server/tomcat/port",
  "value": {
    "cert": "my-cert",
    "private-key": "my private key"
  }
}
```

---

### 2 - Set Name Value
```
PUT /v1/data/:name
```

| Name | Description |
| ---- | ----------- |
| name | Full path |

##### Request Body
`Content-Type: application/json`

| Name | Type | Description |
| ---- | ---- | ----------- |
| value | JSON Object | Any valid JSON object |

##### Sample Request

`PUT /v1/data/full/path/to/name`

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
| path | string | Full path |
| value | JSON Object | Any valid JSON object |

##### Response Codes
| Code | Description |
| ---- | ----------- |
| 200 | Call successful - name value was added |
| 400 | Bad Request |
| 401 | Not Authorized |
| 415 | Unsupported Media Type |
| 500 | Server Error |

---

### 3 - Generate password/certificate

```
POST /v1/data/:name
```

| Name | Description |
| ---- | ----------- |
| name | Full path |

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
| name | string | Full path  |
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
  "name": "/mypasswd",
  "value":"49cek4ow75ev5zw4t3v3"
}
```
###### Certificate
``` 
{
  "id": "some_id",
  "name":"/mycert",
  "value": {
    "ca" : "---- Root CA Certificate ----",
    "certificate": "---- Generated Certificate. Signed by rootCA ----",
    "private_key": "---- Private key for the Generated certificate ----"
  }
}
```

### 4 - Delete Name
```
DELETE /v1/data/:name
```

| Name | Description |
| ---- | ----------- |
| name | Full path |


##### Sample Request

`DELETE /v1/data/full/path/to/name`

##### Response Codes
| Code | Description |
| ---- | ----------- |
| 204 | Call successful - name was deleted |
| 400 | Bad Request |
| 401 | Not Authorized |
| 404 | Not Found |
| 500 | Server Error |
