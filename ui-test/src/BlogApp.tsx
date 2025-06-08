// A sample React component for a blog application, for proxy testing purposes

import React from "react";

const BlogApp: React.FC = () => {
  return (
    <div className="blog-app">
      <header id="heading" className="app-header">
        <h1>My Blog</h1>
      </header>
      <main className="app-content">
        <article className="blog-post">
          <h2>Welcome to My Blog</h2>
          <p>This is a sample blog post for testing purposes.</p>
        </article>
      </main>
      <footer className="app-footer">
        <p>&copy; 2023 My Blog</p>
      </footer>
    </div>
  );
};

export default BlogApp;
