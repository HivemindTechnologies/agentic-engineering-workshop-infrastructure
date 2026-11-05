# 💡 Secrets

**Encrypted in Git, Key Outside the Sandbox.** A SOPS file is safe in the repository.
Whoever reaches its key reads every secret in it.

----

## SOPS

- Encrypts the values, keeps the keys: a diff stays readable
- The key comes from a backend: age, GPG, a cloud KMS, Vault
- age: one key file, works offline

----

## Who holds the key

- The key file: not granted, so not readable
- `SOPS_AGE_KEY`: readable, unless `deny_vars` lists it
- GPG: same for its key files; also keep out its agent socket

---

## 🛠️ Secrets

**A secret encrypted with SOPS. Can the agent read it?**

1. Make a key and a secret: the commands are in `exercises/part1.md`
2. `SOPS_AGE_KEY_FILE=$key sops decrypt secret.sops.yaml` works; then `SOPS_AGE_KEY_FILE=$key just sandbox` and the same: which line of the error says why?
3. `SOPS_AGE_KEY=$(grep AGE-SECRET $key) just sandbox`, decrypt again
4. Make decrypt fail in that sandbox too, with the key still in the variable. Commit `nono.json`

The public key encrypts, the private key decrypts: the agent may hold the first, never the second.

