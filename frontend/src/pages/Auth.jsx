import { useState } from "react";
import API from "../api/api";
import { useNavigate } from "react-router-dom";

export default function Auth() {
  const [isLogin, setIsLogin] = useState(true);
  const [form, setForm] = useState({ name: "", email: "", password: "" });
  const nav = useNavigate();

  const submit = async () => {
    try {
      if (!form.email || !form.password || (!isLogin && !form.name)) {
        alert("Fill all fields");
        return;
      }

      if (isLogin) {
        const res = await API.post("/auth/login", {
          email: form.email,
          password: form.password,
        });

        localStorage.setItem("token", res.data.token);
        nav("/projects");
      } else {
        await API.post("/auth/register", form);
        alert("Registered! Please login.");
        setIsLogin(true);
      }
    } catch (err) {
      alert(err.response?.data?.error || "Error");
    }
  };

  return (
    <div style={{ padding: 20 }}>
      <h2>{isLogin ? "Login" : "Register"}</h2>

      {!isLogin && (
        <input
          placeholder="Name"
          value={form.name}
          onChange={(e) => setForm({ ...form, name: e.target.value })}
        />
      )}

      <br />

      <input
        placeholder="Email"
        value={form.email}
        onChange={(e) => setForm({ ...form, email: e.target.value })}
      />

      <br />

      <input
        type="password"
        placeholder="Password"
        value={form.password}
        onChange={(e) => setForm({ ...form, password: e.target.value })}
      />

      <br />

      <button onClick={submit}>Submit</button>

      <p onClick={() => setIsLogin(!isLogin)} style={{ cursor: "pointer" }}>
        Switch to {isLogin ? "Register" : "Login"}
      </p>
    </div>
  );
}