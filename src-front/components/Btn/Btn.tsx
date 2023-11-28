import "./Btn.scss";
import { ReactNode } from "react";
import { Icon } from "@iconify-icon/react";
import { Link } from "react-router-dom";

type Status = "info" | "success" | "warning" | "danger" | "gray" | "gray-black";

type PropsBtn = {
  href?: string;
  icon?: string;
  status?: Status;
  color?: string;
  borderRadius?: string;
  onClick?: () => void;
  children: ReactNode;
  target?: string;
  disabled?: boolean;
};

export default function Btn({
  href,
  icon,
  children,
  status,
  color,
  borderRadius,
  target,
  disabled = false,
  onClick,
}: PropsBtn) {
  const Content = (
    <>
      <span>{children}</span>
      {icon && <Icon icon={icon} />}
    </>
  );

  if (onClick) {
    return (
      <button onClick={onClick} className="Btn-btn" data-status={status} style={{color: color, borderRadius: borderRadius}} disabled={disabled}>
        {Content}
      </button>
    );
  }

  return (
    <Link to={href ? href : ''} target={target} className="Btn-btn" data-status={status} style={{color: color, borderRadius: borderRadius}}>
      {Content}
    </Link>
  );
}
