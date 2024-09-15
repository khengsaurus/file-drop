"use client";

import { saveAs } from "file-saver";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { Button } from "../../../components";
import { get } from "../../utils";

interface FilePageProps {
  params: { fileKey: string };
}

export default function FilePage({ params }: FilePageProps) {
  const { fileKey } = params;
  const [file, setFile] = useState({ url: "", fileName: "" });
  const router = useRouter();

  useEffect(() => {
    get(`${process.env.API_BASE_PATH}/object-record/${fileKey}`)
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
      <Button onClick={downloadFile}>Download {file?.fileName}</Button>
    </main>
  );
}
