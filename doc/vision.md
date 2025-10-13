# Technical vision of the project

## Architecture and services

1. **API Gateway** — The only public HTTP point (REST API) for user interaction. Does REST→gRPC mapping, validates JWT locally using cached public keys from Auth Service. Implements rate limiting.
2. **Auth Service** — Registration, login, JWT token generation, and JWKS endpoint for public key distribution.
3. **User Service** — User profiles management and user search.
4. **Social Service** — Friend relationships management (requests, blocking, friends list).
5. **Chat Service** — Direct messaging and chat management. Publishes events to message bus.
6. **Notification Service** — Consumes events from Kafka message bus and delivers notifications to users via SSE (Server-Sent Events).

**Communication:**
- All inter-service communications use **gRPC**.
- Service-to-service authentication via **service tokens** (or mTLS in production).
- Asynchronous events via **Kafka** message bus.

## Services contracts

### Auth Service

| RPC          | Request             | Response                                    | Purpose                           | Errors                             |
| ------------ | ------------------- | ------------------------------------------- | --------------------------------- | ---------------------------------- |
| Register     | { email, password } | { user_id }                                 | Register new user                 | ALREADY_EXISTS, INVALID_ARGUMENT   |
| Login        | { email, password } | { access_token, refresh_token, user_id }    | Authenticate and get JWT tokens   | UNAUTHENTICATED, INVALID_ARGUMENT  |
| Refresh      | { refresh_token }   | { access_token, refresh_token, user_id }    | Refresh JWT tokens                | UNAUTHENTICATED, INVALID_ARGUMENT  |
| GetPublicKeys| { }                 | { keys: [PublicKey { kid, alg, use, n, e }] } | Get public keys for JWT validation | —                                  |

**Notes:**
- JWT tokens use RS256 (asymmetric) algorithm
- Access tokens: short-lived (15 min), Refresh tokens: long-lived (30 days)
- Gateway calls `GetPublicKeys` on startup and caches them (refresh every 5-10 min)
- `PublicKey` contains JWK (JSON Web Key) fields: `kid` (key ID), `alg` (algorithm), `n` (modulus), `e` (exponent)

---

### User Service

| RPC                  | Request                                    | Response                      | Purpose                         | Errors                        |
| -------------------- | ------------------------------------------ | ----------------------------- | ------------------------------- | ----------------------------- |
| CreateProfile        | { user_id, nickname, bio?, avatar_url? }   | UserProfile                   | Create profile for user         | ALREADY_EXISTS, INVALID_ARG   |
| UpdateProfile        | { user_id, nickname?, bio?, avatar_url? }  | UserProfile                   | Update user profile             | ALREADY_EXISTS, NOT_FOUND     |
| GetProfileByID       | { user_id }                                | UserProfile                   | Get profile by user_id          | NOT_FOUND                     |
| GetProfilesByIDs     | { user_ids: [] }                           | { profiles: [UserProfile] }   | Batch get profiles              | —                             |
| GetProfileByNickname | { nickname }                               | UserProfile                   | Get profile by nickname         | NOT_FOUND                     |
| SearchByNickname     | { query, cursor?, limit }                  | { profiles: [], next_cursor? }| Search profiles matching query  | —                             |

**Notes:**
- `nickname` is unique, format: `^[a-z0-9_]{3,20}$`
- Profile creation happens after successful registration (user_id comes from Auth Service)

---

### Social Service

