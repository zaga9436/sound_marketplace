package service

const platformCommissionPercent int64 = 10

func engineerPayoutAmount(grossAmount int64) int64 {
	if grossAmount <= 0 {
		return 0
	}
	commission := grossAmount * platformCommissionPercent / 100
	payout := grossAmount - commission
	if payout < 0 {
		return 0
	}
	return payout
}
