# Code Development Conventions

> **Purpose:** This document defines coding standards and best practices for all developers and AI assistants working on the go-chat project.

---

## 🎯 Core Principles

1. **Iterative Development** — Develop in small sequential steps, gradually improving functionality. Don't try to do everything at once.
2. **KISS (Keep It Simple, Stupid)** — If there is a simple, effective solution, use it rather than overengineering.
3. **Fail Fast** — Errors should be caught as early as possible. Handle errors explicitly; don't hide them.
4. **Readability First** — Code readability is paramount. Write maintainable, self-documenting code.
5. **Explicit is Better Than Implicit** — Be clear about intentions and dependencies.
6. **Complex is Better Than Complicated** — Complexity from requirements is acceptable; unnecessary complication is not.
7. **Performance Matters, But Not at All Costs** — Don't sacrifice readability or maintainability for premature optimization.
8. **Special Cases Aren't Special** — Consistency trumps convenience. Avoid one-off solutions.
9. **Single Technology Stack** — Don't add new dependencies if existing ones solve the problem.
10. **Now is Better Than Never** — Ship iteratively, but avoid rushing into bad decisions.
11. **If Hard to Explain, It's a Bad Idea** — Simple solutions are easier to maintain and debug.

---

## 🔧 Go-Specific Conventions

### Naming
- **Variables/Functions:** Use `camelCase` for private, `PascalCase` for public/exported
- **Packages:** Short, lowercase, single-word names (e.g., `auth`, `user`, `chat`)
- **Interfaces:** Name by behavior, often with `-er` suffix (e.g., `Repository`, `Validator`, `Handler`)
- **Acronyms:** Keep consistent case (e.g., `userID`, `httpClient`, `URL`)

### Error Handling
- **Always check errors** — Never ignore returned errors (`if err != nil`)
- **Wrap errors with context** — Use `fmt.Errorf("context: %w", err)` to add context while preserving the original error
- **Return early** — Handle errors immediately and return early to avoid deep nesting
- **Don't panic** — Reserve `panic` only for truly unrecoverable errors (e.g., startup configuration failures)
- **Custom errors** — Use custom error types for domain-specific errors that need special handling

### Code Structure
- **Package organization:**
  ```
  service-name/
    ├── cmd/              # Application entry points
    ├── internal/         # Private application code
    │   ├── domain/       # Business entities and logic
    │   ├── repository/   # Database access layer
    │   ├── service/      # Business logic services
    │   └── handler/      # gRPC handlers
    ├── pkg/              # Public libraries (if needed)
    └── proto/            # Protocol definitions (.proto files for API and events)
  ```
- **Small functions** — Functions should do one thing well (ideally < 50 lines)
- **Small interfaces** — Prefer many small interfaces over large ones (1-3 methods)
- **Dependency injection** — Pass dependencies explicitly via constructors, avoid global state

### Concurrency
- **Use channels for communication** — Don't communicate by sharing memory; share memory by communicating
- **Always handle goroutine lifecycle** — Use `context.Context` for cancellation, ensure goroutines can be stopped
- **Use `sync.WaitGroup`** — For waiting on multiple goroutines to complete
- **Avoid goroutine leaks** — Always ensure goroutines terminate (especially with infinite loops)

### Comments
- **Public APIs must have godoc comments** — Start with the name of what you're documenting
- **Explain "why", not "what"** — Code should be self-explanatory; comments explain reasoning
- **Keep comments up-to-date** — Outdated comments are worse than no comments
- **TODO comments** — Format: `// TODO(username): description` for tracking pending work

---

## 🏗️ Microservices Conventions

### Service Communication
- **gRPC for all inter-service communication** — No HTTP between services
- **Service tokens in metadata** — Include service identity in `grpc.metadata`
- **Timeout contexts** — Always set request timeouts (e.g., 5-30 seconds)
- **Retry with backoff** — Implement exponential backoff for transient failures
- **Circuit breakers** — Prevent cascading failures (consider using libraries like `gobreaker`)

### API Design
- **Versioning required** — All protobuf packages must be versioned (`api.v1`, `api.v2`)
- **Backwards compatibility** — Never break existing clients; add new fields, deprecate old ones
- **Cursor-based pagination** — Use cursor pagination for list operations (not limit-offset)
- **Idempotency keys** — Support idempotency for write operations when possible
- **Standard error codes** — Use gRPC status codes consistently (e.g., `NOT_FOUND`, `PERMISSION_DENIED`)

