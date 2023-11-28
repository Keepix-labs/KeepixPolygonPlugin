import { ReactNode } from "react";

export type LayoutProps = {
  children: ReactNode;
};

export type PopinProps = {
  title: string;
  close?: () => void;
  children: ReactNode;
};

export type LogoProps = {
  text?: Boolean;
};

export type BannerProps = {
  status: "info" | "success" | "warning" | "danger";
  children: ReactNode;
};

export type AppBaseProps = {
  title: string;
  subTitle?: string;
  icon?: string;
  color?: string;
  children: ReactNode;
  right?: ReactNode;
};
