import LockedTokens from 0x8d0e87b65159ae63

pub fun main(account: Address): Address? {

	let pubAccount = getAccount(account)
	let lockedAccountInfoCap = pubAccount.getCapability<&LockedTokens.TokenHolder{LockedTokens.LockedAccountInfo}>(LockedTokens.LockedAccountInfoPublicPath)

	if let lockedAccountInfoRef = lockedAccountInfoCap.borrow() {
		return lockedAccountInfoRef.getLockedAccountAddress()
	}

	return nil
}
