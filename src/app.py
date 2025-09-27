import os
import yaml
import datetime
import shutil
import json
import threading
import time
from flask import Flask, render_template, request, redirect, url_for, flash
from configuration_manager import ConfigurationManager
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler

app = Flask(__name__)
app.secret_key = 'supersecretkey'  # Üretimde bu değiştirilmelidir

# Ortamı belirle (örneğin, bir ortam değişkeninden)
APP_ENV = os.getenv('APP_ENV', 'development')

# Global ConfigurationManager örneği
config_manager = ConfigurationManager(env=APP_ENV)

def get_config_manager():
    """Global ConfigurationManager örneğini döndürür."""
    return config_manager

class ConfigChangeHandler(FileSystemEventHandler):
    def __init__(self, manager):
        self.manager = manager

    def on_modified(self, event):
        print(f"Event detected: {event.event_type} on {event.src_path}, is_directory={event.is_directory}")
        # Dosya adı ve uzantısını kontrol et
        if not event.is_directory and event.src_path.endswith(('.yaml', '.json')):
            # Yedekleme dosyalarından gelen olayları yoksay
            if ".bak" in event.src_path:
                print(f"Ignoring backup file modification: {event.src_path}")
                return

            print(f"Yapılandırma dosyası değiştirildi: {os.path.basename(event.src_path)}. Yeniden yükleniyor...")
            self.manager.reload()


def start_watcher():
    event_handler = ConfigChangeHandler(config_manager)
    observer = Observer()
    observer.schedule(event_handler, config_manager.config_path, recursive=True)
    observer.start()
    print(f"'{config_manager.config_path}' dizini izleniyor...")
    try:
        while True:
            time.sleep(1)
    except:  # Daemon thread will be terminated abruptly, so catch everything
        observer.stop()
    observer.join()

@app.route('/')
def index():
    try:
        cm = get_config_manager()
        # Birleştirilmiş tam yapılandırmayı göster
        content = yaml.dump(cm.config, default_flow_style=False, sort_keys=False)
        return render_template('index.html', config_content=content, env=APP_ENV)
    except Exception as e:
        flash(f"Yapılandırma yüklenirken hata oluştu: {e}", "error")
        return render_template('index.html', config_content=f"# Hata: {e}", env=APP_ENV)

