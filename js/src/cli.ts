#!/usr/bin/env node
import { getBinaryPath } from './platform';

try {
  const binPath = getBinaryPath();
  require('child_process').execFileSync(binPath, process.argv.slice(2), { stdio: 'inherit' });
} catch (error) {
  if (error && typeof error === 'object' && 'status' in error && typeof error.status === 'number') {
    process.exit(error.status);
  }

  throw error;
}
