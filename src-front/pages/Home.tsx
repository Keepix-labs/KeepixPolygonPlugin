import { useEffect, useState } from "react";
import Btn from "../components/Btn/Btn";
import Field from "../components/Field/Field";
import "./Home.scss";
import { safeFetch } from "../lib/utils";
import { KEEPIX_API_URL, PLUGIN_API_SUBPATH } from "../constants";
import { useQuery } from "@tanstack/react-query";
import { getPluginStatus, getPluginSyncProgress, getPluginWallet, postResync } from "../queries/api";
import Sprites from "../components/Sprites/Sprites";
import BigLoader from "../components/BigLoader/BigLoader";
import BannerAlert from "../components/BannerAlert/BannerAlert";
import BigLogo from "../components/BigLogo/BigLogo";
import { Icon } from "@iconify-icon/react";
import FAQ from "../components/Faq/Faq";
import Progress from "../components/Progress/Progress";
import { Node } from "../components/Node/Node";
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
  const [stakeRplDisplay, setStakeRplDisplay] = useState(false);
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
    enabled: statusQuery.data?.NodeState === 'NodeStarted'
  });


  //syncProgressQuery?.data?.IsSynced === true

  console.log('statusQuery', statusQuery, loading);

  return (
    <div className="AppBase-content">
      {(!statusQuery?.data || loading) && (
        <BigLoader title="" full={true}></BigLoader>
      )}
      {statusQuery?.data && statusQuery.data?.NodeState === 'NoState' && (
        <BannerAlert status="danger">Error with the Plugin "{statusQuery.data?.NodeState}" please Reinstall.</BannerAlert>
      )}
      {statusQuery?.data && statusQuery.data?.NodeState === 'NodeInstalled' && (
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
        && statusQuery.data?.NodeState === 'NodeStarted'
        && walletQuery.data?.Wallet === undefined && (<>
          setup wallet
      </>)}

      {statusQuery?.data
        && !syncProgressQuery?.data
        && statusQuery.data?.NodeState === 'NodeStarted'
        /*&& walletQuery.data?.Wallet !== undefined*/ && (
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
        && statusQuery.data?.NodeState === 'NodeStarted'
        /*&& walletQuery.data?.Wallet !== undefined*/ && (
        <BigLoader title="Estimation: 1 hour to several days." disableLabel={true} full={true}>
          <div className="state-title">
                <strong>{`Bor Sync Progress:`}</strong>
                <Progress percent={Number(syncProgressQuery?.data.borSyncProgress)} description={syncProgressQuery?.data.executionSyncProgressStepDescription ?? ''}></Progress>
                <strong>{`Heimdall Sync Progress:`}</strong>
                <Progress percent={Number(syncProgressQuery?.data.heimdallSyncProgress)} description={syncProgressQuery?.data.borSyncProgress ?? ''}></Progress>
                {syncProgressQuery?.data.heimdallSyncProgress}
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
                    const result = await postResync({ bor: "true", heimdall: "false" });
                    console.log('NICEE', result);
                  } catch (e) {
                    console.log(e);
                  }
                  setLoading(false);
                }}
              >Re-sync Bor</Btn>
              <Btn
                status="warning"
                onClick={async () => {
                  setLoading(true);
                  try {
                    const result = await postResync({ bor: "false", heimdall: "true" });
                    console.log('NICEE', result);
                  } catch (e) {
                    console.log(e);
                  }
                  setLoading(false);
                }}
              >Re-sync Heimdall</Btn>
          </div>
        </BigLoader>
      )}

      {/* Register the node to RocketPool */}
      {statusQuery?.data && syncProgressQuery?.data
        && syncProgressQuery?.data?.IsSynced === true
        && statusQuery.data?.NodeState === 'NodeStarted'
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
      
      {/* stake Rpl */}
      {stakeRplDisplay
        && statusQuery?.data && syncProgressQuery?.data
        && syncProgressQuery?.data?.IsSynced === true
        && statusQuery.data?.NodeState === 'NodeStarted'
        && walletQuery.data?.Wallet !== undefined
        && statusQuery.data?.IsRegistered === true
        && (<>
          <RplStaking wallet={walletQuery.data?.Wallet} status={statusQuery?.data} backFn={() => { setStakeRplDisplay(false); }}></RplStaking>
      </>)}
      <Sprites></Sprites>
    </div>
  );
}
