# dsl_bypass_ultra/web/dashboard.py (GÖREV 5 Hali)

from flask import Flask, render_template, jsonify
import sys
import os

# Proje kökünü yola ekleyerek 'core' modülüne erişimi sağlıyoruz
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from core.performance_monitor import PerformanceMonitor

app = Flask(__name__)

# Bu global değişken, ana betik tarafından monitör nesnesiyle doldurulacak.
monitor_instance: PerformanceMonitor = None

def set_monitor(monitor: PerformanceMonitor):
    """
    Sets the global monitor instance for the web app.
    """
    global monitor_instance
    monitor_instance = monitor
    print("WEB: PerformanceMonitor instance has been injected into the web dashboard.")

@app.route('/')
def index():
    """
    Serves the main dashboard page.
    """
    return render_template('index.html', title='DSL Bypass Ultra')

@app.route('/api/status')
def api_status():
    """
    API endpoint to get the latest DSL status from the PerformanceMonitor.
    """
    if monitor_instance:
        status = monitor_instance.get_latest_status()
        if status:
            return jsonify(status)

    return jsonify({
        "status": "Initializing...",
        "snr_margin_down": 0,
        "data_rate_down": 0,
        "attenuation_down": 0
    })

def run_dashboard(host='0.0.0.0', port=8080, debug=False):
    """
    Runs the Flask web server.
    """
    if monitor_instance is None:
        print("WEB-ERROR: PerformanceMonitor has not been set! API will not work.")

    print(f"WEB: Starting web dashboard at http://{host}:{port}")
    app.run(host=host, port=port, debug=debug)

if __name__ == '__main__':
    print("This script is not meant to be run directly. Run main.py instead.")