
# Find PRs by owner, author, and branch

The optional branch filter is useful for excluding PRs that are not based on the main branch.

## Install

```bash
go install github.com/qrkourier/gh-pr-finder@main
```

## Run

```bash
GITHUB_TOKEN=$(gh auth status -t | awk '/Token/{print $3}') \
    $(go env GOPATH)/bin/gh-pr-finder \
        --owners openziti,openziti-test-kitchen,openziti-terraform-modules,netfoundry \
        --authors qrkourier \
        --branches main,release-next
```

Output

```text
https://github.com/openziti/helm-charts/pull/187
https://github.com/openziti/ziti-console/pull/289
https://github.com/openziti/ziti-console/pull/284
https://github.com/openziti/ziti-doc/pull/818
https://github.com/openziti/ziti-doc/pull/814
https://github.com/openziti/ziti/pull/1800
https://github.com/openziti/ziti-builder/pull/31
https://github.com/openziti/ziti-tunnel-sdk-c/pull/796
https://github.com/openziti/zrok/pull/553
https://github.com/openziti/helm-charts/pull/160
https://github.com/openziti-terraform-modules/terraform-k8s-openziti-controller/pull/4
https://github.com/openziti-terraform-modules/terraform-lke-ziti/pull/36
```
