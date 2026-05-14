import { useEffect, useState } from "react";
import NotificationList from "./NotificationList";
import "./App.css";

type Item = {
  ID: string;
  Type: string;
  Message: string;
  Timestamp: string;
};

export default function App() {
  const [items, setItems] = useState<Item[]>([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState("All");

  useEffect(() => {
    fetch("http://localhost:8080/api/v1/priority-inbox")
      .then((r) => r.json())
      .then((j) => {
        if (j && j.data) setItems(j.data);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Failed to fetch notifications", err);
        setLoading(false);
      });
  }, []);

  const types = ["All", "Placement", "Result", "Event"];
  const filtered = items.filter((i) => filter === "All" || i.Type === filter);

  return (
    <div className="app-root">
      <header className="header">
        <h1>Notifications</h1>
        <div className="controls">
          <select value={filter} onChange={(e) => setFilter(e.target.value)}>
            {types.map((t) => (
              <option key={t} value={t}>
                {t}
              </option>
            ))}
          </select>
        </div>
      </header>
      <main>
        {loading ? (
          <div className="loading">Loading...</div>
        ) : (
          <NotificationList items={filtered} />
        )}
      </main>
    </div>
  );
}
