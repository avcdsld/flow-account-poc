transaction(
    newPublicKeyAsHexString: String,
) {
    prepare(signer: AuthAccount) {
        signer.keys.add(
            publicKey: PublicKey(publicKey: newPublicKeyAsHexString.decodeHex(), signatureAlgorithm: SignatureAlgorithm.ECDSA_secp256k1),
            hashAlgorithm: HashAlgorithm.SHA3_256,
            weight: 999.0
        )
    }
}
