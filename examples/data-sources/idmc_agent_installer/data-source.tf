
# Linux
data "idmc_agent_installer" "linux" {
  platform = "linux64"
}

# Windows
data "idmc_agent_installer" "windows" {
  platform = "win64"
}
