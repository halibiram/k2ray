import numpy as np

class AnalyticsEngine:
    def __init__(self, monitor):
        self.monitor = monitor

    def analyze_speed_trend(self):
        """
        Analyzes the trend of downstream and upstream speeds.
        Returns a simple trend indicator ('stable', 'improving', 'degrading').
        """
        history = list(self.monitor.history)
        if len(history) < 10:
            return {"downstream_trend": "insufficient_data", "upstream_trend": "insufficient_data"}

        timestamps = np.array([p['timestamp'] for p in history])
        down_speeds = np.array([p['downstream_rate'] for p in history])

        try:
            # Using a simple linear regression slope to determine trend
            down_slope = np.polyfit(timestamps - timestamps[0], down_speeds, 1)[0]
        except (np.linalg.LinAlgError, TypeError):
            down_slope = 0


        def get_trend(slope):
            # This threshold is arbitrary and may need tuning.
            # It represents a change of 10 kbps per second.
            if abs(slope) < 10:
                return "stable"
            elif slope > 0:
                return "improving"
            else:
                return "degrading"

        return {
            "downstream_trend": get_trend(down_slope),
        }

    def predict_future_performance(self):
        """
        Placeholder for a more complex performance prediction model.
        Extrapolates the next data point based on the current trend.
        """
        history = list(self.monitor.history)
        if len(history) < 10:
            return {"predicted_downstream_in_1m": 0}

        timestamps = np.array([p['timestamp'] for p in history])
        down_speeds = np.array([p['downstream_rate'] for p in history])

        try:
            # Fit a line and predict the next point
            model = np.polyfit(timestamps - timestamps[0], down_speeds, 1)
            predict = np.poly1d(model)

            # Predict 60 seconds into the future
            future_point = (timestamps[-1] - timestamps[0]) + 60
            predicted_speed = predict(future_point)
        except (np.linalg.LinAlgError, TypeError):
            predicted_speed = down_speeds[-1]


        return {"predicted_downstream_in_1m": round(predicted_speed, 2)}