from abc import ABC, abstractmethod

class ModemInterface(ABC):
    """
    Abstract base class defining the interface for modem interactions.
    This serves as a stand-in for the deliverable from TASK 1.
    """

    @abstractmethod
    async def connect(self, host, username, password):
        """Establishes a connection to the modem."""
        pass

    @abstractmethod
    async def disconnect(self):
        """Closes the connection to the modem."""
        pass

    @abstractmethod
    async def get_dsl_status(self):
        """Retrieves the current DSL status and metrics."""
        pass

    @abstractmethod
    async def execute_raw_command(self, command):
        """Executes a raw command on the modem's shell or API."""
        pass

    @abstractmethod
    async def manipulate_line_length(self, target_length_m):
        """Simulates a shorter line length by adjusting DSL parameters."""
        pass

    @abstractmethod
    async def spoof_snr_values(self, target_snr_db):
        """Spoofs the Signal-to-Noise Ratio (SNR) values."""
        pass

    @abstractmethod
    async def adjust_attenuation(self, target_db):
        """Adjusts the line attenuation values reported by the modem."""
        pass

    @abstractmethod
    async def optimize_bit_loading(self):
        """Optimizes the bit-loading table for better performance."""
        pass