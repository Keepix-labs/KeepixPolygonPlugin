export type JaugeProps = {
  max: number;
  value: number;
  unit?: string;
};

export type RightProps = {
  icon: string;
  name: string;
  max: number;
  value: number;
  unit?: string;
};

export type MonitoringInformation = {
  cpu: string;
  memory: {
    free: string;
    used: string;
    total: string;
  };
};

export type MonitoringMeasure = {
  name: string;
  max: number;
  value: number;
  fixed: number;
  unit: string;
  unitMax?: string;
  icon: string;
};
