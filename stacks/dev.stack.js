const opts = {
  org: 'moonwalker',
  common_name: 'comet'
}

const stage = stack('dev', { opts })

metadata({
  description: 'Development environment for Comet testing',
  owner: 'platform-team',
  tags: ['dev', 'testing', 'non-prod'],
  custom: {
    server_type: 'cx33',
    vcpu: 4,
    ram_gb: 8,
    storage_gb: 80,
    monthly_cost_eur: 5.49
  }
})

backend('gcs', {
  bucket: 'mw-tf-state',
  prefix: `${opts.org}-${opts.common_name}/stacks/{{ .stack }}/{{ .component }}`
})

const p1 = component('project', 'test/modules/project', {
  name: `${opts.common_name}-${stage.name}-p1`
})

const p2 = component('project2', 'test/modules/project', {
  name: `${opts.common_name}-{{ .stack }}-p2`
})

const vpc = component('vpc', 'test/modules/vpc', {
  name: `${p1.name}-vpc`,
  id: `${p2.id}-vpc`
})

const gke = component('gke', 'test/modules/gke', {
  inputs: {
    project_id: p1.id
  }
})

append('providers', [`data "google_client_config" "default" {}`])

const k8s = {
  // alias: 'main',
  host: gke.kube_host,
  cluster_ca_certificate: gke.kube_cert,
  token: 'data.google_client_config.default.access_token'
}

const helm = {
  kubernetes: k8s
}

const metsrv = component('metsrv', 'test/modules/kubernetes', {
  providers: {
    // google: {},
    kubernetes: k8s,
    helm: helm
  }
})

kubeconfig({
  current: 0,
  clusters: [
    {
      context: 'cluster-1-gke-eu',
      host: '1.2.3.4',
      cert: 'Zm9vYmFy',
      exec_command: 'gke-gcloud-auth-plugin',
      exec_args: [
        'kubernetes',
        'cluster',
        'kubeconfig',
        'exec-credential',
        '--version=v1beta1',
        '--context=default'
      ]
    }
  ]
})
