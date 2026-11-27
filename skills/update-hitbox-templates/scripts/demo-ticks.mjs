// Prints the number of playback ticks of a CS2 demo, read from the file info message its header points to.
// Prints 0 when the demo has none, which happens with demos whose recording was cut short.
// Usage: node demo-ticks.mjs <demo.dem>
import fs from 'node:fs';

function readVarint(buffer, position) {
  let value = 0;
  let shift = 0;
  for (;;) {
    const byte = buffer[position++];
    value += (byte & 0x7f) * 2 ** shift;
    if ((byte & 0x80) === 0) {
      return [value, position];
    }
    shift += 7;
  }
}

function readPlaybackTicks(path) {
  const fd = fs.openSync(path, 'r');
  try {
    const header = Buffer.alloc(16);
    fs.readSync(fd, header, 0, 16, 0);
    const fileInfoOffset = header.readInt32LE(8);
    const size = fs.fstatSync(fd).size;
    if (header.toString('latin1', 0, 7) !== 'PBDEMS2' || fileInfoOffset <= 0 || fileInfoOffset >= size) {
      return 0;
    }

    // The message starts with its command, its tick and the size of its body. The tick is the last one of the demo,
    // which is enough here and is there even when the body is compressed.
    const buffer = Buffer.alloc(Math.min(64, size - fileInfoOffset));
    fs.readSync(fd, buffer, 0, buffer.length, fileInfoOffset);
    const [, tickPosition] = readVarint(buffer, 0);
    const [tick] = readVarint(buffer, tickPosition);

    return tick;
  } finally {
    fs.closeSync(fd);
  }
}

const path = process.argv[2];
if (path === undefined) {
  console.error('Usage: node demo-ticks.mjs <demo.dem>');
  process.exit(1);
}

console.log(readPlaybackTicks(path));
