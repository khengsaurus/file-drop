import { NextRequest } from "next/server";
import { get, serverUrl } from "../../utils";

export async function GET(req: NextRequest) {
  const fileKey = req.nextUrl.searchParams.get("fileKey");
  if (!fileKey) {
    return new Response(null, { status: 400 });
  }

  return get(`${serverUrl}/download/${fileKey}`);
}
