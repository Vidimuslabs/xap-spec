# Image provenance — verifying the XAP node container

Vidimus signs the `xap-server` container image with [cosign](https://docs.sigstore.dev/)
so you can prove, before you run it, that the image you pulled was published by
Vidimus and has not been altered. This is non-repudiable: the signature could only
have been produced by the holder of the Vidimus image-signing key.

## Public key

`xap-server-cosign.pub` in this directory is the Vidimus image-signing public key.
The private half is held by Vidimus (a KMS/HSM-backed key for generally-available
builds).

## Verify an image

```sh
cosign verify \
  --key xap-server-cosign.pub \
  ghcr.io/vidimuslabs/xap-server@sha256:<digest>
```

A successful verification prints the signed claims and exits 0; verify that the
`docker-manifest-digest` in the output matches the digest you are about to run.

Pin and run by **digest**, not by tag, so the bytes you verified are the bytes you
run:

```sh
docker run ghcr.io/vidimuslabs/xap-server@sha256:<digest>
```

## Notes

- Signatures are currently published **without** a public transparency-log entry
  (`--tlog-upload=false`), so verification does not require reaching public Sigstore
  infrastructure — appropriate for air-gapped and pre-GA deployments. Add
  `--insecure-ignore-tlog=true` to `cosign verify` while no tlog entry exists.
- Every image also carries an SBOM. Supply-chain hardening (SLSA provenance
  attestation, KMS-backed signing key) lands with the GA build.
