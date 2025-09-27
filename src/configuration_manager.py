import os
import yaml
import json
import threading
from typing import Any, Dict
from jsonschema import validate, ValidationError

class ConfigurationManager:
    def __init__(self, config_path: str = 'config', env: str = 'development'):
        self.config_path = config_path
        self.env = env
        self._lock = threading.Lock()
        self.schema = self._load_schema()
        self.config = self._load_config()

    def reload(self):
        """Reloads the configuration from the files."""
        with self._lock:
            print("Configuration change detected. Reloading...")
            self.config = self._load_config()
            print("Configuration reloaded successfully.")

    def _load_schema(self) -> Dict[str, Any]:
        """Loads the configuration schema."""
        schema_path = os.path.join(self.config_path, 'schema.json')
        if not os.path.exists(schema_path):
            raise FileNotFoundError("Configuration schema 'schema.json' not found.")
        with open(schema_path, 'r') as f:
            return json.load(f)

    def _validate_config(self, config: Dict[str, Any]):
        """Validates the configuration against the schema."""
        try:
            validate(instance=config, schema=self.schema)
        except ValidationError as e:
            raise ValueError(f"Configuration validation failed: {e.message}")

    def _load_config(self) -> Dict[str, Any]:
        """Loads, merges, and validates configuration from YAML files."""
        base_config = self._read_yaml(os.path.join(self.config_path, 'default.yaml'))
        env_config = self._read_yaml(os.path.join(self.config_path, f'{self.env}.yaml'))

        merged_config = self._deep_merge(base_config, env_config)
        self._override_with_env_vars(merged_config)

        self._validate_config(merged_config)

        return merged_config

    def _read_yaml(self, file_path: str) -> Dict[str, Any]:
        """Reads a YAML file and returns its content."""
        if not os.path.exists(file_path):
            return {}
        with open(file_path, 'r') as f:
            return yaml.safe_load(f) or {}

    def _deep_merge(self, base: Dict, override: Dict) -> Dict:
        """Deeply merges two dictionaries."""
        merged = base.copy()
        for key, value in override.items():
            if isinstance(value, dict) and key in merged and isinstance(merged[key], dict):
                merged[key] = self._deep_merge(merged[key], value)
            else:
                merged[key] = value
        return merged

    def _override_with_env_vars(self, config: Dict, prefix: str = 'APP'):
        """Overrides configuration with environment variables."""
        for key, value in config.items():
            env_var_name = f"{prefix}_{key.upper()}"
            if isinstance(value, dict):
                self._override_with_env_vars(value, env_var_name)
            else:
                env_value = os.getenv(env_var_name)
                if env_value is not None:
                    config[key] = self._cast_env_value(env_value)

    def _cast_env_value(self, value: str) -> Any:
        """Casts environment variable string to appropriate type."""
        if value.lower() in ['true', 'false']:
            return value.lower() == 'true'
        if value.isdigit():
            return int(value)
        try:
            return float(value)
        except ValueError:
            return value

    def get(self, key: str, default: Any = None) -> Any:
        """Retrieves a configuration value by key."""
        keys = key.split('.')
        value = self.config
        for k in keys:
            if isinstance(value, dict) and k in value:
                value = value[k]
            else:
                return default
        return value

    def __getitem__(self, key: str) -> Any:
        return self.get(key)