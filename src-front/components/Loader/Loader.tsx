import "./Loader.scss";
import { Icon } from "@iconify/react";

export default function Loader({ label = 'Loading...' }: any) {
  return (
    <div className="Loader">
      <Icon icon="svg-spinners:180-ring-with-bg" />
      <span>{label}</span>
    </div>
  );
}
