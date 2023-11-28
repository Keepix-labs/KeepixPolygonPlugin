import "./style.scss";

import { useEffect, useState } from "react";
import { Icon } from '@iconify-icon/react';
import Logo from "../Logo/Logo";

export default function BigLogo({
  full = false,
  children
}: any) {
  return (
    <div className={`transfer card card-default`} style={{ height: full ? '100vh' : undefined }}>
        <div className="state" style={{ height: full ? '100vh' : undefined }}>
            <div className="state-logo">
            <Logo text={false} />
            </div>
            <div className="state-title">
                {children}
            </div>
        </div>
    </div>
  );
}