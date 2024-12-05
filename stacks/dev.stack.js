const opts = {
  org: 'moonwalker',
  common_name: 'comet'
}

const stage = stack('dev')

backend('gcs', {
  bucket: 'mw-tf-state',
  prefix: `${opts.org}-${opts.common_name}/stacks/{{ .stack }}/{{ .component }}`
})

const p1 = component('project', 'test/module', {
  name: `${opts.common_name}-${stage.name}`
})

const p2 = component('project2', 'test/module', {
  name: `${opts.common_name}-{{ .stack }}`
})

const vpc = component('vpc', 'test/module', {
  name: `${p1.name}-vpc`,
  id: `${p2.id}-vpc`
})
