import json
import asyncio
import sys
import os
from flask import Flask, jsonify, render_template_string

# Hack to allow imports from parent directory
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from core.performance_monitor import PerformanceMonitor
from core.analytics_engine import AnalyticsEngine
from engines.keenetic_engine import KeeneticDSLEngine

app = Flask(__name__)

# Global instances for demonstration
keenetic_engine = KeeneticDSLEngine()
monitor = PerformanceMonitor(keenetic_engine)
analytics = AnalyticsEngine(monitor)


HTML_TEMPLATE = """
<!DOCTYPE html>
<html>
<head>
    <title>K2Ray Monitoring Dashboard</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        body { font-family: Arial, sans-serif; background-color: #f0f2f5; }
        h1, h2 { color: #333; }
        .container { display: flex; flex-wrap: wrap; }
        .card { background-color: white; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); margin: 10px; padding: 20px; }
        .chart-container { flex: 2; min-width: 400px; }
        .logs-container { flex: 1; min-width: 300px; }
        .analytics-container { flex: 1; min-width: 300px;}
        pre { background-color: #f4f4f4; border: 1px solid #ddd; padding: 10px; max-height: 200px; overflow-y: auto; white-space: pre-wrap; word-wrap: break-word; }
        .trend { font-weight: bold; }
        .trend.improving { color: green; }
        .trend.degrading { color: red; }
        .trend.stable { color: orange; }
    </style>
</head>
<body>
    <h1>K2Ray Real-Time Monitoring</h1>
    <div class="container">
        <div class="card chart-container">
            <h2>Downstream/Upstream Rate (kbps)</h2>
            <canvas id="speedChart"></canvas>
        </div>
        <div class="card analytics-container">
            <h2>Analytics</h2>
            <p>Downstream Trend: <span id="downstreamTrend" class="trend"></span></p>
            <p>Predicted Downstream in 1 min: <span id="prediction"></span> kbps</p>
            <hr/>
            <h2>Latest Data Point</h2>
            <pre id="dataLog">{}</pre>
        </div>
        <div class="card logs-container">
            <h2>Alerts</h2>
            <pre id="alertsLog">[]</pre>
        </div>
    </div>

    <script>
        const ctx = document.getElementById('speedChart').getContext('2d');
        const speedChart = new Chart(ctx, {
            type: 'line',
            data: { labels: [], datasets: [ { label: 'Downstream Rate', data: [], borderColor: 'rgb(75, 192, 192)', tension: 0.1, fill: false }, { label: 'Upstream Rate', data: [], borderColor: 'rgb(255, 99, 132)', tension: 0.1, fill: false } ] },
            options: { scales: { y: { beginAtZero: false } }, animation: { duration: 500 } }
        });

        async function updateDashboard() {
            try {
                const response = await fetch('/api/data');
                const data = await response.json();

                // Update Chart
                const label = new Date(data.timestamp * 1000).toLocaleTimeString();
                if (speedChart.data.labels.length > 30) {
                    speedChart.data.labels.shift();
                    speedChart.data.datasets.forEach((ds) => ds.data.shift());
                }
                speedChart.data.labels.push(label);
                speedChart.data.datasets[0].data.push(data.downstream_rate);
                speedChart.data.datasets[1].data.push(data.upstream_rate);
                speedChart.update();

                // Update Logs & Analytics
                document.getElementById('dataLog').textContent = JSON.stringify(data.latest_point, null, 2);
                document.getElementById('alertsLog').textContent = JSON.stringify(data.alerts, null, 2);

                const trendEl = document.getElementById('downstreamTrend');
                trendEl.textContent = data.analytics.downstream_trend;
                trendEl.className = 'trend ' + data.analytics.downstream_trend;

                document.getElementById('prediction').textContent = data.analytics.predicted_downstream_in_1m;


            } catch (error) {
                console.error('Failed to fetch data:', error);
            }
        }

        setInterval(updateDashboard, 2000);
        window.onload = updateDashboard;
    </script>
</body>
</html>
"""

@app.route('/')
def index():
    return render_template_string(HTML_TEMPLATE)

@app.route('/api/data')
def api_data():
    if not keenetic_engine.connected:
        try:
            loop = asyncio.get_event_loop()
        except RuntimeError:
            loop = asyncio.new_event_loop()
            asyncio.set_event_loop(loop)

        if not keenetic_engine.connected:
            loop.run_until_complete(keenetic_engine.connect("localhost", "admin", "pass"))
            loop.run_until_complete(keenetic_engine.get_dsl_status())

    data_point = monitor.collect_metrics()
    monitor.detect_connection_issues()

    # Run analytics
    speed_trend = analytics.analyze_speed_trend()
    prediction = analytics.predict_future_performance()

    response_data = {
        "latest_point": data_point,
        "alerts": monitor.alerts,
        "analytics": {**speed_trend, **prediction}
    }

    return jsonify(response_data)


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5001, debug=True)