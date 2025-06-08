import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import App from "./App.tsx";
import { BrowserRouter, Route } from "react-router";
import { Routes } from "react-router";
import BlogApp from "./BlogApp.tsx";

createRoot(document.getElementById("root")!).render(
  <BrowserRouter>
    <StrictMode>
      <Routes>
        <Route path="/" element={<App />} />
        <Route path="/blog" element={<BlogApp />} />
      </Routes>
    </StrictMode>
  </BrowserRouter>,
);
