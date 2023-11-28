import "./NewMiniPool.scss";
import { Icon } from "@iconify-icon/react";
import Btn from "../Btn/Btn";
import Field from "../Field/Field";
import { KEEPIX_API_URL, PLUGIN_API_SUBPATH } from "../../constants";
import { safeFetch } from "../../lib/utils";
import Web3 from "web3";
import { useQuery } from "@tanstack/react-query";
import { getMinipoolMinimumStakeRplAmounts, getPluginNodeInformation, postPluginCreateMiniPool } from "../../queries/api";
import BigLoader from "../BigLoader/BigLoader";
import Loader from "../Loader/Loader";
import Popin from "../Popin/Popin";
import { useState } from "react";
import BannerAlert from "../BannerAlert/BannerAlert";

export const NewMiniPool = ({ wallet, status, minipools, backFn }: any) => {
    const [open, setPopinOpen] = useState(false);
    const [loading, setLoading] = useState(false);
    const [postResult, setPostResult] = useState<any>(undefined);
    // const web3 = new Web3();

    const nodeInformationQuery = useQuery({
        queryKey: ["getNodeInformation"],
        queryFn: getPluginNodeInformation,
        refetchInterval: 10000,
        enabled: status?.NodeState === 'NODE_RUNNING'
      });

    const miniPoolMinimumStakeRplAmountsQuery = useQuery({
        queryKey: ["getMinipoolMinimumStakeRplAmounts"],
        queryFn: getMinipoolMinimumStakeRplAmounts,
        refetchInterval: 60000,
        enabled: status?.NodeState === 'NODE_RUNNING'
    });

    const sendCreateMiniPool = async (body: any) => {
        setPopinOpen(true);
        setLoading(true);
        setPostResult(undefined);
        const result = await postPluginCreateMiniPool(body);
        setLoading(false);
        setPostResult(result);
    };

    return (<>
        <div className="card card-default">
            <div className="home-row-full" >
                <Btn
                status="gray-black"
                color="white"
                onClick={async () => { backFn(); }}
                >Back</Btn>
            </div>
            <header className="AppBase-header">
                <div className="AppBase-headerIcon icon-app">
                <Icon icon="ion:rocket" />
                </div>
                <div className="AppBase-headerContent">
                <h1 className="AppBase-headerTitle">Create one new MiniPool</h1>
                </div>
            </header>
            <div className="home-row-full" >
                <Field
                    status="gray-black"
                    title="Wallet Address"
                    icon="ion:wallet"
                    color="white"
                >{ wallet }</Field>
            </div>
            <div className="home-row-full" >
                {!nodeInformationQuery.data && (<Loader></Loader>)}
                {nodeInformationQuery.data && (
                    <Field
                        status="gray-black"
                        title="Wallet ETH Balance"
                        icon="mdi:ethereum"
                        color="white"
                    >{ nodeInformationQuery.data.node.ethWalletBalance } ETH</Field>
                )}
            </div>
            <div className="home-row-full" >
                {!nodeInformationQuery.data && (<Loader></Loader>)}
                {nodeInformationQuery.data && (
                    <Field
                        status="gray-black"
                        title="Wallet RPL Balance"
                        icon="ion:rocket"
                        color="white"
                    >{ nodeInformationQuery.data.node.rplWalletBalance } RPL</Field>
                )}
            </div>
            <div className="home-row-full" >
                {!nodeInformationQuery.data && (<Loader></Loader>)}
                {nodeInformationQuery.data && (
                    <Field
                        status="gray-black"
                        title="Node Staking RPL Balance"
                        icon="material-symbols:lock"
                        color="white"
                    >{ nodeInformationQuery.data.node.nodeRPLStakedBalance } RPL</Field>
                )}
            </div>
            {(!nodeInformationQuery.data || !miniPoolMinimumStakeRplAmountsQuery.data) && (<Loader></Loader>)}
            {nodeInformationQuery.data
                && miniPoolMinimumStakeRplAmountsQuery.data
                && (<>
                <div className="card card-default">
                    {nodeInformationQuery.data.node.ethWalletBalance < 8 && (<>
                        <div className="home-row-full" >
                                <Field
                                    status="danger"
                                    title="Require"
                                    color="white"
                                >You currently have less than 8 ETH in your wallet. Please add more ETH to your wallet before creating new MiniPool</Field>
                        </div>
                    </>)}
                    {nodeInformationQuery.data.node.ethWalletBalance < 8.1 && (<>
                        <div className="home-row-full" >
                                <Field
                                    status="warning"
                                    title="Warning"
                                    color="white"
                                >You currently have less than 8.1 ETH in your wallet. Please add at least 0.1 ETH for fees.</Field>
                        </div>
                    </>)}
                    {parseFloat(nodeInformationQuery.data.node.nodeRPLStakedBalance) < parseFloat(miniPoolMinimumStakeRplAmountsQuery.data.minimumSmallPoolStake) && (<>
                        <div className="home-row-full" >
                            <Field
                                status="warning"
                                title="Warning"
                                color="white"
                            >You don't have enough RPL tokens to create a minipool of 8 ETH. You are missing {parseFloat(miniPoolMinimumStakeRplAmountsQuery.data.minimumSmallPoolStake) - parseFloat(nodeInformationQuery.data.node.nodeRPLStakedBalance)} Staked RPL tokens.</Field>
                        </div>
                    </>)}
                    {parseFloat(nodeInformationQuery.data.node.nodeRPLStakedBalance) < parseFloat(miniPoolMinimumStakeRplAmountsQuery.data.minimumMiniPoolStake) && (<>
                        <div className="home-row-full" >
                            <Field
                                status="warning"
                                title="Warning"
                                color="white"
                            >You don't have enough RPL tokens to create a minipool of 16 ETH. You are missing {parseFloat(miniPoolMinimumStakeRplAmountsQuery.data.minimumMiniPoolStake) - parseFloat(nodeInformationQuery.data.node.nodeRPLStakedBalance)} Staked RPL tokens.</Field>
                        </div>
                    </>)}
                    <div className="home-row-2" >
                        <Btn
                            icon="fe:plus"
                            status="gray-black"
                            color="white"
                            onClick={async () => {
                                await sendCreateMiniPool({ SmallPool: true });
                            }}
                            disabled={nodeInformationQuery.data.node.ethWalletBalance < 8 || parseFloat(nodeInformationQuery.data.node.nodeRPLStakedBalance) < parseFloat(miniPoolMinimumStakeRplAmountsQuery.data.minimumSmallPoolStake)}
                        >Create 8 ETH MiniPool</Btn>
                        <Btn
                            icon="fe:plus"
                            status="gray-black"
                            color="white"
                            onClick={async () => {
                                await sendCreateMiniPool({ SmallPool: false });
                            }}
                            disabled={nodeInformationQuery.data.node.ethWalletBalance < 16 || parseFloat(nodeInformationQuery.data.node.nodeRPLStakedBalance) < parseFloat(miniPoolMinimumStakeRplAmountsQuery.data.minimumMiniPoolStake)}
                        >Create 16 ETH MiniPool</Btn>
                    </div>
                </div>
            </>)}
        </div>
        {open && (
        <>
          <Popin
            title="MiniPool Creation"
            close={() => {
              setPopinOpen(false);
            }}
          >
            {loading === true && (
                <Loader></Loader>
            )}
            {postResult !== undefined && postResult.result !== true && (<BannerAlert status="danger">MiniPool Creation failed. StackTrace: {postResult.stdOut}</BannerAlert>)}
            {postResult !== undefined && postResult.result === true && (<BannerAlert status="success">MiniPool Created With Success</BannerAlert>)}
          </Popin>
        </>
      )}
    </>);
}