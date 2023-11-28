import { KEEPIX_API_URL, PLUGIN_API_SUBPATH } from "../constants";

// Plugin
export const getPluginStatus = async () =>
  request<any>({
    url: `${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/status`,
    method: 'GET',
    name: "getPluginStatus",
    parser: (data: any) => { return JSON.parse(data.result); } 
  });

export const getPluginWallet = async () => 
  request<any>({
    url: `${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/wallet-fetch`,
    method: 'GET',
    name: "getPluginWallet",
    parser: (data: any) => { return JSON.parse(data.result); } 
  });

export const getPluginSyncProgress = async () => 
  request<any>({
    url: `${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/sync-state`,
    method: 'GET',
    name: "getPluginSyncProgress",
    parser: (data: any) => { return JSON.parse(data.result); } 
  });

export const getPluginMiniPools = async () => 
  request<any>({
    url: `${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/minipool-fetch`,
    method: 'GET',
    name: "getPluginMiniPools",
    parser: (data: any) => { return JSON.parse(data.result); } 
  });

export const getPluginNodeInformation = async () => 
  request<any>({
    url: `${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/node-fetch`,
    method: 'GET',
    name: "getPluginNodeInformation",
    parser: (data: any) => { return JSON.parse(data.result); } 
  });

export const postPluginCreateMiniPool = async (body: any) => 
  request<any>({
    url: `${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/create-pool`,
    method: 'POST',
    name: "postPluginCreateMiniPool",
    body: body,
    parser: (data: any) => { return { result: JSON.parse(data.result), stdOut: data.stdOut }; } 
  });

export const postPluginMiniPoolExit = async (body: any) => 
  request<any>({
    url: `${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/exit-minipool`,
    method: 'POST',
    name: "postPluginMiniPoolExit",
    body: body,
    parser: (data: any) => { return { result: JSON.parse(data.result), stdOut: data.stdOut }; } 
  });

export const postPluginMiniPoolClose = async (body: any) => 
  request<any>({
    url: `${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/close-minipool`,
    method: 'POST',
    name: "postPluginMiniPoolClose",
    body: body,
    parser: (data: any) => { return { result: JSON.parse(data.result), stdOut: data.stdOut }; } 
  });

export const postPluginStakeRpl = async (amount: any) => 
  request<any>({
    url: `${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/stake-rpl`,
    method: 'POST',
    name: "postPluginStakeRpl",
    body: { Amount: amount },
    parser: (data: any) => { return { result: JSON.parse(data.result), stdOut: data.stdOut }; } 
  });

export const postPluginUnStakeRpl = async (amount: any) => 
  request<any>({
    url: `${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/unstake-rpl`,
    method: 'POST',
    name: "postPluginUnStakeRpl",
    body: { Amount: amount },
    parser: (data: any) => { return { result: JSON.parse(data.result), stdOut: data.stdOut }; } 
  });

export const getMinipoolMinimumStakeRplAmounts = async () => 
  request<any>({
    url: `${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/fetch-minimum-pool-stake-rpl-amounts`,
    method: 'GET',
    name: "getMinipoolMinimumStakeRplAmounts",
    parser: (data: any) => { return JSON.parse(data.result); } 
  });

// Functions
async function request<T>(options: any) {
  if (options.method === undefined) {
    options.method = 'GET';
  }
  const response: Response = await fetch(options.url, {
    method: options.method,
    headers: {
      "Content-Type": "application/json",
    },
    body: options.method === 'POST' && options.body !== undefined ? JSON.stringify(options.body): undefined
  });

  if (!response.ok) {
    throw new Error(`${options.name} call failed.`);
  }

  const data: T = await response.json();

  if (options.parser !== undefined) {
    return options.parser(data);
  }
  return data;
}
