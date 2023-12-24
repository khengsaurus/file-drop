"use client";

import Link from "next/link";
import { isDev, post, uploadFile } from "./utils";
import "./globals.css";
import Dropzone from "react-dropzone";
import { useState } from "react";

export default function Home() {
  const [uploading, setUploading] = useState(false);

  async function handleFileDrop(file: File) {
    if (!file) return;

    setUploading(true);
    handleFileUpload(file)
      .then((fileKey) => postDownloadRecord(file.name, fileKey))
      .then(async (res) => {
        const resData = await res.json();
        console.log(resData);
      })
      .catch()
      .finally(() => setUploading(false));
  }

  return (
    <main>
      Home
      <br />
      <br />
      <Dropzone
        onDrop={(files) => handleFileDrop(files?.[0])}
        disabled={uploading}
        multiple={false}
      >
        {({ getRootProps, isDragActive }) => (
          <div
            {...getRootProps()}
            className="drop-zone column-center"
            style={{ borderColor: isDragActive ? "gainsboro" : "" }}
          >
            {isDragActive
              ? "Release to upload"
              : "Drag a single file here to upload"}
          </div>
        )}
      </Dropzone>
      <br />
      <br />
      <Link href="/help">Help</Link>
    </main>
  );
}

async function handleFileUpload(file: File): Promise<string> {
  let fileKey = "";
  return post("/api/file", {})
    .then(async (res) => {
      const { key, url } = (await res.json()) || {};
      if (!key || !url || res.status !== 200) {
        throw new Error("Failed to retrieve presigned URL");
      }
      fileKey = key;
      return url;
    })
    .then((url) => uploadFile(file, url, isDev))
    .then((res) => {
      if (res?.status !== 200) {
        throw new Error("Failed to upload file to presigned URL");
      }
      return fileKey;
    });
}

async function postDownloadRecord(fileName: string, key: string) {
  return post("/api/record", { fileName, key });
}
