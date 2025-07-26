# Authentication Service

Authentication service should include authentication verification that include JWT token 

## Functionality

- Service can receive and process the request of making new user
- Admin can change the role of user
- Service can receive login credentails and issue JTW token

<!-- Note:
Client should send loging credentials to authentication service and service should provide JTW token.

Auth service should validate the token. ()

POST request for getting token credentials  - service respond and send back a JWT token and 200 status

Authentization: Beareer <JWT> - json body

With JTW token presnet in request in broker service, broker service should check the JWT token validity. -->