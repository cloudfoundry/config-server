## Sample Request / Responses

### Generate Certificate

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



### Generate Password

- Request

    ```
curl -k -X POST https://localhost:8080/v1/data/passwd -d '{"type":"password"}'
```

- Response:

    Status **201** for first time request. Status **200** for subsequent requests.

    ```
{
    "path": "passwd",
    "value":"49cek4ow75ev5zw4t3v3"
}
```


### Set/Update Key
- Request
 
    ```
curl -k -X PUT https://localhost:8080/v1/data/sample-key -d '{"value":"blah"}'
```

- Response:
  - Status **204**
  - No Body

### Get Value for existing key
- Request
 
    ```
curl -kv -X GET https://localhost:8080/v1/data/sample-key
```

- Response

    ```
{"path":"sample-key","value":"blah"}
```

### Get Value for key that does not exist
- Request

    ```
curl -kv -X GET https://localhost:8080/v1/data/derp
```

- Response
  - Status **404** Not Found
  - No Body 