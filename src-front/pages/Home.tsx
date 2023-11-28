import { useEffect, useState } from "react";
import Btn from "../components/Btn/Btn";
import Field from "../components/Field/Field";
import "./Home.scss";
import { safeFetch } from "../lib/utils";
import { KEEPIX_API_URL, PLUGIN_API_SUBPATH } from "../constants";
import { useQuery } from "@tanstack/react-query";
import { getPluginMiniPools, getPluginStatus, getPluginSyncProgress, getPluginWallet } from "../queries/api";
import Sprites from "../components/Sprites/Sprites";
import BigLoader from "../components/BigLoader/BigLoader";
import BannerAlert from "../components/BannerAlert/BannerAlert";
import BigLogo from "../components/BigLogo/BigLogo";
import { Icon } from "@iconify-icon/react";
import FAQ from "../components/Faq/Faq";
import Progress from "../components/Progress/Progress";
import { MiniPool } from "../components/MiniPool/MiniPool";
import { Node } from "../components/Node/Node";
import { NewMiniPool } from "../components/NewMiniPool/NewMiniPool";
import { RplStaking } from "../components/RplStaking/RplStaking";


const faqSyncProgress: any[] = [
  {
    title: "Does the window have to stay open?",
    desc: "No, it's not necessary, but if you're on your own computer, please leave it running."
  },
  {
    title: "Why does the Execution sync remain at zero percent as the Consensus advances?",
    desc: "The Consensus will be 100% before the Execution node. This is a normal situation, so don't hesitate to wait several hours before doing anything."
  },
  {
    title: "Is the progress of the Consensus node blocked at 0%?",
    desc: "Please wait several hours to be sure of the blockage, once the problem has been confirmed please select re-sync. This will restart synchronization from zero. If this happens again, please repeat the same procedure. Ethereum clients have several synchronization problems that are resolved after several re-syncs.",
  },
  {
    title: "Is the progress of the Execution node blocked at 99.99%?",
    desc: "At the end of the Execution node's synchronization, it waits for the concensus to send it all the information. The concensus may be delayed, so it needs to check the data again, which may take some time, so please be patient at this stage. If nothing happens after 24 hours, please uninstall and reinstall the plugin."
  }
];

