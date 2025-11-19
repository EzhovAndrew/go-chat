# Service Layer and Repository Interface Implementation Plan

## Overview

Implement the complete service layer for all microservices, including:
- Repository interface definitions for data access abstraction
- Concrete service implementations with business logic
- Cross-service gRPC client interfaces
- Proper dependency injection, error handling, and validation
- Full testability via interface-based design

## Current State Analysis

**Completed:**
- Domain entities and value objects
- Domain errors
- Service interfaces
- Handler implementations (using service interfaces)

**To Be Implemented:**
- Repository interfaces (data access layer abstraction)
- Concrete service implementations
- Cross-service client interfaces (for gRPC calls)
- Unit tests for service implementations

## Architecture & Design Principles

### Layer Responsibilities

```
Handler Layer (exists)
    ↓ depends on
Service Layer (to implement)
    ↓ depends on
Repository Layer (interfaces only, to define)
    ↓ will depend on
Database (future iteration)
```

### Key Design Decisions

1. **Repository Pattern**: Encapsulate all data access behind repository interfaces
2. **Dependency Injection**: Services receive dependencies via constructors
3. **Interface-First**: Define interfaces for testability (mock repositories, mock external services)
4. **Error Translation**: Repositories return domain errors, services add business logic validation
5. **Context Propagation**: All methods accept `context.Context` for cancellation and tracing
6. **Transaction Management**: Repository interfaces will expose transaction primitives when needed

## Implementation Structure by Service

### 1. Auth Service

**Repository Interfaces to Create:**

`internal/repository/user_repository.go`:
```go
type UserRepository interface {
    // Create stores a new user
    Create(ctx context.Context, user *domain.User) error
    
    // GetByID retrieves a user by ID
    GetByID(ctx context.Context, userID domain.UserID) (*domain.User, error)
    
    // GetByEmail retrieves a user by email
    GetByEmail(ctx context.Context, email string) (*domain.User, error)
    
    // EmailExists checks if an email is already registered
    EmailExists(ctx context.Context, email string) (bool, error)
}
```

`internal/repository/refresh_token_repository.go`:
```go
type RefreshTokenRepository interface {
    // Create stores a new refresh token
    Create(ctx context.Context, token *domain.RefreshToken) error
    
    // GetByToken retrieves a refresh token by its value
    GetByToken(ctx context.Context, token string) (*domain.RefreshToken, error)
    
    // Revoke marks a refresh token as revoked
    Revoke(ctx context.Context, tokenID string) error
    
    // DeleteExpired removes all expired tokens for a user
    DeleteExpired(ctx context.Context, userID domain.UserID) error
}
```

**Service Implementations to Create:**

`internal/service/auth_service_impl.go`:
```go
type authService struct {
    userRepo         repository.UserRepository
    refreshTokenRepo repository.RefreshTokenRepository
    tokenService     TokenService  // For generating tokens
}

func NewAuthService(
    userRepo repository.UserRepository,
    refreshTokenRepo repository.RefreshTokenRepository,
    tokenService TokenService,
) AuthService {
    return &authService{...}
}

// Implements:
// - Register: validate email format, check uniqueness, hash password, create user
// - Login: fetch user by email, compare password, generate tokens, store refresh token
// - Refresh: validate token, check expiration/revocation, rotate tokens
```

`internal/service/token_service_impl.go`:
```go
type tokenService struct {
    privateKey *rsa.PrivateKey
    publicKey  *rsa.PublicKey
    keyID      string
}

func NewTokenService(privateKeyPath string) (TokenService, error) {
    // Load RSA keys
}

// Implements:
// - GenerateTokenPair: create JWT access token (15 min) + refresh token (30 days)
// - GetPublicKeys: return JWK format public keys for Gateway validation
// - ValidateAccessToken: parse and verify JWT signature
```

**Key Considerations:**
- Password hashing uses existing `utils.HashPassword` and `utils.ComparePassword`
- Token generation needs RSA key management (load from files or env)
- Refresh tokens stored hashed in database for security
- Token rotation: old refresh token invalidated, new one issued

### 2. Users Service

**Repository Interfaces to Create:**

