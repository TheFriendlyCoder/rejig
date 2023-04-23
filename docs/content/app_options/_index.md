+++
title = "Application Options"
weight = 4
+++

**Rejigger** makes use of a [YAML](https://yaml.org) formatted configuration file stored on the local machine to find available templates and to customize its behavior. This file is expected to be named `.rejig` or `.rejig.yml` and it should be located in the home folder of the user running the tool. An example of what the app options file should look like is shown below:

```yaml
templates:
  - type: git
    source: git@github.com:TheFriendlyCoder/rejigger.git
    subdir: samples/simple
    name: sample
inventories:
  - type: git
    source: https://github.com/TheFriendlyCoder/rejigger.git
    namespace: demo
options:
  theme: dark
```

## Templates

This subsection contains a list of 0 or more definitions of the [templates](../templates) currently available to the application. As you create your own templates, or are given links to templates shared by your friends and workmates, you can add them here so they can be used by the application. Each template reference supports the following properties:

* `type` - (required) describes the protocol to be used when accessing the template definition. May be one of the following values:
  * `git` - indicates the template is stored in a remote Git repository
  * `local` - indicates the template is stored on the local file system
* `source` - (required) provides either the path to the template (if `type` is `local`) or the URL to the remote repository (if `type` is `git`).
* `subdir` - (optional) provides a relative path within the `source` location where the template definition exists. If not provided, the application will assume the template definition will be stored in the root folder.
* `name` - (required) this is a friendly, easy to remember name you give to the template. It is used when referring to the template on the command line, like when using a template to create a new project using the `create` command. It must be unique across all the templates in your options file.

{{% notice note "Note" %}} 
Make sure to use the URL you would use for checking out a remote template from a Git repository using a git client, and not the URL for the landing page for the GitHub / Gitlab / Bitbucket project (ie: "https://github.com/TheFriendlyCoder/rejigger.git" and not "https://github.com/TheFriendlyCoder/rejigger")
{{% /notice %}}

{{% notice warning "Warning" %}}
**Rejigger** currently only supports pulling templates from authenticated remote Git repos using the SSH protocol. To make sure this works as expected you will need to use the `git` or `ssh` URL instead of the `http` or `https` URL when you define your template source (ie: "git@github.com:TheFriendlyCoder/rejigger.git" instead of "https://github.com/TheFriendlyCoder/rejigger.git"). In addition, you will need to make sure your git client is properly configured to authenticate to the remote repository using an SSH key (ie: typically stored in `~/.ssh/id_rsa`)
{{% /notice %}}

## Inventories

This subsection contains a list of 0 or [inventories](../inventories) of templates available to the application. As the number of templates you create grows, you may want to group them together into a single location to make them easier to manage. This section in your options file allows you to point to the central location for batches of templates. Each element in this section supports the following options:

* `type` - (required) describes the protocol to be used when accessing the inventory definition. May be one of the following values:
  * `git` - indicates the inventory is stored in a remote Git repository
  * `local` - indicates the inventory is stored on the local file system
* `source` - (required) provides either the path to the template (if `type` is `local`) or the URL to the remote repository (if `type` is `git`). Unlike templates, inventory definitions are assumed to be located in the root folder of the source location, with each template being stored in a sub-folder.
* `namespace` - (required) similar to the template `name`, this is a friendly identifier you give to the inventory to make it easy to reference within the application. **NOTE:** templates stored within an inventory need to be referenced by their namespace name followed by a period separator, as in "MyNamespace.MyTemplate".

## Options

This section allows you to customize the behavior of the **Rejigger** application. It supports the following properties, all of which are optional:

* `theme` - allows you to adjust the color scheme used by the application when displaying text content on the console. When not defined, the tool defaults to monochrome output. The supported values for this property are:
  * `dark` - colors that work well for a dark background with lighter colored text
  * `light` - colors that work well for a light background with darker colored text