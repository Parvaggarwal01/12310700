# Stage 1

## Core Actions

The notification platform must support the following core actions for logged-in students:

1. **Fetch Notifications:** Retrieve a paginated list of notifications, with the ability to filter by notification type (Event, Result, Placement).
2. **Mark as Read:** Update the status of a specific notification from unread to read.
3. **Real-Time Updates:** Receive new notifications instantly without needing to refresh the page.

---

## REST API Endpoints & Contracts

### 1. Fetch Notifications

Retrieves a list of notifications for the authenticated user.

- **Endpoint:** `GET /api/v1/notifications`
- **Headers:**
  - `Authorization`: `Bearer <Access_Token>`
- **Query Parameters:**
  - `page` (integer, default: 1): The page number.
  - `limit` (integer, default: 10): Items per page.
  - `type` (string, optional): Filter by type (`Event`, `Result`, `Placement`).
- **Response (200 OK):**

  ```json
  {
    "data": [
      {
        "id": "uuid-string",
        "type": "Placement",
        "message": "CSX Corporation hiring",
        "timestamp": "2026-04-22T17:51:18Z",
        "isRead": false
      }
    ],
    "meta": {
      "currentPage": 1,
      "totalPages": 5,
      "totalItems": 45
    }
  }
  ```

  ### 2. Mark Notification as Read

  Called when a user clicks on or views an unread notification. I'm using PATCH since we are only doing a partial update on the resource.
  - **Endpoint:** `PATCH /api/v1/notifications/:id/read`

- **Headers:**
  - `Authorization`: `Bearer <Access_Token>`

- **JSON Response (200 OK):**

  ```json
  {
    "success": true,
    "message": "status updated"
  }
  ```

  ### 3. Real-Time Implementation

  To push live notifications to the browser, I'll design this using Server-Sent Events (SSE) rather than WebSockets.
  - **Endpoint:** `GET /api/v1/notifications/stream`

- **Headers:**
  - `Accept`: `text/event-stream`, `Authorization`: `Bearer <Access_Token>`
