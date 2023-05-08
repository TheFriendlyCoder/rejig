# Templates

## Overview

Templates in **Rejigger** are simply specially formatted folders containing a [YAML](https://yaml.org) file in the root folder named `.rejig.yml` (referred to as the template manifest) and a set of example source files that will be used to generate new projects from it.

## Creating a template

You can create a simple template as follows

1. create a new folder for your template in your home folder, say `~/mytemplate`
2. create a new file named `.rejig.yml` in the empty folder with the following content:
   ```yaml
   versions:
     schema: 1.0
     rejigger: 0.0.1
     template: 1.0
   template:
     args:
       - name: project_name
         description: Name of the source code project
   ```
   For details on the contents of this file, see [here](../tmpl).
3. Create another file named `README.md` alongside the manifest file, and put the following content in it:

   ```markdown
   Welcome to the {{project_name}} project!
   ```

That's it! Congratulations, you've just made your first template!

## Using your new template

To test out your newly created template, do the following:

1. Add the following lines to a file named `.rejig` in your home folder:
   ```yaml
   templates:
      - type: local
        source: ~/mytemplate
        name: MyTemplate
   ```
2. Next, run the following command:
   ```
   rejig create ./MyNewProject MyTemplate
   ```
3. You should be prompted to enter a value for the `project_name` argument. Type any value that you like (ie: CoolApp)
4. Once the operation completes, you should see a new folder named `MyNewProject`. In that folder you should see a file named `README.md` with the contents "Welcome to the CoolApp project"!