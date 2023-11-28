import "./style.scss";

import { useEffect, useState } from "react";
import { Icon } from '@iconify-icon/react';
import Logo from "../Logo/Logo";

type Data = {
  state: "IN_PROGRESS" | "DONE";
};

export default function BigLoader({
  title,
  label = 'Loading',
  full = false,
  disableLabel = false,
  children
}: any) {
  return (
    <div className={`transfer card card-default`} style={{ height: full ? '100vh' : undefined }}>
        <div className="state" style={{ height: full ? '100vh' : undefined }}>
            <div className="state-logo">
            <Logo text={false} />
            </div>
            {(
              <div className="state-title">
                  <span>{title}</span>
                  {disableLabel === true ? '' : <strong>{label} <Icon icon="svg-spinners:3-dots-scale" /></strong>}
              </div>
            )}
            {children}
        </div>
    </div>
  );
}