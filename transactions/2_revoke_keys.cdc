transaction(
    keyIndex1: Int,
    keyIndex2: Int,
) {
    prepare(signer: AuthAccount) {
        signer.keys.revoke(keyIndex: keyIndex1)
        signer.keys.revoke(keyIndex: keyIndex2)
    }
}
