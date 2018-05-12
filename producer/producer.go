// Copyright 2017~2022 The Bottos Authors
// This file is part of the Bottos Chain library.
// Created by Rocket Core Team of Bottos.

//This program is free software: you can distribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.

//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.

//You should have received a copy of the GNU General Public License
// along with bottos.  If not, see <http://www.gnu.org/licenses/>.

/*
 * file description:  producer entry
 * @Author:
 * @Date:   2017-12-06
 * @Last Modified by:
 * @Last Modified time:
 */
package producer

import (
	"fmt"

	"github.com/bottos-project/core/chain"
	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/config"
	"github.com/bottos-project/core/role"
)

type Reporter struct {
	core     chain.BlockChainInterface
	roleIntf role.RoleInterface
	state    ReportState
}
type ReporterRepo interface {
	Woker(Trxs []*types.Transaction) *types.Block
	VerifyTrxs(Trxs []*types.Transaction) error
	IsReady() bool
}

func New(b chain.BlockChainInterface, roleIntf role.RoleInterface) ReporterRepo {
	stat := ReportState{0, "", "", false, 0, false}
	return &Reporter{core: b, roleIntf: roleIntf, state: stat}
}

func (p *Reporter) Woker(trxs []*types.Transaction) *types.Block {

	now := common.NowToSeconds()
	slot := p.roleIntf.GetSlotAtTime(now)
	//scheduledTime := p.roleIntf.GetSlotTime(slot)
	accountName, err1 := p.roleIntf.GetCandidateBySlot(slot)
	if err1 != nil {
		return nil // errors.New("report Block failed")
	}
	block, err := p.reportBlock(p.state.ScheduledTime, accountName, trxs)
	if err != nil {
		return nil // errors.New("report Block failed")
	}

	//fmt.Println("brocasting block", block)
	return block
} /*TODO
func (p *Reporter) IsValid(block *blockchain.BlockData, receivedAt int64) (valid bool, err error) {
	slot := p.getSlotAt(receivedAt)
	requiredTS := p.GetCurrentSlotStart(receivedAt)

	position := p.GetPosition(block.Signer)
	if position != slot {
		return false, fmt.Errorf("not in required slot: %d, producer position = %d", slot, position)
	}
	diff := requiredTS - block.GetTimestamp()
	if diff < 0 {
		diff = -diff
	}
	valid = diff < Epsilon
	if !valid {
		err = fmt.Errorf("incorrect timestamp: %d. Required: %d ± %v", block.GetTimestamp(), requiredTS, Epsilon)
	}
	return
}*/
func (p *Reporter) StartTag() error {
	//p.core.

	return nil

}
func (p *Reporter) VerifyTrxs(trxs []*types.Transaction) error {

	return nil
}

//func reportBlock(reportTime time.Time, reportor role.Delegate) *types.Block {
func (p *Reporter) reportBlock(blockTime uint64, accountName string, trxs []*types.Transaction) (*types.Block, error) {
	head := types.NewHeader()
	head.PrevBlockHash = p.core.HeadBlockHash().Bytes()
	head.Number = p.core.HeadBlockNum() + 1
	head.Timestamp = blockTime
	head.Delegate = []byte(accountName)
	block := types.NewBlock(head, trxs)
	//TODO
	block.Header.DelegateSign = block.Sign("123").Bytes()
	// If this block is last in a round, calculate the schedule for the new round
	if block.Header.Number%config.BLOCKS_PER_ROUND == 0 {
		newSchedule := p.roleIntf.ElectNextTermDelegates()
		fmt.Println("next term delgates", newSchedule)

		currentState, err := p.roleIntf.GetCoreState()
		if err != nil {
			return nil, err
		}
		block.Header.DelegateChanges = common.Filter(currentState.CurrentDelegates, newSchedule)
		//fmt.Println("DelegateChanges", block.Header.DelegateChanges)
	}

	return block, nil
}
