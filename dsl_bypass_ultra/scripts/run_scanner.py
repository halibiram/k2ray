#!/usr/bin/env python
# ==============================================================================
# Helper Script: Run Network Scanner
# ==============================================================================
#
# Bu betik, NetworkScanner'ı çalıştırır ve bulunan ilk modem IP'sini
# standart çıktıya yazdırır. `install.sh` tarafından kullanılır.
#
# ==============================================================================

import sys
import os

# Proje kök dizinini Python yoluna ekle
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from utils.network_scanner import NetworkScanner

def main():
    """
    Ağ tarayıcısını çalıştırır ve sonucu yazdırır.
    """
    print("... (Python) Searching for Keenetic modem on the network ...")

    # Simülasyon ortamı için, ortak ağ geçitlerini dene
    # Gerçek bir ortamda bu aralık daha geniş olabilir.
    scanner = NetworkScanner(network_prefix="192.168.1.")

    # Keenetic modemler genellikle 80 (HTTP) veya 443 (HTTPS) portlarını kullanır
    modem_ip = scanner.find_modem(ports=[80, 443])

    if modem_ip:
        # Başarı durumunda, IP adresini standart çıktıya yazdır.
        # install.sh bu çıktıyı yakalayacak.
        print(modem_ip)

if __name__ == "__main__":
    main()