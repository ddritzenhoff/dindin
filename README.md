This is a slack bot to organize a dinner rotation between my friends and me :)

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
  isProd: false
  dev:
    channelID: ""
  prod:
    channelID: ""
```
