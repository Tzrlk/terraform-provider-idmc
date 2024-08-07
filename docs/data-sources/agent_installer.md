---
# Had to mess with description rendering because it was failing with a weird error.
page_title: "idmc_agent_installer Data Source - idmc"
subcategory: ""
description: |-
  https://docs.informatica.com/integration-cloud/b2b-gateway/current-version/rest-api-reference/platform-rest-api-version-2-resources/agent.html
---

# idmc_agent_installer (Data Source)

https://docs.informatica.com/integration-cloud/b2b-gateway/current-version/rest-api-reference/platform-rest-api-version-2-resources/agent.html

## Example Usage

```terraform
# Linux
data "idmc_agent_installer" "linux" {
  platform = "linux"
}

# Windows
data "idmc_agent_installer" "windows" {
  platform = "windows"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `platform` (String) Platform of the Secure Agent machine. Must be one of the following values:
win64
linux64

### Read-Only

- `checksum_download_url` (String) The URL of the CRC-32 SHA256 package checksum.
- `download_url` (String) The URL of the latest Secure Agent installer package.
- `install_token` (String) Token needed to install and register a Secure Agent.
