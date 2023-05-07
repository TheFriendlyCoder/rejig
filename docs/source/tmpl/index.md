# Templates

**Rejigger** templates are simply folders which contain a [YAML](https://yaml.org) formatted manifest file describing the template and its various options, along with any other files and folders which define the content of the project it creates. This section describes how to define a template and how the application manipulates the template to generate new projects from it.

**Example**

Suppose you have a folder which contains the following files:

```
.
├── .rejig.yml
├── {{project_name}}
│   └── version.txt
│       .... more files and folders here for the template  
├── README.md
└── .... more files here for the template
```

This folder is considered a **template** simply by the fact that it contains a file named `.rejig.yml`. This file, referred to as the **template manifest** or just the **manifest**, contains instructions on how **Rejigger** should process the rest of the contents of the folder when creating new projects. Details on the contents of these manifest files can be found [here](manifest.md).

You could then add this new template to your [application options](../app_options/index.md#templates) by doing something like:

```yaml
templates:
  - type: local
    source: /path/to/my/template
    name: myFirstTemplate
```

and then use of your new template by doing something like:

```
rejig ProjDir myFirstTemplate
```