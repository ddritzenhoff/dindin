**TODO**
- Improve logging functionality
- Improve error-handling
- Dockerize
- 'Who's eating on Xday' post:
    - The post is live for 24 hours.
    - If someone likes the post, check to see if that person is within the database. If no, add him. If yes, add meals eaten.
- 'Cooking next week' post:
    - publishes a list with the worst meals eaten to meals cooked ratios.
        - If 0 meals eaten and 0 meals cooked, don't include.
        - If >0 meals eaten and 0 meals cooked, be at the top.
        - If >0 meals eaten and >0 meals cooked, include the highest ratios first.
- Assign next week cooks
    - Pass in arguments in the order of day-of-week for a min of 1 arg and a max of 6 args.
- Get upcoming cooks.

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
