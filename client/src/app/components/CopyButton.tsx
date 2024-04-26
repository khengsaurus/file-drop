import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import CheckCircleOutlineIcon from "@mui/icons-material/CheckCircleOutline";
import { useRef, useState } from "react";

interface CopyButtonProps {
  text: string;
}

export default function CopyButton({ text }: CopyButtonProps) {
  const [hover, setHover] = useState(false);
  const [showChecked, setShowChecked] = useState(false);
  const checkedTimerRef = useRef<NodeJS.Timeout>();

  function handleCopy() {
    setShowChecked(true);
    navigator.clipboard.writeText(text);
    clearTimeout(checkedTimerRef.current);
    checkedTimerRef.current = setTimeout(() => setShowChecked(false), 2000);
  }

  const Icon = showChecked ? CheckCircleOutlineIcon : ContentCopyIcon;

  return (
    <span
      className="copy-button"
      onMouseEnter={() => setHover(true)}
      onMouseLeave={() => setHover(false)}
      onClick={handleCopy}
    >
      <Icon style={hover ? { fill: "rgb(100, 160, 220)" } : undefined} />
    </span>
  );
}
