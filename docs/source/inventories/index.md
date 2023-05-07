# Inventories

**Rejigger** inventories are simply groupings or master lists of [templates](../tmpl) which are stored together to make them easier to locate and manage. Inventories have several benefits over managing templates individually including:

* making it easy to locate the source for the templates (ie: co-location)
* making it easy to create / add new templates (ie: tools configured to use the inventory will automatically recognize them without further intervention)
* making it easier to share templates with friends and coworkers (ie: one URL to share an entire set of templates)

To create an inventory you simply need to create a new folder containing a file named `.rejig.inv.yml` which we refer to as a **template manifest** and populate it with the names and locations of the individual templates contained within the inventory. For more details on the contents of this manifest file see [here](manifest.md).


**Example**

Suppose you make regular use of two programming languages: GO and Python, and you want to create one template for each type of project. This can be done by creating a folder with two sub-folders, one for each of your templates, and adding an inventory manifest to the root resulting in a folder structure something like this:

```linenums="0"
.
├── .rejig.inv.yml
├── golang
│   └── .rejig.yml
│       .... more files here for the GO template
├── python
│   └── .rejig.yml
│       .... more files here for the Python template
```

* `.rejig.inv.yml` - the [inventory manifest](manifest.md) with links to the templates
* `golang` - folder containing the files and folders for the GO template
* `golang/.rejig.yml` - [template manifest](../tmpl/manifest.md) for the GO template
* `python` - folder containing the files and folders for the Python template
* `python/.rejig.yml` - [template manifest](../tmpl/manifest.md) for the python template

In this example, the `.rejig.inv.yml` file would contain the following:

```yaml
templates:
  - source: ./golang
    name: go
  - source: ./python
    name: python
```

You could then add this new inventory to your [application options](../app_options/index.md#inventories) by doing something like:

```yaml
inventories:
  - type: local
    source: /path/to/my/inventory
    namespace: my
```

and then use of your new templates by doing something like:

```linenums="0"
rejig Proj1Dir my.go
rejig Proj2Dir my.python
```

!!! note
    **Rejigger** prepends each inventory template with the name of the inventory to avoid naming conflicts between inventories. So, for example, two different inventories may provide a template named `python` and they would be referred to as `inventory1.python` and `inventory2.python` respectfully (assuming their respective inventories were named `inventory1` and `inventory2`)