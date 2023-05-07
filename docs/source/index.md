Do you find it tedious to start a new programming project, and having to remember how to set up the structure for it? Trying to think back to the last project you made, recalling how you set up your favorite CI tool, all of the build tools, document generators, test frameworks, dependencies, etc. giving yourself a migraine in the process? Well that is where **Rejigger** comes to the rescue.

Using **Rejigger** you can keep track of your favorite project setups by defining a set of patterns in special template projects, which can then be applied quickly and easily when creating new projects in the future using simple commands like the following:

```linenums="0"
rejig create . MyTemplate
```

**Rejigger** provides a framework for organizing these "template" projects in a way that makes them easy to find and easy to customize when generating new projects. You can store your favorite templates in folders on your local machine, save them on shared folders or publish them to Git. This makes it easy to share your favorite templates with your friends and coworkers.

In a professional setting, **Rejigger** allows your build engineers or infrastructure teams to develop patterns for how projects should be structured to work efficiently with specific tools and technologies. That way dev teams can focus on writing code instead of trying to figure out the nuances of how the different build and automation tools work.

To get a feel for how **Rejigger** works, check out the [quickstart guide](getting_started).