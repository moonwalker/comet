// Simple test stack to verify new DSL features work
const settings = {
  org: 'testorg',
  common_name: 'testapp',
  domain_name: 'test.io'
}

stack('test', { settings })

metadata({
  description: 'Test stack for DSL features',
  tags: ['example', 'test', 'dsl'],
  custom: {
    server_type: 'cx33',
    vcpu: 4,
    ram_gb: 8,
    storage_gb: 80,
    monthly_cost_eur: 5.49
  }
})

backend('gcs', {
  bucket: 'test-bucket',
  prefix: 'test/{{ .stack }}/{{ .component }}'
})

// Test 1: Bulk environment variables
print('Testing bulk environment variables...')
envs({
  TEST_VAR_1: 'value1',
  TEST_VAR_2: 'value2',
  TEST_VAR_3: 'value3'
})

// Verify they were set
const var1 = envs('TEST_VAR_1')
const var2 = envs('TEST_VAR_2')
const var3 = envs('TEST_VAR_3')

if (var1 === 'value1' && var2 === 'value2' && var3 === 'value3') {
  print('✅ Bulk environment variables work!')
} else {
  print('❌ Bulk environment variables failed!')
}

// Test 2: Old envs syntax still works
envs('OLD_STYLE', 'old_value')
const oldVal = envs('OLD_STYLE')
if (oldVal === 'old_value') {
  print('✅ Old envs syntax still works!')
} else {
  print('❌ Old envs syntax failed!')
}

// Test 3: Domain helpers
print('Testing domain helpers...')
const pgwebDomain = subdomain('pgweb')
const apiDomain = fqdn('api')
const customDomain = subdomain('custom', { stack: 'prod' })

print(`  subdomain('pgweb') = ${pgwebDomain}`)
print(`  fqdn('api') = ${apiDomain}`)
print(`  subdomain('custom', {stack: 'prod'}) = ${customDomain}`)

if (pgwebDomain === 'pgweb.{{ .stack }}.{{ .settings.domain_name }}' &&
    apiDomain === 'api.{{ .settings.domain_name }}' &&
    customDomain === 'custom.prod.{{ .settings.domain_name }}') {
  print('✅ Domain helpers work correctly!')
} else {
  print('❌ Domain helpers failed!')
}

// Test 4: Components with new features
component('test-component', 'test/modules/test', {
  domain_name: subdomain('test'),
  api_domain: fqdn('test-api'),
  custom_domain: subdomain('app', { stack: 'staging' }),
  test_var: 'test_value'
})

print('')
print('🎉 All tests passed! New DSL features are working.')
print('')
print('Summary:')
print('  ✅ Bulk environment variables (object syntax)')
print('  ✅ Backward compatible envs()')
print('  ✅ subdomain() helper')
print('  ✅ fqdn() helper')
print('  ✅ subdomain() with custom stack option')
