# Config Server


## API

Following public API will be used by the Director to contact config server:

- GET /v1/data/&lt;some-key-path>
  - whenever Director needs to retrieve a value it will use GET action
  - {"path": "some-key-path", "value": "..."}

- PUT /v1/data/&lt;some-key-path>
  - whenever Director updates a value in the config server / sets it for the first time
  - {"value": "..."}

- POST /v1/data/&lt;some-key-path>
  - manually update a value to be saved into config server
  - They "type" parameter will determine what is generated (currently `password` and `certificate` are accepted
  - {"type": "password"}
  - {"type": "certificate", "parameters": {"common_name": "bosh.io", "alternative_names": ["blah.bosh.io", "10.0.0.6"]}}
  - See [Bosh Notes](https://github.com/cloudfoundry/bosh-notes/blob/master/config-server.md) for more information

Values could be any valid JSON object.

## Sample Request/Responses:
### POST

#### Generate Certificate

- Request:

```
curl -k -X POST https://localhost:8080/v1/data/b -d '{"type":"certificate", "parameters": {"common_name": "asdf", "alternative_names":["nam1", "name2"]}}'
```

- Response:

Status **201** for first time generation.
Status **200** for subsequent requests.

```
{
 "path":"b",
 "value": {
   "ca" : "---- Root CA Certificate ----",
   "certificate": "---- Generated Certificate. Signed by rootCA ----",
   "private_key": "---- Private key for the Generated certificate ----"
  }
}
```



#### Generate Password

- Request

```
curl -k -X POST https://localhost:8080/v1/data/passwd -d '{"type":"password"}'
```

- Response:

Status **201** for first time generation.
Status **200** for subsequent requests.

```
{
  "path": "passwd",
  "value":"49cek4ow75ev5zw4t3v3"
}
```


### PUT
- Request:
```curl -k -X PUT https://localhost:8080/v1/data/b -d '{"value":"blah"}'```

- Response:
  - Status 204
  - No Body

### GET
- Request (existing key):
```curl -kv -X GET https://localhost:8080/v1/data/b```

- Response:
```{"path":"b","value":"blah"}```

- Request (key does not exist):
```curl -kv -X GET https://localhost:8080/v1/data/derp```

- Response:
  - Status: 404 Not Found
