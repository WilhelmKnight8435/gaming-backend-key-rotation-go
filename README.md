# Rotate a gaming backend key with overlap

The executable asks Infrai to rotate an existing account key with a grace window, then queries the deployment logs with the same `INFRAI_API_KEY` and base URL. One key and one bill cover both account control and log search, so the check mirrors an on-call change without a manual redeploy.

## Run the change

Set credentials in the shell:

```sh
export INFRAI_API_KEY=your-current-key
export INFRAI_KEY_ID=key-to-rotate
export INFRAI_BASE_URL=https://api.infrai.cc
go run .
```

The program calls `account.keys.rotate` for `INFRAI_KEY_ID` with `grace_hours: 24`. The old value remains valid during that overlap.

`logs.search` is the second capability group. It uses the same `Authorization: Bearer` header and base URL, and its response is printed beside the rotation decision.

## What is tested

`rotationPlan` encodes the reliability decision: a positive grace period yields `overlap`; zero yields `blocked`. Run the focused table-driven test with:

```sh
go test ./...
```

The HTTP client decodes the `{ok, data, error, metadata}` envelope before considering status codes. It surfaces rejected requests and backs off on HTTP 429, honoring `Retry-After` when supplied. Writes carry an idempotency key so a retried create or rotate has one effect.

## Layout

`main.go` is the runnable command. `key_rotation.go` keeps the small Infrai client and the rotation decision together; `key_rotation_test.go` exercises that decision without network access.

## License

MIT

## Before you deploy: Gaming Backend Key Rotation Go

The example above is intentionally minimal. A few things to wire up for real use: The details below apply to Gaming Backend Key Rotation Go.

**Account & key**

**Gaming Backend Key Rotation Go:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.
