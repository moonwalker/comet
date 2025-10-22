// @ts-check
/// <reference path="../index.d.ts" />

/**
 * Example stack showing how TypeScript definitions improve developer experience
 *
 * TypeScript definitions are auto-generated when you first run Comet.
 * Just add the reference comment above and enjoy IDE autocomplete!
 *
 * Benefits:
 * 1. Autocomplete for all Comet functions (stack, component, secrets, etc.)
 * 2. Inline documentation as you type
 * 3. Type checking to catch errors before running
 * 4. Better refactoring support in VS Code
 *
 * Note: This is still a .js file - no TypeScript compilation needed!
 */

// ============================================================================
// Configuration
// ============================================================================

const opts = {
  org: 'acme-corp',
  domain: 'acme.io',
  region: 'us-east-1',
  environment: 'production'
}

// TypeScript knows this returns a Stack object
const myStack = stack('production', { opts })

metadata({
  description: 'TypeScript support with autocomplete and type checking',
  tags: ['example', 'typescript', 'ide']
})

// ============================================================================
// Backend Configuration
// ============================================================================

// Autocomplete suggests: 'gcs', 's3', 'azurerm', 'local', etc.
backend('gcs', {
  bucket: 'acme-terraform-state',
  prefix: `${opts.org}/stacks/{{ .stack }}/{{ .component }}`
})

// ============================================================================
// Secrets Configuration
// ============================================================================

// Autocomplete shows available options: defaultProvider, defaultPath
secretsConfig({
  defaultProvider: 'sops',
  defaultPath: 'secrets.enc.yaml'
})

// ============================================================================
// Environment Variables
// ============================================================================

// Old way - still works
envs('AWS_REGION', opts.region)

// New bulk way - TypeScript validates the object structure
envs({
  AWS_REGION: opts.region,
  AWS_ACCESS_KEY_ID: secret('aws/access_key_id'),
  AWS_SECRET_ACCESS_KEY: secret('aws/secret_access_key'),
  ENVIRONMENT: opts.environment
})

// ============================================================================
// User-Defined Helper Functions (with JSDoc types)
// ============================================================================

/**
 * Generate a subdomain for the current stack
 * @param {string} name - Subdomain name
 * @returns {string} Full domain template string
 */
function subdomain(name) {
  return `${name}.{{ .stack }}.${opts.domain}`
}

/**
 * Create a standardized Kubernetes application component
 * @param {string} name - Application name
 * @param {Object} config - Additional configuration
 * @param {number} [config.replicas=2] - Number of replicas
 * @param {string} [config.image] - Container image
 * @param {string} [config.domain] - Custom domain
 * @returns {ComponentProxy} Component proxy for output references
 */
function k8sApp(name, config = {}) {
  return component(name, 'modules/k8s-app', {
    namespace: opts.org,
    domain: config.domain || subdomain(name),
    region: opts.region,
    replicas: config.replicas || 2,
    image: config.image || `${opts.org}/${name}:latest`,
    environment: opts.environment,
    // Spread any additional config
    ...config
  })
}

/**
 * Setup AWS credentials from secrets
 * @param {string} [role='default'] - AWS role name
 */
function setupAWS(role = 'default') {
  envs({
    AWS_ACCESS_KEY_ID: secret(`aws/${role}/access_key_id`),
    AWS_SECRET_ACCESS_KEY: secret(`aws/${role}/secret_access_key`),
    AWS_REGION: opts.region
  })
}

/**
 * Validate required environment variable exists
 * @param {string} name - Environment variable name
 * @returns {string} The environment variable value
 * @throws {Error} If environment variable is missing
 */
function requireEnv(name) {
  const value = envs(name)
  if (!value) {
    throw new Error(`Missing required environment variable: ${name}`)
  }
  return value
}

// ============================================================================
// Infrastructure Components
// ============================================================================

// TypeScript provides autocomplete for component() parameters
const vpc = component('vpc', 'modules/vpc', {
  cidr: '10.0.0.0/16',
  name: `${opts.org}-${myStack.name}-vpc`,
  region: opts.region
})

// TypeScript knows vpc is a ComponentProxy with dynamic properties
const database = component('database', 'modules/rds', {
  vpc_id: vpc.id,              // Autocomplete suggests vpc properties
  subnet_ids: vpc.subnet_ids,  // TypeScript knows these return template strings
  instance_class: 'db.t3.medium',
  allocated_storage: 100,
  engine: 'postgres',
  engine_version: '15.3',
  username: 'admin',
  password: secret('database/master_password')
})

// Append provider configuration
// TypeScript knows this takes (string, string[])
append('providers', [
  'data "aws_caller_identity" "current" {}'
])

// Using our helper functions
setupAWS('production')

// Deploy multiple apps using our factory
const api = k8sApp('api', {
  replicas: 5,
  database_url: database.endpoint,
  secret_key: secret('api/secret_key')
})

const worker = k8sApp('worker', {
  replicas: 3,
  database_url: database.endpoint,
  queue_url: secret('worker/queue_url')
})

const admin = k8sApp('admin', {
  replicas: 1,
  database_url: database.endpoint,
  domain: subdomain('admin')
})

// Deploy a batch of microservices
const services = ['auth', 'billing', 'notifications', 'analytics']

services.forEach(serviceName => {
  k8sApp(serviceName, {
    replicas: 2,
    database_url: database.endpoint
  })
})

// ============================================================================
// Kubernetes Configuration
// ============================================================================

// TypeScript validates the kubeconfig structure
kubeconfig({
  current: 0,
  clusters: [
    {
      context: `${opts.org}-${opts.environment}-gke`,
      host: 'https://kubernetes.example.com',
      cert: secret('kubernetes/ca_cert'),
      exec_command: 'gke-gcloud-auth-plugin',
      exec_args: [
        '--version=v1beta1',
        `--context=${opts.environment}`
      ]
    }
  ]
})

// ============================================================================
// What TypeScript Catches
// ============================================================================

// Uncomment these to see TypeScript errors in VS Code:

// @ts-expect-error - Wrong number of arguments
// stack()

// @ts-expect-error - backend expects 2 arguments
// backend('gcs')

// @ts-expect-error - secretsConfig expects an object
// secretsConfig('invalid')

// @ts-expect-error - component expects 3 arguments
// component('name', 'source')

// TypeScript won't catch everything (because it's still JavaScript)
// but it helps catch obvious mistakes before running comet!
