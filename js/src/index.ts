import { exec } from 'node:child_process';
import fs from 'node:fs/promises';
import { getBinaryPath } from './platform';
import { DemoSource, ExportFormat } from './constants';

export type Options = {
  demoPath: string;
  outputFolderPath: string;
  format: ExportFormat;
  source?: DemoSource;
  analyzePositions?: boolean;
  minify?: boolean; // JSON only
  onStart?: (command: string) => void;
  onStdout?: (data: string) => void;
  onStderr?: (data: string) => void;
  onEnd?: (exitCode: number) => void;
  executablePath?: string;
};

export async function analyzeDemo({
  demoPath,
  outputFolderPath,
  format,
  source,
  analyzePositions,
  minify,
  onStart,
  onStdout,
  onStderr,
  onEnd,
  executablePath,
}: Options): Promise<void> {
  await fs.mkdir(outputFolderPath, { recursive: true });

  return new Promise<void>((resolve, reject) => {
    const binPath = executablePath ?? getBinaryPath();
    const args: string[] = [
      `"${binPath}"`,
      `-demo-path="${demoPath}"`,
      `-output="${outputFolderPath}"`,
      `-format="${format}"`,
    ];
    if (source) {
      args.push(`-source="${source}"`);
    }
    if (analyzePositions) {
      args.push(`-positions="${analyzePositions}"`);
    }
    if (minify) {
      args.push('-minify');
    }
    const command = args.join(' ');
    if (onStart) {
      onStart(command);
    }

    const child = exec(command, { windowsHide: true, maxBuffer: undefined });
    if (onStdout) {
      child.stdout?.on('data', (data: string) => {
        onStdout(data);
      });
    }

    if (onStderr) {
      child.stderr?.on('data', (data: string) => {
        onStderr(data);
      });
    }

    child.on('exit', (code: number) => {
      if (onEnd) {
        onEnd(code);
      }
      if (code === 0) {
        resolve();
      } else {
        reject();
      }
    });
  });
}

export * from './constants';