`internal/repository/profile_repository.go`:
```go
type ProfileRepository interface {
    // Create stores a new user profile
    Create(ctx context.Context, profile *domain.UserProfile) error
    
    // Update modifies an existing profile
    Update(ctx context.Context, profile *domain.UserProfile) error
    
    // GetByID retrieves a profile by user ID
    GetByID(ctx context.Context, userID domain.UserID) (*domain.UserProfile, error)
    
    // GetByIDs batch retrieves multiple profiles
    GetByIDs(ctx context.Context, userIDs []domain.UserID) ([]*domain.UserProfile, error)
    
    // GetByNickname retrieves a profile by nickname
    GetByNickname(ctx context.Context, nickname string) (*domain.UserProfile, error)
    
    // SearchByNickname searches profiles with cursor pagination
    // Returns profiles and next cursor (empty if no more results)
    SearchByNickname(ctx context.Context, query, cursor string, limit int32) ([]*domain.UserProfile, string, error)
    
    // NicknameExists checks if nickname is taken
    NicknameExists(ctx context.Context, nickname string) (bool, error)
}
```

**Service Implementation to Create:**

`internal/service/user_service_impl.go`:
```go
type userService struct {
    profileRepo repository.ProfileRepository
}

func NewUserService(profileRepo repository.ProfileRepository) UserService {
    return &userService{profileRepo: profileRepo}
}

// Implements:
// - CreateProfile: validate nickname format (regex ^[a-z0-9_]{3,20}$), check uniqueness
// - UpdateProfile: validate fields, handle partial updates (empty = no change)
// - GetProfileByID: direct fetch
// - GetProfilesByIDs: batch fetch (handles empty list gracefully)
// - GetProfileByNickname: direct fetch
// - SearchByNickname: delegate to repo with pagination
```

**Key Considerations:**
- Nickname validation: lowercase alphanumeric + underscore, 3-20 chars
- UpdateProfile: only update non-empty fields (partial update pattern)
- SearchByNickname: cursor pagination for scalability
- Batch operations for efficiency (GetProfilesByIDs)

### 3. Social Service

**Repository Interfaces to Create:**

`internal/repository/friend_request_repository.go`:
```go
type FriendRequestRepository interface {
    // Create stores a new friend request
    Create(ctx context.Context, request *domain.FriendRequest) error
    
    // GetByID retrieves a friend request by ID
    GetByID(ctx context.Context, requestID domain.RequestID) (*domain.FriendRequest, error)
    
    // GetPendingRequest checks if pending request exists between users
    GetPendingRequest(ctx context.Context, requesterID, targetID domain.UserID) (*domain.FriendRequest, error)
    
    // ListByTargetID lists pending requests for a user (pagination)
    ListByTargetID(ctx context.Context, targetID domain.UserID, cursor string, limit int32) ([]*domain.FriendRequest, string, error)
    
    // UpdateStatus updates request status (accepted/declined)
    UpdateStatus(ctx context.Context, requestID domain.RequestID, status domain.FriendRequestStatus) error
}
```

`internal/repository/friendship_repository.go`:
```go
type FriendshipRepository interface {
    // Create creates a friendship (bidirectional)
    Create(ctx context.Context, friendship *domain.Friendship) error
    
    // Delete removes a friendship
    Delete(ctx context.Context, userID1, userID2 domain.UserID) error
    
    // Exists checks if friendship exists between two users
    Exists(ctx context.Context, userID1, userID2 domain.UserID) (bool, error)
    
    // ListFriends lists all friends for a user (pagination)
    ListFriends(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]domain.UserID, string, error)
}
```

`internal/repository/block_repository.go`:
```go
type BlockRepository interface {
    // Create blocks a user
    Create(ctx context.Context, block *domain.Block) error
    
    // Delete unblocks a user
    Delete(ctx context.Context, blockerID, blockedID domain.UserID) error
    
    // Exists checks if a block exists
    Exists(ctx context.Context, blockerID, blockedID domain.UserID) (bool, error)
    
    // IsBlocked checks if either user has blocked the other (bidirectional check)
    IsBlocked(ctx context.Context, userID1, userID2 domain.UserID) (bool, error)
}
```

**Service Implementations to Create:**

`internal/service/friend_request_service_impl.go`:
```go
type friendRequestService struct {
    requestRepo    repository.FriendRequestRepository
    friendshipRepo repository.FriendshipRepository
    blockRepo      repository.BlockRepository
    // Future: event publisher for Kafka
}

// Implements:
// - SendFriendRequest: validate not self, check not blocked, check not already friends, check no pending request
// - ListRequests: delegate to repo
// - AcceptFriendRequest: verify target user owns request, create friendship, update request status
// - DeclineFriendRequest: verify ownership, update status
```

