import { Button as MuiButton, type ButtonProps } from "@mui/material";

export default function Button(props: ButtonProps) {
  return (
    <MuiButton
      variant="contained"
      style={{ fontSize: "14px", textTransform: "none" }}
      size="small"
      {...props}
    />
  );
}
