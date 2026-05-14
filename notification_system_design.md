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

---


# Stage 3: Query Analysis & Performance Optimization

# Original Query

```sql
SELECT *
FROM notifications
WHERE studentID = 1042
  AND isRead = false
ORDER BY createdAt ASC;
```


# 1. Is This Query Accurate? Why Is It Slow?

## Accuracy

Yes, the query is functionally accurate and will return the correct unread notifications for that specific student. However, it is extremely slow because, without a proper index, the database engine must perform a Sequential Scan (full table scan). It is forcing the database to read through all 5,000,000 rows to find the handful of rows that match studentID = 1042 and isRead = false.

# 2. What would you change and what would be the likely computation cost?

I would create a Index on the columns used in the WHERE and ORDER BY clauses:
```sql
CREATE INDEX idx_student_unread_notifications
ON notifications (studentID, isRead, createdAt);
```

## Computation Cost Shift:
- **Before Index**: $O(N)$ where $N$ is 5,000,000.
- ***After Index**: $O(\log N)$

# 3. What would you change and what would be the likely computation cost?
No, this is bad advice for a production system.

**Why**: Every time a new notification is inserted, the database must also update every single index. This will severely degrade write performance.


# 4 Recent Placement Notifications Query
To find all students who received a "Placement" notification in the last 7 days, we need to join the notifications table with the students table and filter by the enum and timestamp.

```sql
SELECT DISTINCT s.id, s.name, s.email, s.roll_no
FROM students s
JOIN notifications n ON s.id = n.studentID
WHERE n.notificationType = 'Placement'
  AND n.createdAt >= NOW() - INTERVAL '7 days';
```

---

# Stage 4

## Performance Analysis
Fetching notifications directly from the PostgreSQL database on every single page load for 50,000 students is an anti-pattern. Even with optimal indexing, the sheer volume of concurrent read connections and queries will exhaust database connection pools, spike CPU/Memory usage, and result in a sluggish user experience.

To solve this, we must decouple the read-heavy traffic from the primary database.

## Proposed Solutions & Tradeoffs

### 1. Implement a Caching Layer (Redis)
Instead of querying PostgreSQL on page load, we introduce Redis, an in-memory key-value store. When a user requests their notifications, the API checks Redis first. If the data is there (Cache Hit), it returns immediately. If not (Cache Miss), it queries PostgreSQL, stores the result in Redis with a Time-To-Live (TTL), and then returns it.
*   **Implementation Strategy:** Cache the "Top 50" or "Unread Only" notifications per user.
*   **Pros:** Sub-millisecond read latency. Massively reduces the load on the primary SQL database.
*   **Tradeoffs:**
    *   *Complexity:* Adds another infrastructure component to manage.
    *   *Cost:* High-memory Redis instances can be expensive.

### 2. Server-Sent Events (SSE) & Local State Management
We already designed an SSE stream in Stage 1. We can leverage this to fundamentally change how the frontend behaves. Instead of fetching data *every* time a user navigates pages, the React app fetches the initial state once upon login. From then on, it maintains the notification list in local state (e.g., Redux or Context API) and updates it dynamically via the SSE stream.
*   **Implementation Strategy:** Fetch once -> Listen to SSE -> Update local UI state.
*   **Pros:** Drops backend API calls for page-loads to nearly zero after the initial session start. Provides the best UX since updates appear instantly.
*   **Tradeoffs:**
    *   *Client Memory:* The browser must hold the notification state.
    *   *Connection Management:* The frontend must robustly handle reconnects and state-resyncing if the SSE connection drops (e.g., user switches internet networks or wakes up a sleeping laptop).

## Final Recommendation
For a production environment of this scale, a hybrid approach is best: Use **Redis** to cache the initial payload to ensure the first page load is blazing fast, and rely on **SSE + Local State** to handle subsequent updates so the user doesn't need to refresh the page to see new data.


---


# Stage 5

## Analysis of the Proposed Implementation

**1. What shortcomings do you observe with this implementation?**
The current `for` loop is synchronous, blocking, and sequential.
*   **Latency:** If `send_email` takes just 500ms per student, processing 50,000 students will take over 6 hours. The HR user's browser will time out long before this finishes.
*   **Lack of Fault Tolerance:** There is no retry mechanism. If an external API blips, the notification fails permanently.
*   **Resource Exhaustion:** Keeping a single thread open and making 50,000 sequential DB/Network calls will likely exhaust memory or connection pools.

**2. Logs indicate that the 'send_email' call failed for 200 students midway. What now?**
With the current sequential design, an unhandled exception during the `send_email` call will crash the entire loop.
*   The system is now in an inconsistent state.
*   We don't know exactly which students received the email and which didn't without manually auditing the logs.
*   The remaining students in the array will not receive their notifications.

**3. Should the process of saving to DB as well as sending the email happen together? Why or why not?**
No, they should be decoupled.
*   **Different Failure Domains:** Saving to our own internal PostgreSQL database is fast and highly reliable. Sending an email relies on an external third-party API (like SendGrid or AWS SES) which is slow and prone to rate-limiting or network timeouts.
*   **Isolation:** If the email API goes down, it shouldn't prevent the in-app notification (DB save + real-time push) from succeeding. They must scale and fail independently.

## Redesign for Reliability and Speed

To make this fast and reliable, we need an **Event-Driven Architecture** using a Message Broker (like RabbitMQ, Apache Kafka, or Redis Pub/Sub) and Background Workers.

1.  When HR clicks "Notify All", the main API simply drops 50,000 messages into a queue and immediately responds with a "202 Accepted".
2.  Independent worker services consume these messages asynchronously.
3.  We split the queues: one for Emails, one for In-App notifications. This decouples the processes.
4.  Workers have built-in retry logic (exponential backoff) and Dead Letter Queues (DLQ) for permanent failures.

### Revised Pseudocode

```go

type NotificationPayload struct {
    StudentID string `json:"student_id"`
    Message   string `json:"message"`
}

func NotifyAllHandler(studentIDs []string, message string) string {
    // Quickly publish events to a message broker
    for _, id := range studentIDs {
        payload := NotificationPayload{StudentID: id, Message: message}

        PublishToQueue("email_queue", payload)
        PublishToQueue("in_app_queue", payload)
    }

    // Return immediately
    return "Notifications are processing in the background"
}

// --- 2. Background Worker: Email Service ---
func EmailWorker() {
    // Continuously consume messages from the broker
    for payload := range ConsumeFromQueue("email_queue") {
        err := SendEmail(payload.StudentID, payload.Message)

        if err != nil {
            if IsNetworkOrRateLimitError(err) {
                // Temporary failure: Requeue with exponential backoff
                RetryWithBackoff("email_queue", payload)
            } else {
                // Permanent failure: Send to Dead Letter Queue for manual review
                SendToDeadLetterQueue("email_dlq", payload)
            }
        }
    }
}

// --- 3. Background Worker: In-App Service ---
func InAppWorker() {
    for payload := range ConsumeFromQueue("in_app_queue") {
        err := SaveToDB(payload.StudentID, payload.Message)

        if err != nil {
            RetryWithBackoff("in_app_queue", payload)
            continue // Skip real-time push if DB save fails
        }

        // Trigger Server-Sent Event (SSE) to update the user's frontend instantly
        PushToApp(payload.StudentID, payload.Message)
    }
}

```