`internal/service/friendship_service_impl.go`:
```go
type friendshipService struct {
    friendshipRepo repository.FriendshipRepository
}

// Implements:
// - ListFriends: delegate to repo
// - RemoveFriend: verify friendship exists, delete
```

`internal/service/block_service_impl.go`:
```go
type blockService struct {
    blockRepo      repository.BlockRepository
    friendshipRepo repository.FriendshipRepository
}

// Implements:
// - BlockUser: create block, remove friendship if exists
// - UnblockUser: delete block
```

`internal/service/relationship_service_impl.go`:
```go
type relationshipService struct {
    friendshipRepo repository.FriendshipRepository
    requestRepo    repository.FriendRequestRepository
    blockRepo      repository.BlockRepository
}

// Implements:
// - CheckRelationship: check blocked first, then friends, then pending request, else none
```

**Key Considerations:**
- Friendship is bidirectional: stored once but queried both ways
- Block takes precedence: cannot send friend requests to blocked users
- AcceptFriendRequest: creates Friendship + updates request status (may need transaction)
- BlockUser: automatically removes friendship if exists
- Permission checks: only target user can accept/decline their requests
- Self-action prevention: cannot befriend/block yourself

### 4. Chat Service

**Repository Interfaces to Create:**

`internal/repository/chat_repository.go`:
```go
type ChatRepository interface {
    // Create stores a new chat
    Create(ctx context.Context, chat *domain.Chat) error
    
    // GetByID retrieves a chat by ID
    GetByID(ctx context.Context, chatID domain.ChatID) (*domain.Chat, error)
    
    // GetByParticipants finds chat between two users
    GetByParticipants(ctx context.Context, userID1, userID2 domain.UserID) (*domain.Chat, error)
    
    // ListUserChats lists all chats for a user (pagination)
    ListUserChats(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]*domain.Chat, string, error)
    
    // IsParticipant checks if user is a participant
    IsParticipant(ctx context.Context, chatID domain.ChatID, userID domain.UserID) (bool, error)
}
```

`internal/repository/message_repository.go`:
```go
type MessageRepository interface {
    // Create stores a new message
    Create(ctx context.Context, message *domain.Message) error
    
    // GetByID retrieves a message by ID
    GetByID(ctx context.Context, messageID domain.MessageID) (*domain.Message, error)
    
    // ListByChatID lists messages in a chat (pagination, newest first)
    ListByChatID(ctx context.Context, chatID domain.ChatID, cursor string, limit int32) ([]*domain.Message, string, error)
    
    // CheckIdempotency checks if message with idempotency key exists
    CheckIdempotency(ctx context.Context, idempotencyKey string) (*domain.Message, error)
}
```

**Cross-Service Client Interface:**

`internal/client/social_client.go`:
```go
// SocialClient wraps gRPC calls to Social Service
type SocialClient interface {
    // CheckRelationship calls SocialService.CheckRelationship
    CheckRelationship(ctx context.Context, userID, targetID domain.UserID) (domain.RelationshipStatus, error)
}
```

**Service Implementations to Create:**

`internal/service/chat_service_impl.go`:
```go
type chatService struct {
    chatRepo     repository.ChatRepository
    socialClient client.SocialClient
}

// Implements:
// - CreateDirectChat: check relationship (must be friends), check not blocked, check chat doesn't exist, create
// - GetChat: verify requester is participant
// - ListUserChats: delegate to repo
// - ListChatMembers: get chat, verify requester is participant, return participant IDs
```

`internal/service/message_service_impl.go`:
```go
type messageService struct {
    messageRepo repository.MessageRepository
    chatRepo    repository.ChatRepository
    // Future: Kafka producer for event publishing
}

// Implements:
// - SendMessage: validate text non-empty, check idempotency, verify sender is participant, create message
// - ListMessages: verify requester is participant, delegate to repo
// - StreamMessages: placeholder for future gRPC streaming implementation
```

**Key Considerations:**
- CreateDirectChat: calls Social Service gRPC to verify friendship
- Must handle Social Service unavailability gracefully (circuit breaker future consideration)
- Idempotency: check idempotency_key before creating duplicate messages
- Permission checks: all operations verify user is chat participant
- Future: Kafka event publishing for `message.sent` events

