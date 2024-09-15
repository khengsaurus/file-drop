"use client";

import { useEffect, useRef } from "react";
import "../globals.css";
import { Upload, UrlInput } from "./components";
import { getNewClientToken } from "./utils";

export default function Home() {
  const newTokenTimerRef = useRef(false);

  useEffect(() => {
    if (
      document.cookie.indexOf("X-FD-Client=") == -1 &&
      !newTokenTimerRef.current
    ) {
      newTokenTimerRef.current = true;
      getNewClientToken()
        .catch(console.error)
        .finally(() => (newTokenTimerRef.current = false));
    }

    return () => {
      newTokenTimerRef.current = false;
    };
  }, []);

  return (
    <main>
      <Upload />
      <UrlInput />
      <div className="home-info">Links will be available for 1 day</div>
    </main>
  );
}
