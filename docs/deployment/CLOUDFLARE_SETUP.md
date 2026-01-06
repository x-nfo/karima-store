# Cloudflare Tunnel & Zero Trust Setup Guide

This guide explains how to configure Cloudflare Tunnel and Zero Trust to secure the Karima Store infrastructure with a unified domain (`*.karimasyari.com`).

## 1. Prerequisites
- A Cloudflare account.
- The domain `karimasyari.com` must be active on Cloudflare (authoritative DNS).
- (Optional) `ks-backend.cloud` can be removed or redirected later.

## 2. Create Cloudflare Tunnel
1.  Go to **Cloudflare Dashboard** > **Zero Trust**.
2.  Navigate to **Networks** > **Tunnels**.
3.  Click **Create a Tunnel**.
4.  Select **Cloudflared**.
5.  Name the tunnel: `karima-production-vps`.
6.  **Install Connector**:
    - Choose **Docker** as the environment.
    - Copy the **Token** (the long string after `--token` in the command).
    - **Do NOT run the command yet.** Just save the token.
    - Paste this token into your `.env.production` file as `TUNNEL_TOKEN`.

## 3. Configure Public Hostnames (Ingress)
In the Tunnel configuration page (or "Public Hostname" tab), add the following mappings.
**Note**: The service `http://<service_name>:<port>` refers to the Docker Compose service names and internal ports.

| Domain (Public) | Path (Optional) | Service (Internal) | Notes |
| :--- | :--- | :--- | :--- |
| `karimasyari.com` | `/` | `http://frontend:3000` | (If frontend is on VPS) * |
| `api.karimasyari.com` | `/` | `http://backend:8080` | Main API |
| `auth.karimasyari.com` | `/` | `http://kratos:4433` | Kratos Public API |
| `admin.karimasyari.com` | `/` | `http://frontend-admin:3000` | (If Admin UI is on VPS) * |
| `media.karimasyari.com` | `/` | `http://minio:9000` | (If using Minio/Local) ** |
| `monitor.karimasyari.com` | `/grafana` | `http://grafana:3000` | Grafana Dashboard |
| `monitor.karimasyari.com` | `/prometheus` | `http://prometheus:9090` | Prometheus |

*> **Note on Frontend**: If your frontend (Next.js/Astro) is hosted on Vercel/Netlify, you do **NOT** route `karimasyari.com` through this tunnel. You point the A record to Vercel/Netlify directly. This tunnel is mainly for the **Backend** and **Auth** on the VPS.

**If Frontend is external (Vercel), only map:**
- `api.karimasyari.com` -> `http://backend:8080`
- `auth.karimasyari.com` -> `http://kratos:4433`
- `monitor.karimasyari.com` -> `http://grafana:3000` (plus others)

## 4. configure Zero Trust Access (Security)
Protect sensitive subdomains so only YOU can access them.

1.  Go to **Zero Trust** > **Access** > **Applications**.
2.  **Add an Application** > **Self-hosted**.
3.  **Application Name**: `Karima Monitor`.
4.  **Subdomain**: `monitor` . `karimasyari.com`.
5.  **Identity Providers**: Enable your preferred login method (e.g., One-time PIN to your email).
6.  **Policies**:
    - Name: `Admin Only`
    - Action: `Allow`
    - Include: `Email` = `your-email@example.com`
7.  **Save**.

**Repeat for:**
- `admin.karimasyari.com` (if hosted on VPS)
- `monitor.karimasyari.com`

**Do NOT enable Access for:**
- `api.karimasyari.com` (Must be public for the frontend to call)
- `auth.karimasyari.com` (Public auth endpoints)

## 5. SSL/TLS
Cloudflare handles SSL automatically.
- Ensure **SSL/TLS** mode in Cloudflare Dashboard is set to **Full** or **Strict**.

## 6. Update DNS
When you save the Public Hostnames in the Tunnel config, Cloudflare **automatically** creates the CNAME records for you. You do not need to manually edit DNS records.

## 7. Firewall (VPS)
Once the Tunnel is working (status: Healthy), you can **CLOSE** all inbound ports on your VPS firewall except SSH (22).
- Deny 80/tcp
- Deny 443/tcp
- Deny 8080/tcp
- Deny 4433/tcp
- Allow 22/tcp (SSH)
