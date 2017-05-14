    
1. Create: POST /v1-domains/domains 
    - input `{domainName: "foo.com",projectid: "1a1"}` I will get the accountid from token
    - output `{domainName: "foo.com",domainName: "foo.com"}`
    

1. List: GET /v1-domains/domains
1. Get by ID: GET /v1-domains/domains/:id
1. Have the user retry validation: POST /v1-domains/domain/:id?action=validate
1. Delete: DELETE /v1-domains/domains/id
1. Errors have specific fields, {type: "error", status: 422, code: "DomainAlreadyInUse", message: "domain.com has already been validated by another account"}


- Error Message `{Type: "error",status: "401", message: "error message"}`    