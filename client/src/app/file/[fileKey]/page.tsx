"use client";

import { useEffect } from "react";
import { get } from "../../utils";

interface FilePageProps {
  params: { fileKey: string };
}

export default function FilePage(props: FilePageProps) {
  const { params } = props;
  const { fileKey } = params;

  useEffect(() => {
    get("/api/record", { fileKey }).then(async (res) => {
      const resData = await res.json();
      console.log(resData);
      return;
    });
  }, [fileKey]);

  return (
    <main>
      File
      <br />
      <br />
      {fileKey}
    </main>
  );
}
