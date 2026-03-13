<p align="center">
  <img alt="Oura CLI logo" src="https://em-content.zobj.net/source/apple/325/ring_1f48d.png" height="128px" />
  <p align="center">A terminal UI for your Oura Ring data</p>
</p>

<hr>

<p align="center">
<a href="https://github.com/hagelstam/oura-cli/releases/latest"><img src="https://img.shields.io/github/release/hagelstam/oura-cli.svg?style=for-the-badge" alt="Release"></a>
<a href="/LICENSE"><img src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge" alt="Software License"></a>
<a href="https://github.com/hagelstam/oura-cli/actions/workflows/build.yml"><img src="https://img.shields.io/github/actions/workflow/status/hagelstam/oura-cli/build.yml?style=for-the-badge" alt="Build status"></a>
</a>
<a href="https://goreportcard.com/report/github.com/hagelstam/oura-cli"><img src="https://goreportcard.com/badge/github.com/hagelstam/oura-cli?style=for-the-badge" alt="GoReportCard"></a>
</p>

## Install

```bash
go install github.com/hagelstam/oura-cli@latest
```

## Features

- **Sleep:** browse daily sleep scores and contributors

## Usage

Run `oura --help` for the full list of commands and flags.

> [!TIP]
> Generate a token at [cloud.ouraring.com/personal-access-tokens](https://cloud.ouraring.com/personal-access-tokens).

## Under the hood

- [cobra](https://github.com/spf13/cobra) for the CLI
- [bubbletea](https://github.com/charmbracelet/bubbletea) for the TUI
- [lipgloss](https://github.com/charmbracelet/lipgloss) for the styling