### 5. Notifications Service

**Repository Interfaces to Create:**

`internal/repository/notification_repository.go`:
```go
type NotificationRepository interface {
    // Create stores a new notification
    Create(ctx context.Context, notification *domain.Notification) error
    
    // GetByUserID lists notifications for a user (pagination)
    GetByUserID(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]*domain.Notification, string, error)
    
    // MarkAsRead marks a notification as read
    MarkAsRead(ctx context.Context, notificationID domain.NotificationID) error
    
    // GetByID retrieves a notification by ID
    GetByID(ctx context.Context, notificationID domain.NotificationID) (*domain.Notification, error)
}
```

**Service Implementation to Create:**

`internal/service/notification_service_impl.go`:
```go
type notificationService struct {
    notificationRepo repository.NotificationRepository
}

// Implements:
// - GetNotifications: delegate to repo
// - MarkAsRead: verify notification belongs to user, mark as read
```

**Key Considerations:**
- This service is primarily for notification history
- Real-time delivery via Kafka consumer + SSE (separate from this iteration)
- Permission check: ensure user owns notification before marking as read

## Cross-Cutting Concerns

### 1. Error Handling Pattern

All service methods should follow this pattern:

```go
func (s *service) Operation(ctx context.Context, ...) (*Type, error) {
    // 1. Input validation (business rules)
    if err := validateInput(...); err != nil {
        return nil, domain.ErrInvalidXXX
    }
    
    // 2. Business logic checks (using repositories)
    exists, err := s.repo.Exists(ctx, ...)
    if err != nil {
        return nil, fmt.Errorf("check existence: %w", err)
    }
    if exists {
        return nil, domain.ErrAlreadyExists
    }
    
    // 3. Main operation
    entity := &domain.Entity{...}
    if err := s.repo.Create(ctx, entity); err != nil {
        return nil, fmt.Errorf("create entity: %w", err)
    }
    
    return entity, nil
}
```

### 2. Validation Helpers

Create validation utilities where needed:

`users/internal/service/validation.go`:
```go
func validateNickname(nickname string) error {
    matched, _ := regexp.MatchString(`^[a-z0-9_]{3,20}$`, nickname)
    if !matched {
        return domain.ErrInvalidNickname
    }
    return nil
}
```

### 3. ID Generation

Services generate IDs using UUIDs:

```go
import "github.com/google/uuid"

chatID := domain.NewChatID(uuid.New().String())
```

### 4. Testing Strategy

Each service implementation requires:

**Unit Tests Structure:**
```go
// Example: chat/internal/service/chat_service_test.go
type mockChatRepository struct {
    createFunc             func(ctx context.Context, chat *domain.Chat) error
    getByParticipantsFunc  func(ctx context.Context, u1, u2 domain.UserID) (*domain.Chat, error)
    // ... other methods
}

type mockSocialClient struct {
    checkRelationshipFunc func(ctx context.Context, u1, u2 domain.UserID) (domain.RelationshipStatus, error)
}

func TestCreateDirectChat_ValidFriends_ReturnsChat(t *testing.T) {
    // Arrange: mock returns friends status, no existing chat
    // Act: call CreateDirectChat
    // Assert: chat created successfully
}

func TestCreateDirectChat_NotFriends_ReturnsError(t *testing.T) {
    // Arrange: mock returns non-friend status
    // Act: call CreateDirectChat
    // Assert: returns ErrUsersNotFriends
}
```

**Test Coverage Requirements:**
- Happy path for each method
- Each domain error scenario
- Edge cases (empty lists, nil values, etc.)
- Cross-service call failures (for services with external dependencies)

### 5. Constructor Patterns

All service constructors follow this pattern:

```go
func NewXXXService(
    repo1 repository.Repo1,
    repo2 repository.Repo2,
    externalClient client.ExternalClient,
) XXXService {
    return &xxxService{
        repo1:          repo1,
        repo2:          repo2,
        externalClient: externalClient,
    }
}
```

## Implementation Order

1. **Auth Service** - Foundation, no external dependencies
   - UserRepository interface
   - RefreshTokenRepository interface
   - authService implementation
   - tokenService implementation
   - Unit tests

