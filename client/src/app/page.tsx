"use client";

import { Upload, UrlInput } from "./components";
import "../globals.css";
import { useEffect } from "react";
import { getCurrentClientToken, getNewClientToken } from "./utils";

export default function Home() {
  useEffect(() => {
    if (!getCurrentClientToken()) {
      getNewClientToken()
        .then((res) => res.json())
        .then((data) => {
          window.localStorage.setItem("client-token", data.token || "");
        })
        .catch(console.error);
    }
  }, []);

  return (
    <main>
      <Upload />
      <UrlInput />
      <div className="home-info">Links will be available for 1 day</div>
    </main>
  );
}
