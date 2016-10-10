## Config Server



Job            | Status  
:------------: | -------------
Unit tests     | ![Resize icon][badge-test-unit] 
Integration tests  | ![Resize icon][badge-test-integration]
BOSH Release   | ![Resize icon][badge-bosh-release]
Go Formatting| ![Resize icon][badge-test-gofmt]   

[badge-test-unit]: https://main.bosh-ci.cf-app.com/api/v1/teams/main/pipelines/config-server/jobs/test-unit/badge

[badge-test-integration]: https://main.bosh-ci.cf-app.com/api/v1/teams/main/pipelines/config-server/jobs/test-integration/badge

[badge-bosh-release]: https://main.bosh-ci.cf-app.com/api/v1/teams/main/pipelines/config-server/jobs/bosh-release/badge

[badge-test-gofmt]: https://main.bosh-ci.cf-app.com/api/v1/teams/main/pipelines/config-server/jobs/test-gofmt/badge

__CI :__ <https://main.bosh-ci.cf-app.com/teams/main/pipelines/config-server> <br>


## API 

- GET /v1/data/&lt;some-key-path>
  
  Used by Director to retrieve a value from the config server. Value can an be any valid JSON object.
    
	```
  {"path": "some-key-path", "value": "..."}  
  ```

- PUT /v1/data/&lt;some-key-path>

  Used by Director to set or update a value in the config server. Value can an be any valid JSON object.
  ```
  {"value": "..."}
  ```

- POST /v1/data/&lt;some-key-path>

	Used by Director to generate a value in the config server. The `type` parameter will determine what is generated. Currently `password` and `certificate` are accepted.
  
  ```
  {"type": "password"}
  ```
  ```
  {"type": "certificate", "parameters": {"common_name": "bosh.io", "alternative_names": ["blah.bosh.io", "10.0.0.6"]}}
  ```
  	
See [Sample Requests & Responses](https://github.com/cloudfoundry/config-server/blob/master/docs/sample-requests-responses.md)

See [Bosh Notes](https://github.com/cloudfoundry/bosh-notes/blob/master/config-server.md) for more information
