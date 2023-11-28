import "./Field.scss";
import { ReactNode } from "react";
import { Icon } from "@iconify-icon/react";
import { Link } from "react-router-dom";

type Status = "info" | "success" | "warning" | "danger" | "gray" | "gray-black";

type PropsBtn = {
  href?: string;
  target?: string;
  title?: string;
  icon?: string;
  status?: Status;
  color?: string;
  onClick?: () => void;
  userSelect?: 'none' | 'all' | 'text';
  children: ReactNode;
};

export default function Field({
  href,
  target,
  icon,
  title,
  children,
  status,
  color,
  onClick,
  userSelect = 'none',
}: PropsBtn) {
  const Content = (
    <>
      {icon && !icon.startsWith("http") && <Icon icon={icon} />}
      {icon && icon.startsWith("http") && <img src={icon} width="16px" />}
      {title && <h4>{title}:</h4>}
      <span>{href ? <><a href={href} target={target}>{children}</a></> : children}</span>
    </>
  );

  if (onClick) {
    return (
      <button onClick={onClick} className="Field-field" data-status={status} style={{color: color, userSelect: userSelect}} disabled>
        {Content}
      </button>
    );
  }

  return (
    <button className="Field-field" data-status={status} style={{color: color, userSelect: userSelect }}>
      {Content}
    </button>
  );
}
