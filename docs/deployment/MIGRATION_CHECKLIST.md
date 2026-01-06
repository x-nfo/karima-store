# Migration Checklist: Zero Downtime Arch Update

Follow this checklist to switch from the old split-domain architecture to the new Unified Domain + Tunnel architecture.

## Phase 1: Preparation (Local)
- [ ] **Generate Tunnel Token**: Create the tunnel in Cloudflare and get the token.
- [ ] **Commit Changes**: Git add/commit/push the updated `docker-compose.prod.yml` and docs.
  ```bash
  git add .
  git commit -m "chore: update architecture to cloudflare tunnel"
  git push origin main
  ```

## Phase 2: Preparation (VPS)
- [ ] **SSH into VPS**.
- [ ] **Pull Latest Code**:
  ```bash
  cd /path/to/karima_store
  git pull origin main
  ```
- [ ] **Update `.env.production`**:
  - Add `TUNNEL_TOKEN=ey...`
  - Update `CORS_ORIGIN` to `https://karimasyari.com,https://www.karimasyari.com` (and others)
  - Update `KRATOS_` URLs to `https://auth.karimasyari.com...`
  ```bash
  nano .env.production
  ```
- [ ] **Verify Environment**:
  ```bash
  # Check if env vars are loaded correctly (dry run)
  podman-compose -f docker-compose.prod.yml config
  ```

## Phase 3: Switch Over
- [ ] **Stop Old Services**:
  ```bash
  podman-compose -f docker-compose.prod.yml down
  ```
- [ ] **Start New Stack (with Tunnel)**:
  ```bash
  podman-compose -f docker-compose.prod.yml up -d
  ```
- [ ] **Verify Tunnel Connection**:
  - Check Cloudflare Dashboard: Tunnel status should be **Healthy**.
  - Check logs:
    ```bash
    podman logs karima_tunnel
    ```

## Phase 4: Validation (IMPORTANT)
- [ ] **Frontend**: access `https://karimasyari.com`. Does it load?
- [ ] **API**: `curl https://api.karimasyari.com/health`.
- [ ] **Auth**: Try to Register/Login. Check browser network tab.
  - Cookies should now be on `.karimasyari.com` (or `auth.karimasyari.com` if not explicitly set to wildcard).
  - No more 3rd party cookie warnings.
- [ ] **Monitor**: Access `https://monitor.karimasyari.com`.
  - Should prompt for Cloudflare Access (Email login).

## Phase 5: Cleanup
- [ ] **Close Firewall Ports**:
  - Once everything works via the Tunnel, BLOCK port 8080 and 4433 on the VPS Firewall (UFW/IPTables/Cloud Provider Firewall).
  - ONLY Port 22 (SSH) needs to be open.
