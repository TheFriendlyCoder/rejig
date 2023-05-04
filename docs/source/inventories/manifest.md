# Inventory Manifest

Every **Rejigger** inventory must contain a [YAML](https://yaml.org) formatted text file in the root / parent folder of the inventory named `.rejig.inv.yml`. This file contains links and descriptions for one or more [templates](../tmpl). This file has the following structure:


```yaml
templates:
  - source: ./template1
    name: MyFirstTemplate
  - source: ./template2
    name: MySecondTemplate
  # ... more template references here
```

## templates

This is the primary node containing the list of templates managed by the inventory. Each list item under this section requires the following properties:

!!! tip
    Currently, inventories must contain all template sources under the same file system location / path as the inventory itself. So for inventories stored on disk, each template definition must be stored as sub-folders under the root inventory folder. For inventories stored in Git repositories, each template definition must be stored as sub-folders under the root of the same Git repository.

* `source` - relative path to the template in the inventory. This folder is expected to contain the entire definition for the template, as defined by the [template definition](../tmpl)
* `name` - each template within the inventory must be given a unique name. This allows users to uniquely identify and reference each template.

!!! tip
    Template names within the **same inventory** must be unique, but templates from **different inventories** may use the same names. **Rejigger** prepends the name of the inventory to the template name to avoid naming conflicts, similar to the way a namespace works in most programming languages.

