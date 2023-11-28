import "./Node.scss";
import { Icon } from "@iconify-icon/react";
import Btn from "../Btn/Btn";
import Field from "../Field/Field";
import { KEEPIX_API_URL, PLUGIN_API_SUBPATH } from "../../constants";
import { safeFetch } from "../../lib/utils";
import Web3 from "web3";
import { useQuery } from "@tanstack/react-query";
import { getPluginNodeInformation } from "../../queries/api";
import BigLoader from "../BigLoader/BigLoader";
import Loader from "../Loader/Loader";
import Logo from "../Logo/Logo";

export const Node = ({ node, wallet, status, minipools }: any) => {

    const nodeInformationQuery = useQuery({
        queryKey: ["getNodeInformation"],
        queryFn: getPluginNodeInformation,
        refetchInterval: 10000,
        enabled: status?.NodeState === 'NODE_RUNNING'
      });

    const lastBlockNumberQuery = useQuery({
        queryKey: ["getLatestBlockNumber"],
        queryFn: async () => {
            const web3 = new Web3(nodeInformationQuery.data.node.rpcUrl);
            const latest = await web3.eth.getBlockNumber();
            return latest.toString();
        },
        enabled: nodeInformationQuery?.data?.node.rpcUrl !== undefined
    });

    return (<>
        <div className="card card-default">
            <header className="AppBase-header">
                <div className="AppBase-headerIcon icon-app">
                {/* <Icon icon="logos:ethereum-color" /> */}
                <Logo text={false} width="50px"></Logo>
                </div>
                <div className="AppBase-headerContent">
                <h1 className="AppBase-headerTitle">Your Ethereum Node</h1>
                <div className="AppBase-headerSubtitle">Information & actions</div>
                </div>
            </header>
            <div className="home-row-full" >
                <Field
                icon="pajamas:status-health"
                status="success"
                title="Status"
                >{status?.NodeState}</Field>
            </div>
            <div className="home-row-full" >
                <Field
                    status="gray-black"
                    title="Wallet Address"
                    icon="ion:wallet"
                    color="white"
                    userSelect="text"
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
                        userSelect="text"
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
                        userSelect="text"
                    >{ nodeInformationQuery.data.node.rplWalletBalance } RPL</Field>
                )}
            </div>
            <div className="home-row-full">
                <Field
                        status="gray-black"
                        title="Total Borrowed"
                        icon="ri:parent-line"
                        color="white"
                        userSelect="text"
                    >{ minipools.reduce((acc: any, x: any) => acc + (parseFloat(x['RP-deposit'])), 0) } ETH</Field>
            </div>
            <div className="home-row-full">
                <Field
                        status="gray-black"
                        title="Total ETH Staked"
                        icon="ri:parent-line"
                        color="white"
                        userSelect="text"
                    >{ minipools.reduce((acc: any, x: any) => acc + parseFloat(x['Node-deposit']), 0) } ETH</Field>
            </div>
            <div className="home-row-full">
                <Field
                        status="gray-black"
                        title="APY Estimation"
                        icon="material-symbols-light:rewarded-ads-sharp"
                        color="white"
                        userSelect="text"
                    >... %</Field>
            </div>
            <div className="home-row-full">
                    <Field
                        status="gray-black"
                        title="Latest Block"
                        icon="clarity:block-line"
                        color="white"
                        userSelect="text"
                    >
                        {!lastBlockNumberQuery.data && (<Icon icon="svg-spinners:180-ring-with-bg" />)}
                        {lastBlockNumberQuery.data && (lastBlockNumberQuery.data)}
                    </Field>
            </div>
            <div className="home-row-full" >
                {!nodeInformationQuery.data && (<Loader></Loader>)}
                {nodeInformationQuery.data && (<Field
                    icon="material-symbols:link"
                    status="info"
                    title="Your Own RPC URL"
                    color="white"
                    href={nodeInformationQuery.data.node.rpcUrl}
                    target="_blank"
                >{nodeInformationQuery.data.node.rpcUrl}</Field>)}
            </div>
            <div className="home-row-full" >
                {
                    (status?.NodeState !== 'NODE_STOPPED' && status?.NodeState !== 'NO_STATE') ?
                    <Btn
                    icon="material-symbols:stop"
                    status="gray-black"
                    color="red"
                    onClick={async () => { await safeFetch(`${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/stop`) }}
                    >Stop</Btn>
                    :
                    <Btn
                    icon="mdi:play"
                    status="warning"
                    onClick={async () => { await safeFetch(`${KEEPIX_API_URL}${PLUGIN_API_SUBPATH}/start`) }}
                    >Start</Btn>
                }
            </div>
            {/* <div className="card card-default">
                {!nodeInformationQuery.data && (<Loader></Loader>)}
                {nodeInformationQuery.data && (<>
                    <div className="home-row-full" >
                        <Field
                            status="gray"
                            title="RPL Staked"
                            icon="ion:rocket"
                            color="white"
                        >{ nodeInformationQuery.data.node.nodeRPLStakedBalance }</Field>
                    </div>
                    <div className="home-row-full" >
                        <Field
                            status="gray"
                            title="Effective RPL Staked"
                            icon="ion:rocket"
                            color="white"
                        >{ nodeInformationQuery.data.node.nodeRPLStakedEffectiveBalance }</Field>
                    </div>
                    <div className="home-row-full" >
                        <Field
                            status="gray"
                            title="RPL Staked Percentage of Borrowed ETH"
                            icon="ion:rocket"
                            color="white"
                        >{ nodeInformationQuery.data.node.nodeRPLStakedBorrowedETHPercentage } %</Field>
                    </div>
                    <div className="home-row-full" >
                        <Field
                            status="gray"
                            title="RPL Staked Percentage of Bonded ETH"
                            icon="ion:rocket"
                            color="white"
                        >{ nodeInformationQuery.data.node.nodeRPLStakedBondedETHPercentage } %</Field>
                    </div>
                    <div className="home-row-2" >
                        <Field
                            status="gray"
                            title="Minimum RPL Needed"
                            icon="ion:rocket"
                            color="white"
                        >{ nodeInformationQuery.data.node.nodeMinimumRPLStakeNeeded } RPL</Field>
                        <Field
                            status="gray"
                            title="Maximum RPL Stake Possible"
                            icon="ion:rocket"
                            color="white"
                        >{ nodeInformationQuery.data.node.nodeMaximumRPLStakePossible } RPL</Field>
                    </div>
                </>)}
            </div> */}
        </div>
    </>);
}