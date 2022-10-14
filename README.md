**TODO**
- Improve logging functionality
- Dockerize

Sample config.yaml:
```
http:
  rest:
    host: "localhost"
    port: 7777
  grpc:
    host: "localhost"
    port: 7778
database:
  name: "dindin.db"

slack:
  botSigningKey: ""
  appID: ""
  clientID: ""
  clientSecret: ""
  signingSecret: ""
  dev:
    channelID: ""
  prod:
    channelID: ""
```
