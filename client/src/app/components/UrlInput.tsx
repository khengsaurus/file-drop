import { useState } from "react";
import { Button } from "../../components";
import { getUrlUrl, isValidUrl, post, serverUrl } from "../utils";
import CopyButton from "./CopyButton";

export default function UrlInput() {
  const [link, setLink] = useState("");
  const [uploading, setUploading] = useState(false);
  const [uploadedUrlKeys, setUploadedUrlKeys] = useState<string[]>([]);

  function handleLink() {
    if (!isValidUrl(link)) {
      console.error("Invalid URL");
      return;
    }
    setUploading(true);
    postUrl(link)
      .then(async (res) => {
        const resData = await res.json();
        if (resData?.key) {
          setUploadedUrlKeys((keys) => [resData.key, ...keys.slice(0, 4)]);
          setLink("");
        }
      })
      .finally(() => setUploading(false));
  }

  const fileUrlPrefix =
    typeof window !== "undefined" ? getUrlUrl(window.location) : "";

  return (
    <>
      <div className="url-shortener">
        <input
          placeholder="URL"
          disabled={uploading}
          value={link}
          onChange={(e) => setLink(e?.target?.value)}
        />
        <Button
          onClick={handleLink}
          variant={uploading ? "outlined" : "contained"}
        >
          Shorten
        </Button>
      </div>
      {uploadedUrlKeys.length > 0 && (
        <div className="uploaded-list">
          Your shortened URL(s)
          <ul>
            {uploadedUrlKeys.map((key) => {
              const url = `${fileUrlPrefix}/${key}`;

              return (
                <li key={key}>
                  {url}
                  <CopyButton text={url} />
                </li>
              );
            })}
          </ul>
        </div>
      )}
    </>
  );
}

async function postUrl(url: string) {
  return post(`${serverUrl}/api/url`, { url });
}
