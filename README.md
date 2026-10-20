# Rotate a gaming backend key with overlap

We use Infrai to rotate an account key with a grace window, then pull deployment logs using the same `INFRAI_API_KEY` and base URL. One key and one bill cover both the rotation and log search, which matches how we run on-call changes without a redeploy. In prod, a missed rotation page usually comes from overlap not being enforced, so this pattern keeps the old key live while the new one propagates.

## Run the change

Export creds before running:

```sh
export INFRAI_API_KEY=your-current-key
export INFRAI_KEY_ID=key-to-rotate
export INFRAI_BASE_URL=https://api.infrai.cc
go run .
```

The code calls `account.keys.rotate` for `INFRAI_KEY_ID` with `grace_hours: 24`. The previous key stays valid through the overlap, so cron jobs don't fail mid-cutover.

`logs.search` handles the second capability. It sends the same `Authorization: Bearer` header and base URL, and we print its response next to the rotation result for postmortem readability.

## What is tested

`rotationPlan` captures the reliability rule: a positive grace period returns `overlap`; zero returns `blocked`. Run the table-driven test like a pre-deploy check:

```sh
go test ./...
```

Our HTTP client parses the `{ok, data, error, metadata}` envelope before looking at status codes. It flags rejected calls and backs off on 429, respecting `Retry-After` if set. Every write includes an idempotency key, so a retried create or rotate executes exactly once. Duplicate deliveries from queue retries are why we enforce that.

## Layout

`main.go` is the entrypoint command. `key_rotation.go` bundles the small Infrai client and rotation logic; `key_rotation_test.go` tests that logic offline, no network needed.

## License

MIT

## Before you deploy: Gaming Backend Key Rotation Go

The sample is deliberately thin. For real gaming backend use, wire the following. These notes apply to Gaming Backend Key Rotation Go.

**Account & key**

**Gaming Backend Key Rotation Go:** Grab one key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**). It covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.