export default function HomePage() {
  const [loading, setLoading] = useState(false);
  const [newMiniPoolDisplay, setNewMiniPoolDisplay] = useState(false);
  const [stakeRplDisplay, setStakeRplDisplay] = useState(false);
  const [manageMiniPoolsDisplay, setManageMiniPoolsDisplay] = useState(false);
  
  const walletQuery = useQuery({
    queryKey: ["getPluginWallet"],
    queryFn: getPluginWallet
  });

  const statusQuery = useQuery({
    queryKey: ["getPluginStatus"],
    queryFn: async () => {
      if (walletQuery.data === undefined) { // wallet check required before anything
        await walletQuery.refetch();
      }
      return getPluginStatus();
    },
    refetchInterval: 2000
  });

  const syncProgressQuery = useQuery({
    queryKey: ["getPluginSyncProgress"],
    queryFn: getPluginSyncProgress,
    refetchInterval: 5000,
    enabled: statusQuery.data?.NodeState === 'NODE_RUNNING'
  });

  const miniPoolsQuery = useQuery({
    queryKey: ["getPluginMiniPools"],
    queryFn: async () => {
      let pools: any[] = [];
      try {
        const strRocketPoolMiniPools = await getPluginMiniPools();
        if (strRocketPoolMiniPools.pools != undefined && strRocketPoolMiniPools.pools != '') {
          let splitInformations = strRocketPoolMiniPools.pools.split("\n\n").slice(0, -2);

          if (splitInformations.length >= 1) {
              let nextsIsFinalizedMiniPools = false;
              let nextsIsPrelaunchMiniPools = false;
              // let numberOfPools = parseInt(splitInformations[0]);
              for (let i = 0; i < splitInformations.length; i++) {

                if (splitInformations[i].includes("finalized minipool")
                  || splitInformations[i].includes("Staking minipool")) {
                  if (splitInformations[i + 1] !== undefined && splitInformations[i + 2] !== undefined) {
                    nextsIsFinalizedMiniPools = true;
                    nextsIsPrelaunchMiniPools = false;
                    let numberOfPools = parseInt(splitInformations[i]);
                    i++;
                    continue ;
                  }
                }

                if (splitInformations[i].includes("Prelaunch minipool")) {
                  if (splitInformations[i + 1] !== undefined && splitInformations[i + 2] !== undefined) {
                    nextsIsFinalizedMiniPools = false;
                    nextsIsPrelaunchMiniPools = true;
                    let numberOfPools = parseInt(splitInformations[i]);
                    i++;
                    continue ;
                  }
                }

                if (splitInformations[i].trim() !== '') {
                  console.log(splitInformations[i]);
                  const poolInfos = splitInformations[i];
                  const poolData = poolInfos.split("\n").reduce((acc: any, x: any) => {
                    const line = x.split(":", 2);
                    if (line.length == 2) {
                      acc[line[0].trim().replace(/ /gm, '-')] = line[1].trim();
                    }
                    return acc;
                  }, {});
                  pools.push({
                    ... poolData,
                    Finalized: nextsIsFinalizedMiniPools ? true : false,
                    Prelaunch: nextsIsPrelaunchMiniPools ? true : false
                  });
                }
              }
          }
        }
      } catch (e) {
        console.log('failed parse or fetch pools', e);
      }
      console.log(pools);
      return pools;
    },
    refetchInterval: 10000,
    enabled: statusQuery.data?.NodeState === 'NODE_RUNNING' && statusQuery.data?.IsRegistered === true && syncProgressQuery?.data?.IsSynced === true
  });

  //syncProgressQuery?.data?.IsSynced === true

  return (
    <div className="AppBase-content">
      {(!statusQuery?.data || loading) && (
        <BigLoader title="" full={true}></BigLoader>
      )}
      {statusQuery?.data && statusQuery.data?.NodeState === 'NO_STATE' && (
        <BannerAlert status="danger">Error with the Plugin "{statusQuery.data?.NodeState}" please Reinstall.</BannerAlert>
      )}
      {statusQuery?.data && statusQuery.data?.NodeState === 'NODE_STOPPED' && (
        <BigLogo full={true}>
          <Btn
            status="warning"
            onClick={async () => {
              setLoading(true);
              await safeFetch(`${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/start`);
              setLoading(false);
            }}
          >Start</Btn>
        </BigLogo>
      )}
      {statusQuery?.data
        && statusQuery.data?.NodeState === 'NODE_RUNNING'
        && walletQuery.data?.Wallet === undefined && (<>
          setup wallet
      </>)}

      {statusQuery?.data
        && !syncProgressQuery?.data
        && statusQuery.data?.NodeState === 'NODE_RUNNING'
        && walletQuery.data?.Wallet !== undefined && (
        <BigLoader title="Estimation: 1 to 10 minutes." label="Retrieving synchronization information" full={true}>
          <Btn
                status="danger"
                onClick={async () => { await safeFetch(`${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/stop`) }}
              >Stop</Btn>
        </BigLoader>
      )}

      {statusQuery?.data
        && syncProgressQuery?.data
        && syncProgressQuery?.data?.IsSynced === false
        && statusQuery.data?.NodeState === 'NODE_RUNNING'
        && walletQuery.data?.Wallet !== undefined && (
        <BigLoader title="Estimation: 1 hour to several days." disableLabel={true} full={true}>
          <div className="state-title">
                <strong>{`Execution Sync Progress:`}</strong>
                <Progress percent={Number(syncProgressQuery?.data.executionSyncProgress)} description={syncProgressQuery?.data.executionSyncProgressStepDescription ?? ''}></Progress>
                <strong>{`Consensus Sync Progress:`}</strong>
                <Progress percent={Number(syncProgressQuery?.data.consensusSyncProgress)} description={syncProgressQuery?.data.consensusSyncProgressStepDescription ?? ''}></Progress>
                {/* <strong><Icon icon="svg-spinners:3-dots-scale" /></strong> */}
          </div>
          <FAQ questions={faqSyncProgress}></FAQ>
          <div className="home-row-2" >
            <Btn
                status="danger"
                onClick={async () => { await safeFetch(`${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/stop`) }}
              >Stop</Btn>
            <Btn
                status="danger"
                onClick={async () => { await safeFetch(`${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/restart`) }}
              >Restart</Btn>
          </div>
          <div className="home-row-2" >
              <Btn
                status="warning"
                onClick={async () => {
                  setLoading(true);
                  try {
                    const resultEth2 = await safeFetch(`${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/resync-eth1`);

                    console.log('NICEE', resultEth2);
                  } catch (e) {
                    console.log(e);
                  }
                  setLoading(false);
                }}
              >Re-sync Execution</Btn>
              <Btn
                status="warning"
                onClick={async () => {
                  setLoading(true);
                  try {
                    const resultEth2 = await safeFetch(`${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/resync-eth2`);

                    console.log('NICEE', resultEth2);
                  } catch (e) {
                    console.log(e);
                  }
                  setLoading(false);
                }}
              >Re-sync Consensus</Btn>
          </div>
        </BigLoader>
      )}

      {/* Register the node to RocketPool */}
      {statusQuery?.data && syncProgressQuery?.data
        && syncProgressQuery?.data?.IsSynced === true
        && statusQuery.data?.NodeState === 'NODE_RUNNING'
        && walletQuery.data?.Wallet !== undefined
        && statusQuery.data?.IsRegistered === false && (<>
        <BigLoader title="Node Ready." disableLabel={true} full={true}>
          <Btn
                status="warning"
                onClick={async () => {
                  setLoading(true);
                  try {
                    const resultRegister = await safeFetch(`${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/register-node`);
                  } catch (e) {
                    console.log(e);
                  }
                  setLoading(false);
                }}
              >Register My Node to RocketPool</Btn>
        </BigLoader>
      </>)}
      
      {/* Has one or more Pool */}
      {!newMiniPoolDisplay
        && !stakeRplDisplay
        && !manageMiniPoolsDisplay
        && statusQuery?.data && syncProgressQuery?.data
        && syncProgressQuery?.data?.IsSynced === true
        && statusQuery.data?.NodeState === 'NODE_RUNNING'
        && walletQuery.data?.Wallet !== undefined
        && statusQuery.data?.IsRegistered === true
        && miniPoolsQuery?.data
        && miniPoolsQuery?.data?.length >= 0 && (<>
        <Node wallet={walletQuery.data?.Wallet} minipools={miniPoolsQuery?.data ?? []} status={statusQuery?.data} ></Node>
        <div className="card card-default">
          <header className="AppBase-header">
              <div className="AppBase-headerIcon icon-app">
              <Icon icon="ion:rocket" />
              </div>
              <div className="AppBase-headerContent">
              <h1 className="AppBase-headerTitle">Actions</h1>
              </div>
          </header>
          <div className="home-row-full" >
            <Btn
              icon="material-symbols:lock"
              status="gray-black"
              color="white"
              onClick={async () => { setStakeRplDisplay(true); }}
            >Manage RPL Staking</Btn>
          </div>
          <div className="home-row-full" >
            <Btn
              icon="system-uicons:grid-small"
              status="gray-black"
              color="white"
              onClick={async () => { setManageMiniPoolsDisplay(true); }}
              disabled={miniPoolsQuery?.data?.length === 0}
            >Manage My MiniPools ({miniPoolsQuery?.data?.length})</Btn>
          </div>
          <div className="home-row-full" >
            <Btn
              icon="fe:plus"
              status="gray-black"
              color="white"
              onClick={async () => { setNewMiniPoolDisplay(true); }}
            >Create One New MiniPool</Btn>
          </div>
        </div>
      </>)}

      {/* Manage MiniPools */}
      {manageMiniPoolsDisplay
        && statusQuery?.data && syncProgressQuery?.data
        && syncProgressQuery?.data?.IsSynced === true
        && statusQuery.data?.NodeState === 'NODE_RUNNING'
        && walletQuery.data?.Wallet !== undefined
        && statusQuery.data?.IsRegistered === true
        && miniPoolsQuery?.data && (<>
          <div className="card card-default">
            <div className="home-row-full" >
                  <Btn
                  status="gray-black"
                  color="white"
                  onClick={async () => { setManageMiniPoolsDisplay(false); }}
                  >Back</Btn>
            </div>
            <header className="AppBase-header">
                <div className="AppBase-headerIcon icon-app">
                <Icon icon="ion:rocket" />
                </div>
                <div className="AppBase-headerContent">
                <h1 className="AppBase-headerTitle">Manage my MiniPools</h1>
                </div>
            </header>
            { miniPoolsQuery?.data.map((pool: any, index, array) => <MiniPool key={index} index={index + 1} total={array.length} pool={pool} wallet={walletQuery.data?.Wallet} ></MiniPool>)}
          </div>
      </>)}

      {/* Add new MiniPool */}
      {newMiniPoolDisplay
        && statusQuery?.data && syncProgressQuery?.data
        && syncProgressQuery?.data?.IsSynced === true
        && statusQuery.data?.NodeState === 'NODE_RUNNING'
        && walletQuery.data?.Wallet !== undefined
        && statusQuery.data?.IsRegistered === true
        && miniPoolsQuery?.data && (<>
          <NewMiniPool wallet={walletQuery.data?.Wallet} minipools={miniPoolsQuery?.data ?? []} status={statusQuery?.data} backFn={() => { setNewMiniPoolDisplay(false); }}></NewMiniPool>
      </>)}

      {/* stake Rpl */}
      {stakeRplDisplay
        && statusQuery?.data && syncProgressQuery?.data
        && syncProgressQuery?.data?.IsSynced === true
        && statusQuery.data?.NodeState === 'NODE_RUNNING'
        && walletQuery.data?.Wallet !== undefined
        && statusQuery.data?.IsRegistered === true
        && miniPoolsQuery?.data && (<>
          <RplStaking wallet={walletQuery.data?.Wallet} minipools={miniPoolsQuery?.data ?? []} status={statusQuery?.data} backFn={() => { setStakeRplDisplay(false); }}></RplStaking>
      </>)}
      <Sprites></Sprites>
    </div>
  );
}
