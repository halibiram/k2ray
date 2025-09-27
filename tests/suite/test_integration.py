import unittest
import os
import subprocess
import yaml
import sys
import tempfile
import shutil

# Add project root to the Python path to allow importing project modules
ROOT_DIR = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..'))
sys.path.append(ROOT_DIR)

class TestIntegration(unittest.TestCase):
    """
    Tests the integration between the deployment script (deploy.py)
    and the configuration it generates.
    """

    def setUp(self):
        """Set up the test environment by creating a temporary directory."""
        self.test_dir = tempfile.mkdtemp()
        self.scripts_dir = os.path.join(ROOT_DIR, "scripts")
        self.deploy_script_path = os.path.join(self.scripts_dir, "deploy.py")

    def tearDown(self):
        """Clean up the test environment by removing the temporary directory."""
        shutil.rmtree(self.test_dir)

    def test_keenetic_config_generation(self):
        """
        Verify that the deploy script generates the correct config for 'keenetic'.
        """
        lab_setup = "keenetic"
        target_speed = 150
        config_file_path = os.path.join(self.test_dir, "config.yaml")

        # Run the deploy script, directing its output to the temporary directory
        result = subprocess.run(
            [
                "python3",
                self.deploy_script_path,
                "--lab-setup", lab_setup,
                "--target-speed", str(target_speed),
                "--output-dir", self.test_dir
            ],
            capture_output=True, text=True
        )

        self.assertEqual(result.returncode, 0, f"Deploy script failed with error: {result.stderr}")
        self.assertTrue(os.path.exists(config_file_path), "Config file was not created.")

        # Load and validate the generated YAML
        with open(config_file_path, 'r') as f:
            config = yaml.safe_load(f)

        self.assertIsNotNone(config)
        self.assertEqual(config.get('lab_setup', {}).get('name'), lab_setup)
        self.assertEqual(config.get('lab_setup', {}).get('target_speed_mbps'), target_speed)
        self.assertEqual(config['v2ray']['outbounds'][0]['settings']['domainStrategy'], 'AsIs')

if __name__ == '__main__':
    unittest.main()