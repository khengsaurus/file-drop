"use client";

import { Button } from "@mui/material";
import { saveAs } from "file-saver";
import { useEffect, useState } from "react";
import { get, serverUrl } from "../../utils";
import { useRouter } from "next/navigation";

interface FilePageProps {
  params: { fileKey: string };
}

export default function FilePage({ params }: FilePageProps) {
  const { fileKey } = params;
  const [file, setFile] = useState({ url: "", fileName: "" });
  const router = useRouter();

  useEffect(() => {
    get(`${serverUrl}/api/object-record/${fileKey}`)
      .then((res) => {
        if (res?.status === 404) {
          router.push("/");
          throw new Error("Resource not found");
        } else {
          return res.json();
        }
      })
      .then(setFile)
      .catch(console.error);
  }, [fileKey, router]);

  function downloadFile() {
    if (!file?.url) {
      return;
    }
    fetch(file.url)
      .then((res) => res.blob())
      .then((blob) => saveAs(blob, file.fileName))
      .catch(console.error);
  }

  return (
    <main>
      <Button
        variant="contained"
        onClick={downloadFile}
        style={{ fontSize: "15px", textTransform: "none" }}
        size="small"
      >
        Download {file?.fileName}
      </Button>
    </main>
  );
}