### Data Management
- **Database per service** — Each service owns its data; no direct database access across services
- **Migrations as code** — Use migration tools (e.g., `goose`)
- **Migrations managed separately** — Migrations should NOT be part of application code; manage them in a separate directory structure with dedicated deployment pipeline
- **Atomic migrations** — Each migration should be reversible and atomic
- **No shared databases** — Services communicate via APIs, not shared databases

### Observability
- **Structured logging** — Use JSON-formatted logs with consistent fields
  - Required fields: `timestamp`, `level`, `service`, `trace_id`, `message`
- **Distributed tracing** — Use OpenTelemetry for request tracing across services
- **Metrics** — Expose Prometheus metrics (requests, latency, errors)
- **Request ID propagation** — Pass `trace_id` through gRPC metadata and log it

### Event-Driven
- **Kafka for async events** — Use Kafka for event-driven communication
- **Event schema versioning** — Version event schemas, maintain compatibility
- **Idempotent consumers** — Kafka consumers must handle duplicate messages
- **Dead letter queues** — Route failed messages to DLQ for investigation

---

## ✅ Code Quality Standards

### Linting
- **golangci-lint required** — All code must pass `golangci-lint run`
- **Pre-commit hooks** — Run linters before commits
- **No warnings allowed** — Code must be warning-free before merge

### Testing
- **Unit tests required** — All business logic must have unit tests
- **Test coverage:** Aim for **80%+ coverage** for critical paths
- **Table-driven tests** — Use table-driven tests for multiple scenarios
- **Integration tests** — Each service should have integration tests for gRPC endpoints
- **Test naming:** `TestFunctionName_Scenario_ExpectedResult`
  ```go
  func TestCreateUser_ValidInput_ReturnsUser(t *testing.T) { ... }
  func TestCreateUser_DuplicateEmail_ReturnsError(t *testing.T) { ... }
  ```
- **Mocking:** Use interfaces for dependencies, mock external services

### Code Review
- **All changes require review** — No direct commits to `main`
- **Review checklist:**
  - [ ] Code follows conventions
  - [ ] Tests included and passing
  - [ ] Error handling is correct
  - [ ] No unnecessary complexity
  - [ ] Documentation updated (if needed)
  - [ ] No security vulnerabilities
- **Small PRs preferred** — Keep PRs focused and reviewable (< 1000 lines when possible)

---

## 📝 Git Conventions

### Commit Messages
Follow **Conventional Commits** format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `refactor`: Code refactoring (no functional changes)
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `chore`: Maintenance tasks (dependencies, configs)
- `perf`: Performance improvements

**Examples:**
```
feat(auth): add JWT token refresh endpoint

Implements refresh token rotation with RS256 signing.
Tokens expire after 15 minutes and refresh tokens after 30 days.

Closes #42
```

```
fix(chat): prevent duplicate messages with idempotency keys

Added idempotency_key field to SendMessage RPC to prevent
duplicate message creation on retries.
```

### Branch Naming
```
<type>/<short-description>
```

**Examples:**
- `feat/user-profile-search`
- `fix/auth-token-expiration`
- `refactor/chat-service-repository`
- `docs/update-api-contracts`

### Pull Requests
- **Title:** Use commit message format
- **Description template:**
  ```markdown
  ## What
  Brief description of changes
  
  ## Why
  Problem being solved or feature being added
  
  ## How
  Technical approach
  
  ## Testing
  How to test these changes
  
  ## Checklist
  - [ ] Tests added/updated
  - [ ] Documentation updated
  - [ ] Linter passing
  - [ ] No breaking changes (or documented)
  ```

---

## 🔒 Security Best Practices

1. **Never commit secrets** — Use environment variables or secret management systems
2. **Validate all inputs** — Sanitize and validate user inputs at service boundaries
3. **Use parameterized queries** — Prevent SQL injection
4. **Rate limiting** — Implement rate limiting at Gateway level
5. **Principle of least privilege** — Services should only have necessary permissions
6. **Audit logging** — Log security-relevant events (login, permission changes)

---

## 🚀 Deployment & Operations

1. **12-Factor App principles** — Follow 12-factor methodology
2. **Health checks** — Every service must expose `/health` or gRPC health check
3. **Graceful shutdown** — Handle SIGTERM gracefully, drain connections
4. **Configuration via environment** — No hardcoded configs in code
5. **Docker for containerization** — Each service has its own Dockerfile
6. **Resource limits** — Set memory and CPU limits for containers

---

## 📚 Additional Resources

- [Effective Go](https://go.dev/doc/effective_go)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [gRPC Best Practices](https://grpc.io/docs/guides/performance/)
- [Microservices Patterns](https://microservices.io/patterns/)
