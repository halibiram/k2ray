import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import ConfigEditor from './ConfigEditor.vue'
import { v4 as uuidv4 } from 'uuid'

// A helper to create a complete default structure for merging
const createDefaultConfigData = (protocol) => {
    const wsSettings = { path: '/', headers: {} };
    const grpcSettings = { serviceName: '' };
    const transport = { net: 'tcp', tls: 'none', wsSettings, grpcSettings };

    switch (protocol) {
        case 'vmess':
            return { v: '2', add: '', port: 443, id: uuidv4(), aid: 0, type: 'none', host: '', path: '', ...transport };
        case 'vless':
            return { id: uuidv4(), add: '', port: 443, encryption: 'none', flow: '', ...transport };
        case 'trojan':
            return { server: '', server_port: 443, password: '', sni: '', ...transport };
        case 'shadowsocks':
            return { server: '', server_port: 8388, password: '', method: 'aes-256-gcm' };
        default:
            return {};
    }
};

describe('ConfigEditor.vue', () => {
  let wrapper;

  beforeEach(() => {
    wrapper = mount(ConfigEditor, {
      props: {
        isOpen: true,
        config: null, // Start in "create" mode
      },
      global: {
        stubs: {
          TransitionRoot: { template: '<div><slot /></div>' },
          TransitionChild: { template: '<div><slot /></div>' },
          Dialog: { template: '<div><slot /></div>' },
          DialogPanel: { template: '<div><slot /></div>' },
          DialogTitle: { template: '<h3><slot /></h3>' },
        }
      }
    });
  });

  it('renders the form for creating a new config', () => {
    expect(wrapper.find('h3').text()).toContain('Create Configuration');
    expect(wrapper.find('#name').exists()).toBe(true);
    expect(wrapper.find('#protocol').exists()).toBe(true);
  });

  it('switches protocol fields correctly', async () => {
    // Default is VMess
    expect(wrapper.find('#vmess-aid').exists()).toBe(true);

    // Switch to VLESS
    await wrapper.find('#protocol').setValue('vless');
    expect(wrapper.find('#vless-flow').exists()).toBe(true);
    expect(wrapper.find('#vmess-aid').exists()).toBe(false);

    // Switch to Trojan
    await wrapper.find('#protocol').setValue('trojan');
    expect(wrapper.find('#trojan-sni').exists()).toBe(true);
    expect(wrapper.find('#vless-flow').exists()).toBe(false);

    // Switch to Shadowsocks
    await wrapper.find('#protocol').setValue('shadowsocks');
    expect(wrapper.find('#ss-method').exists()).toBe(true);
    expect(wrapper.find('#trojan-sni').exists()).toBe(false);
  });

  it('shows transport settings for applicable protocols', async () => {
    // VMess should have transport settings
    await wrapper.find('#protocol').setValue('vmess');
    expect(wrapper.find('#transport-net').exists()).toBe(true);

    // VLESS should have transport settings
    await wrapper.find('#protocol').setValue('vless');
    expect(wrapper.find('#transport-net').exists()).toBe(true);

    // Trojan should have transport settings
    await wrapper.find('#protocol').setValue('trojan');
    expect(wrapper.find('#transport-net').exists()).toBe(true);

    // Shadowsocks should NOT have transport settings
    await wrapper.find('#protocol').setValue('shadowsocks');
    expect(wrapper.find('#transport-net').exists()).toBe(false);
  });

  it('shows conditional transport fields correctly', async () => {
    await wrapper.find('#protocol').setValue('vmess');

    // WebSocket path should not be visible by default
    expect(wrapper.find('#ws-path').exists()).toBe(false);

    // Select WebSocket network
    await wrapper.find('#transport-net').setValue('ws');
    expect(wrapper.find('#ws-path').exists()).toBe(true);

    // gRPC service name should not be visible
    expect(wrapper.find('#grpc-service-name').exists()).toBe(false);

    // Select gRPC network
    await wrapper.find('#transport-net').setValue('grpc');
    expect(wrapper.find('#grpc-service-name').exists()).toBe(true);
    expect(wrapper.find('#ws-path').exists()).toBe(false);
  });

  it('loads an existing config for editing and merges with defaults', async () => {
    const existingConfig = {
      id: 1,
      name: 'My Old WS Server',
      protocol: 'vmess',
      // This config is "old" and doesn't have grpcSettings
      config_data: JSON.stringify({
        add: 'test.com',
        port: 1234,
        id: 'existing-uuid',
        net: 'ws',
        wsSettings: { path: '/my-path' }
      }),
    };

    await wrapper.setProps({ config: existingConfig });

    expect(wrapper.find('h3').text()).toContain('Edit Configuration');
    expect(wrapper.find('#name').element.value).toBe('My Old WS Server');

    // Check that basic fields are populated
    const vm = wrapper.vm as any;
    expect(vm.form.config_data.add).toBe('test.com');
    expect(vm.form.config_data.port).toBe(1234);

    // Check that transport settings are populated
    expect(vm.form.config_data.net).toBe('ws');
    expect(wrapper.find('#ws-path').element.value).toBe('/my-path');

    // Check that a field missing from the old config (grpcSettings) exists in the model
    // This confirms the deep merge worked
    expect(vm.form.config_data.grpcSettings).toBeDefined();
    expect(vm.form.config_data.grpcSettings.serviceName).toBe('');
  });

  it('emits a save event with the correct form data', async () => {
    await wrapper.find('#name').setValue('New Test Server');
    await wrapper.find('#protocol').setValue('trojan');
    await wrapper.find('#trojan-server').setValue('my-trojan.com');
    await wrapper.find('#trojan-port').setValue(8443);
    await wrapper.find('#trojan-password').setValue('secret');
    await wrapper.find('#transport-net').setValue('grpc');
    await wrapper.find('#grpc-service-name').setValue('my-grpc-service');

    await wrapper.find('form').trigger('submit');

    const emitted = wrapper.emitted('save');
    expect(emitted).toHaveLength(1);
    const payload = emitted[0][0];

    expect(payload.name).toBe('New Test Server');
    expect(payload.protocol).toBe('trojan');
    expect(payload.config_data.server).toBe('my-trojan.com');
    expect(payload.config_data.net).toBe('grpc');
    expect(payload.config_data.grpcSettings.serviceName).toBe('my-grpc-service');
  });
});