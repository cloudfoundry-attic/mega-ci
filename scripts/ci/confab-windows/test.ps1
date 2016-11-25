trap {
  write-error $_
  exit 1
}

$env:GOPATH = Join-Path -Path $PWD "gopath"
mkdir -path $env:GOPATH
$env:PATH = $env:GOPATH + "/bin;C:/go/bin;" + $env:PATH

# Move vendor directory into GOPATH to eliminate long file problems
mv ./consul-release/src/confab/vendor $env:GOPATH/src
mkdir -path $env:GOPATH/src/github.com/cloudfoundry-incubator
mv ./consul-release $env:GOPATH/src/github.com/cloudfoundry-incubator/consul-release
cd $env:GOPATH/src/github.com/cloudfoundry-incubator/consul-release

if ((Get-Command "go.exe" -ErrorAction SilentlyContinue) -eq $null)
{
  Write-Host "Installing Go 1.6.3!"
  Invoke-WebRequest https://storage.googleapis.com/golang/go1.6.3.windows-amd64.msi -OutFile go.msi

  $p = Start-Process -FilePath "msiexec" -ArgumentList "/passive /norestart /i go.msi" -Wait -PassThru

  if($p.ExitCode -ne 0)
  {
    throw "Golang MSI installation process returned error code: $($p.ExitCode)"
  }
  Write-Host "Go is installed!"
}

go.exe install github.com/onsi/ginkgo/ginkgo
if ($LastExitCode -ne 0)
{
    Write-Error $_
    exit 1
}

ginkgo.exe -r src/confab -race -randomizeAllSpecs -skipPackage vendor
if ($LastExitCode -ne 0)
{
    Write-Error $_
    exit 1
}

Exit 0

