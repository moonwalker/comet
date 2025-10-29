// Example: Using static token authentication in kubeconfig()
//
// This demonstrates how to use token-based authentication instead of
// exec-based credential providers. Useful for CI/CD pipelines and
// simplified environments.

stack({
  name: 'token-auth-example',
  backend: 'local'
})

// Example 1: Token-based authentication
kubeconfig({
  current: 0,
  clusters: [
    {
      context: 'my-cluster-token',
      host: 'https://kubernetes.example.com',
      cert: 'LS0tLS1CRUdJTi...', // Base64-encoded CA cert
      token: 'eyJhbGciOiJSUzI1NiIsImtpZCI6Ij...' // Static bearer token
    }
  ]
})

// Example 2: Exec-based authentication (existing behavior)
kubeconfig({
  current: 0,
  clusters: [
    {
      context: 'my-cluster-exec',
      host: 'https://kubernetes.example.com',
      cert: 'LS0tLS1CRUdJTi...',
      exec_command: 'doctl',
      exec_args: ['kubernetes', 'cluster', 'kubeconfig', 'exec-credential', '--version=v1beta1', 'my-cluster']
    }
  ]
})

// Example 3: Using with secrets (recommended for tokens)
// const cluster = secret('op://vault/k8s-cluster/config')
//
// kubeconfig({
//   current: 0,
//   clusters: [
//     {
//       context: cluster.context,
//       host: cluster.host,
//       cert: cluster.cert,
//       token: cluster.token  // Token from secrets provider
//     }
//   ]
// })

// Example 4: Multiple clusters with mixed authentication
kubeconfig({
  current: 0,
  clusters: [
    {
      context: 'prod-cluster',
      host: 'https://prod.k8s.example.com',
      cert: 'LS0tLS1CRUdJTi...',
      exec_command: 'gcloud',
      exec_args: ['container', 'clusters', 'get-credentials', 'prod-cluster']
    },
    {
      context: 'ci-cluster',
      host: 'https://ci.k8s.example.com',
      cert: 'LS0tLS1CRUdJTi...',
      token: 'service-account-token-here' // Static token for CI
    }
  ]
})
