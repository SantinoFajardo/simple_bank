<h1>Quick doc to explain how test the create user endpoint</h1>
This can be harder than other endpoints since we need handle grpc call and the background worker.
<h2>Testing simple gRPC endpoints</h2>

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

<h2>Testing gRPC endpoints that need</h2>

On the `testCases` struct we need add a new prop `buildContext` that will contain the logic to include the authorization permission to test the endpoint
Doing this we can setup differentes `context metadata`, would be helpful if we want to test other cases, for example `unauthorized` cases

Now we can use the returned context of this function inside the testing call functions to test the authentication flow on the endpoint:

```go
ctx := tc.buildContext(t, server.tokenMaker)
res, err := server.UpdateUser(ctx, tc.body)
```

Here is a example of how we can build a success context to test the happy cases.

```go
buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
	accessToken, _, err := tokenMaker.CreateToken(user.Username, time.Minute)
	require.NoError(t, err)
	bearerToken := fmt.Sprintf("%s %s", authorizatioBearerKey, accessToken)
	md := metadata.MD{
		authorizationHeader: []string{bearerToken},
	}
	return metadata.NewIncomingContext(context.Background(), md)
},
```
