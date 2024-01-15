import "./Staking.scss";
import { Icon } from "@iconify-icon/react";
import Btn from "../Btn/Btn";
import Field from "../Field/Field";
import { KEEPIX_API_URL, PLUGIN_API_SUBPATH } from "../../constants";
import { safeFetch } from "../../lib/utils";
import Web3 from "web3";
import { useQuery } from "@tanstack/react-query";
import { getMinipools, postPluginStake, postPluginUnstake, postPluginReward } from "../../queries/api";
import BigLoader from "../BigLoader/BigLoader";
import Loader from "../Loader/Loader";
import Popin from "../Popin/Popin";
import Pool from "../Pool/Pool";
import { useState } from "react";
import BannerAlert from "../BannerAlert/BannerAlert";
import { Input } from "../Form/Form";

export const Staking = ({ wallet, ethBalance, maticBalance, backFn }: any) => {
    const [open, setPopinOpen] = useState(false);
    const [selectedPool, setPool] = useState(null);
    const [operation, setOperation] = useState<"stake" | "unstake" | "reward" | undefined>(undefined); // ["stake", "unstake"
    const [loading, setLoading] = useState(false);
    const [postResult, setPostResult] = useState<any>(undefined);
    const [maticAmount, setMaticAmount] = useState<number>(0);

    const poolsQuery = useQuery({
        queryKey: ["getMinipools"],
        queryFn: getMinipools,
        refetchInterval: 60000,
        enabled: true
    });

    const sendStake = async (body: any) => {
        setPopinOpen(true);
        setLoading(true);
        setPostResult(undefined);
        const result = await postPluginStake(body.amount, body.pool);
        setLoading(false);
        setPostResult({result: result.result !== false});
    };

    const sendUnstake = async (body: any) => {
        setPopinOpen(true);
        setLoading(true);
        setPostResult(undefined);
        const result = await postPluginUnstake(body.amount, body.pool);
        setLoading(false);
        setPostResult({result: result.result !== false});
    }

    const sendReward = async (body: any) => {
        setPopinOpen(true);
        setLoading(true);
        setPostResult(undefined);
        const result = await postPluginReward(body.pool);
        setLoading(false);
        setPostResult({result: result.result !== false});
    }
    
    const formattedETHBalance = Web3.utils.fromWei(ethBalance, "ether");
    const formattedMATICBalance = Web3.utils.fromWei(maticBalance, "ether");

    const onStake = (pool: any) => {
        setOperation("stake");
        setPool(pool);
    }

    const onUnstake = (pool: any) => {
        setOperation("unstake");
        setPool(pool);
    }

    const onReward = (pool: any) => {
        setOperation("reward");
        setPool(pool);
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
                <h1 className="AppBase-headerTitle">Manage MATIC Staking</h1>
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
                <Field
                        status="gray-black"
                        title="Wallet ETH Balance"
                        icon="mdi:ethereum"
                        color="white"
                    >{ formattedETHBalance } ETH</Field>
            </div>
            <div className="home-row-full" >
                    <Field
                        status="gray-black"
                        title="Wallet MATIC Balance"
                        icon="ion:rocket"
                        color="white"
                    >{ formattedMATICBalance } MATIC</Field>
            </div>
            <div>
            {poolsQuery.data?.length > 0 ? (
                poolsQuery.data.map((pool : any, index: number) => (
                    <Pool key={index} pool={pool} onStake={onStake} onUnstake={onUnstake} onReward={onReward}/>
                ))
            ) : (
            <div>No pools available.</div>
            )}
</div>

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
            {postResult !== undefined && postResult.result !== true && (<BannerAlert status="danger">Transaction failed. StackTrace: {postResult.stdOut}</BannerAlert>)}
            {postResult !== undefined && postResult.result === true && (<BannerAlert status="success">Transaction succeeded</BannerAlert>)}
          </Popin>
        </>
      )}
      {selectedPool && operation === "stake" && (
        <>
        <Popin
          title="Stake MATIC"
          close={() => {
              setPool(null);
          }}
        >
          <h2>Pool: {selectedPool.name}</h2>
          Current balance: {formattedMATICBalance} MATIC
          <div className="home-row-full">
                      <Input
                          label="Amount you want Stake"
                          name="maticAmount"
                          icon="material-symbols:edit"
                          required={true}
                          >
                          <input
                              onChange={(event: any) => {
                                  setMaticAmount(Number(event.target.value));
                              }}
                              type="number"
                              id="maticAmount"
                              defaultValue={0}
                              placeholder="0.000"
                          />
                      </Input>
                  </div>
          <Btn
              status="warning"
              disabled={maticAmount === 0 || maticAmount > parseFloat(formattedMATICBalance)}
              onClick={async () => {
                  const wei = Web3.utils.toWei(maticAmount.toString(), "ether");
                  setPool(null);
                  await sendStake({ amount: wei, pool: selectedPool.contractAddress });
              }}
          >Stake</Btn>
        </Popin>
      </>
      )}
      {selectedPool && operation === "unstake" && (
        <>
          <Popin
            title="Unstake MATIC"
            close={() => {
                setPool(null);
            }}
          >
            <h2>Pool: {selectedPool.name}</h2>
            Current stake: {selectedPool.userStake} MATIC
            <div className="home-row-full">
                        <Input
                            label="Amount you want Unstake"
                            name="maticAmount"
                            icon="material-symbols:edit"
                            required={true}
                            >
                            <input
                                onChange={(event: any) => {
                                    setMaticAmount(Number(event.target.value));
                                }}
                                type="number"
                                id="maticAmount"
                                placeholder="0.000"
                            />
                        </Input>
                    </div>
            <Btn
                status="warning"
                disabled={maticAmount === 0 || maticAmount > parseFloat(selectedPool.userStake)}
                onClick={async () => {
                    const wei = Web3.utils.toWei(maticAmount.toString(), "ether");
                    setPool(null);
                    await sendUnstake({ amount: wei, pool: selectedPool.contractAddress });
                }}
            >Unstake</Btn>
          </Popin>
        </>
      )}
      {selectedPool && operation === "reward" && (
        <>
          <Popin
            title="Claim MATIC rewards"
            close={() => {
                setPool(null);
            }}
          >
            <h2>Pool: {selectedPool.name}</h2>
            Current reward: {selectedPool.userReward} MATIC
            <Btn
                status="warning"
                onClick={async () => {
                    setPool(null);
                    await sendReward({ pool: selectedPool.contractAddress });
                }}
            >Claim</Btn>
          </Popin>
        </>
      )}
    </>);
}