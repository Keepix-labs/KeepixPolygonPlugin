import "./Form.scss";
import { Icon } from "@iconify-icon/react";
import { CheckboxProps, FormProps, InputProps } from "../../model/Form";
import { useEffect, useState } from "react";
import Btn from "../Btn/Btn";
import Popin from "../Popin/Popin";
import Loader from "../Loader/Loader";
import BannerAlert from "../BannerAlert/BannerAlert";
import { useQuery } from "@tanstack/react-query";

/** Form container */

export const Form = ({ children }: FormProps) => {
  return <form className="Form-form">{children}</form>;
};

/** Input wrapper with label */

export const Input = ({
  label,
  small,
  name,
  icon,
  required,
  children,
}: InputProps) => {
  return (
    <div>
      <label className="Form-label" htmlFor={name}>
        {label}: {required && <abbr title="Required">*</abbr>}
        {small && <small>{small}</small>}
      </label>
      <div className="Form-input">
        {icon && <Icon icon={icon} />}
        {children}
        <span></span>
      </div>
    </div>
  );
};

/** Radio / Checkbox */

export const Checkbox = ({
  label,
  name,
  type = "checkbox",
  enabled,
  disabled,
  defaultValue = false,
  style = {},
  onChange = () => {},
}: CheckboxProps) => {
  const [selected, setSelected] = useState(defaultValue);

  useEffect(() => {
    onChange(selected);
  }, [selected]);
  return (
    <div className="Form-checkbox">
      <div className="Form-label">{label}:</div>
      <input
        onClick={() => {
          setSelected(!selected);
        }}
        className={selected ? "active" : ""}
        id={name}
        name={name}
        type={type}
      />
      <label style={style} className="Form-checkboxCase" htmlFor={name}>
        {disabled && (
          <span className="Form-checkboxCaseDisabled">{disabled}</span>
        )}
        {enabled && <span className="Form-checkboxCaseEnabled">{enabled}</span>}
      </label>
    </div>
  );
};
