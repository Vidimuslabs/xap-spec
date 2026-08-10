# Image provenance — verifying the XAP node container

Vidimus signs the `xap-server` container image with [cosign](https://docs.sigstore.dev/)
so you can prove, before you run it, that the image you pulled was published by
Vidimus and has not been altered. This is non-repudiable: the signature could only
have been produced by the holder of the Vidimus image-signing key.

## Access to the image

The image is distributed to evaluators rather than published for anonymous pull.
`ghcr.io/vidimuslabs/xap-server` requires authentication: an anonymous
`docker pull` or `cosign verify` against it returns `401 Unauthorized`. To request
access, email <hello@vidimuslabs.com> with the GitHub username you want granted —
access is a read role on the package, and you then authenticate once with
`docker login ghcr.io` before running the commands below.

Everything needed to *audit* this procedure is public regardless: the signing key
is in this directory and the commands are below, so you can see exactly what a
verification checks before asking anyone for anything.

**This page is about image provenance, not receipt verification — and only the
former needs the image.** Image provenance answers "did Vidimus publish these
bytes"; it is a supply-chain question about our container. Receipt verification
answers "was this action actually authorized", which is the protocol's own claim,
and it requires no container, no credentials, and no access to any enforcement
point. You can exercise it in full today:
[xap-go](https://github.com/Vidimuslabs/xap-go) replays every conformance vector
in this repository with `xap vectors run`, and the same verifier runs in your
browser at [vidimuslabs.com/verify](https://www.vidimuslabs.com/verify). Nothing
about XAP's independent verifiability depends on being able to pull this image.

## Public key

`xap-server-cosign.pub` in this directory is the Vidimus image-signing public key.
The private half is held by Vidimus (a KMS/HSM-backed key for generally-available
builds).

## Verify an image

```sh
cosign verify \
  --key xap-server-cosign.pub \
  --insecure-ignore-tlog=true \
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
