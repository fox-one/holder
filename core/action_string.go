// Code generated by "stringer -type Action -trimprefix Action"; DO NOT EDIT.

package core

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ActionSysWithdraw-1]
	_ = x[ActionSysProperty-2]
	_ = x[ActionProposalMake-3]
	_ = x[ActionProposalShout-4]
	_ = x[ActionProposalVote-5]
	_ = x[ActionPoolDonate-6]
	_ = x[ActionPoolGain-7]
	_ = x[ActionVaultLock-8]
	_ = x[ActionVaultRelease-9]
	_ = x[ActionPoolPardon-10]
	_ = x[ActionPoolPardonAll-11]
}

const _Action_name = "SysWithdrawSysPropertyProposalMakeProposalShoutProposalVotePoolDonatePoolGainVaultLockVaultReleasePoolPardonPoolPardonAll"

var _Action_index = [...]uint8{0, 11, 22, 34, 47, 59, 69, 77, 86, 98, 108, 121}

func (i Action) String() string {
	i -= 1
	if i < 0 || i >= Action(len(_Action_index)-1) {
		return "Action(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Action_name[_Action_index[i]:_Action_index[i+1]]
}