@app.route('/save', methods=['POST'])
def save():
    try:
        # Kaydetmeden önce, yeni yapılandırmanın geçerli olduğundan emin olmalıyız.
        new_content_str = request.form['config']
        new_config_data = yaml.safe_load(new_content_str)

        if not isinstance(new_config_data, dict):
            raise ValueError("Sağlanan yapılandırma geçerli bir YAML nesnesi değil.")

        # Geçici bir ConfigurationManager ile doğrula
        # Bu, `default.yaml`'ı yükler ve yeni verilerle birleştirir, sonra doğrular.
        temp_manager = ConfigurationManager(env=APP_ENV)

        # Yeni verileri varsayılanın üzerine yazarak tam bir birleştirilmiş yapılandırma oluşturun.
        # Bu, kullanıcının tam yapılandırmayı düzenlemesine olanak tanır.
        # Ancak, yalnızca `default.yaml`'dan farklı olan değişiklikleri kaydetmek daha iyidir.
        # Şimdilik, basitlik adına, kullanıcı tarafından sağlanan her şeyi ortama özgü dosyaya kaydediyoruz.

        # Tam birleştirilmiş yapılandırmayı doğrula
        # Önce temel yapılandırmayı al
        base_config = temp_manager._read_yaml(os.path.join(temp_manager.config_path, 'default.yaml'))
        # Düzenlenen verileri temel yapılandırmayla birleştir
        merged_for_validation = temp_manager._deep_merge(base_config, new_config_data)
        # Birleştirilmiş yapılandırmayı doğrula
        temp_manager._validate_config(merged_for_validation)

        # Doğrulama başarılı olursa, yalnızca ortama özgü yapılandırma dosyasını kaydet.
        # Kullanıcının gönderdiği YAML'ı, default'tan farklı olanlarla sınırlamak en iyisi olacaktır,
        # ancak şimdilik tüm düzenlenmiş içeriği kaydediyoruz.
        env_config_path = os.path.join(temp_manager.config_path, f'{temp_manager.env}.yaml')

        # Yedekleme oluştur
        backup_path = None
        if os.path.exists(env_config_path):
            timestamp = datetime.datetime.now().strftime("%Y-%m-%dT%H-%M-%S")
            backup_path = f"{env_config_path}.{timestamp}.bak"
            shutil.copy2(env_config_path, backup_path)
            flash(f"'{os.path.basename(env_config_path)}' yedeklendi -> '{os.path.basename(backup_path)}'", "info")

        # Değişiklik geçmişini kaydet
        if backup_path:
            history_path = os.path.join(temp_manager.config_path, 'history.json')
            history = []
            if os.path.exists(history_path):
                with open(history_path, 'r') as f:
                    history = json.load(f)

            history.insert(0, {
                "timestamp": timestamp,
                "env": temp_manager.env,
                "backup_file": os.path.basename(backup_path)
            })

            with open(history_path, 'w') as f:
                json.dump(history, f, indent=4)

        # Kullanıcının gönderdiği YAML'ı, default'tan farklı olanlarla sınırlamak için bir fark oluşturun
        # Bu, `development.yaml` dosyasını temiz tutar.
        default_data = temp_manager._read_yaml(os.path.join(temp_manager.config_path, 'default.yaml'))

        def get_diff(d1, d2):
            diff = {}
            for k, v2 in d2.items():
                if k not in d1:
                    diff[k] = v2
                else:
                    v1 = d1[k]
                    if isinstance(v2, dict) and isinstance(v1, dict):
                        sub_diff = get_diff(v1, v2)
                        if sub_diff:
                            diff[k] = sub_diff
                    elif v1 != v2:
                        diff[k] = v2
            return diff

        diff_to_save = get_diff(default_data, new_config_data)

        with open(env_config_path, 'w') as f:
            yaml.dump(diff_to_save, f, default_flow_style=False, sort_keys=False)

        flash(f"'{temp_manager.env}.yaml' başarıyla kaydedildi ve doğrulandı!", "success")
    except Exception as e:
        flash(f"Yapılandırma kaydedilirken hata oluştu: {e}", "error")

    return redirect(url_for('index'))

@app.route('/history')
def history():
    config_manager = get_config_manager()
    history_path = os.path.join(config_manager.config_path, 'history.json')
    history = []
    if os.path.exists(history_path):
        with open(history_path, 'r') as f:
            history = json.load(f)
    return render_template('history.html', history=history, env=APP_ENV)

@app.route('/rollback/<timestamp>')
def rollback(timestamp):
    try:
        config_manager = get_config_manager()
        history_path = os.path.join(config_manager.config_path, 'history.json')
        history = []
        if os.path.exists(history_path):
            with open(history_path, 'r') as f:
                history = json.load(f)

        entry_to_restore = next((item for item in history if item["timestamp"] == timestamp), None)

        if not entry_to_restore:
            flash("Geri alınacak kayıt bulunamadı.", "error")
            return redirect(url_for('history'))

        env = entry_to_restore['env']
        backup_filename = entry_to_restore['backup_file']

        backup_path = os.path.join(config_manager.config_path, backup_filename)
        env_config_path = os.path.join(config_manager.config_path, f'{env}.yaml')

        if not os.path.exists(backup_path):
            flash(f"Yedekleme dosyası '{backup_filename}' bulunamadı.", "error")
            return redirect(url_for('history'))

        shutil.copy2(backup_path, env_config_path)
        flash(f"'{env}.yaml' başarıyla '{backup_filename}' içeriğine geri döndürüldü.", "success")

    except Exception as e:
        flash(f"Geri alma sırasında hata oluştu: {e}", "error")

    return redirect(url_for('history'))


if __name__ == '__main__':
    # Arka plan iş parçacığında dosya izleyiciyi başlat
    watcher_thread = threading.Thread(target=start_watcher, daemon=True)
    watcher_thread.start()

    # Flask uygulamasını ana iş parçacığında çalıştır
    # Not: Flask'in kendi yeniden yükleyicisini (use_reloader=False) devre dışı bırakıyoruz
    # çünkü bu, watchdog gözlemcimizle çakışabilir.
    app.run(debug=True, port=5001, use_reloader=False)