2. **Users Service** - No external dependencies
   - ProfileRepository interface
   - userService implementation
   - Unit tests

3. **Social Service** - Complex business logic
   - All repository interfaces
   - All service implementations
   - Unit tests

4. **Chat Service** - Depends on Social Service
   - Repository interfaces
   - SocialClient interface
   - Service implementations
   - Unit tests

5. **Notifications Service** - Simple, independent
   - NotificationRepository interface
   - notificationService implementation
   - Unit tests

## Potential Pitfalls & Mitigations

### Pitfall 1: Transaction Management
**Issue**: Operations like "accept friend request" need to atomically create friendship + update request status.

**Mitigation**: Repository interfaces will expose transaction methods in next iteration:
```go
type TransactionManager interface {
    WithReadCommitted(ctx context.Context, fn func(ctx context.Context) error) error
}
```

The `WithReadCommitted` method name explicitly indicates the isolation level guarantees (READ COMMITTED), making the transactional behavior clear to developers.

For now, document which operations need transactions in code comments.

### Pitfall 2: Circular Dependencies Between Services
**Issue**: Chat depends on Social, but careful not to create reverse dependency.

**Mitigation**: Clear dependency direction:
- Chat → Social (allowed)
- Social never calls Chat
- Use Kafka events for reverse flow (Social publishes, Chat consumes if needed)

### Pitfall 3: Context Cancellation Handling
**Issue**: Long-running operations may not respect context cancellation.

**Mitigation**: Repository implementations (next iteration) must check `ctx.Err()` before expensive operations.

### Pitfall 4: Nil Pointer Dereferences
**Issue**: Services receive nil dependencies or nil return values.

**Mitigation**:
- Constructors validate dependencies are not nil
- Always check repo return values for nil before dereferencing
- Return empty slices instead of nil for list operations

### Pitfall 5: Pagination Cursor Encoding
**Issue**: Cursor format not standardized, security concerns.

**Mitigation**:
- Cursor encodes only timestamp (createdAt) in base64 format
- UUID identifiers eliminate need for additional ID in cursor
- Repository implementations will handle encoding/decoding
- Never expose raw database values in cursors

### Pitfall 6: gRPC Client Initialization
**Issue**: Chat service needs Social Service gRPC client, but connection management is unclear.

**Mitigation**: 
- Define SocialClient interface now
- Implementation (gRPC connection setup) done in `cmd/main.go` in next iteration
- Service only depends on interface, fully testable with mocks

## Files to Create

### Auth Service
- `internal/repository/user_repository.go`
- `internal/repository/refresh_token_repository.go`
- `internal/service/auth_service_impl.go`
- `internal/service/token_service_impl.go`
- `internal/service/auth_service_test.go`
- `internal/service/token_service_test.go`

### Users Service
- `internal/repository/profile_repository.go`
- `internal/service/user_service_impl.go`
- `internal/service/validation.go`
- `internal/service/user_service_test.go`

### Social Service
- `internal/repository/friend_request_repository.go`
- `internal/repository/friendship_repository.go`
- `internal/repository/block_repository.go`
- `internal/service/friend_request_service_impl.go`
- `internal/service/friendship_service_impl.go`
- `internal/service/block_service_impl.go`
- `internal/service/relationship_service_impl.go`
- `internal/service/friend_request_service_test.go`
- `internal/service/friendship_service_test.go`
- `internal/service/block_service_test.go`
- `internal/service/relationship_service_test.go`

### Chat Service
- `internal/repository/chat_repository.go`
- `internal/repository/message_repository.go`
- `internal/client/social_client.go`
- `internal/service/chat_service_impl.go`
- `internal/service/message_service_impl.go`
- `internal/service/chat_service_test.go`
- `internal/service/message_service_test.go`

### Notifications Service
- `internal/repository/notification_repository.go`
- `internal/service/notification_service_impl.go`
- `internal/service/notification_service_test.go`

## Summary

This plan provides a complete blueprint for implementing the service layer while maintaining clean architecture principles. Each service will have:

1. Clear separation of concerns (repository for data, service for business logic)
2. Interface-based design for testability
3. Proper error handling and validation
4. Cross-service communication via interfaces
5. Comprehensive unit test coverage

The implementation will be iterative, starting with simpler services (Auth, Users) and progressing to more complex ones (Social, Chat) that have cross-service dependencies.

