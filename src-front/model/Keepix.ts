export type KeepixInformation = {
  version: string;
  latestVersion?: string;
};

export type KeepixUpdateState = {
  success: boolean;
  description: string;
};

export type KeepixActionType = "UPDATE";

export type KeepixSettings = {
  "keepix-name": string;
};
