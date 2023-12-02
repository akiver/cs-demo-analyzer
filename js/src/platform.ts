import path from 'node:path';

function getBinarySubpath() {
  const supportedPlatforms: Record<string, string> = {
    'darwin-x64': 'bin/darwin-x64/csda',
    'darwin-arm64': 'bin/darwin-arm64/csda',
    'linux-x64': 'bin/linux-x64/csda',
    'linux-arm64': 'bin/linux-arm64/csda',
    'win32-x64': 'bin/windows-x64/csda.exe',
  };

  const platformKey = `${process.platform}-${process.arch}`;
  if (!supportedPlatforms[platformKey]) {
    throw new Error(`Unsupported platform: ${platformKey}`);
  }

  return supportedPlatforms[platformKey];
}

export function getBinaryPath() {
  return path.join(__dirname, getBinarySubpath());
}
