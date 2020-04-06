# karabiner-config

Generate karabiner-Elements configuration file with simple syntax.

![](https://github.com/uphy/karabiner-config/workflows/Release/badge.svg)

[karabiner-Elements](https://github.com/tekezo/Karabiner) is a great keyboard customization tool but its configuration file is a little redundant and hard to maintain.

`karabiner-config` generates the configuration file from more simple and readable configuration file.

- YAML format
- Commonize conditions between multiple manipulators (e.g., `frontmost_application_if`, `frontmost_application_unless`)
- Shorten frequently used configurations
  - Set variable
  - Condition based on variable
  - Condition based on frontmost application
- In `manipulator` configuration switch

## Installation

### Binary

See [release](https://github.com/uphy/karabiner-config/releases) page.

### Build

```sh
$ git clone https://github.com/uphy/karabiner-config
$ cd karabiner-config
$ make install
```

## Getting Started

Create new rule file.
This example defines keybinding `Control-m` to start a new line.

`sample/instruction/01_c-m.yml`

```yaml
title: Control+M to newline
maintainers:
  - uphy
rules:
  - description: Control+H to newline
    manipulators:
      - from:
          key: C-m
        to:
          - key: return_or_enter
```

Then import this rule to the existing karabiner config.

**Note**: This command replaces your karabiner config's `Default profile` directly.  Please back up first.

```sh
$ karabiner-config sample/instruction/01_c-m.yml ~/.config/karabiner/karabiner.json
```

Also you can watch changes of the yaml file and automatically update the karabiner config with `-w` option.

```sh
$ karabiner-config -w sample/instruction/01_c-m.yml ~/.config/karabiner/karabiner.json
```

In the terminal, most of shell supports the `Control-m` by default.
So, disable the keybinding.

`sample/instruction/01_c-m.yml`

```diff
        - from:
            key: C-m
          to:
            - key: return_or_enter
+         conditions:
+           - app_unless:
+               identifiers:
+                 - ^com\\.apple\\.Terminal$
```

After save this rule file, your karabiner config will be automatically updated.

## Rule file

The rule file structure is same as Karabiner-Elements.

But there are some improvements.

### YAML format

### Commonize conditions between multiple manipulators

### Shorten frequently used configurations

#### Keybindings

#### Set variable

#### Condition based on variable

#### Condition based on frontmost application

### In `manipulator` configuration switch