---
sidebar_position: 7
---

# Kubernetes Integration

Comet can generate kubeconfig files to connect to your Kubernetes clusters. This is useful when your infrastructure includes managed Kubernetes services like GKE, EKS, AKS, or DigitalOcean Kubernetes.

## Basic Usage

The `kubeconfig()` function configures Kubernetes cluster access in your stack:

```javascript title="stacks/production.stack.js"
stack('production', { settings })

backend('gcs', {
  bucket: 'my-terraform-state',
  prefix: 'comet/{{ .stack }}/{{ .component }}'
})

// Define your cluster
const cluster = component('gke', 'modules/gke', {
  cluster_name: 'production-cluster',
  region: 'us-central1'
})

// Configure kubeconfig
kubeconfig({
  current: 0,
  clusters: [{
    context: 'production-gke',
    host: cluster.endpoint,
    cert: cluster.ca_certificate,
    exec_command: 'gke-gcloud-auth-plugin',
    exec_args: []
  }]
})
```

## Authentication Methods

Comet supports two authentication methods for Kubernetes clusters:

### Token-Based Authentication

Use static bearer tokens for simple authentication. Best for CI/CD pipelines and service accounts:

```javascript
kubeconfig({
  current: 0,
  clusters: [{
    context: 'my-cluster',
    host: 'https://kubernetes.example.com',
    cert: 'LS0tLS1CRUdJTi...',  // Base64-encoded CA cert
    token: 'dop_v1_...'          // Static bearer token
  }]
})
```

**Pros:**
- Simple setup, no external tools needed
- Works well in CI/CD environments
- No need to install cloud CLI tools

**Cons:**
- Less secure than exec-based auth
- Tokens may expire and need rotation
- Not recommended for interactive use

### Exec-Based Authentication (Recommended)

Use cloud provider CLI tools for dynamic credential generation:

```javascript
kubeconfig({
  current: 0,
  clusters: [{
    context: 'my-cluster',
    host: 'https://kubernetes.example.com',
    cert: 'LS0tLS1CRUdJTi...',
    exec_command: 'doctl',
    exec_args: [
      'kubernetes',
      'cluster',
      'kubeconfig',
      'exec-credential',
      '--version=v1beta1',
      'my-cluster'
    ]
  }]
})
```

**Pros:**
- More secure - credentials refresh automatically
- Uses your existing cloud authentication
- Recommended for interactive use

**Cons:**
- Requires cloud CLI tools installed
- May add latency to kubectl commands
- More complex setup

## Using with Secrets

Store sensitive cluster information in secrets:

```javascript
// Configure secrets
secretsConfig({
  defaultProvider: 'sops',
  defaultPath: 'secrets.enc.yaml'
})

// Get cluster config from secrets
const cluster = secret('k8s/production')

kubeconfig({
  current: 0,
  clusters: [{
    context: cluster.context,
    host: cluster.host,
    cert: cluster.cert,
    token: cluster.token
  }]
})
```

**Example secrets file:**

```yaml
# secrets.enc.yaml (encrypted with SOPS)
k8s:
  production:
    context: "production-cluster"
    host: "https://1.2.3.4"
    cert: "LS0tLS1CRUdJTi..."
    token: "dop_v1_..."
```

## Multiple Clusters

Configure access to multiple clusters in the same stack:

```javascript
kubeconfig({
  current: 0,  // Index of default cluster
  clusters: [
    {
      context: 'production-primary',
      host: 'https://prod-us.k8s.example.com',
      cert: 'LS0tLS1CRUdJTi...',
      exec_command: 'gcloud',
      exec_args: ['container', 'clusters', 'get-credentials', 'prod-us']
    },
    {
      context: 'production-backup',
      host: 'https://prod-eu.k8s.example.com',
      cert: 'LS0tLS1CRUdJTi...',
      exec_command: 'gcloud',
      exec_args: ['container', 'clusters', 'get-credentials', 'prod-eu']
    }
  ]
})
```

The `current` field specifies which cluster is the default (0-indexed).

## Cloud Provider Examples

### Google Kubernetes Engine (GKE)

