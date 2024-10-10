# Todo

## Basic functional requirements

- [x] The APIs should be RESTful
- [x] The admin endpoints should be authenticated
  - Using simple authentication via API key
- [x] Invite tokens expire after 7 days
- [x] A public endpoint for validating invite token
  - Redeem endpoint

### Nice to have (functional)

- [x] The invite token validatation logic needs to be throttled (limit the requests coming from a specific client)
  - Using a simple rate-limit middleware
- [x] An admin can get an overview of active and inactive tokens
  - Not sure what active/inactive means, but we can list the tokens at `/admin/tokens`

## Basic non-functional requirements

- [x] Design and document the APIs that will facilitate the workflow outlined above
- [x] Develop the API in Go
- [x] Use any framework or library that will help you develop the solution faster
  - Using labstack/echo for routing
- [x] Make sure your code is well-formatted, clean and follow best practices
- [x] Separate concerns
- [x] Write testable code
- [x] Use in-memory storage for tokens

### Nice to have (non-fuctional)

- [x] Document the APIs in swagger or similar tool
- [ ] Write functional code
  - Not sure what this means
- [x] Test, all levels of them
- [x] Use an actual DB
  - Using postgres
- [x] Provide deployment instructions
