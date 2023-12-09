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
import { Staking } from "../components/Staking/Staking";


const faqSyncProgress: any[] = [
  {
    title: "Does the window have to stay open?",
    desc: "No, it's not necessary, but if you're on your own computer, please leave it running."
  },
  {
    title: "Why Erigon is not syncing while Heimdall is syncing?",
    desc: "The Erigon node needs to wait for the Heimdall node to be fully synchronized before starting its synchronization."
  },
  {
    title: "Heimdall finished syncing but Erigon is still at 0%?",
    desc: "Resync Erigon and wait for a few minutes.",
  }
];

export default function HomePage() {
  const [loading, setLoading] = useState(false);
  const [stakeDisplay, setStakeDisplay] = useState(true);
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

  const canStart = () => {
    const acceptedStates = ['NodeInstalled', 'StartingNode', "StartingHeimdall", "StartingErigon"];
    return acceptedStates.includes(statusQuery?.data?.NodeState ?? '');
  }

  return (
    <div className="AppBase-content">
      {(!statusQuery?.data || loading) && (
        <BigLoader title="" full={true}></BigLoader>
      )}
      {statusQuery?.data && statusQuery.data?.NodeState === 'NoState' && (
        <BannerAlert status="danger">Error with the Plugin "{statusQuery.data?.NodeState}" please Reinstall.</BannerAlert>
      )}
      {statusQuery?.data && canStart() && (
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
                <strong>{`Erigon Sync Progress:`}</strong>
                <Progress percent={Number(syncProgressQuery?.data.erigonSyncProgress)} description={syncProgressQuery?.data.erigonStepDescription ?? ''}></Progress>
                <strong>{`Heimdall Sync Progress:`}</strong>
                <Progress percent={Number(syncProgressQuery?.data.heimdallSyncProgress)} description={syncProgressQuery?.data.heimdallStepDescription ?? ''}></Progress>
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
                    const result = await postResync({ erigon: "true", heimdall: "false" });
                    console.log('NICEE', result);
                  } catch (e) {
                    console.log(e);
                  }
                  setLoading(false);
                }}
              >Re-sync Erigon</Btn>
              <Btn
                status="warning"
                onClick={async () => {
                  setLoading(true);
                  try {
                    const result = await postResync({ erigon: "false", heimdall: "true" });
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
      
      {/* stake MATIC */}
      {stakeDisplay
        && statusQuery?.data && syncProgressQuery?.data
        /*&& syncProgressQuery?.data?.IsSynced === true*/
        && statusQuery.data?.NodeState === 'NodeStarted'
        /*&& walletQuery.data?.Wallet !== undefined*/
        /*&& statusQuery.data?.IsRegistered === true*/
        && (<>
          <Staking wallet={walletQuery.data?.Wallet} status={statusQuery?.data} backFn={() => { setStakeDisplay(false); }}></Staking>
      </>)}
      <Sprites></Sprites>
    </div>
  );
}
