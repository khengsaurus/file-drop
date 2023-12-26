"use client";

import { Button } from "@mui/material";
import { saveAs } from "file-saver";
import { useEffect, useState } from "react";
import { get } from "../../utils";

interface FilePageProps {
  params: { fileKey: string };
}

export default function FilePage(props: FilePageProps) {
  const { params } = props;
  const [file, setFile] = useState({ url: "", fileName: "" });
  const { fileKey } = params;

  useEffect(() => {
    get("/api/record", { fileKey })
      .then((res) => res.json())
      .then(setFile)
      .catch(console.error);
  }, [fileKey]);

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
      <Button variant="outlined" onClick={downloadFile}>
        Download {file?.fileName}
      </Button>
    </main>
  );
}
