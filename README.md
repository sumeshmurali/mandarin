# mandarin
Free to use fake/mock api server suite for developers, product managers and for load testing

**Note: Under development. Not ready for production use yet**

## Example Configuration
```yaml
name: Mock Server
description: Mock server for testing
config:
  port: 80
endpoints:
  /:
    description: Root endpoint
    request_config:
      allowed_methods:
        - GET
        - POST
    response_config:
      status_code: 200
      headers:
        Content-Type: application/json
      body: "{\"message\": \"Hello, World!\"}"
    ratelimit_config:
      ratelimit: 1 # requests per second
      ratelimit_type: global
```

# Docker Usage

1. Clone the repo
2. Run `docker build . -t mandarin`
3. Create a folder called `config` and copy the example config file above to the folder
4. Run the docker image using `docker run --volume ./config/:/mandarin/ -p 8080:80 mandarin:latest`