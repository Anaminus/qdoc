# QuickDoc
The `qdoc` command converts documentation within a Lua script into a Markdown
file.

## Usage
```bash
qdoc [-ext EXT] [-o DIR] [GLOB...]
```


Each non-flag argument is a glob matching a number of input files. If no files
are specified, then all files with the `.lua` extension in the working directory
are matched.

The content of each matched file, assumed to have the Lua syntax, is parsed for
comments. Adjacent line comments are merged into a single group, while block
comments always have their own group.

Groups may contain directives of the form "@name" to control how they are
interpreted and converted. A directive divides the group into chunks; the name
of the directive determines the chunk type, while the rest of the chunk is the
content. An optional ":" may follow the directive to separate it from the
content as needed. The leading and trailing spaces of the content are trimmed.
Multiple directives of the same type can be specified, in which case the content
of each are concatenated together with newlines as separators, to form a single
chunk.

The following directives are processed:

- `@sec`

    Describes a section. Only the first line of the chunk is processed.

    The content determines the text of the section heading. Additionally, it is
    a dot-separated list of names, which determines the level and location of
    the section in the generated Markdown file. For example, section "Foo.Bar"
    will appear under section "Foo". If there is no parent section, then the
    section is placed in the root.

    If the content ends with a path wrapped in curly braces, that is used as the
    path instead, while the text before is used as the heading. This can be used
    to finely control the path, or to make it different from the heading. For
    example:

        Description of Foo.Bar {Foo.Bar}

    Each section defines a footnote that links to the heading, using the path as
    the name, allowing the section to be linked to more easily. For example:

        See [section Foo.Bar][Foo.Bar] for more information.

- `@ord`

    Describes the priority of the section. Only the first line of the chunk is
    processed.

    The content is a signed integer, which determines the order of the section
    among its siblings, sorted by ascending priority. By default, sections have
    a priority of 0. Sections that have the same priority are sorted by path.

- `@def`

    Describes a type definition. The content is enclosed in a fenced code block
    that appears at the top of the section.

    This directive is optional.

- `@doc`

    Describes the content of the section. This is general markdown text.

Other directives are ignored. An entire group is ignored if it does not contain
a valid `@sec` directive.

An output file is not written if its content is empty or only contains spacing.

The `-o` flag is a path to a directory to which output files will be written. If
`-o` is not specified, then an output file is written to the same directory as
the input file.

To determine the location of an output file, the extension of each matched file
is replaced with ".md", which can be changed with the `-ext` flag. An input file
is skipped if the output and input are the same file.

## Installation
1. [Install Go](https://golang.org/doc/install)
2. [Install Git](http://git-scm.com/downloads)
3. Using a shell with Git (such as Git Bash), run the following command:

```
go get github.com/anaminus/qdoc
```

If Go was installed correctly, this will install qdoc to `$GOPATH/bin`, which
will allow it to be run directly from a shell.
