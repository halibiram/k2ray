import time
from collections import deque

class PerformanceMonitor:
    def __init__(self, keenetic_engine):
        self.engine = keenetic_engine
        self.history = deque(maxlen=100)  # Store last 100 data points
        self.alerts = []

    def collect_metrics(self):
        """Collects current metrics from the Keenetic engine."""
        status = self.engine.dsl_status
        performance = self.engine.performance_metrics

        data_point = {
            "timestamp": time.time(),
            "downstream_rate": status.get("downstream_rate", 0),
            "upstream_rate": status.get("upstream_rate", 0),
            "snr_margin": status.get("snr_margin", 0),
            "attenuation": status.get("attenuation", 0),
            "success_count": performance.get("success", 0),
            "failure_count": performance.get("failure", 0),
            "line_status": status.get("line_status", "Down")
        }
        self.history.append(data_point)
        return data_point

    def track_speed_improvements(self):
        """Analyzes history to track speed improvements."""
        if len(self.history) < 2:
            return 0, 0

        initial_rate = self.history[0]["downstream_rate"]
        current_rate = self.history[-1]["downstream_rate"]

        improvement = current_rate - initial_rate
        return improvement, current_rate

    def detect_connection_issues(self):
        """Detects connection issues based on recent metrics."""
        if not self.history:
            return False

        last_point = self.history[-1]

        if last_point["line_status"] != "Up":
            self.generate_alert("Line is down!")
            return True

        if last_point["failure_count"] > self.history[0]["failure_count"]:
            self.generate_alert("Detected a new failure in operations.")
            return True

        return False

    def generate_alert(self, message):
        """Generates an alert for the system."""
        print(f"ALERT: {message}")
        self.alerts.append({"message": message, "timestamp": time.time()})
        # This is a hook for a more advanced alerting system (GÖREV 4)