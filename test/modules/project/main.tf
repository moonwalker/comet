resource "time_static" "project_ts" {}

locals {
  project_id = var.id != "" ? var.id : "${var.name}-${time_static.project_ts.unix}"
}
