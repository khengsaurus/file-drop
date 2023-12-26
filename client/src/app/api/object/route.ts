import { post, serverUrl } from "../../utils";

export async function POST() {
  const res = await post(`${serverUrl}/api/object`, {});

  const resData = await res.json();
  if (!resData?.key || !resData?.url) {
    return new Response(null, { status: 500 });
  }

  return Response.json(resData);
}
