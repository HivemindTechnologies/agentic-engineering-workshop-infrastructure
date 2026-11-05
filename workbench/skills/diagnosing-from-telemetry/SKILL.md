---
name: diagnosing-from-telemetry
description: Use when a service is slow, failing, restarting or out of memory, or when someone reports a problem with a running service
---

# Diagnosing From Telemetry

First turn the complaint into one measurable quantity over one time window.
Touch nothing before that.

## The Ladder

Climb it in order. Each rung answers one question and hands the next a narrower
search space. The recipes below run against this cluster as they are.

1. **Metric: what and when.** Which quantity left its normal range, and at
   which minute?
2. **Event and log: what happened at that minute.** Kubernetes events, and the
   service's own lines around the change.
3. **Trace: where inside a request.** Which span carries the time or the error.
4. **Profile: which function.** CPU, heap, goroutines, mutex.

Each skipped rung leaves a guess in the diagnosis. When a rung yields nothing,
say so in the report; an empty result is evidence too.

## Recipes

Query with `curl` and `jq`, and filter before the output reaches you: a raw
answer costs thousands of tokens and says less.

**Metrics**, Prometheus on `localhost:9090`. Names come from the running
system; a query against an invented name returns an empty result, which reads
exactly like health.

```sh
curl -s localhost:9090/api/v1/label/__name__/values | jq -r '.data[]' | grep -i http_server

# Now: requests per second by route and status. Swap in any query below.
q='sum by (http_route, http_response_status_code) (rate(http_server_request_duration_seconds_count{service_name="scoring"}[5m]))'
curl -s localhost:9090/api/v1/query --data-urlencode "query=$q" | jq -c '.data.result[] | [.metric, .value[1]]'

# When: the same query over the last hour, one point per 5 minutes.
curl -s localhost:9090/api/v1/query_range --data-urlencode "query=$q" \
  --data-urlencode start=$(( $(date +%s) - 3600 )) --data-urlencode end=$(date +%s) --data-urlencode step=300 \
  | jq -r '.data.result[] | (.metric | tostring), (.values[] | "\(.[0] | todate)  \(.[1])")'
```

Useful queries for `scoring`:

- p99 latency: `histogram_quantile(0.99, sum by (le, http_route) (rate(http_server_request_duration_seconds_bucket{service_name="scoring"}[5m])))`
- memory: `container_memory_working_set_bytes{k8s_container_name="scoring"}`
- CPU: `container_cpu_usage{k8s_container_name="scoring"}`
- goroutines: `go_goroutine_count{service_name="scoring"}`

**Kubernetes**, read-only kubeconfig: restarts, kills and scheduling.

```sh
kubectl -n demo get pods -o wide
kubectl -n demo get events --sort-by=.lastTimestamp | tail -20
kubectl -n demo get pods -o json | jq -c '.items[] | {pod: .metadata.name, restarts: [.status.containerStatuses[].restartCount], last: [.status.containerStatuses[].lastState]}'
```

**Logs**, Loki on `localhost:3100`. Scoring logs only through OpenTelemetry:
`kubectl logs` for it is empty. Kubernetes events are a second stream,
`service_name="kubernetes-events"`. Count first, read lines second.

```sh
curl -s -G localhost:3100/loki/api/v1/query \
  --data-urlencode 'query=sum by (detected_level) (count_over_time({service_name="scoring"}[15m]))' \
  | jq -c '.data.result[] | [.metric, .value[1]]'
curl -s -G localhost:3100/loki/api/v1/query \
  --data-urlencode 'query=sum by (k8s_event_reason, k8s_object_name) (count_over_time({service_name="kubernetes-events"}[1h]))' \
  | jq -c '.data.result[] | [.metric, .value[1]]'
# The lines themselves, newest first, at most 5:
curl -s -G localhost:3100/loki/api/v1/query_range \
  --data-urlencode 'query={service_name="scoring"} | detected_level != "info"' --data-urlencode limit=5 \
  | jq -r '.data.result[] | .values[][1]'
```

**Traces**, Tempo on `localhost:3200`. Search with TraceQL, then read one trace
span by span.

```sh
curl -s -G localhost:3200/api/search \
  --data-urlencode 'q={resource.service.name="scoring" && duration > 100ms}' --data-urlencode limit=5 \
  | jq -c '.traces[]? | {traceID, rootTraceName, durationMs}'
curl -s localhost:3200/api/v2/traces/<traceID> \
  | jq -c '.trace.resourceSpans[].scopeSpans[].spans[] | {name, ms: (((.endTimeUnixNano|tonumber) - (.startTimeUnixNano|tonumber)) / 1e6)}'
```

**Profiles**, Pyroscope on `localhost:4040`: the functions with the most CPU
time of their own over the last 15 minutes. For heap or goroutines, swap the
profile type; the list of types comes from the first command.

```sh
curl -s localhost:4040/querier.v1.QuerierService/ProfileTypes -H 'content-type: application/json' -d '{}' | jq -r '.profileTypes[].ID'
curl -s -G localhost:4040/pyroscope/render \
  --data-urlencode 'query=process_cpu:cpu:nanoseconds:cpu:nanoseconds{service_name="scoring"}' \
  --data-urlencode from=now-15m --data-urlencode until=now --data-urlencode format=json \
  | jq -r '.flamebearer as $f | [$f.levels[] as $l | range(0; $l | length; 4) as $i | {n: $f.names[$l[$i+3]], self: $l[$i+2]}]
      | group_by(.n) | map({n: .[0].n, self: (map(.self) | add)}) | sort_by(-.self) | .[:10][] | "\(.self)\t\(.n)"'
```

## Evidence Format

Every claim carries its query and its answer, so another person can rerun it
and get the same number:

```
claim:  scoring answers 0.2 requests/s since 14:22, down from 5
query:  sum(rate(http_server_request_duration_seconds_count{service_name="scoring"}[5m]))
value:  0.2
```

A claim without a query below it is a hypothesis and must be labelled one.

## One Hypothesis at a Time

State it as: *X causes Y, because Z*. Then name a second observation that must
hold if X is true, and check that one too. A hypothesis that predicts nothing
new is not testable.

## Distinguish the Layer

A container that dies and a process that misbehaves look alike from outside.
Separate them: the platform reports what it did to the container (restarts,
kills, throttling), the process reports what it was doing (goroutines, heap,
spans, frames). A cause on one layer often shows as a symptom on the other.

## Fix the Cause

Raising a limit, adding a retry, or restarting on a timer removes the symptom
and keeps the cause. If you propose one of those, say plainly that it is
mitigation and name the real fix.

## Close the Loop

Rerun the same queries after the fix. The report shows before and after from
the same query.
