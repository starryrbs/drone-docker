# drone-docker

drone docker插件，用于构建docker镜像，push docker镜像

解决官方[drone-docker](http://plugins.drone.io/drone-plugins/drone-docker/)不能很好支持缓存的问题

## 使用方式

```yaml

kind: pipeline
name: default

steps:
- name: docker
  image: plugins/docker
  environment:
    /var/lib/docker:/var/lib/docker
  settings:
    dockerfile: src/Frontend/entrypoint/Dockerfile
    context: src/Frontend/entrypoint
    username: kevinbacon
    password: pa55word
    repo: foo/bar
    tags: latest
    registry: registry.xxx.com
```

