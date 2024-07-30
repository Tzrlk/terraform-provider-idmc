run "data" {

  ## Linux #####################################################################

  assert {
    error_message = "Linux platform query should be configured correctly."
    condition     = data.idmc_agent_installer.linux.platform == "linux64"
  }

  assert {
    error_message = "Linux download url should not be null or empty."
    condition     = (data.idmc_agent_installer.linux.download_url != null &&
    trimspace(data.idmc_agent_installer.linux.download_url) != "")
  }

  assert {
    error_message = "Linux checksum download url should not be null or empty."
    condition     = (data.idmc_agent_installer.linux.checksum_download_url != null &&
    trimspace(data.idmc_agent_installer.linux.checksum_download_url) != "")
  }

  assert {
    error_message = "Linux install token should not be null or empty."
    condition     = (data.idmc_agent_installer.linux.install_token != null &&
    trimspace(data.idmc_agent_installer.linux.install_token) != "")
  }

  ## Windows ###################################################################

  assert {
    error_message = "Windows platform query should be configured correctly."
    condition     = data.idmc_agent_installer.windows.platform == "win64"
  }

  assert {
    error_message = "Windows download url should not be null or empty."
    condition     = (data.idmc_agent_installer.windows.download_url != null &&
    trimspace(data.idmc_agent_installer.windows.download_url) != "")
  }

  assert {
    error_message = "Windows checksum download url should not be null or empty."
    condition     = (data.idmc_agent_installer.windows.checksum_download_url != null &&
    trimspace(data.idmc_agent_installer.windows.checksum_download_url) != "")
  }

  assert {
    error_message = "Windows install token should not be null or empty."
    condition     = (data.idmc_agent_installer.windows.install_token != null &&
    trimspace(data.idmc_agent_installer.windows.install_token) != "")
  }

}
