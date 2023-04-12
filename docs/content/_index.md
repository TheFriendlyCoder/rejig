+++
archetype = "home"
title = 'Rejigger'
+++

Do you find it tedious starting a new programming project, and having to remember how to set up the structure for your project? Trying to think back to the last project you made, recalling how you set up your favorite CI tool, all of the build tools, document generators, test frameworks, dependencies, etc. giving yourself a migraine in the process? Well that is where **Rejigger** comes to the rescue.

Using **Rejigger** you can create a set of example projects that define your favorite project structures, and then apply those structures quickly and easy when creating new projects in the future, with a simple command something like:

```shell
rejig create . MyTemplate
```

**Rejigger** provides a framework for organizing these "template" projects in a way that makes them easy to organize and easy to customize when generating enw projects. You can store your favorite templates in folders on your local machine, or save them on shared folders or even in Git repositories, making them easy to find and organize. You can even share your favorite templates with your friends or coworkers.

In a professional setting, **Rejigger** allows your build engineers or infrastructure teams to develop patterns for how projects should be structured to work efficiently with specific tools and technologies so that the teams they work with can focus on writing code instead of trying to figure out the nuances of how the different build and automation tools work.

To get a feel for how **Rejigger** works, check out the [quickstart guide](getting_started/quickstart).