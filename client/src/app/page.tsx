"use client";

import { Upload } from "./components";
import "../globals.css";

export default function Home() {
  return (
    <main>
      <Upload />
      <div className="home-info">Links will have a TTL of 1 day</div>
    </main>
  );
}
