# Wiring Checklist for New Features

When adding new services, handlers, or routes, use this checklist to ensure complete wiring.

## New Service Checklist

When creating a new service (e.g., `AdminService`):

- [ ] **Define interface** in `service/service.go` or dedicated file
- [ ] **Create implementation** in `service/{name}.go`
- [ ] **Add to Services struct** in `service/service.go`:
  ```go
  type Services struct {
      // ...existing services...
      NewService NewServiceInterface  // <- ADD THIS
  }
  ```
- [ ] **Initialize in main.go** (or NewServices if no external deps):
  ```go
  services.NewService = service.NewNewService(repo, otherDeps)
  ```
- [ ] **Add required config** to `config.go` if service needs env vars
- [ ] **Add env vars to Lambda** in `infrastructure/backend/lambda-api.tf`
- [ ] **Write unit tests** for service methods
- [ ] **Write integration test** that verifies service is accessible

## New Handler Checklist

When creating a new handler:

- [ ] **Create handler struct** with service dependency
- [ ] **Create constructor** `NewXxxHandler(service XxxService)`
- [ ] **Implement handler methods** for each endpoint
- [ ] **Create handler instance in main.go**:
  ```go
  xxxHandler := handlers.NewXxxHandler(services.Xxx)
  ```
- [ ] **Write unit tests** with mocked service

## New Route Checklist

When adding new routes:

- [ ] **Create route registration function** (e.g., `RegisterXxxRoutes`)
- [ ] **Call registration function in main.go**:
  ```go
  handlers.RegisterXxxRoutes(e, xxxHandler)
  ```
- [ ] **Add route to API Gateway** in `infrastructure/backend/api-gateway.tf`
- [ ] **Write smoke test** that verifies route returns expected status:
  ```go
  // integration_test.go
  func TestAdminRoutes_Accessible(t *testing.T) {
      resp := httptest.NewRecorder()
      req := httptest.NewRequest("GET", "/api/v1/admin/users?search=test", nil)
      req.Header.Set("Authorization", "Bearer "+adminToken)

      router.ServeHTTP(resp, req)

      // Should NOT be 404 - route must be registered
      assert.NotEqual(t, 404, resp.Code, "Route not registered!")
  }
  ```

## New Frontend Route Checklist

When adding new frontend pages:

- [ ] **Create page component** in `routes/path/page.tsx`
- [ ] **Check routing pattern** - is it file-based or code-based?
  - Look at `main.tsx` - if routes are imported there, it's code-based
- [ ] **For code-based routing (this project):**
  - [ ] Import page component in `main.tsx`
  - [ ] Create route with `createRoute()`
  - [ ] Add to `routeTree.addChildren([...])`
- [ ] **Add navigation link** if needed (sidebar, header, etc.)

## Verification Steps

After completing all checklists:

1. **Build succeeds**: `go build ./...`
2. **Unit tests pass**: `go test ./...`
3. **Integration test passes**: Route returns expected status (not 404)
4. **Manual smoke test**: Can access endpoint via curl or browser

## Why This Matters

Unit tests with mocks can pass even when:
- Service isn't added to Services struct
- Handler isn't created in main.go
- Routes aren't registered
- Env vars are missing from Lambda

Only integration/smoke tests that hit actual endpoints catch these wiring issues.
