# Part 1: Exercises

One section per 🛠️ slide. Four sessions of 50 minutes: §1 and §2, §3 to §6,
§7 to §9, §10 and §11.

## Before You Start

Everything runs in `workbench/`; nothing gets installed on your system.

- `nix develop`: a shell with every tool, pinned
- `just`: lists every command; `justfile` shows what each one runs

## 1. Sandbox

Box a command you run every day, and see what the box stops. Run
everything inside `nix develop`:

```sh
nono run --allow . -- CMD
```

NixOS: add `--read /nix/store`, where every program lives.

1. **Run:** `ls ~`, `env`, `curl https://example.com`: what got through?
2. **Network:** block it, then allow one domain:
   `nono run --allow . --block-net -- curl https://example.com`,
   `nono run --allow . --allow-domain example.com -- curl https://example.com`
3. **Why:** ask `nono why` about one denial, e.g.
   `nono why --path ~/.ssh --op read --allow .`
4. **Yours:** `terraform plan`, `kubectl get pods` or your own script: what
   broke, and which flag fixes it?

## 2. Setup

### Cluster

```sh
just up        # first run pulls images: a few minutes
just urls      # Argo CD, Grafana (admin/admin), ...
just signals   # a metric, a log line and a trace of scoring
just status    # pods, and which commit Argo CD runs: compare `git log -1`
```

`.kubeconfig` here is read-only, for you and the agent. The admin kubeconfig
and the push key live in `~/.local/state/agentic-workshop-infrastructure/`,
outside the sandbox; `nix develop .#admin` uses the admin kubeconfig.
`kubectl auth whoami` tells you which one you are.

### Sandbox

```sh
export SECRET_TOKEN=fake-token
just sandbox
```

`just sandbox` is the agent's sandbox, as a shell. Guess, then run:

```sh
echo "$SECRET_TOKEN"
ls ~/.ssh
curl -sI https://google.de | head -1
cat .kubeconfig
cat ~/.local/state/agentic-workshop-infrastructure/admin.kubeconfig
cat ~/.local/state/agentic-workshop-infrastructure/push-key
kubectl -n demo delete pod -l app=scoring
```

For each command: which layer answered? nono, Kubernetes RBAC, or nothing.

Find your `delete` in the cluster's audit log:

```sh
exit
just audit | grep agent
```

**Optional:** find any way out of `just sandbox`: a file, a socket, a
credential, a command that runs outside the box. Post it in the chat.

### Claude

In the terminal where you exported `SECRET_TOKEN`:

1. `just claude`, then `/login` once: the login lands in `.claude-home/`,
   apart from yours. No browser? Open the printed link, paste the code.
2. `/status`: which model? `/context`: what fills the window before your
   first message? Which tools does it have? `/` lists its commands.
3. Ask it to read `$SECRET_TOKEN`: does it see it?
4. `Shift+Tab` cycles the modes: what changes?

## 3. Secrets

In `nix develop`, outside the sandbox. A workshop-only key, so your real one
stays untouched:

```sh
key=~/.local/state/agentic-workshop-infrastructure/age.key
age-keygen -o $key                    # prints your public key, age1…; again? rm $key first
printf 'password: hunter2\n' > secret.yaml
sops encrypt --age "$(age-keygen -y $key)" secret.yaml > secret.sops.yaml
rm secret.yaml
cat secret.sops.yaml                  # keys readable, values encrypted
```

Decrypt, then the same from the sandbox: once with the key file, once with
the key in a variable:

```sh
SOPS_AGE_KEY_FILE=$key sops decrypt secret.sops.yaml
SOPS_AGE_KEY_FILE=$key just sandbox
sandbox $ sops decrypt secret.sops.yaml   # which line of the error says why?
sandbox $ exit
SOPS_AGE_KEY=$(grep AGE-SECRET $key) just sandbox
sandbox $ sops decrypt secret.sops.yaml
```

Now make decrypt fail in that sandbox too, with the key still in the
variable: change `nono.json`, not the command. Commit `nono.json`: the
stricter profile stays, and the push in GitOps needs a clean tree.

The public key encrypts, the private key decrypts: the agent may hold the
first, never the second.

## 4. Variance

Two terminals, `just claude` in each, the same prompt:

```text
Write a HorizontalPodAutoscaler for scoring. Print the YAML, do not edit files.
```

Compare the YAML, then the notes: what did only one session warn about?
`/clear` in both, and the same prompt plus:

```text
Interview me relentlessly about every aspect of this plan until we reach a shared understanding.
```

Give both the same answers and compare again. Which YAML would you merge,
and why?

## 5. GitOps

Ask `just claude` for a change in `deploy/`, one you would make at work, e.g.:

```text
Run scoring with two replicas.
```

