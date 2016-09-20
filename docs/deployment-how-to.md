# How to Deploy Config Server

### Snippet from a Deployment Manifest 

```
    jobs:
    - name: config_server
      properties:
        port: 8080
        store: "memory" 
        jwt:
          verification_key: |
            -----BEGIN PUBLIC KEY-----
            {{Replace with public key that will validate the JWT token}}
            -----END PUBLIC KEY-----
        ca:
          certificate: |
            -----BEGIN CERTIFICATE-----
            {{Replace with the CA cert that will be used to sign generated certificates}}
            -----END CERTIFICATE-----
          private_key: |
            -----BEGIN RSA PRIVATE KEY-----
            {{Replace with the private key that will be used to sign generated certificates}}
            -----END RSA PRIVATE KEY-----
        ssl:
          certificate: |
            -----BEGIN CERTIFICATE-----
            {{Replace with the certificate that will be used for TLS connections}}
            -----END CERTIFICATE-----
          private_key: |
            -----BEGIN RSA PRIVATE KEY-----
            {{Replace with the private key that will be used for TLS connections}}
            -----END RSA PRIVATE KEY-----

```
