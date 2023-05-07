# Template Manifest

In order for **Rejigger** to recognize a folder as containing a template definition, it requires a [YAML](https://yaml.org) formatted manifest file named `.rejig.yml` in the root folder. This file contains a description of the template and various options it supports. An example of a complete manifest file showing all supported options is shown below:

```yaml
versions:
  schema: 1.0
  rejigger: 0.0.1
  template: 1.0
template:
  args:
    - name: project_name
      description: Name of the source code project
    - name: version
      description: Initial version number for the project
  exclusions:
    - "docs/.*"
```

Each sub-section of the manifest file is described in more detail in the following sections.

## Versions

The first section of the manifest file is defined by the key `versions`. It defines the versions of several key ingredients that are needed to properly process the template. The required properties are as follows:

* `schema` - used to tell **Rejigger** what contents to expect in the rest of the manifest file. Currently this should always be set to the value `1.0` because we have not yet required any breaking changes to the format.
* `rejigger` - indicates the minimum version of the *Rejigger* application that is required to process the contents of this template. If a user attempts to process a template with an older version of the application, the generation logic will fail.
* `template` - indicates the edition or revision of the template itself. This can be used to track changes to a template over time, and potentially add custom logic to the template to support migrating existing projects from older versions to newer versions without having to completely rebuild the project from scratch.

The values for each of these version fields is expected to conform to the format and interpretation of the [semantic versioning](https://semver.org) standards.

## Template

The values in this section are read by the **Rejigger** application and provide it with instructions on how to customize the projects generated from the template contents.

### Args

This optional subsection defines a list of customizable arguments or parameters that can be defined by a user and injected into various places during the generation process. *Rejigger* will use the definitions in this section to prompt the user for values which are relevant for the project they are creating. The app will then look for references to these args throughout the template content (files, folders, etc.) and replace references to them with the values provided by the user.

* `name` - unique variable name to save the custom value provided by the user. This name may then be used anywhere in the template to act as a placeholder for the value provided during generation. When used within the template contents, the argument name must be wrapped in dual curly-braces as in `{{ my_arg_name }}`. This allows the app to distinguish between references to template argument names and plain text which may exist in the template with the same sequence of characters.
* `description` - contains a very short explanation of the purpose of the argument. This text is displayed to the user by the application so they can understand better how the value they are providing will be used within the template. The text should be as short as possible so it doesn't occupy unnecessary space when shown on the console, while still providing enough detail for the user to understand what it is for.

**Example: simple file substitution**
Suppose you have a file in your template named "project_version.prop" and within that file you want to inject a custom version number which is to be defined by the user of your template. To accomplish this you could put the value `{{version}}` in the .prop file, and add the following definition to your manifest file:

```yaml
template:
  args:
    - name: version
      description: initial version for your application
```

When a user then tries to apply your template to create a new project (ie: using the `create` operation) they will be prompted for the following:

```
initial version for your application(version): 
```

Then, whatever value the user enters in this prompt will be used as a replacement value for the `{{version}}` token. So in our example case we would see a file named "project_version.prop" created in our new project folder, and the contents of the file would be updated to contain the value provided by the user.

**Example: simple folder substitution**

Suppose you have a folder in your template which needs to be named the same thing as the project generated from the template. This is a common occurrence with projects based on Java, Python and several other languages. You can have **Rejigger** automatically rename folders and files in the same way that it replaces argument tokens within the contents of files. Consider the following example template:

```yaml
template:
  args:
    - name: project
      description: name of the project being generated
```

When the user applies this template to create a new project, they will be prompted as follows:

```
name of the project being generated(project):
```

In order to have the template create a new folder with whatever project name the user provides, you would then create a folder in your template named something like `{{project}}`.

```
.
└── {{project}}
```

So if the user said the name of their new project was "ProcessLib", the newly created project would have a folder with that name created in place of the templated one.

```
.
└── ProcessLib
```

!!! note
    The reason **Rejigger** uses double curly-braces to delineate argument tokens is to ensure that text within your template does not get accidentally replaced with values provided by the user. For example, suppose your template needs to create a file named "project.props" inside a folder named after the actual project name selected by the user. To ensure the folder whose name must be replaced by the `project` argument, but the file within it keeps its exact name "project.props" you can simply define your template contents as:

    ```
    .
    └── {{project}}
        └── project.props
    ```

    This will result in a generated output that looks something like:
    
    ```
    .
    └── ProcessLib
        └── project.props
    ```

### Exclusions

This property allows you to provide a list of 0 or more [regular expressions](https://en.wikipedia.org/wiki/Regular_expression) which define files and folders which should be ignored by **Rejigger** when generating new projects from the template. This can be helpful if there are support files, documentation, or other such content stored in the same repository as your template but which should not be included in new projects.

!!! note
    It is helpful to use tools like [regex101](https://regex101.com) to test your regular expressions to make sure they work as you expect. Just make sure to select the "golang" language preferences in the tool to ensure you are using a validator that is compatible with **Rejigger**.
