import unittest
import os
import shutil
import yaml
from src.configuration_manager import ConfigurationManager

class TestConfigurationManager(unittest.TestCase):

    def setUp(self):
        """Her testten önce geçici bir yapılandırma dizini ve dosyaları oluşturur."""
        self.test_dir = 'test_config'
        os.makedirs(self.test_dir, exist_ok=True)

        self.default_config = {
            'app': {'name': 'TestApp'},
            'database': {'host': 'localhost', 'port': 5432, 'user': 'default', 'password': 'password'}
        }
        with open(os.path.join(self.test_dir, 'default.yaml'), 'w') as f:
            yaml.dump(self.default_config, f)

        self.dev_config = {
            'database': {'user': 'dev_user', 'password': 'dev_password'}
        }
        with open(os.path.join(self.test_dir, 'development.yaml'), 'w') as f:
            yaml.dump(self.dev_config, f)

        self.schema = {
            "type": "object",
            "properties": {
                "app": {"type": "object", "properties": {"name": {"type": "string"}}},
                "database": {
                    "type": "object",
                    "properties": {
                        "host": {"type": "string"},
                        "port": {"type": "integer"},
                        "user": {"type": "string"},
                        "password": {"type": "string"}
                    },
                    "required": ["host", "port", "user", "password"]
                }
            },
            "required": ["app", "database"]
        }
        with open(os.path.join(self.test_dir, 'schema.json'), 'w') as f:
            import json
            json.dump(self.schema, f)

    def tearDown(self):
        """Her testten sonra geçici dizini kaldırır."""
        shutil.rmtree(self.test_dir)
        if 'APP_DATABASE_PORT' in os.environ:
            del os.environ['APP_DATABASE_PORT']

    def test_load_and_merge_config(self):
        """Yapılandırmaların doğru bir şekilde yüklenip birleştirildiğini test eder."""
        cm = ConfigurationManager(config_path=self.test_dir, env='development')
        self.assertEqual(cm.get('app.name'), 'TestApp')
        self.assertEqual(cm.get('database.host'), 'localhost')
        self.assertEqual(cm.get('database.user'), 'dev_user') # dev, default'u geçersiz kılar

    def test_env_override(self):
        """Ortam değişkenlerinin yapılandırmayı geçersiz kıldığını test eder."""
        os.environ['APP_DATABASE_PORT'] = '5433'
        cm = ConfigurationManager(config_path=self.test_dir, env='development')
        self.assertEqual(cm.get('database.port'), 5433)

    def test_validation_success(self):
        """Geçerli bir yapılandırmanın doğrulamadan geçtiğini test eder."""
        try:
            ConfigurationManager(config_path=self.test_dir, env='development')
        except ValueError:
            self.fail("Geçerli yapılandırma ile doğrulama başarısız oldu.")

    def test_validation_failure(self):
        """Geçersiz bir yapılandırmanın doğrulama hatası verdiğini test eder."""
        # Şemayı ihlal eden bir yapılandırma oluşturun (gerekli bir alanı kaldırın)
        invalid_config = self.default_config.copy()
        del invalid_config['database']['password']

        with open(os.path.join(self.test_dir, 'default.yaml'), 'w') as f:
            yaml.dump(invalid_config, f)

        # Geliştirme yapılandırmasının bu alanı sağlamadığından emin olun
        with open(os.path.join(self.test_dir, 'development.yaml'), 'w') as f:
            yaml.dump({}, f)

        # Yapılandırma Yöneticisi'nin bir ValueError yükseltmesini bekleyin
        with self.assertRaises(ValueError) as context:
            ConfigurationManager(config_path=self.test_dir, env='development')

        self.assertIn("validation failed", str(context.exception))

    def test_get_nested_key(self):
        """İç içe anahtarların doğru bir şekilde alındığını test eder."""
        cm = ConfigurationManager(config_path=self.test_dir, env='development')
        self.assertEqual(cm.get('database.user'), 'dev_user')

    def test_get_default_value(self):
        """Mevcut olmayan bir anahtar için varsayılan değerin döndürüldüğünü test eder."""
        cm = ConfigurationManager(config_path=self.test_dir, env='development')
        self.assertEqual(cm.get('non.existent.key', 'default_val'), 'default_val')

if __name__ == '__main__':
    unittest.main()