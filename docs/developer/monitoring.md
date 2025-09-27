# Monitoring and Alerting Guide

This document provides a comprehensive guide to setting up monitoring, logging, and alerting for the K2Ray application.

## 1. Structured Logging

The application uses `zerolog` for structured, high-performance logging. All logs are output in JSON format, making them easy to parse, search, and analyze with log management tools like Fluentd, Logstash, or Grafana Loki.

### Configuration

Logging behavior can be controlled via the following environment variables:

-   `LOG_LEVEL`: Sets the minimum log level to record. Can be `debug`, `info`, `warn`, or `error`. Defaults to `info`.
-   `LOG_PATH`: Specifies a file path for log output. If set, logs will be written to this file with automatic rotation. If not set, logs are written to `stderr`.
-   `REMOTE_LOG_URL`: A TCP address (e.g., `log-aggregator:5000`) to stream logs to a remote service.

### Example Log Entry

```json
{
  "level": "info",
  "service": "k2ray",
  "method": "POST",
  "path": "/api/v1/auth/login",
  "status_code": 200,
  "client_ip": "::1",
  "latency": 1.234567,
  "body_size": 123,
  "error": "",
  "message": "Request completed",
  "time": "2023-10-27T10:00:00Z"
}
```

## 2. Prometheus Metrics

The application exposes a wide range of metrics in a Prometheus-compatible format. These metrics provide insights into the application's health, performance, and behavior.

### Endpoint

The metrics are available at the following unauthenticated endpoint:

-   **URL**: `http://<k2ray-host>:8080/metrics`

### Key Metrics

-   `k2ray_app_info`: Static information about the application build.
-   `k2ray_http_requests_total`: A counter for all HTTP requests, labeled by method, path, and status code.
-   `k2ray_http_request_duration_seconds`: A histogram of request latencies.
-   `k2ray_user_logins_total`: A counter for user login attempts, labeled by status (`success` or `failure`).
-   `k2ray_system_cpu_usage_percent`: Current CPU utilization of the host machine.
-   `k2ray_system_memory_usage_bytes`: Memory usage, labeled by type (`total`, `used`, `free`).
-   `k2ray_system_disk_usage_bytes`: Disk usage for the root filesystem, labeled by type (`total`, `used`, `free`).

### Prometheus Configuration Example

Below is a sample `prometheus.yml` configuration to scrape metrics from the K2Ray application.

```yaml
# prometheus.yml

global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'k2ray'
    static_configs:
      - targets: ['localhost:8080'] # Replace with your k2ray instance address

rule_files:
  - 'alerting.rules.yml'
```

## 3. Alerting with Alertmanager

Prometheus can be configured to fire alerts based on predefined rules. These alerts are then sent to Alertmanager, which can route them to various notification channels like email, Slack, or PagerDuty.

### Example Alerting Rules

Create a file named `alerting.rules.yml` with the following content:

```yaml
# alerting.rules.yml

groups:
  - name: k2ray_alerts
    rules:
      - alert: K2RayInstanceDown
        expr: up{job="k2ray"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "K2Ray instance is down"
          description: "The K2Ray instance at {{ $labels.instance }} has been down for more than 1 minute."

      - alert: K2RayHighRequestLatency
        expr: histogram_quantile(0.95, sum(rate(k2ray_http_request_duration_seconds_bucket[5m])) by (le)) > 0.5
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High request latency detected"
          description: "The 95th percentile of request latency is above 0.5s."

      - alert: K2RayHighErrorRate
        expr: (sum(rate(k2ray_http_requests_total{status=~"5.."}[5m])) by (job) / sum(rate(k2ray_http_requests_total[5m])) by (job)) * 100 > 5
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High HTTP 5xx error rate"
          description: "More than 5% of requests are failing with a 5xx error."

      - alert: K2RayHighCPUUsage
        expr: k2ray_system_cpu_usage_percent > 80
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "High CPU usage detected"
          description: "CPU usage on the host running K2Ray is above 80%."
```

## 4. Grafana Dashboard

Grafana is an excellent tool for visualizing the metrics collected by Prometheus. You can import the JSON model below to create a pre-configured dashboard for K2Ray.

### Dashboard JSON Model

*(To use this, go to Grafana -> Dashboards -> Import and paste the JSON content.)*

