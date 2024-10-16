import chai from "chai";
import { ethers } from "hardhat";
import { solidity } from "ethereum-waffle";

import { deployContracts } from "../test-utils";
import {
  getSignerAddresses,
  makeCheckpoint,
  signHash,
  makeTxBatchHash,
  examplePowers
} from "../test-utils/pure";

chai.use(solidity);
const { expect } = chai;


async function runTest(opts: {}) {


  // Prep and deploy contract
  // ========================
  const signers = await ethers.getSigners();
  const gravityId = ethers.utils.formatBytes32String("foo");
  // This is the power distribution on the Cosmos hub as of 7/14/2020
  let powers = examplePowers();
  let validators = signers.slice(0, powers.length);
  const powerThreshold = 6666;
  const {
    gravity,
    testERC20,
    checkpoint: deployCheckpoint
  } = await deployContracts(gravityId, validators, powers, powerThreshold);


  // Transfer out to Cosmos, locking coins
  // =====================================
  await testERC20.functions.approve(gravity.address, 1000);
  await expect(gravity.functions.sendToSankofa(
    testERC20.address,
    "0xffffffffffffffffffffffffffffffffffffffff",
    1000
  )).to.emit(gravity, 'SendToCosmosEvent').withArgs(
      testERC20.address,
      await signers[0].getAddress(),
      "0x000000000000000000000000ffffffffffffffffffffffffffffffffffffffff",
      1000, 
      2
    );

  expect((await testERC20.functions.balanceOf(gravity.address))[0]).to.equal(1000);
  expect((await gravity.functions.state_lastEventNonce())[0]).to.equal(2);


    
  // Do it again
  // =====================================
  await testERC20.functions.approve(gravity.address, 1000);
  await expect(gravity.functions.sendToSankofa(
    testERC20.address,
    "0xffffffffffffffffffffffffffffffffffffffff",
    1000
  )).to.emit(gravity, 'SendToCosmosEvent').withArgs(
      testERC20.address,
      await signers[0].getAddress(),
      "0x000000000000000000000000ffffffffffffffffffffffffffffffffffffffff",
      1000, 
      3
    );

  expect((await testERC20.functions.balanceOf(gravity.address))[0]).to.equal(2000);
  expect((await gravity.functions.state_lastEventNonce())[0]).to.equal(3);
}

describe("sendToSankofa tests", function () {
  it("works right", async function () {
    await runTest({})
  });
});
