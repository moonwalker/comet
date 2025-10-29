/**
 * Comet TypeScript Definitions
 *
 * These type definitions provide IDE autocomplete and type checking
 * for Comet stack files while still allowing pure JavaScript.
 *
 * Usage in your .stack.js files:
 *
 *   /// <reference types="./types" />
 *
 * Or if you use JSDoc:
 *
 *   // @ts-check
 *   /// <reference types="./types" />
 */

// ============================================================================
// Core Types
// ============================================================================

/**
 * Stack configuration options
 */
export interface StackOptions {
  /** Custom options/settings passed to stack (accessible in templates as {{ .settings }}) */
  [key: string]: any;
}

/**
 * Stack object returned by stack() function
 */
export interface Stack {
  /** Stack name */
  name: string;
  /** Stack options */
  options: StackOptions;
}

/**
 * Backend configuration for Terraform/OpenTofu state
 */
export interface BackendConfig {
  /** Backend-specific configuration */
  [key: string]: any;
}

/**
 * Component configuration
 */
export interface ComponentConfig {
  /** Input variables for the component */
  [key: string]: any;

  /** Provider configuration (optional) */
  providers?: {
    [providerName: string]: ProviderConfig;
  };

  /** Explicit inputs (optional, alternative to root-level config) */
  inputs?: {
    [key: string]: any;
  };
}

/**
 * Provider configuration
 */
export interface ProviderConfig {
  /** Provider alias (optional) */
  alias?: string;
  /** Provider-specific configuration */
  [key: string]: any;
}

/**
 * Component proxy object with dynamic property access
 * Properties return template strings referencing component outputs
 */
export interface ComponentProxy {
  /** Component name */
  name: string;
  /** Access any output from the component as a template reference */
  [output: string]: any;
}

/**
 * Secrets configuration
 */
export interface SecretsConfig {
  /** Default secrets provider (e.g., 'sops', 'op') */
  defaultProvider?: 'sops' | 'op' | string;
  /** Default secrets file path */
  defaultPath?: string;
}

/**
 * Kubeconfig cluster configuration
 */
export interface KubeconfigCluster {
  /** Kubernetes context name */
  context: string;
  /** Cluster API server host */
  host: string;
  /** Base64-encoded cluster CA certificate */
  cert: string;
  /** Static bearer token for authentication (mutually exclusive with exec_command) */
  token?: string;
  /** Exec plugin command (mutually exclusive with token) */
  exec_command?: string;
  /** Exec plugin arguments */
  exec_args?: string[];
}

/**
 * Kubeconfig configuration
 */
export interface Kubeconfig {
  /** Current cluster index */
  current: number;
  /** List of cluster configurations */
  clusters: KubeconfigCluster[];
}

/**
 * Environment variables object
 */
export interface EnvVars {
  [key: string]: string;
}

// ============================================================================
// Global Functions
// ============================================================================

/**
 * Print output to console
 * @param args - Values to print
 */
export function print(...args: any[]): void;

/**
 * Access environment variables as a proxy object
 * @example
 * const token = env.GITHUB_TOKEN
 * const home = env.HOME
 */
export const env: {
  [key: string]: string | undefined;
};

/**
 * Get or set environment variables
 *
 * @example
 * // Get a single variable
 * const value = envs('MY_VAR')
 *
 * // Set a single variable
 * envs('MY_VAR', 'my-value')
 *
 * // Set multiple variables (bulk mode)
 * envs({
 *   AWS_ACCESS_KEY_ID: 'xxx',
 *   AWS_SECRET_ACCESS_KEY: 'yyy',
 *   AWS_REGION: 'us-east-1'
 * })
 */
export function envs(key: string): string | undefined;
export function envs(key: string, value: string): string;
export function envs(vars: EnvVars): void;

/**
 * Get a secret value using the full provider URI
 *
 * @param ref - Full secret reference (e.g., 'sops://secrets.enc.yaml#/path/to/secret')
 * @returns The decrypted secret value
 *
 * @example
 * const apiKey = secrets('sops://secrets.enc.yaml#/datadog/api_key')
 * const password = secrets('op://vault/item/field')
 */
export function secrets(ref: string): any;

/**
 * Configure default secrets provider and path
 *
 * @param config - Secrets configuration object
 *
 * @example
 * secretsConfig({
 *   defaultProvider: 'sops',
 *   defaultPath: 'secrets.enc.yaml'
 * })
 */
export function secretsConfig(config: SecretsConfig): void;

