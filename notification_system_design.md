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


# Stage 2: Database Design & Scalability

## Database Selection

For persistent storage, I recommend **PostgreSQL**.

### Why PostgreSQL?

1. **Relational Integrity**
   Notifications are inherently tied to specific users (students). A relational database enforces strict foreign key constraints, ensuring data consistency.

2. **ACID Compliance**
   We need guarantees that when a notification is marked as `read`, that state change is reliably saved and not lost.

3. **Advanced Indexing**
   PostgreSQL handles large datasets exceptionally well when properly indexed, which is critical for a high-read environment like a notification feed.

---

# Database Schema

```sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE notif_type AS ENUM (
    'Event',
    'Result',
    'Placement'
);

CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    roll_no VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    student_id INT NOT NULL,

    type notif_type NOT NULL,

    message TEXT NOT NULL,

    is_read BOOLEAN DEFAULT false,

    created_at TIMESTAMP WITH TIME ZONE
    DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_student
        FOREIGN KEY(student_id)
        REFERENCES students(id)
        ON DELETE CASCADE
);
```

---

# Performance Optimization

## 1. Indexing

As the notification table grows into millions of rows, queries can become slow.

### Problem
Fetching unread notifications for a user may require scanning large amounts of data.

### Solution
Use a composite index:

```sql
CREATE INDEX idx_notifications_student_read_created
ON notifications(student_id, is_read, created_at DESC);
```

---

## 2. Table Partitioning

### Problem
A single massive table increases:
- Query latency
- Index maintenance overhead
- Insert delays

### Solution
Partition the table by date (monthly or yearly).

Example:

```sql
CREATE TABLE notifications_2026_05
PARTITION OF notifications
FOR VALUES FROM ('2026-05-01')
TO ('2026-06-01');
```

---

## 3. Data Archiving / TTL Strategy

### Problem
Old notifications consume unnecessary storage.

### Solution
Archive notifications older than 6 months.

Possible approaches:
- Move old rows to an archive table
- Store archived data in object storage like S3

Example:

```sql
INSERT INTO notifications_archive
SELECT *
FROM notifications
WHERE created_at < NOW() - INTERVAL '6 months';

DELETE FROM notifications
WHERE created_at < NOW() - INTERVAL '6 months';
```

---

## 4. Redis Caching Layer

### Problem
Repeated reads for active users increase database load.

### Solution
Cache frequently accessed notifications in Redis.

Example:
- Cache top 50 unread notifications per user
- Refresh cache when a new notification arrives


---

# REST API Queries

## 1. Fetch Notifications

### Features
- Pagination support
- Optional filtering by notification type
- Sorted by latest notifications first

```sql

SELECT
    id,
    type,
    message,
    is_read,
    created_at
FROM notifications
WHERE student_id = $1


ORDER BY created_at DESC
LIMIT 10 OFFSET 0;
```

---

## 2. Mark Notification as Read

```sql
UPDATE notifications
SET is_read = true
WHERE id = $1
AND student_id = $2;
```

