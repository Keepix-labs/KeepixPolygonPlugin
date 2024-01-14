import "./style.scss";

import { useEffect, useState } from "react";
import { Icon } from '@iconify-icon/react';
import Web3 from "web3";
import Btn from "../Btn/Btn";


export default function Pool({
  pool,
  onStake,
  onUnstake,
  onReward
}: any) {
  const statusColor = pool.currentState !== "HEALTHY" ? "red" : "green";
  return (
    <div className="pool card card-default">
      <div className="details">
          <img src={pool.logoUrl} alt="" />
          <label className="h3">{pool.name !== "" ? pool.name : "Anonymous"}</label>
          <label className="h3" style={{color: statusColor}}>{pool.currentState}</label>
          <label>Accepts stakers: {pool.delegationEnabled? "Yes": "No"}</label>
          <label>Commission: {pool.commissionPercent}%</label>
          <label>Performance: {pool.performanceIndex}%</label>
          <label>Stake: {(pool.totalStaked / 10**18).toFixed(0)} MATIC</label>
          <label>Your stake: {parseFloat(pool.userStake).toFixed(2)} MATIC</label>
          <label>Your reward: {parseFloat(pool.userReward).toFixed(2)} MATIC</label>
          <label>Min stake: {parseFloat(pool.minStake).toFixed(2)} MATIC</label>
          <div className="home-row-3">
            <Btn onClick={() => onStake(pool)}>Stake</Btn>
            {parseFloat(pool.userStake) > 0 && <Btn onClick={() => onUnstake(pool)}>Unstake</Btn>}
            {parseFloat(pool.userReward) > 0 && <Btn onClick={() => onReward(pool)}>Claim reward</Btn>}
          </div>
      </div>
    </div>
  );
}