| RPC                  | Request                      | Response                                     | Purpose                        | Errors                                    |
| -------------------- | ---------------------------- | -------------------------------------------- | ------------------------------ | ----------------------------------------- |
| SendFriendRequest    | { target_user_id }           | FriendRequest                                | Send friend request            | INVALID_ARGUMENT, ALREADY_EXISTS          |
| ListRequests         | { user_id, cursor?, limit }  | { requests: [FriendRequest], next_cursor? }  | List pending friend requests   | —                                         |
| AcceptFriendRequest  | { request_id }               | FriendRequest                                | Accept friend request          | NOT_FOUND, PERMISSION_DENIED              |
| DeclineFriendRequest | { request_id }               | FriendRequest                                | Decline friend request         | NOT_FOUND, PERMISSION_DENIED              |
| RemoveFriend         | { friend_user_id }           | { }                                          | Remove user from friends       | NOT_FOUND                                 |
| ListFriends          | { user_id, cursor?, limit }  | { user_ids: [], next_cursor? }               | List user's friends            | —                                         |
| BlockUser            | { target_user_id }           | { }                                          | Block user                     | NOT_FOUND                                 |
| UnblockUser          | { target_user_id }           | { }                                          | Unblock user                   | NOT_FOUND                                 |
| CheckRelationship    | { user_id, target_user_id }  | { status: FRIEND/BLOCKED/PENDING/NONE }      | Check relationship status      | —                                         |

**Notes:**
- Friend requests are bidirectional checks: cannot send if already friends or blocked
- Blocking automatically removes from friends
- `CheckRelationship` is used by Chat Service to enforce permissions

---

### Chat Service

| RPC              | Request                             | Response                                    | Purpose                         | Errors                                |
| ---------------- | ----------------------------------- | ------------------------------------------- | ------------------------------- | ------------------------------------- |
| CreateDirectChat | { participant_id }                  | { chat_id }                                 | Create 1-on-1 chat              | ALREADY_EXISTS, PERMISSION_DENIED     |
| GetChat          | { chat_id }                         | Chat                                        | Get chat info                   | NOT_FOUND, PERMISSION_DENIED          |
| ListUserChats    | { user_id, cursor?, limit }         | { chats: [Chat], next_cursor? }             | List user's chats               | —                                     |
| ListChatMembers  | { chat_id }                         | { user_ids: [] }                            | Get chat participants           | NOT_FOUND, PERMISSION_DENIED          |
| SendMessage      | { chat_id, text, idempotency_key? } | Message                                     | Send message to chat            | INVALID_ARGUMENT, PERMISSION_DENIED   |
| ListMessages     | { chat_id, cursor?, limit }         | { messages: [Message], next_cursor? }       | Get message history             | PERMISSION_DENIED                     |
| StreamMessages   | { chat_id, since_unix_ms? }         | stream Message                              | Real-time message stream        | PERMISSION_DENIED                     |

**Notes:**
- Chat creation checks friend relationship via `SocialService.CheckRelationship`
- Cannot create chat if users are not friends or if target blocked sender
- `SendMessage` supports idempotency via optional `idempotency_key`
- `StreamMessages` uses gRPC server-side streaming for real-time updates
- Publishes events to Kafka: `message.sent { chat_id, sender_id, recipient_id, message_id, text }`

---

### Notification Service

**Event-Driven Architecture:**
- Purely consumer-based, listens to Kafka topics
- Delivers notifications to connected clients via **SSE (Server-Sent Events)**

**Consumed Kafka Events:**
- `message.sent` — New message notification
- `friend_request.sent` — New friend request notification
- `friend_request.accepted` — Friend request accepted notification

**Client Endpoints:**
- `GET /notifications/stream` — SSE endpoint for real-time notification stream
  - Returns: SSE stream of notification events
  - Requires: JWT authentication via query param or Authorization header

**Optional RPCs (for notification history):**

| RPC              | Request                        | Response                                     | Purpose                     | Errors        |
| ---------------- | ------------------------------ | -------------------------------------------- | --------------------------- | ------------- |
| GetNotifications | { user_id, cursor?, limit }    | { notifications: [Notification], next_cursor?} | Get notification history    | —             |
| MarkAsRead       | { notification_id }            | { }                                          | Mark notification as read   | NOT_FOUND     |

---

### Gateway (API Gateway)

**REST API:** Accepts REST requests, translates to gRPC calls to backend services, returns REST responses.

**Key Responsibilities:**
- JWT validation (local, using cached JWKS from Auth Service)
- Rate limiting per user/IP
- REST ↔ gRPC translation
- Request/response logging and tracing

**Example REST Routes:**

