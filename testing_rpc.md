<h1>Quick doc to explain how test the create user endpoint</h1>
This can be harder than other endpoints since we need handle grpc call and the background worker.
Inside the `rpc_create_user_test.go` will be located the test
In order to test the workers and avoid use the real workers of redis, we can mock the workers as we did with the database.
On the `Makefile` on the `mock` command we can add the next command to mock the workers interface:
<br/><br/>

```makefile
mockgen -package mockwk -destination worker/mocks/distributor.go github.com/santinofajardo/simpleBank/workers TaskDistributor
```

In order to call the mocked gRPC server we can use this calls:

```go
tc.buildStubs(store, taskDistributor)
server := newTestServer(t, store, taskDistributor)

res, err := server.CreateUser(context.Background(), tc.req)
tc.checkResponse(t, res, err)
```
