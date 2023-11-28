export type Plugin = {
  id: string;
  title: string;
  description: string;
  subTitle: string;
  installed: boolean;
  icon: string;
  version: string;
  latestVersion?: string;
  packageName?: string;
  repositoryUrl?: string;
};

export type PluginActionType = "INSTALL" | "UNINSTALL" | "UPDATE";

export type AsyncQueryResponse = { taskId?: string, error?: string };
