import os
import yaml
from flask import Flask, render_template, request, redirect, url_for, flash
from configuration_manager import ConfigurationManager

app = Flask(__name__)
app.secret_key = 'supersecretkey'  # Üretimde bu değiştirilmelidir

# Ortamı belirle (örneğin, bir ortam değişkeninden)
APP_ENV = os.getenv('APP_ENV', 'development')

def get_config_manager():
    """Her istek için yeni bir ConfigurationManager örneği oluşturur."""
    return ConfigurationManager(env=APP_ENV)

@app.route('/')
def index():
    try:
        config_manager = get_config_manager()
        # Birleştirilmiş tam yapılandırmayı göster
        # Not: Kaydetme işlemi yalnızca ortama özgü dosyayı günceller.
        # Bu, varsayılanların üzerine yazmak için en temiz yoldur.
        content = yaml.dump(config_manager.config, default_flow_style=False, sort_keys=False)
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

if __name__ == '__main__':
    app.run(debug=True, port=5001)