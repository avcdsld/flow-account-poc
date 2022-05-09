transaction(
    publicKeyAsHexString1: String,
    publicKeyAsHexString2: String,
    publicKeyAsHexString3: String
) {
    prepare(signer: AuthAccount) {
        let account = AuthAccount(payer: signer)

        account.keys.add(
            publicKey: PublicKey(publicKey: publicKeyAsHexString1.decodeHex(), signatureAlgorithm: SignatureAlgorithm.ECDSA_secp256k1),
            hashAlgorithm: HashAlgorithm.SHA3_256,
            weight: 999.0
        )
        account.keys.add(
            publicKey: PublicKey(publicKey: publicKeyAsHexString2.decodeHex(), signatureAlgorithm: SignatureAlgorithm.ECDSA_secp256k1),
            hashAlgorithm: HashAlgorithm.SHA3_256,
            weight: 1000.0
        )
        account.keys.add(
            publicKey: PublicKey(publicKey: publicKeyAsHexString3.decodeHex(), signatureAlgorithm: SignatureAlgorithm.ECDSA_secp256k1),
            hashAlgorithm: HashAlgorithm.SHA3_256,
            weight: 1.0
        )
    }
}
