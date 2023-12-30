import { useState } from "react";
import Dropzone from "react-dropzone";
import { maxFileSize } from "../consts";
import "../styles/globals.css";
import { post, serverUrl, uploadFile } from "../utils";

const maxFiles = 5;

export default function UploadPage() {
  const [uploading, setUploading] = useState(false);
  const [uploadedFileKeys, setUploadedFileKeys] = useState<string[]>([]);
  const maxUploadsReached = uploadedFileKeys.length > maxFiles;

  async function handleFileDrop(file: File) {
    if (!file || (file.size || 10e6) >= maxFileSize) return;

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

  function renderLabel(isDragActive: boolean) {
    return uploading ? (
      "File uploading..."
    ) : isDragActive ? (
      "Release to upload"
    ) : (
      <div className="column-center">
        Drag here or click to upload
        <br />
        <div className="small-text">Max file size: {maxFileSize / 1e6}MB</div>
      </div>
    );
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
            style={{ borderColor: isDragActive ? "rgb(25, 118, 210)" : "" }}
          >
            {renderLabel(isDragActive)}
          </div>
        )}
      </Dropzone>
      {uploadedFileKeys.length > 0 && (
        <div className="uploaded-list">
          Your uploaded file(s)
          <ul>
            {uploadedFileKeys.map((key) => (
              <li key={key}>
                {window.location.host}/file/{key}
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
  const { size, type } = file;
  if (size >= maxFileSize) throw new Error();

  return post(`${serverUrl}/api/object`, { size, type })
    .then(async (res) => {
      const { key, url } = (await res.json()) || {};
      if (!key || !url || res.status !== 200) {
        throw new Error("Failed to retrieve presigned URL");
      }
      fileKey = key;
      return url;
    })
    .then((url) => uploadFile(file, url))
    .then((res) => {
      if (res?.status !== 200) {
        throw new Error("Failed to upload file to presigned URL");
      }
      return fileKey;
    });
}

async function postDownloadRecord(fileName: string, key: string) {
  return post(`${serverUrl}/api/record`, { fileName, key });
}
