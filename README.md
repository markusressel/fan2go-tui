<h1 align="center">fan2go-tui</h1>
<h4 align="center">Terminal UI for fan2go.</h4>

<div align="center">

[![Programming Language](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white)]()
[![Latest Release](https://img.shields.io/github/release/markusressel/fan2go-tui.svg)](https://github.com/markusressel/fan2go-tui/releases)
[![License](https://img.shields.io/badge/license-AGPLv3-blue.svg)](/LICENSE)

[![asciicast](https://asciinema.org/a/612087.svg)](https://asciinema.org/a/612087)

</div>

# Features

* [x] Visualize detected and configured fans
* [x] Monitor curve values
* [x] Inspect sensor readings

# How to use

## Installation

### Arch Linux ![](https://img.shields.io/badge/Arch_Linux-1793D1?logo=arch-linux&logoColor=white)

```shell
yay -S fan2go-tui-git
```

<details>
<summary>Community Maintained Packages</summary>

None yet

</details>

### Manual

Compile yourself:

```shell
git clone https://github.com/markusressel/fan2go-tui.git
cd fan2go-tui
make build
sudo cp ./bin/fan2go-tui /usr/bin/fan2go-tui
sudo chmod ug+x /usr/bin/fan2go-tui
```

## Configuration

> **Note:**
> The configuration is optional and sane defaults will be used if omitted.

Then configure fan2go-tui by creating a YAML configuration file in **one** of the following locations:

* `/etc/fan2go-tui/fan2go-tui.yaml` (recommended)
* `/home/<user>/.config/fan2go-tui/fan2go-tui.yaml`
* `./fan2go-tui.yaml`

```shell
sudo mkdir -P ~/.config/fan2go-tui
sudo nano ~/.config/fan2go-tui/fan2go-tui.yaml
```

### Example

An example configuration file including more detailed documentation can be found
in [fan2go-tui.yaml](/fan2go-tui.yaml).

# Dependencies

See [go.mod](go.mod)

# License

```
fan2go-tui
Copyright (C) 2023  Markus Ressel

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
```
