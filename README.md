# Config Server


## API

Following public API will be used by the Director to contact config server:

- GET /v1/data/&lt;some-key-path>
  - whenever Director needs to retrieve a value it will use GET action
  - {"path": "some-key-path", "value": "..."}

- PUT /v1/data/&lt;some-key-path>
  - whenever Director generates a value it will be saved into the config server
  - {"value": "..."}

Values could be any valid JSON object.
