# skaffold-buildkit-wrapper

## Usage

### For local docker daemon

Enabling the `useBuildkit` option in skaffold, you can use buildkit in the local docker daemon, but the experimental features, like `--mount=type=secret`, are not available now.

`skaffold-buildkit-wrapper` executes the docker CLI with the arguments of new features. (only the `--secret` option is supported now.)

```sample.Dockerfile
# syntax = docker/dockerfile:experimental

RUN --mount=type=secret,id=gitconfig,dst=/root/.gitconfig \
    git clone https://github.com/daisaru11/sample_private_repo.git
```

```
build:
  artifacts:
    - image: daisaru11/skaffold-buildkit-wrapper-example
      context: .
      custom:
        buildCommand: skaffold-buildkit-wrapper build --cli docker --secret id=gitconfig,src=.gitconfig --build-arg TEST=test1 --file sample.Dockerfile
```

See [example](https://github.com/daisaru11/skaffold-buildkit-wrapper/tree/master/example/docker_buildkit).

### For buildkitd cluster in Kuberenetes

`skaffold-buildkit-wrapper` connects a buildkitd daemon via the `buildkit` command.

Using `--kube-pod-selector` options, the pod to be connected is automatically chosen by the consistent hashing algorithm.

```
apiVersion: skaffold/v1beta15
kind: Config
build:
  local: {}
  artifacts:
    - image: daisaru11/skaffold-buildkit-wrapper-example
      context: .
      custom:
        buildCommand: skaffold-buildkit-wrapper build --cli buildctl --kube-pod-selector app=buildkitd --kube-pod-balancing-hash-key daisaru11/skaffold-buildkit-wrapper-example --secret id=gitconfig,src=.gitconfig --build-arg TEST=test1 --file sample.Dockerfile
```

See [example](https://github.com/daisaru11/skaffold-buildkit-wrapper/tree/master/example/k8s_buildkitd).

## References

- [github.com/moby/buildkit/tree/master/examples/kube-consistent-hash](https://github.com/moby/buildkit/tree/master/examples/kube-consistent-hash)
- [[KubeConEU] Building images efficiently and securely on Kubernetes with BuildKit](https://www.slideshare.net/AkihiroSuda/kubeconeu-building-images-efficiently-and-securely-on-kubernetes-with-buildkit)
- [buildkit as a cluster builder?](https://github.com/GoogleContainerTools/skaffold/issues/2642)
