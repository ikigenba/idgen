# Author a spec with the Codex CLI

This page walks through the process of authoring a spec.

## Requirements

- The Codex CLI, installed and authenticated.
- The `idgen` binary on your path.

## Part 1 — a new project from scratch

### 1. Copy the skills to the repo.

I prefer to keep everything related to a project in the project, and this
includes the related skills.  Start by copying the `.agents` directory
structure from an existing project into the new one.


```sh
cp -r ../idgen/.agents .
```

### 2. Load the spec context

Start `codex` in the new repo, then load the spec contracts.

```
use $spec-shapes
```

### 3. Describe the goal in plain English

Just talk and describe the app, whatever matters to you.

> Let's build a small Go CLI application called foo. It will be a simple echo
> application. Every argument passed to it, separated by whitespace, is printed
> back to stdout on a new line.
>
> The main function should be structured so it passes its arguments to a foo.Run
> function. The idea is that the Run function is then testable. The Run function
> just returns an error. If an error is returned main prints the message to
> stderr and exits with a non-zero exit code.
>
> Create a Makefile with these targets: build, test, run, install, fmt, clean. The
> build target should be the default. The install target should accept the standard
> Unix PREFIX argument so that install directories can be overridden. If no
> PREFIX is provided it should use the standard go install.
>
> $grillme

The agent will interrogate you until it runs out of questions.


### 4. Codify

Then ask it to write the design decisions and build phase breakdown.

```
$codify
```

### 5. Generate the build loop (once per project)

Then, if this is the first run, generate the prompts for ralph. This only needs to be
done the first time; after that, you can commit the project/loop prompts to git.

```
$create-gather-build-verify-prompts
```

### 6. Run it

```sh
project/loops/run
```
