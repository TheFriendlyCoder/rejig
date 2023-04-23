+++
title = "Inventories"
weight = 3
+++

**Rejigger** inventories are simply groupings or master lists of [templates](../templates) which are stored together to make them easier to locate and manage. Inventories have several benefits over managing templates individually including:

* making it easy to locate the source for the templates (ie: co-location)
* making it easy to create / add new templates (ie: tools configured to use the inventory will automatically recognize them without further intervention)
* making it easier to share templates with friends and coworkers (ie: one URL to share an entire set of templates)

To create an inventory you simply need to create a new folder containing a [YAML](https://yaml.org) formatted manifest file named `.rejig.inv.yml` and populate it with the names and locations of the templates contained in the inventory. Then you can just create one sub-folder for each of your templates and add a reference to each one in the inventory.

**Example**

Suppose you make regular use of two programming languages: GO and Python, and you want to create one template for each type of project. This can be done by creating a new sub-folder for each of your templates, and adding an inventory manifest to the parent folder, resulting in a folder structure something like this:

```shell
.
├── .rejig.inv.yml
├── golang
│   └── .rejig.yml
│       .... more files here for the GO template
├── python
│   └── .rejig.yml
│       .... more files here for the Python template
```

* `.rejig.inv.yml` - defines the inventory contents
* `golang` - folder containing the files and folders for the GO template
* `golang/.rejig.yml` - [template manifest](../templates/manifest.md) for the GO template
* `python` - folder containing the files and folders for the Python template
* `python/.rejig.yml` - [template manifest](../templates/manifest.md) for the python template

In this example, the `.rejig.inv.yml` file would contain the following definition:

```yaml
templates:
  - source: ./golang
    name: go
  - source: ./python
    name: python
```

You could then add this new inventory to your [application options](../app_options#inventories) by doing something like:

```yaml
inventories:
  - type: local
    source: /path/to/my/inventory
    namespace: my
```

and then use of your new templates by doing something like:

```shell
rejig Proj1Dir my.go
rejig Proj2Dir my.python
```