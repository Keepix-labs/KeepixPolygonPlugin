import "./MiniPool.scss";
import { Icon } from "@iconify-icon/react";
import Btn from "../Btn/Btn";
import Field from "../Field/Field";
import { safeFetch } from "../../lib/utils";
import { KEEPIX_API_URL, PLUGIN_API_SUBPATH } from "../../constants";
import { useState } from "react";
import Popin from "../Popin/Popin";
import Loader from "../Loader/Loader";
import BannerAlert from "../BannerAlert/BannerAlert";
import { postPluginMiniPoolClose, postPluginMiniPoolExit } from "../../queries/api";

export const MiniPool = ({ index, total, pool, wallet }: any) => {
    const [open, setPopinOpen] = useState(false);
    const [loading, setLoading] = useState(false);
    const [postResult, setPostResult] = useState<any>(undefined);

    const sendExitMiniPool = async (body: any) => {
        setPopinOpen(true);
        setLoading(true);
        setPostResult(undefined);
        const result = await postPluginMiniPoolExit(body);
        setLoading(false);
        setPostResult(result);
    };

    const sendCloseMiniPool = async (body: any) => {
        setPopinOpen(true);
        setLoading(true);
        setPostResult(undefined);
        const result = await postPluginMiniPoolClose(body);
        setLoading(false);
        setPostResult(result);
    };

    return (<>
        <div className="card card-default">
            <header className="AppBase-header">
                <div className="AppBase-headerIcon icon-app">
                <Icon icon="material-symbols:rocket" />
                </div>
                <div className="AppBase-headerContent">
                <h1 className="AppBase-headerTitle">MiniPool ({index} / {total})</h1>
                <div className="AppBase-headerSubtitle">MiniPool {pool['Address']}</div>
                </div>
            </header>
            <div className="home-row-full" >
            <Field
                status="gray-black" color="white"
                title="Minipool Address"
                icon={`${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/view/rocket-small.png`}
                href={`https://rocketscan.io/minipool/${pool['Address']}`}
                target="_blank"
            >{ pool['Address'] }</Field>
            </div>
            <div className="home-row-full" >
            <Field
                status="gray-black" color="white"
                title="Node Operator"
                icon={`${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/view/rocket-small.png`}
                href={`https://rocketscan.io/node/${pool['Delegate-address']}`}
                target="_blank"
            >{ pool['Delegate-address'] }</Field>
            </div>
            <div className="home-row-full" >
            <Field
                status="gray-black" color="white"
                title="Pub Key"
                icon={`${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/view/beaconcha.png`}
                href={`https://beaconcha.in/validator/${pool['Validator-pubkey']}`}
                target="_blank"
            >{ pool['Validator-pubkey'] }</Field>
            </div>
            <div className="home-row-full" >
            <Field
                status="success"
                title="Minipool Status"
                icon="material-symbols:work"
            >{ pool['Prelaunch'] ? 'Prelaunch (Your 8 or 16 ETH deposit will be transferred to be Beacon Chain in 12 hours)' : 'Running' }</Field>
            </div>
            <div className="home-row-full" >
            <Field
                status="gray-black" color="white"
                title="Refund"
                icon="formkit:ethereum"
            >{ pool['Available-refund'] }</Field>
            </div>
            <div className="home-row-full" >
            <Field
                status="gray-black" color="white"
                title="Portion"
                icon="formkit:ethereum"
            >{ pool['Your-portion'] }</Field>
            </div>
            <div className="home-row-2" >
                <Btn
                    icon="material-symbols:stop"
                    status="gray-black"
                    color="red"
                    onClick={async () => { await sendExitMiniPool({ MiniPoolAddress: pool['Address'] }); }}
                    >Exit</Btn>
                <Btn
                    icon="material-symbols:close"
                    status="gray-black"
                    color="orange"
                    onClick={async () => { await sendCloseMiniPool({ MiniPoolAddress: pool['Address'] }); }}
                >Close</Btn>
            </div>
            <div className="home-row-2" style={{ textAlign: "center" }} >
                <div>(Exit staking minipools from the beacon chain)</div>
                <div>(Withdraw any remaining balance from a minipool and close it)</div>
            </div>
        </div>
        {open && (
        <>
          <Popin
            title="MiniPool Task"
            close={() => {
              setPopinOpen(false);
            }}
          >
            {loading === true && (
                <Loader></Loader>
            )}
            {postResult !== undefined && postResult.result !== true && (<BannerAlert status="danger">Task failed. StackTrace: {postResult.stdOut}</BannerAlert>)}
            {postResult !== undefined && postResult.result === true && (<BannerAlert status="success">Task Done With Success</BannerAlert>)}
          </Popin>
        </>
      )}
    </>);
}