```javascript
kubeconfig({
  current: 0,
  clusters: [{
    context: 'gke-production',
    host: cluster.endpoint,
    cert: cluster.ca_certificate,
    exec_command: 'gke-gcloud-auth-plugin',
    exec_args: []
  }]
})
```

### DigitalOcean Kubernetes (DOKS)

```javascript
// With token (simple)
kubeconfig({
  current: 0,
  clusters: [{
    context: 'doks-production',
    host: cluster.endpoint,
    cert: cluster.ca_certificate,
    token: cluster.kube_token
  }]
})

// With doctl (recommended)
kubeconfig({
  current: 0,
  clusters: [{
    context: 'doks-production',
    host: cluster.endpoint,
    cert: cluster.ca_certificate,
    exec_command: 'doctl',
    exec_args: [
      'kubernetes',
      'cluster',
      'kubeconfig',
      'exec-credential',
      '--version=v1beta1',
      cluster.name
    ]
  }]
})
```

### Amazon EKS

```javascript
kubeconfig({
  current: 0,
  clusters: [{
    context: 'eks-production',
    host: cluster.endpoint,
    cert: cluster.ca_certificate,
    exec_command: 'aws',
    exec_args: [
      'eks',
      'get-token',
      '--cluster-name',
      cluster.name,
      '--region',
      'us-west-2'
    ]
  }]
})
```

### Azure Kubernetes Service (AKS)

```javascript
kubeconfig({
  current: 0,
  clusters: [{
    context: 'aks-production',
    host: cluster.endpoint,
    cert: cluster.ca_certificate,
    exec_command: 'kubelogin',
    exec_args: [
      'get-token',
      '--environment',
      'AzurePublicCloud',
      '--server-id',
      cluster.server_id,
      '--client-id',
      cluster.client_id,
      '--tenant-id',
      cluster.tenant_id
    ]
  }]
})
```

## Generating Kubeconfig

After defining your stack, generate the kubeconfig:

```bash
# Generate kubeconfig for the stack
comet kube production gke

# Use the generated kubeconfig
export KUBECONFIG=~/.kube/comet-production-gke
kubectl get nodes
```

This merges the cluster configuration into your kubeconfig file at `~/.kube/config` and sets it as the current context.

## CI/CD Usage

For CI/CD pipelines, use token-based authentication:

```javascript
// Store token in CI secrets
const token = env.K8S_TOKEN

kubeconfig({
  current: 0,
  clusters: [{
    context: 'ci-cluster',
    host: env.K8S_HOST,
    cert: env.K8S_CERT,
    token: token
  }]
})
```

Then in your CI pipeline:

```yaml
# .github/workflows/deploy.yml
- name: Setup kubeconfig
  run: |
    comet kube production cluster
    kubectl apply -f manifests/
  env:
    K8S_TOKEN: ${{ secrets.K8S_TOKEN }}
    K8S_HOST: ${{ secrets.K8S_HOST }}
    K8S_CERT: ${{ secrets.K8S_CERT }}
```

## Best Practices

1. **Use exec-based auth for interactive use** - More secure and credentials refresh automatically
2. **Use token-based auth for CI/CD** - Simpler and doesn't require cloud CLI tools
3. **Store credentials in secrets** - Never commit tokens or certificates to git
4. **Rotate tokens regularly** - If using token-based auth, rotate tokens periodically
5. **Use separate clusters for environments** - Don't share clusters between dev/staging/prod
6. **Limit token permissions** - Create service accounts with minimal required permissions

## Troubleshooting

### "Unable to connect to the server"

Check that your cluster endpoint is correct and accessible:

```bash
# Test connectivity
curl -k https://your-cluster-endpoint

# Check kubeconfig
kubectl config view
```

### "error: You must be logged in to the server"

For exec-based auth, ensure the CLI tool is installed and authenticated:

```bash
# GKE
gcloud auth login
gcloud container clusters get-credentials cluster-name

# DigitalOcean
doctl auth init
doctl kubernetes cluster kubeconfig save cluster-name

# AWS
aws configure
aws eks update-kubeconfig --name cluster-name
```

### Token expired

Rotate the token and update your secrets:

```bash
# Get new token from your cluster
kubectl create token service-account-name

# Update secrets file
# Then regenerate kubeconfig
comet kube production cluster
```
