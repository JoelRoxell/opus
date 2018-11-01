# Opus

## WIP

CI/CD for frontend applications using docker.

TODO:

1. Create a server which can run on a bare machine or in a kubernetes cluster.
2. Create a client which can manage and list artifacts and from server.
3. Server should be able to listen receive triggers (i.e. github push).
4. Allow concurrent build.
5. Push events to external hooks.

First commands:

1. add / remove project.
2. trigger builds.
3. fetch status.
4. deployment management

## Usage WIP

Add project `opus -add <project-name> <source>`

Remove project `opus -delete <project-name>`

Trigger builds `opus -trigger <project-name> <tag|latest> <file-store>`

Status `opus -list <filter>`

Deployment `opus -deploy project-name <tag>`

Rollback `opus -rollback project-name`

Install server-processes to current kube-context `opus -install`

Super simplified initial sketch
[![Image from Gyazo](https://i.gyazo.com/bebc4df2dd3bcf61ea34c205ed3ec5bb.png)](https://gyazo.com/bebc4df2dd3bcf61ea34c205ed3ec5bb)

## Artifact structure

```bash
.
└── project-x
    ├── 1.0.1
    │   ├── 687ae5f03324498e19eba79dae8ce891.png
    │   ├── d834a0da39eba2788afe51cf86f0739a.png
    │   ├── d96cee9a1d15ab252fdd6dbf435472b6.png
    │   ├── f997be6320078b12cfd2552375e56242.png
    │   ├── index.html
    │   └── static
    │       ├── css
    │       │   └── main.2123c96.css
    │       └── js
    │           └── main.11835c96.js
    ├── 1.0.3
    │   ├── 687ae5f03324498e19eba79dae8ce891.png
    │   ├── d834a0da39eba2788afe51cf86f0739a.png
    │   ├── d96cee9a1d15ab252fdd6dbf435472b6.png
    │   ├── f997be6320078b12cfd2552375e56242.png
    │   ├── index.html
    │   └── static
    │       ├── css
    │       │   └── main.21823c23.css
    │       └── js
    │           └── main.1134kc96.js
    ├── 1.0.4
    │   ├── 687ae5f03324498e19eba79dae8ce891.png
    │   ├── d834a0da39eba2788afe51cf86f0739a.png
    │   ├── d96cee9a1d15ab252fdd6dbf435472b6.png
    │   ├── f997be6320078b12cfd2552375e56242.png
    │   ├── index.html
    │   └── static
    │       ├── css
    │       │   └── main.23965c96.css
    │       └── js
    │           └── main.2185da96.js
    ├── 2.0.0
    │   ├── 687ae5f03324498e19eba79dae8ce891.png
    │   ├── d834a0da39eba2788afe51cf86f0739a.png
    │   ├── d96cee9a1d15ab252fdd6dbf435472b6.png
    │   ├── f997be6320078b12cfd2552375e56242.png
    │   ├── index.html
    │   └── static
    │       ├── css
    │       │   └── main.22315b96.css
    │       └── js
    │           └── main.218jd236.js
    └── archive
        ├── 1.0.1.tar.gz
        ├── 1.0.3.tar.gz
        ├── 1.0.4.tar.gz
        └── 2.0.0.tar.gz
```
