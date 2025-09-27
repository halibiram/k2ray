# dsl_bypass_ultra/core/performance_monitor.py

import time
import threading
from collections import deque

class PerformanceMonitor:
    """
    Real-time Performance Monitoring and Analytics System

    This class tracks key DSL metrics over time in a background thread,
    logs performance data, and provides analytics.

    GÖREV 3: Bu sınıf, arka planda çalışan izleme mantığı ile dolduruldu.
    """

    def __init__(self, modem_interface):
        """
        Initializes the performance monitor.

        Args:
            modem_interface (ModemInterface): The interface to get data from.
        """
        self.modem = modem_interface
        self.history = deque(maxlen=100)  # Son 100 kaydı sakla
        self._stop_event = threading.Event()
        self._thread = None
        print("PerformanceMonitor initialized.")

    def _monitor_loop(self, interval):
        """The actual monitoring logic that runs in the background."""
        print("Monitoring thread started.")
        while not self._stop_event.is_set():
            status = self.modem.get_dsl_status()
            if status:
                self.log_status(status)
            else:
                print("MONITOR: Failed to get status from modem.")
            time.sleep(interval)
        print("Monitoring thread stopped.")

    def start_monitoring(self, interval=5):
        """
        Starts monitoring the DSL connection in a separate daemon thread.
        """
        if self._thread is not None and self._thread.is_alive():
            print("MONITOR: Monitoring is already running.")
            return

        print(f"MONITOR: Starting real-time monitoring with {interval}s interval.")
        self._stop_event.clear()
        self._thread = threading.Thread(target=self._monitor_loop, args=(interval,), daemon=True)
        self._thread.start()

    def stop_monitoring(self):
        """
        Stops the monitoring loop gracefully.
        """
        if self._thread is None or not self._thread.is_alive():
            print("MONITOR: Monitoring is not running.")
            return

        print("MONITOR: Stopping real-time monitoring...")
        self._stop_event.set()
        self._thread.join(timeout=5) # Wait for the thread to finish
        print("MONITOR: Stopped.")

    def log_status(self, status):
        """
        Logs the current status and timestamp to the history deque.
        """
        log_entry = {
            "timestamp": time.time(),
            "status": status
        }
        self.history.append(log_entry)
        # print(f"MONITOR: Status logged. Rate: {status['data_rate_down']/1000:.2f} Mbps")

    def get_latest_status(self):
        """
        Returns the most recently logged status.

        Returns:
            dict: The latest status dictionary, or None if history is empty.
        """
        if not self.history:
            return None
        return self.history[-1]['status']

    def get_performance_analytics(self):
        """
        Generates a summary of performance analytics from the history.
        """
        print("Generating performance analytics...")
        if not self.history:
            return {}

        rates = [entry['status']['data_rate_down'] for entry in self.history]

        return {
            "average_down_rate_kbps": sum(rates) / len(rates),
            "max_down_rate_kbps": max(rates),
            "min_down_rate_kbps": min(rates),
            "total_disconnects": 0, # GÖREV 3'te geliştirilecek
            "log_count": len(self.history)
        }

    def generate_alerts(self):
        """
        Checks for conditions that should trigger an alert (e.g., instability).
        """
        print("Checking for alert conditions...")
        latest = self.get_latest_status()
        if latest and latest.get('crc_errors', 0) > 100:
            return [{"type": "high_crc", "value": latest['crc_errors']}]
        return []