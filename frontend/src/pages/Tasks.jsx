import { useEffect, useState } from "react";
import API from "../api/api";
import { useParams } from "react-router-dom";

export default function Tasks() {
  const { projectID } = useParams();

  const [tasks, setTasks] = useState([]);
  const [title, setTitle] = useState("");
  const [priority, setPriority] = useState("medium");
  const [filter, setFilter] = useState("");

  const fetchTasks = async () => {
    let url = `/projects/${projectID}/tasks`;

    if (filter) {
      url += `?status=${filter}`;
    }

    const res = await API.get(url);
    setTasks(res.data.tasks || []);
  };

  useEffect(() => {
    fetchTasks();
  }, [filter]);

  const createTask = async () => {
  try {
    if (!title) return alert("Enter task title");
    console.log(projectID)
    await API.post(`/projects/${projectID}/tasks`, {
      title,
      description: "",
      priority,
    });

    setTitle("");
    fetchTasks();
  } catch (err) {
    console.log(err.response?.data); // 🔥 see real backend error
    alert(err.response?.data?.error || "Failed to create task");
  }
};

  const updateTask = async (id, status) => {
    await API.patch(`/tasks/${id}`, { status });
    fetchTasks();
  };

  const deleteTask = async (id) => {
    await API.delete(`/tasks/${id}`);
    fetchTasks();
  };

  return (
    <div style={{ padding: 20 }}>
      <h2>Tasks</h2>

      {/* Create */}
      <input
        placeholder="Task title"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
      />

        <select value={priority} onChange={(e) => setPriority(e.target.value)}>
  <option value="low">Low</option>
  <option value="medium">Medium</option>
  <option value="high">High</option>
</select>

      <button onClick={createTask}>Add</button>

      <br /><br />

      {/* Filter */}
      <select onChange={(e) => setFilter(e.target.value)}>
        <option value="">All</option>
        <option value="todo">Todo</option>
        <option value="in_progress">In Progress</option>
        <option value="done">Done</option>
      </select>

      <br /><br />

      {/* Tasks */}
      {tasks.length === 0 ? (
        <p>No tasks</p>
      ) : (
        tasks.map((t) => (
          <div key={t.id} style={{ border: "1px solid #ccc", margin: 10, padding: 10 }}>
            <p><b>{t.title}</b></p>
            <p>Status: {t.status}</p>

            <button onClick={() => updateTask(t.id, "todo")}>Todo</button>
            <button onClick={() => updateTask(t.id, "in_progress")}>In Progress</button>
            <button onClick={() => updateTask(t.id, "done")}>Done</button>

            <br />

            <button onClick={() => deleteTask(t.id)}>Delete</button>
          </div>
        ))
      )}
    </div>
  );
}