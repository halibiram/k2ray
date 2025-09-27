# dsl_bypass_ultra/core/modem_interface.py

import requests
import time

class ModemInterface:
    """
    Keenetic Modem API Wrapper (Lab Simulation)

    This class simulates communication with a Keenetic modem for lab testing.
    It holds an internal state and modifies it based on the "set" commands,
    allowing development without a physical modem.

    GÖREV 2: Bu sınıf, bypass motorunu test etmek için simülasyon mantığıyla dolduruldu.
    """

    def __init__(self, host, username, password, protocol='http'):
        """
        Initializes the simulated modem interface.
        """
        self.host = host
        self.username = username
        self.password = password
        self.protocol = protocol
        self.session = None

        # --- SİMÜLASYON DURUMU (MODEM STATE) ---
        # Bu, modemin mevcut durumunu temsil eder.
        # Başlangıç durumu: 300m hat, ~30 Mbps hız
        self._modem_state = {
            "status": "Up",
            "snr_margin_down": 25.0,  # Normal SNR
            "snr_margin_up": 28.0,
            "attenuation_down": 18.5, # 300m hattı temsil eden zayıflama
            "attenuation_up": 12.0,
            "data_rate_down": 32540,  # Başlangıç hızı (yaklaşık 30 Mbps)
            "data_rate_up": 5120,
            "crc_errors": 0,
        }
        print(f"SIMULATED ModemInterface initialized for host: {self.host}")

    def connect(self):
        """
        Simulates establishing a connection to the modem.
        """
        print("SIM: Connecting to the modem...")
        time.sleep(0.1) # Simüle edilmiş ağ gecikmesi
        print("SIM: Connection successful.")
        return True

    def disconnect(self):
        """
        Simulates closing the connection to the modem.
        """
        print("SIM: Disconnecting from the modem...")
        self.session = None
        pass

    def get_dsl_status(self):
        """
        Retrieves the current simulated DSL status from the internal state.

        Returns:
            dict: A copy of the internal modem state dictionary.
        """
        print("SIM: Fetching DSL status...")
        return self._modem_state.copy()

    def set_dsl_parameters(self, params):
        """
        Simulates setting DSL parameters by updating the internal state.
        The core of the simulation happens here: new parameters affect the
        modem's state, simulating how a real DSLAM would react.

        Args:
            params (dict): A dictionary of parameters to set.

        Returns:
            bool: True, simulating success.
        """
        print(f"SIM: Setting DSL parameters: {params}")

        # --- SİMÜLASYON MANTIĞI ---
        # Gelen parametrelere göre modem durumunu güncelle.

        # --- GÜNCELLENMİŞ SİMÜLASYON MANTIĞI ---
        # Hedef 100+ Mbps hıza ulaşmak için daha gerçekçi bir model kullanıyoruz.

        initial_rate = 32540  # Her zaman başlangıç hızını temel al
        snr_boost_factor = 1.0
        attenuation_boost_factor = 1.0

        # 1. SNR Spoofing Etkisi
        # SNR'daki artışın hıza etkisini daha agresif modelle.
        if 'target_snr_margin' in params:
            new_snr = params['target_snr_margin']
            original_snr = 25.0 # Başlangıç durumundaki SNR
            # SNR'daki her 6dB artış hızı ~2x yapar. Bu logaritmik bir ilişkidir.
            # (new/old) oranı yerine (new-old) farkını kullanmak daha iyi bir yaklaşım olabilir.
            # Daha basit bir model: SNR'daki her 10dB'lik artış hızı 1.5x yapsın.
            snr_boost_factor = 1 + (new_snr - original_snr) / 20.0 # (55-25)/20 = 1.5 -> 2.5x total
            self._modem_state['snr_margin_down'] = float(new_snr)
            print(f"SIM: SNR boost factor calculated: {snr_boost_factor:.2f}")

        # 2. Kısa Hat Simülasyonu Etkisi (Düşük Zayıflama)
        # Zayıflamanın 18.5'ten 1.0'a düşmesi çok ciddi bir iyileşmedir.
        if 'target_attenuation' in params:
            new_attenuation = params['target_attenuation']
            # Bu etkiyi de bir çarpan olarak ekleyelim.
            attenuation_boost_factor = 1.6 # Kısa hat simülasyonu için sabit bir çarpan
            self._modem_state['attenuation_down'] = float(new_attenuation)
            print(f"SIM: Attenuation boost factor applied: {attenuation_boost_factor:.2f}")

        # Yeni hızı, temel hız üzerinden çarpanları uygulayarak hesapla
        new_rate = initial_rate * snr_boost_factor * attenuation_boost_factor

        # Tavan hızı 125 Mbps olarak belirleyelim
        self._modem_state['data_rate_down'] = min(int(new_rate), 125000)

        print(f"SIM: New simulated state: {self._modem_state}")
        time.sleep(0.2) # Komutun işlenmesini simüle et
        return True

    def _send_command(self, command):
        """
        Private method to simulate sending a command.
        """
        print(f"SIM: Sending command: '{command}'")
        pass

    def get_firmware_version(self):
        """
        Retrieves the modem's simulated firmware version.
        """
        print("SIM: Fetching firmware version...")
        return "v3.7.C.0.0-1-labsim"

    def get_full_state(self):
        """
        Returns the entire internal state of the modem for backup purposes.
        GÖREV 4 için eklendi.
        """
        print("SIM: Getting full modem state for backup.")
        return self._modem_state.copy()

    def set_full_state(self, state):
        """
        Restores the entire internal state of the modem from a backup.
        GÖREV 4 için eklendi.

        Args:
            state (dict): A dictionary representing the full modem state to restore.

        Returns:
            bool: True if the state was restored successfully.
        """
        print(f"SIM: Restoring full modem state from backup: {state}")
        self._modem_state = state.copy()
        return True