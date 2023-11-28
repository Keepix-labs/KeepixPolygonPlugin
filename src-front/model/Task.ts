export type TaskProgress = {
  percent?: number;
  status:
    | "INSTALLING"
    | "UNINSTALLING"
    | "RUNNING"
    | "STOPPED"
    | "FINISHED"
    | "ERROR"
    | "NONE";
  description: string;
};
