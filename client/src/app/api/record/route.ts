import { NextRequest } from "next/server";
import { get, post, serverUrl } from "../../utils";

export async function GET(req: NextRequest) {
  const fileKey = req.nextUrl.searchParams.get("fileKey");
  if (!fileKey) {
    return new Response(null, { status: 400 });
  }

  const res = await get(`${serverUrl}/api/record/${fileKey}`);

  const resData = await res.json();
  if (!resData?.url) {
    return new Response(null, { status: 500 });
  }

  return Response.json(resData);
}

export async function POST(req: Request) {
  const reqData = await req.json();
  const res = await post(`${serverUrl}/api/record`, reqData);

  const resData = await res.json();
  if (!resData?.key) {
    return new Response(null, { status: 500 });
  }

  return Response.json(resData);
}
