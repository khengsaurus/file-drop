import { post, serverUrl } from "../../utils";

export async function POST(req: Request) {
  const { size, type } = (await req.json()) || {};
  const res = await post(`${serverUrl}/api/object`, { size, type });

  const resData = await res.json();
  if (!resData?.key || !resData?.url) {
    return new Response(null, { status: 500 });
  }

  return Response.json(resData);
}
