const opts = {
  org: 'moonwalker',
  common_name: 'comet'
}

const stage = stack('dev')

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
