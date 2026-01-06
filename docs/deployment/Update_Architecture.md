Subject: Migration Plan: Split Domain to Unified Domain with Cloudflare Tunnel & Zero Trust

Context: I have an existing Microservices E-commerce architecture (Go, Astro, Ory Kratos, Postgres) currently deployed on a VPS using Podman Compose. Current status:

Issue: I am using a split domain strategy (karimasyari.com for frontend, ks-backend.cloud for API/Auth). This is causing cross-site cookie issues for authentication.

Security Risk: Ports 8080 (API) and 4433 (Kratos) are currently exposed directly on the VPS public IP.

OS: Linux VPS.

Orchestrator: Podman Compose (rootless).

Objective: I need a step-by-step migration plan to refactor this into a secure, unified domain architecture using Cloudflare services.

Specific Requirements:

Unified Domain Strategy:

Move ALL public services to *.karimasyari.com (Storefront, API, Auth, Admin).

Use monitor.ks-backend.cloud ONLY for internal infrastructure tools.

Cloudflare Tunnel Implementation:

Replace direct port exposure with cloudflared (Cloudflare Tunnel).

The VPS firewall must be closed (Deny Inbound) except for SSH.

Provide the cloudflared service configuration for docker-compose.prod.yml.

Access Control (Zero Trust):

Protect admin.karimasyari.com with Cloudflare Access.

Protect monitor.ks-backend.cloud with Cloudflare Access.

api.karimasyari.com and auth.karimasyari.com must remain public but proxied via Tunnel.

Monitoring Stack:

Add a monitoring stack (Prometheus + Grafana + Portainer) to the docker-compose.prod.yml.

These should be accessible via monitor.ks-backend.cloud subpaths or subdomains.

Deliverables: Please provide:

Updated docker-compose.prod.yml: Including the backend, kratos, postgres, redis, plus the new tunnel service and monitoring stack.

Updated .env.production template: Reflecting the new domain structure (CORS, Cookie Domain) and Tunnel tokens.

Cloudflare Configuration Guide: Briefly explain how to route the ingress in Cloudflare Dashboard (Public Hostnames).

Migration Checklist: A safe sequence of commands to execute on the VPS to switch over without losing data.