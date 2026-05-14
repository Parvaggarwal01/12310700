import React from "react";

type Item = {
  ID: string;
  Type: string;
  Message: string;
  Timestamp: string;
};

export default function NotificationList({ items }: { items: Item[] }) {
  if (items.length === 0) {
    return <div className="empty-state">No notifications</div>;
  }

  return (
    <div className="container">
      <ul className="notif-list">
        {items.map((it) => (
          <li key={it.ID} className="notif-item">
            <div className="meta">
              <span className="type">{it.Type}</span>
              <span className="timestamp">{it.Timestamp}</span>
            </div>
            <div className="message">{it.Message}</div>
          </li>
        ))}
      </ul>
    </div>
  );
}
