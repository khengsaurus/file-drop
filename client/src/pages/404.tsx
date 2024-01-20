import { Button } from "@mui/material";
import Link from "next/link";
import "../globals.css";

export default function Custom404() {
  return (
    <main>
      <p>The link you entered may be expired</p>
      <br />
      <Button
        variant="contained"
        style={{ fontSize: "15px", textTransform: "none" }}
        size="small"
      >
        <Link href="/">Home</Link>
      </Button>
    </main>
  );
}
