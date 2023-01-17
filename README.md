# Algome

This is the CLI tool to generate readme with structure of solved algorithmic problems of different
sources.

## Installation

### Install the latest version using Go

`go get github.com/malyusha/algome`

### Linux

* Download the latest Linux .gz archive
  from [Releases page](https://github.com/malyusha/algome/releases)
* Extract the archive
* Move binary to desired directory (used in $PATH env variable)
* Set "execute" permissions for the binary

```shell
# Extract the archive
tar xf algome_0.0.1_linux_x86_64.tar.gz
# Move the binary
sudo mv algome /usr/local/bin
# Set permissions
sudo chmod +x /usr/local/bin/algome
```

### MacOS

* Download the latest macOS archive from [Releases page](https://github.com/malyusha/algome/releases)
* Extract the archive
* Move binary to desired directory (used in $PATH env variable)
* Set "execute" permissions for the binary

```shell
# Extract the archive
tar xf algome_0.0.1_darvin_x86_64.tar.gz
# Move the binary
sudo mv algome /usr/local/bin
# Set permissions
sudo chmod +x /usr/local/bin/algome
```

### Windows

Download the latest Windows .gz archive from [Releases page](https://github.com/malyusha/algome/releases)

## Definitions

* `Source provider` - source of problems (e.g. list of problems from LeetCode). Source provider is
  responsible for loading all problems from source. Once loaded, all problems are cached in JSON
  files to prevent unnecessary network requests.

## Configuration

Readme generation is based on JSON configuration provided via file. By default, CLI will look for
the file `algome.conf.json` inside of current directory. If no file is found, then the default
configuration is used.

The default configuration contains following JSON:

```json
{
  "output": "./",
  "structure": {
    "catalog": {
      "map_attr": "id",
      "base_dir": "./"
    }
  },
  "sources": ["leetcode"],
  "problems_cache_dir": "~/.algome/cache",
  "templates_dir": null
}
```

Let's look at each property of the configuration.

| name                 | type       | description                                                                                                                                                                              | default         |
|----------------------|------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------|
| `output`             | `string`   | Output directory for index readme.md. Index readme is an index file for allsources of problems                                                                                           | `./`            |
| `sources`            | `[]string` | Array of source names that are supported. For each of source separated readme is generated inside solutions directory. See [supported sources](#Supported sources for readme generation) | `["leetcode"]`  |
| `problems_cache_dir` | `string`   | Directory, where the cached problems of source are stored as JSON                                                                                                                        | ~/.algome/cache |
| `templates_dir`      | `string`   | Templates directory. Allows to override markdown templates for readme files. See [supported templates](#Override templates).                                                             | null            |

## Usage

There are multiple commands currently available for usage:

### `init`

Initializes configuration file. Not required command. By default, the following configuration is
used:

```json
{
  "structure": {
    "catalog": {
      "base_dir": "./"
    }
  },
  "sources": ["leetcode"]
}
```

**Example**

```
$ algome init
```

### `generate`

Generate readme files from available sources, specified in configuration.

**Example**

```
$ algome [generate]
```

## Supported sources for readme generation

* `leetcode` - LeetCode problems provider.

## Override templates

It's possible to override templates used for readme generation. All templates are parsed using
GoLang text/template package.
To override template, place file with name matching with supported template names inside directory,
then specify that directory in configuration file as the value of property `templates_dir`.

**Template names to override**

* `source_markdown` - template, responsible for problems of single source provider (e.g. LeetCode).
* `index` - index template, which contains list of sources with problems. It's an entrypoint readme
  of your repository.