import Link from "next/link";
import { Button } from "../components";
import "../globals.css";

export default function Custom404() {
  return (
    <main>
      <p>The link you entered may be expired</p>
      <br />
      <Button>
        <Link href="/">Home</Link>
      </Button>
    </main>
  );
}
