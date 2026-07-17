<div align="center">
  <img src="frontend/public/favicon.svg" alt="Nucleus" width="64" />

# Nucleus - Vulnerability Scan Orchestrator

Nucleus is a single, self-contained Go binary that acts as an orchestration manager for scheduled [Nuclei](https://github.com/projectdiscovery/nuclei) vulnerability scans.

It serves an embedded Vue 3 SPA web dashboard, reads/writes scan states to a local SQLite database, and automatically dispatches beautifully formatted HTML email reports via SMTP when scans discover vulnerabilities.

</div>

## Prerequisites
- **Go 1.22+**
- **Node.js & npm** (for building the frontend)
- **Nuclei**: Ensure the `nuclei` CLI is installed and available in the system `$PATH`.

## Building the Application

Nucleus bundles the Vue 3 frontend directly into the Go binary using the `//go:embed` directive. 

### 1. Build the Frontend
```bash
cd frontend
npm install
npm run build
```
*(This places the static assets in the `dist/` directory).*

### 2. Compile the Go Binary
```bash
# From the root directory
go mod tidy
go build -o nucleus ./cmd/nucleus
```

## Download

The binary is availble for download from here:

https://apps.jdbnet.co.uk/nucleus

## Running as a Systemd Service

To keep Nucleus running continuously in the background on your VM, it is recommended to create a systemd service.

1. Create a service file:
```bash
sudo nano /etc/systemd/system/nucleus.service
```

2. Add the following configuration (adjust the `User`, `Group`, `WorkingDirectory`, and `ExecStart` paths to match your environment):

```ini
[Unit]
Description=Nucleus Vulnerability Scan Orchestrator
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/nucleus
ExecStart=/opt/nucleus/nucleus

# SMTP Configuration
Environment="SMTP_HOST=localhost"
Environment="SMTP_PORT=25"
Environment="SMTP_FROM=nucleus@example.com"
Environment="SMTP_TO=admin@example.com"

# Optional: Add authentication if your SMTP server requires it
# Environment="SMTP_USER=username"
# Environment="SMTP_PASS=password"

# Optional: Add basic authentication to protect the web dashboard
# Environment="WEB_USER=admin"
# Environment="WEB_PASS=supersecret"

Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

3. Enable and start the service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable nucleus
sudo systemctl start nucleus
sudo systemctl status nucleus
```

## Dashboard
By default, the web dashboard will be available on `http://<your-vm-ip>:8080`.
