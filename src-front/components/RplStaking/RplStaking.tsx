import "./RplStaking.scss";
import { Icon } from "@iconify-icon/react";
import Btn from "../Btn/Btn";
import Field from "../Field/Field";
import { KEEPIX_API_URL, PLUGIN_API_SUBPATH } from "../../constants";
import { safeFetch } from "../../lib/utils";
import Web3 from "web3";
import { useQuery } from "@tanstack/react-query";
import { getMinipoolMinimumStakeRplAmounts, getPluginNodeInformation, postPluginCreateMiniPool, postPluginStakeRpl, postPluginUnStakeRpl } from "../../queries/api";
import BigLoader from "../BigLoader/BigLoader";
import Loader from "../Loader/Loader";
import Popin from "../Popin/Popin";
import { useState } from "react";
import BannerAlert from "../BannerAlert/BannerAlert";
import { Input } from "../Form/Form";

export const RplStaking = ({ wallet, status, minipools, backFn }: any) => {
    const [open, setPopinOpen] = useState(false);
    const [loading, setLoading] = useState(false);
    const [postResult, setPostResult] = useState<any>(undefined);
    const [rplAmount, setRplAmount] = useState<undefined | number>(undefined);
    const [rplAmountUnstake, setRplAmountUnstake] = useState<undefined | number>(undefined);
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

    const sendStake = async (body: any) => {
        setPopinOpen(true);
        setLoading(true);
        setPostResult(undefined);
        const result = await postPluginStakeRpl(body);
        setLoading(false);
        setPostResult(result);
    };

    const sendUnStake = async (body: any) => {
        setPopinOpen(true);
        setLoading(true);
        setPostResult(undefined);
        const result = await postPluginUnStakeRpl(body);
        setLoading(false);
        setPostResult(result);
    }

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
                <h1 className="AppBase-headerTitle">Manage RPL Staking</h1>
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
            {(!nodeInformationQuery.data || !miniPoolMinimumStakeRplAmountsQuery.data) && (<Loader></Loader>)}
            {nodeInformationQuery.data
                && miniPoolMinimumStakeRplAmountsQuery.data
                && (<>
                <div className="card card-default">
                    <div className="home-row-full" >
                        {!miniPoolMinimumStakeRplAmountsQuery.data && (<Loader></Loader>)}
                        {miniPoolMinimumStakeRplAmountsQuery.data && (
                            <Field
                                status="info"
                                title="Minimum Staking Needed for build 8-ETH Minipool"
                                color="white"
                            >{ miniPoolMinimumStakeRplAmountsQuery.data.minimumSmallPoolStake } RPL</Field>
                        )}
                    </div>
                    <div className="home-row-full" >
                        {!miniPoolMinimumStakeRplAmountsQuery.data && (<Loader></Loader>)}
                        {miniPoolMinimumStakeRplAmountsQuery.data && (
                            <Field
                                status="info"
                                title="Minimum Staking Needed for build 16-ETH Minipool"
                                color="white"
                            >{ miniPoolMinimumStakeRplAmountsQuery.data.minimumMiniPoolStake } RPL</Field>
                        )}
                    </div>
                    <div className="home-row-full">
                        <Field
                                status="gray-black"
                                title="Stake Balance"
                                color="white"
                            >{ nodeInformationQuery.data.node.nodeRPLStakedBalance } RPL</Field>
                    </div>
                    <div className="home-row-full">
                        <Input
                            label="Amount you want Stake"
                            name="rplAmount"
                            icon="material-symbols:edit"
                            required={true}
                            >
                            <input
                                onChange={(event: any) => {
                                    setRplAmount(Number(event.target.value));
                                }}
                                type="number"
                                id="rplAmount"
                                defaultValue=""
                                max={nodeInformationQuery.data.node.rplWalletBalance}
                                placeholder="0.000"
                            />
                        </Input>
                    </div>
                    <div className="home-row-full">
                        <Btn
                            icon="fe:plus"
                            status="gray-black"
                            color="white"
                            disabled={rplAmount === undefined}
                            onClick={async () => {
                                await sendStake(rplAmount);
                            }}
                        >Deposit</Btn>
                    </div>
                </div>
                <div className="card card-default">
                    <div className="home-row-full">
                        <Field
                                status="gray-black"
                                title="Withdrawable Balance"
                                color="white"
                            >{ nodeInformationQuery.data.node.nodeRPLStakedBalance } RPL</Field>
                    </div>
                    <div className="home-row-full">
                        <Input
                            label="Amount you want Unstake"
                            name="rplAmountUnstake"
                            icon="material-symbols:edit"
                            required={true}
                            >
                            <input
                                onChange={(event: any) => {
                                    setRplAmountUnstake(Number(event.target.value));
                                }}
                                type="number"
                                id="rplAmountUnstake"
                                defaultValue=""
                                max={nodeInformationQuery.data.node.nodeRPLStakedBalance}
                                placeholder="0.000"
                            />
                        </Input>
                    </div>
                    <div className="home-row-full">
                        <Btn
                            icon="fe:plus"
                            status="gray-black"
                            color="white"
                            disabled={rplAmountUnstake === undefined}
                            onClick={async () => {
                                await sendUnStake(rplAmountUnstake);
                            }}
                        >Unstake</Btn>
                    </div>
                </div>
            </>)}
        </div>
        {open && (
        <>
          <Popin
            title="Task Progress"
            close={() => {
              setPopinOpen(false);
            }}
          >
            {loading === true && (
                <Loader></Loader>
            )}
            {postResult !== undefined && postResult.result !== true && (<BannerAlert status="danger">Stake failed. StackTrace: {postResult.stdOut}</BannerAlert>)}
            {postResult !== undefined && postResult.result === true && (<BannerAlert status="success">Staked With Success</BannerAlert>)}
          </Popin>
        </>
      )}
    </>);
}