**Authentication:**
* `POST /v1/auth/register` → `AuthService.Register`
* `POST /v1/auth/login` → `AuthService.Login`
* `POST /v1/auth/refresh` → `AuthService.Refresh`

**User Profiles:**
* `POST /v1/profile` → `UserService.CreateProfile`
* `PUT /v1/profile` → `UserService.UpdateProfile`
* `GET /v1/profile/{id}` → `UserService.GetProfileByID`
* `GET /v1/users/search?q=nickname` → `UserService.SearchByNickname`

**Social:**
* `POST /v1/friends/request` → `SocialService.SendFriendRequest`
* `GET /v1/friends/requests` → `SocialService.ListRequests`
* `POST /v1/friends/requests/{id}/accept` → `SocialService.AcceptFriendRequest`
* `POST /v1/friends/requests/{id}/decline` → `SocialService.DeclineFriendRequest`
* `GET /v1/friends` → `SocialService.ListFriends`
* `DELETE /v1/friends/{id}` → `SocialService.RemoveFriend`
* `POST /v1/users/{id}/block` → `SocialService.BlockUser`

**Chat:**
* `POST /v1/chats` → `ChatService.CreateDirectChat`
* `GET /v1/chats` → `ChatService.ListUserChats`
* `GET /v1/chats/{id}` → `ChatService.GetChat`
* `POST /v1/chats/{id}/messages` → `ChatService.SendMessage`
* `GET /v1/chats/{id}/messages` → `ChatService.ListMessages`

## Important Notes

### API Design
1. **API Versioning:** All APIs (gRPC and REST) must be versioned (e.g., `/v1/...`, package `api.v1`)
2. **Pagination:** All list operations use cursor-based pagination (better than offset for databases)
   - Request: `{ cursor?: string, limit: int }`
   - Response: `{ items: [], next_cursor?: string }`
   - Cursor encodes the position (e.g., timestamp + id for time-ordered data)

### Data Management
3. **Database per Service:** Each microservice (except Gateway and Notification) has its own database
   - Auth Service: User credentials, refresh tokens
   - User Service: User profiles
   - Social Service: Friend relationships, blocks
   - Chat Service: Chats, messages
   - Notification Service (optional): Notification history

### Security & Authentication
4. **JWT Validation:** Gateway validates JWT locally using public keys from Auth Service
   - Gateway calls `AuthService.GetPublicKeys` on startup and periodically (every 5-10 min)
   - Caches public keys in memory to avoid calling Auth Service on every request
   - Reduces latency and Auth Service load
5. **Service-to-Service Auth:** Services authenticate via service tokens in gRPC metadata
   - Production: Consider mTLS for zero-trust security
6. **Rate Limiting:** Gateway implements rate limiting per user and per IP
   - Prevents abuse and ensures fair resource usage

### Reliability & Scalability
7. **Idempotency:** Critical operations support idempotency keys:
   - `SendMessage` with optional `idempotency_key`
   - `SendFriendRequest` idempotent by nature (check existing before creating)
8. **Event-Driven Architecture:** Kafka message bus for asynchronous events
   - Decouples services (Chat Service doesn't directly call Notification Service)
   - Enables horizontal scaling and replay capabilities
9. **Real-time Communication:**
   - **gRPC streaming** for `ChatService.StreamMessages` (server-to-Gateway)
   - **SSE** for Gateway-to-client notifications (simpler than WebSocket for one-way push)

### Service Dependencies
10. **Cross-Service Calls:**
    - Gateway → Auth Service (`GetPublicKeys` on startup and periodically)
    - Gateway → All services (REST to gRPC translation)
    - Chat Service → Social Service (`CheckRelationship` before chat creation)
    - Chat Service → Kafka (publish `message.sent` events)
    - Social Service → Kafka (publish `friend_request.*` events)
    - Notification Service → Kafka (consume events)

### Observability
11. **Tracing & Logging:** All services should support distributed tracing (OpenTelemetry)
    - Trace requests across service boundaries
    - Helps debug latency and errors in microservice architecture
