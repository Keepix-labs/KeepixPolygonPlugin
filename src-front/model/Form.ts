import { ReactNode } from "react";

export type FormProps = {
  children: ReactNode;
};

export type InputProps = {
  label: string;
  small?: string;
  name?: string;
  icon?: string;
  required?: boolean;
  children: ReactNode;
};

export type Checkbox = "radio" | "checkbox";

export type CheckboxProps = {
  label: string;
  name: string;
  defaultValue?: boolean;
  enabled?: string;
  disabled?: string;
  type?: Checkbox;
  onChange?: (event: any) => void;
  style?: any
};