/**
 * Get a secret value using shorthand notation
 * Uses the configured default provider and path
 *
 * @param path - Secret path (e.g., 'datadog/api_key' or 'datadog.api_key')
 * @returns The decrypted secret value
 *
 * @example
 * // Configure defaults first
 * secretsConfig({ defaultProvider: 'sops', defaultPath: 'secrets.enc.yaml' })
 *
 * // Then use shorthand
 * const apiKey = secret('datadog/api_key')
 * const appKey = secret('datadog.app_key')  // dot notation also works
 */
export function secret(path: string): any;

/**
 * Define a stack with name and options
 *
 * @param name - Stack name (usually environment like 'dev', 'staging', 'production')
 * @param options - Stack configuration options
 * @returns Stack object
 *
 * @example
 * const opts = {
 *   org: 'mycompany',
 *   region: 'us-east-1'
 * }
 *
 * const myStack = stack('production', { opts })
 */
export function stack(name: string, options?: StackOptions): Stack;

/**
 * Configure the Terraform/OpenTofu backend
 *
 * @param type - Backend type ('gcs', 's3', 'azurerm', 'local', etc.)
 * @param config - Backend-specific configuration
 *
 * @example
 * backend('gcs', {
 *   bucket: 'my-terraform-state',
 *   prefix: 'stacks/{{ .stack }}/{{ .component }}'
 * })
 *
 * backend('s3', {
 *   bucket: 'my-terraform-state',
 *   key: 'stacks/{{ .stack }}/{{ .component }}/terraform.tfstate',
 *   region: 'us-east-1'
 * })
 */
export function backend(type: string, config: BackendConfig): void;

/**
 * Define an infrastructure component
 *
 * @param name - Component name (unique within stack)
 * @param source - Path to Terraform module (relative or absolute)
 * @param config - Component configuration (inputs and providers)
 * @returns Component proxy object for referencing outputs
 *
 * @example
 * const vpc = component('vpc', 'modules/vpc', {
 *   cidr: '10.0.0.0/16',
 *   name: 'my-vpc'
 * })
 *
 * // Reference vpc outputs in other components
 * const subnet = component('subnet', 'modules/subnet', {
 *   vpc_id: vpc.id,  // Creates template reference
 *   cidr: '10.0.1.0/24'
 * })
 */
export function component(
  name: string,
  source: string,
  config: ComponentConfig
): ComponentProxy;

/**
 * Append raw Terraform code to generated files
 *
 * @param type - File type to append to (e.g., 'providers', 'variables', 'outputs')
 * @param lines - Array of Terraform code lines to append
 *
 * @example
 * append('providers', [
 *   'data "google_client_config" "default" {}'
 * ])
 *
 * append('outputs', [
 *   'output "cluster_endpoint" {',
 *   '  value = google_container_cluster.primary.endpoint',
 *   '}'
 * ])
 */
export function append(type: string, lines: string[]): void;

/**
 * Configure Kubernetes access for the stack
 *
 * @param config - Kubeconfig configuration
 *
 * @example
 * kubeconfig({
 *   current: 0,
 *   clusters: [
 *     {
 *       context: 'gke-production',
 *       host: 'https://1.2.3.4',
 *       cert: 'LS0tLS1CRU...',
 *       exec_command: 'gke-gcloud-auth-plugin',
 *       exec_args: ['--version=v1beta1']
 *     }
 *   ]
 * })
 */
export function kubeconfig(config: Kubeconfig): void;

// ============================================================================
// Template Functions (available in Go templates)
// ============================================================================

/**
 * Template context available in template strings
 *
 * @example
 * component('app', 'modules/app', {
 *   name: '{{ .stack }}-app',           // Current stack name
 *   region: '{{ .settings.region }}',   // Stack settings
 *   component: '{{ .component }}'       // Current component name
 * })
 */
export interface TemplateContext {
  /** Current stack name */
  stack: string;
  /** Stack settings/options */
  settings: any;
  /** Current component name */
  component: string;
  /** Stack options (alias for settings) */
  opts: any;
}

// ============================================================================
// User-Defined Helpers (Examples)
// ============================================================================

/**
 * Example: Domain helper function
 * Users can create their own helper functions like this
 *
 * @example
 * function subdomain(name: string): string {
 *   return `${name}.{{ .stack }}.{{ .settings.domain }}`
 * }
 *
 * component('api', 'modules/app', {
 *   domain: subdomain('api')  // api.production.example.com
 * })
 */

/**
 * Example: Component factory function
 * Users can create factories for common component patterns
 *
 * @example
 * function k8sApp(name: string, config: Partial<ComponentConfig> = {}): ComponentProxy {
 *   return component(name, 'modules/k8s-app', {
 *     namespace: 'default',
 *     replicas: 2,
 *     ...config
 *   })
 * }
 *
 * const api = k8sApp('api', { replicas: 5 })
 * const worker = k8sApp('worker', { replicas: 3 })
 */
