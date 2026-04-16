import { useEffect, useState } from "react";
import API from "../api/api";
import { useNavigate } from "react-router-dom";

export default function Projects() {
  const [projects, setProjects] = useState([]);
  const [name, setName] = useState("");
  const [loading, setLoading] = useState(false);
  const nav = useNavigate();

  const fetchProjects = async () => {
    try {
      setLoading(true);
      const res = await API.get("/projects");
      setProjects(res.data.projects || []);
    } catch (err) {
      console.log(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProjects();
  }, []);

  const createProject = async () => {
    if (!name) return alert("Enter project name");

    await API.post("/projects", { name });
    setName("");
    fetchProjects();
  };

  const logout = () => {
    localStorage.removeItem("token");
    nav("/");
  };

  return (
    <div style={{ padding: 20 }}>
      <h2>Projects</h2>

      <button onClick={logout}>Logout</button>

      <br /><br />

      <input
        placeholder="Project name"
        value={name}
        onChange={(e) => setName(e.target.value)}
      />

      <button onClick={createProject}>Create</button>

      <br /><br />

      {loading ? (
        <p>Loading...</p>
      ) : projects.length === 0 ? (
        <p>No projects yet</p>
      ) : (
        projects.map((p) => (
          <div
            key={p.id}
            onClick={() => nav(`/projects/${p.id}`)}
            style={{ border: "1px solid #ccc", padding: 10, margin: 10, cursor: "pointer" }}
          >
            {p.name}
          </div>
        ))
      )}
    </div>
  );
}