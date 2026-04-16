import { BrowserRouter, Routes, Route } from "react-router-dom";
import Auth from "./pages/Auth";
import Projects from "./pages/Projects";
import Tasks from "./pages/Tasks";

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Auth />} />
        <Route path="/projects" element={<Projects />} />
        <Route path="/projects/:projectID" element={<Tasks />} />
      </Routes>
    </BrowserRouter>
  );
}