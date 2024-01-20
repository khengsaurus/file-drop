"use client";

import { Upload, UrlInput } from "./components";
import "../globals.css";

export default function Home() {
  return (
    <main>
      <Upload />
      <UrlInput />
      <div className="home-info">Links will have a TTL of 1 day</div>
    </main>
  );
}
