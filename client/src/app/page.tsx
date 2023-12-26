"use client";

import { useState } from "react";
import Dropzone from "react-dropzone";
import "./globals.css";
import { fileUrl, isDev, post, uploadFile } from "./utils";

const maxFiles = 0;

export default function Home() {
  const [uploading, setUploading] = useState(false);
  const [uploadedFileKeys, setUploadedFileKeys] = useState<string[]>([]);
  const maxUploadsReached = uploadedFileKeys.length > maxFiles;

  async function handleFileDrop(file: File) {
    if (!file) return;

    setUploading(true);
    handleFileUpload(file)
      .then((fileKey) => postDownloadRecord(file.name, fileKey))
      .then(async (res) => {
        const resData = await res.json();
        if (resData?.key) {
          setUploadedFileKeys((keys) => [resData.key, ...keys]);
        }
      })
      .catch()
      .finally(() => setUploading(false));
  }

  return (
    <main>
      <Dropzone
        onDrop={(files) => handleFileDrop(files?.[0])}
        disabled={uploading || maxUploadsReached}
        multiple={false}
      >
        {({ getRootProps, isDragActive }) => (
          <div
            {...getRootProps()}
            className="drop-zone column-center"
            style={{ borderColor: isDragActive ? "gainsboro" : "" }}
          >
            {uploading
              ? "File uploading..."
              : isDragActive
              ? "Release to upload"
              : "Drag a file here to upload"}
          </div>
        )}
      </Dropzone>
      {uploadedFileKeys.length > 0 && (
        <div className="uploaded-list">
          Your uploaded file(s)
          <ul>
            {uploadedFileKeys.map((key) => (
              <li key={key}>
                {fileUrl}/{key}
              </li>
            ))}
          </ul>
        </div>
      )}
    </main>
  );
}

async function handleFileUpload(file: File): Promise<string> {
  let fileKey = "";
  return post("/api/object", {})
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
