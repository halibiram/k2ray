# ==============================================================================
# DSL Bypass Ultra - Core Engine Test Script
# ==============================================================================
#
# Bu betik, GÖREV 2'de geliştirilen çekirdek bypass motorunun
# işlevselliğini test etmek için kullanılır.
#
# 1. Simüle edilmiş bir modem arayüzü başlatır.
# 2. Parametre manipülatörünü optimizasyon kurallarıyla yükler.
# 3. DSLAM Spoofer'ı kullanarak bypass işlemini tetikler.
# 4. İşlem öncesi ve sonrası sonuçları karşılaştırır.
#
# Çalıştırma: Projenin ana dizinindeyken `python run_bypass_test.py`
# ==============================================================================

# Gerekli modüllerin yolunu sisteme ekle
import sys
sys.path.append('./dsl_bypass_ultra')

from core.modem_interface import ModemInterface
from core.parameter_manipulator import ParameterManipulator
from core.dslam_spoofer import DslamSpoofer
from core.security_manager import SecurityManager

def main():
    """Ana test fonksiyonu (GÖREV 4 için güncellendi)"""
    print("=================================================")
    print("🚀 LAUNCHING DSL BYPASS ULTRA - SECURITY TEST 🚀")
    print("=================================================")

    # Adım 1: Çekirdek bileşenleri başlat
    print("\n[PHASE 1] Initializing core components...")
    modem = ModemInterface(host="192.168.1.1", username="admin", password="sim_password")
    manipulator = ParameterManipulator(config_path="dsl_bypass_ultra/config/optimization_rules.json")
    security = SecurityManager(modem_interface=modem)
    spoofer = DslamSpoofer(modem_interface=modem, param_manipulator=manipulator, security_manager=security)

    initial_state = modem.get_full_state()
    print("\n✅ Core components initialized successfully.")

    # Adım 2: Yedekleme yap
    print("\n[PHASE 2] Performing initial configuration backup...")
    backup_file = security.backup_configuration()
    if not backup_file:
        print("❌ TEST FAILED: Could not create initial backup.")
        return
    print(f"-> Backup created at: {backup_file}")

    # Adım 3: Bypass işlemini çalıştır
    print("\n[PHASE 3] Executing bypass with 'max_speed' profile...")
    bypass_status = spoofer.execute_bypass(profile='max_speed')
    if not bypass_status:
        print("❌ TEST FAILED: Bypass sequence did not complete.")
        return

    # Adım 4: Bypass sonucunu doğrula
    print("\n[PHASE 4] Verifying bypass results...")
    if bypass_status['data_rate_down'] > 100000:
        print("✅ SUCCESS: Data rate successfully boosted above 100 Mbps.")
    else:
        print("❌ TEST FAILED: Data rate did not reach target after bypass.")
        return

    # Adım 5: Yedekten geri yükle
    print("\n[PHASE 5] Restoring configuration from backup...")
    restore_success = security.restore_configuration(backup_file)
    if not restore_success:
        print("❌ TEST FAILED: Restore operation failed.")
        return
    print("-> Restore operation successful.")

    # Adım 6: Geri yükleme sonucunu doğrula
    print("\n[PHASE 6] Verifying restored state...")
    restored_state = modem.get_full_state()

    # Başlangıç ve geri yüklenmiş durumları karşılaştır
    # Hız ve SNR gibi anahtar değerlerin eşleşip eşleşmediğini kontrol et
    if (restored_state['data_rate_down'] == initial_state['data_rate_down'] and
        restored_state['snr_margin_down'] == initial_state['snr_margin_down']):
        print("✅ SUCCESS: Modem state successfully restored to initial values.")
    else:
        print("❌ TEST FAILED: Modem state does not match initial state after restore.")
        print(f"  Initial:  {initial_state}")
        print(f"  Restored: {restored_state}")
        return

    print("\n=================================================")
    print("✅✅✅ SECURITY FEATURES TEST PASSED ✅✅✅")
    print("=================================================")

if __name__ == "__main__":
    main()