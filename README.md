# Config Server

- CI: <https://main.bosh-ci.cf-app.com/teams/main/pipelines/config-server>
- [API Docs](docs/api.md)

See [config-server-release](https://github.com/cloudfoundry/config-server-release) for the release repo for config-server
See [bosh-notes](https://github.com/cloudfoundry/bosh-notes/blob/master/proposals/config-server.md) for more information

## Contributing

All PRs and code should be against the `develop` branch.  Before submitting PRs, you should make sure that tests are passing properly.

### Unit tests

Before running the unit test on your local machine, you should install
`golangci-lint` following its documentation, see
<https://github.com/golangci/golangci-lint>. (Indeed, the way the
[lint](bin/lint) script installs `golangci-lint` is not suitable for running
on a local machine.)

Then run the unit tests using the dedicated script. 

```bash
bin/test-unit
```

### Integration tests

There are 3 distinct flavors of integration tests. Run the relevant one, or
all, depending on the changes you've done.

```bash
bin/test-integration memory     # <- uses an in-memory database
bin/test-integration mysql      # <- uses a MySQL database
bin/test-integration postgresql # <- uses a PostgreSQL database
```
