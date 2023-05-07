# Inventories

## Overview

Inventories are simply [YAML](https://yaml.org) files containing lists of templates which are stored together in a common location for easy access. Inventories may be stored in different locations including on the local file system as well as in Git repositories. The recommended best practice is to store inventories in Git so they can be more easily shared. 

## Creating a Git-based inventory

Creating a new inventory is as simple as follows:

1. Create a new Git repository for your templates
2. Create a sub-folder under the root of the repository for your new template (ie: `./sample`) and place the contents of your template here
3. In the root folder of the repository, create a file named `.rejig.inv.yml` and provide a name and path for your new template using the following syntax:

    ```yaml
    templates:
      - name: sample1
        source: ./sample
    ```

4. Repeat this process to add as many templates to your repository as you like

## Using an inventory

To make use of the new inventory you've just created, simply add a section to the app configuration file `~/.rejig` that looks like this: 

```yaml
inventories:
  - namespace: <NamespaceName>
    type: git
    source: <repoURL>
```

replacing `<repoURL>` with the URL to your new Git repo. Now you can make use of templates in your new inventory by referring to them by the namespace name `<NamespaceName>`. So, for example, if you named your inventory namespace `personal` and you've created a template in it named `python`, you can create a new project using that template by running: 

```
rejig create ./projdir personal.python
```

!!! note
   Templates within an inventory are referenced by their namespace name AND template name using the format `<namespace>.<template>`.