In a second terminal, in `nix develop`, review and deploy. What did the
hooks run? How long until Argo CD runs your commit?

```sh
git show
git push cluster
just status
```

Now skip Git: scale to three, first with the agent's kubeconfig, then as
admin. Watch for a minute: what happened each time, and who did it?

```sh
kubectl -n demo scale deployment scoring --replicas=3   # this shell
nix develop .#admin
kubectl -n demo scale deployment scoring --replicas=3   # admin shell
kubectl -n demo get deployment scoring -w               # Ctrl+C after a minute
```

## 6. Rollback

Ship a change `just check` cannot catch: in `deploy/demo/scoring.yaml`, point
the readiness probe at `/ready`, a path scoring does not serve. Commit,
`git push cluster`, `just status`: are the old pods still serving? Then ask
`just claude`:

```text
scoring's new pod never gets ready. Find out why and roll back.
```

Review and push, until every scoring pod is ready.

## 7. Spec

Every pod can reach scoring. Only the customers, namespace `traffic`, should.

1. Probe scoring from another namespace, in `nix develop .#admin` (starting a
   pod needs admin). Today it prints `reachable`:

   ```sh
   kubectl -n default \
     run probe -q --rm -i --restart=Never --image=alpine:3.24.2 -- \
     sh -c 'sleep 2; wget -T 3 -qO /dev/null http://scoring.demo:8080/healthz && echo reachable'
   ```

   `sleep 2` gives kubectl time to attach, so the answer prints once.

2. Ask a fresh `just claude`:

   ```text
   Only the customers (namespace traffic) may reach scoring. Interview me about this, one question at a time, then write deploy/demo/SPEC.md. Each acceptance criterion names the command that checks it. Do not implement yet.
   ```

3. Review the spec before any code: what else talks to scoring, or does
   scoring talk to?
4. Then:

   ```text
   Implement deploy/demo/SPEC.md and commit.
   ```

5. `git push cluster`, `just status` until synced, then prove each criterion
   with its command. The probe from step 1 now fails: `wget: download timed out`.

## 8. Quality Gate

1. `kube-linter checks list` has a check for pods that no NetworkPolicy
   selects. Have `just claude` enable it in `.kube-linter.yaml`, make
   `just check` pass, and commit.
2. Every other workload gets a policy, or an ignore annotation with its
   reason. Is each reason true?
3. Delete scoring's NetworkPolicy: `just check` must fail. Restore it and
   push:

   ```sh
   git grep -l 'kind: NetworkPolicy' deploy/   # scoring's selects app: scoring
   just check                                   # after the delete: fails
   git restore deploy/
   git push cluster
   ```

## 9. Skills

With `just claude`, write `skills/working-in-a-git-worktree/SKILL.md`: start
a change in a worktree under `.worktrees/` (inside `workbench/`: the sandbox
writes nowhere else), commit there, merge into `main`, clean up. Its
description says when to use it. Review: would you work this way yourself?
Commit, then a fresh `just claude` (it copies `skills/` at start):

```text
Run scoring with three replicas.
```

Does it load the skill and work in a worktree? Push the result.

## 10. Plugin

Build it by hand or with `just claude`:

```text
my-plugin/
  .claude-plugin/plugin.json          { "name": "my-plugin", "version": "0.1.0" }
  skills/diagnosing-from-telemetry/   copied from skills/
  skills/working-in-a-git-worktree/   copied from skills/, if you wrote it
```

1. `claude plugin validate my-plugin`
2. `just skills=false claude --plugin-dir my-plugin`, then `/skills`: which
   skills come from the plugin? Ask for a change: does it work in a worktree?

**Optional:** make it a marketplace, push, and let a teammate install it.

## 11. Debugging Agent

One session with everything from today: the code, the cluster, the telemetry
skill (`just claude` loads `skills/`), and the scoring pod's debug MCP. No
incident: the normal customer traffic is enough.

1. In a second terminal, `just debug`. It port-forwards the scoring pod's
   debug MCP to 127.0.0.1:7070 with the admin kubeconfig, until Ctrl-C. The
   agent's kubeconfig cannot: only you open this door.
2. Start the session with the debug MCP:

   ```sh
   just claude --mcp-config setup/debug-mcp.json
   ```

3. Check with `/skills` that `diagnosing-from-telemetry` is listed, and
   with `/mcp` that `scoring-debug` is connected.
4. Ask:

   ```text
   Which tools do you have?
   ```

5. Then:

   ```text
   Customers ask: do large events score slower? Answer with proof: a file and line, a query, or a tool call for each claim.
   ```

6. Review the answer. Which claim came from the code, which from telemetry,
   which from the pod? Where would the fix go, and which check would prove
   it?
