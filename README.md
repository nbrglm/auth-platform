# NAP – NBRGLM Auth Platform

**NAP (NBRGLM Auth Platform)** is a self-hostable, multi-tenant authentication platform offering secure session management, tenant isolation, and a customizable UI — all without any paywalls or vendor lock-in.

Built to be simple to use, easy to deploy, and flexible for both startups and individual developers who want full control over their authentication system.

> ⚠️ Codebase is currently being cleaned up and will be pushed soon. Stay tuned!

## Features

- Multi-Tenant Architecture  
- Email/Password Authentication  
- Secure Session Management  
- Fully Self-Hostable  
- Customizable UI & REST API  
- (Coming Soon) MFA Support  
- (Coming Soon) Admin Panel  
- (Planned) Social Login via OIDC  
- (Planned) Act as OIDC Provider

## Documentation

Official documentation is a work in progress and will be available soon at:

**https://docs.nbrglm.com**

## Why NAP?

Most authentication solutions today are either closed-source, expensive, or unnecessarily complex. NAP was built to offer:

- A **fully self-hostable**, production-ready auth platform  
- No paywalls, no feature gates — just clean APIs and UI  
- Simple, fast setup for developers who value autonomy  
- Built with modern tech and scalable design from the start

## Tech Stack

- **Backend:** Go  
- **Frontend:** HTML, TailwindCSS, Alpine.js  
- **Database:** PostgreSQL  
- **Cache/Session Store:** Redis (used for CSRF, session storage, etc.)  
- **Security:** Vault (for secret injection)  
- **Containerization:** Docker  
- **Orchestration Ready:** Kubernetes (K3s, etc.)

## Installation

**Currently:** Docker Compose

You can run NAP locally using Docker Compose. While Kubernetes support is coming via Helm, it's adaptable manually if needed. Running the Go binary directly is not recommended for production.

```bash
# Optional: clone into a dedicated folder
mkdir -p "$HOME/NBRGLM" && cd "$HOME/NBRGLM"
git clone https://github.com/nbrglm/auth-platform.git
cd auth-platform

# Start with Docker Compose
docker-compose up -d

task dev
```

> **Note:** On Windows:
> * Command Prompt: Use `%USERPROFILE%` instead of `$HOME`.
> * PowerShell: Use `$env:USERPROFILE` instead of `$HOME`.
>
> The first run will take some time as it builds the Docker images. Subsequent runs will be faster.

## Security

- All secrets (e.g., JWT keys, email credentials) are mounted as files inside the container at runtime.
- HashiCorp Vault is supported for securely managing and injecting secrets into containers.
- No sensitive information is hardcoded or written to disk.
- Stateless architecture allows for scaling, HA, and safe orchestration in production.

## Roadmap

- [ ] Admin Panel with Tenant Management  
- [ ] MFA (TOTP / Email-based)  
- [ ] OIDC Login (Google, GitHub, etc.)  
- [ ] Act as OIDC Provider for external services  
- [ ] Webhooks & Audit Logs
- [ ] Rate Limiting & Throttling  
- [ ] Email Provider Integrations (SMTP, Mailgun, etc.)  
- [ ] CLI Tool for Tenant & Config Management  
- [ ] Helm Chart for Kubernetes deployment  

## License

Licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0).  
You are free to use, modify, and distribute the software with proper attribution.

## Contributing

We're currently open to **feature suggestions**, feedback, and discussions via issues.

The platform is still under active development, so **pull requests are not being accepted at the moment**.

All contributions in the future will require commits to be signed off under the [Developer Certificate of Origin (DCO)](https://developercertificate.org/).

Stay tuned — contribution support (including PRs) will open up soon!
