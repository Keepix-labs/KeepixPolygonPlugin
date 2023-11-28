import "./Popin.scss";
import { Icon } from "@iconify-icon/react";
import { PopinProps } from "../../model/Props";

const Popin = ({ title, close, children }: PopinProps) => {
  return (
    <div className="Popin">
      <div className="Popin-mask" onClick={close}></div>
      <div className="Popin-content card">
        <div className="Popin-close" onClick={close}>
          <Icon icon="material-symbols:close" />
        </div>
        <div className="Popin-title">{title}</div>
        {children}
      </div>
    </div>
  );
};

export default Popin;
