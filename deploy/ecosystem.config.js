// PM2 ecosystem config for the Sportz44 backend.
//
// The API server runs BOTH the REST/WebSocket API AND the live match listener
// (as an in-process goroutine), so only ONE process is managed here.
//
// Usage:
//   pm2 start deploy/ecosystem.config.js
//   pm2 save                 # persist the process list across reboots
//   pm2 startup              # enable PM2 to start on boot (run the printed cmd)
//
// Env vars are loaded from backend/.env at runtime by the app itself
// (godotenv), so no secrets are hardcoded here.

module.exports = {
  apps: [
    {
      name: 'sportz44-api',
      cwd: './backend',
      script: './bin/api',
      instances: 1,
      exec_mode: 'fork',
      autorestart: true,
      max_memory_restart: '300M',
      env: {
        NODE_ENV: 'production',
      },
      // Graceful shutdown: the API handles SIGTERM/SIGINT for a clean exit
      // (stops the match listener goroutine, then shuts down the HTTP server).
      kill_timeout: 10000,
      listen_timeout: 5000,
      out_file: '/var/log/sportz44/api.out.log',
      error_file: '/var/log/sportz44/api.err.log',
      merge_logs: true,
      time: true,
    },
  ],
};
