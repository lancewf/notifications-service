pkg_name=notifications-service
pkg_origin=lancewf
pkg_version="0.1.0"
pkg_maintainer="Lance Finfrock <lancewf@gmail.com>"
pkg_description="Notifications Service"
pkg_license=("Chef-MLSA")
pkg_upstream_url="https://github.com/lancewf/notifications-service"
pkg_deps=(
  core/glibc
)
pkg_build_deps=(
  core/gcc
)
pkg_exports=(
  [port]=service.port
  [host]=service.host
)
pkg_exposes=(port)

pkg_bin_dirs=(bin)
pkg_scaffolding=core/scaffolding-go
scaffolding_go_base_path=github.com/lancewf

scaffolding_go_build_deps=(
 github.com/icrowley/fake
 github.com/go-chef/chef
 github.com/satori/go.uuid
 github.com/sirupsen/logrus
 github.com/spf13/cobra
 github.com/spf13/viper
)
