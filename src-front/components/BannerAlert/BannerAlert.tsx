import "./BannerAlert.scss";
import { Icon } from "@iconify-icon/react";
import { BannerProps } from "../../model/Props";

export default function BannerAlert({ status, children }: BannerProps) {
  const statusIcons = {
    info: "ph:info",
    success: "ph:check-circle",
    warning: "ph:warning",
    danger: "ph:x-circle",
  };

  const icon = statusIcons[status] || "ph:question";

  return (
    <div className="BannerAlert-main" data-status={status}>
      <Icon icon={icon} />
      <span>{children}</span>
    </div>
  );
}
