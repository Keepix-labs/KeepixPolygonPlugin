import "./Progress.scss";
import { Icon } from '@iconify-icon/react';

type Props = {
  percent?: number;
  description?: string;
}

export default function Progress({ percent = 0, description = '' }: Props) {
  return (
    <div className={`progress ${percent === 100 ? 'complete' : percent}`}>
      <Icon className="progressIcon" icon={percent === 100 ? "ph:check-circle-duotone" : "svg-spinners:6-dots-scale-middle"} />
      <div className="progressBar" style={{ width: `${percent}%` }}></div>
      <span>{percent}%{description !== '' ? ` (${description})` : ''}</span>
    </div>
  )
}