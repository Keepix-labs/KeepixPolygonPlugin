import { KEEPIX_API_URL, PLUGIN_API_SUBPATH } from "../constants";

// Plugin
export const getPluginStatus = async () =>
  request<any>({
    url: `${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/status`,
    method: 'GET',
    name: "getPluginStatus",
    parser: (data: any) => { return data.result; } 
  });

export const getPluginWallet = async () =>
  request<any>({
    url: `${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/wallet-fetch`,
    method: 'GET',
    name: "getPluginWallet",
    parser: (data: any) => { return data.result; } 
  });

export const getPluginSyncProgress = async () => 
  request<any>({
    url: `${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/sync-state`,
    method: 'GET',
    name: "getPluginSyncProgress",
    parser: (data: any) => { return data.result; } 
  });

export const getMinipools = async () => 
  request<any>({
    url: `${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/pools-fetch`,
    method: 'GET',
    name: "getMinipools",
    parser: (data: any) => { return data.result; }
  });

export const getPluginNodeInformation = async () => 
  request<any>({
    url: `${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/node-fetch`,
    method: 'GET',
    name: "getPluginNodeInformation",
    parser: (data: any) => { return data.result; } 
  });

  export const postResync = async (body: any) => 
  request<any>({
    url: `${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/resync`,
    method: 'POST',
    name: "postResync",
    body: body,
    parser: (data: any) => { return { result: data.result, stdOut: data.stdOut }; } 
  });

export const postPluginStake = async (amount: any, address: string) => 
  request<any>({
    url: `${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/stake`,
    method: 'POST',
    name: "postPluginStake",
    body: { amount: amount, address: address },
    parser: (data: any) => { return { result: data.result, stdOut: data.stdOut }; } 
  });

export const postPluginUnstake = async (amount: any, address: string) => 
  request<any>({
    url: `${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/unstake`,
    method: 'POST',
    name: "postPluginUnstake",
    body: { amount: amount, address: address },
    parser: (data: any) => { return { result: data.result, stdOut: data.stdOut }; } 
  });

export const postPluginReward = async (address: string) => 
request<any>({
  url: `${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/reward`,
  method: 'POST',
  name: "postPluginReward",
  body: { address: address },
  parser: (data: any) => { return { result: data.result, stdOut: data.stdOut }; } 
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