```json
{
  "__inputs": [
    {
      "name": "DS_PROMETHEUS",
      "label": "Prometheus",
      "description": "",
      "type": "datasource",
      "pluginId": "prometheus",
      "pluginName": "Prometheus"
    }
  ],
  "__requires": [
    {
      "type": "grafana",
      "id": "grafana",
      "name": "Grafana",
      "version": "8.0.0"
    },
    {
      "type": "datasource",
      "id": "prometheus",
      "name": "Prometheus",
      "version": "1.0.0"
    }
  ],
  "annotations": {
    "list": [
      {
        "builtIn": 1,
        "datasource": "-- Grafana --",
        "enable": true,
        "hide": true,
        "iconColor": "rgba(0, 211, 255, 1)",
        "name": "Annotations & Alerts",
        "type": "dashboard"
      }
    ]
  },
  "editable": true,
  "gnetId": null,
  "graphTooltip": 0,
  "id": null,
  "links": [],
  "panels": [
    {
      "gridPos": { "h": 8, "w": 12, "x": 0, "y": 0 },
      "id": 2,
      "options": {
        "showThresholdLabels": false,
        "showThresholdMarkers": true,
        "stat": "last",
        "textMode": "auto",
        "valueAndName": "value_and_name"
      },
      "pluginVersion": "8.5.2",
      "targets": [
        {
          "datasource": "${DS_PROMETHEUS}",
          "expr": "sum(rate(k2ray_http_requests_total[5m]))",
          "legendFormat": "Requests per second",
          "refId": "A"
        }
      ],
      "title": "Total Requests (5m rate)",
      "type": "stat"
    },
    {
      "gridPos": { "h": 8, "w": 12, "x": 12, "y": 0 },
      "id": 4,
      "options": {
        "showThresholdLabels": false,
        "showThresholdMarkers": true,
        "stat": "last",
        "textMode": "auto",
        "valueAndName": "value_and_name"
      },
      "pluginVersion": "8.5.2",
      "targets": [
        {
          "datasource": "${DS_PROMETHEUS}",
          "expr": "histogram_quantile(0.95, sum(rate(k2ray_http_request_duration_seconds_bucket[5m])) by (le))",
          "legendFormat": "95th Percentile Latency",
          "refId": "A"
        }
      ],
      "title": "Request Latency (p95)",
      "type": "stat"
    },
    {
      "gridPos": { "h": 8, "w": 24, "x": 0, "y": 8 },
      "id": 6,
      "options": {
        "legend": {
          "displayMode": "list",
          "placement": "bottom"
        },
        "tooltip": {
          "mode": "multi"
        }
      },
      "targets": [
        {
          "datasource": "${DS_PROMETHEUS}",
          "expr": "sum(rate(k2ray_http_requests_total[1m])) by (status)",
          "legendFormat": "Status {{status}}",
          "refId": "A"
        }
      ],
      "title": "Requests by Status Code",
      "type": "timeseries"
    },
    {
      "gridPos": { "h": 8, "w": 12, "x": 0, "y": 16 },
      "id": 8,
      "options": {
        "legend": {
          "displayMode": "list",
          "placement": "bottom"
        },
        "tooltip": {
          "mode": "multi"
        }
      },
      "targets": [
        {
          "datasource": "${DS_PROMETHEUS}",
          "expr": "k2ray_system_cpu_usage_percent",
          "legendFormat": "CPU Usage",
          "refId": "A"
        }
      ],
      "title": "CPU Usage",
      "type": "timeseries"
    },
    {
      "gridPos": { "h": 8, "w": 12, "x": 12, "y": 16 },
      "id": 10,
      "options": {
        "legend": {
          "displayMode": "list",
          "placement": "bottom"
        },
        "tooltip": {
          "mode": "multi"
        }
      },
      "targets": [
        {
          "datasource": "${DS_PROMETHEUS}",
          "expr": "k2ray_system_memory_usage_bytes{type='used'} / k2ray_system_memory_usage_bytes{type='total'} * 100",
          "legendFormat": "Memory Usage (%)",
          "refId": "A"
        }
      ],
      "title": "Memory Usage",
      "type": "timeseries"
    }
  ],
  "schemaVersion": 35,
  "style": "dark",
  "tags": ["k2ray"],
  "templating": {
    "list": []
  },
  "time": {
    "from": "now-1h",
    "to": "now"
  },
  "timepicker": {},
  "timezone": "",
  "title": "K2Ray Application Dashboard",
  "uid": "k2ray-dashboard",
  "version": 1
}
```