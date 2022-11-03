'use strict';

const { randomUUID } = require('crypto');
const fs = require('fs').promises;
const path = require('path');

function isUUID(value) {
  return /^[0-9A-Fa-f]{8}(-[0-9A-Fa-f]{4}){4}[0-9A-Fa-f]{8}$/.test(value);
}

// uuid is assumed to pass isUUID()
function uuidToDocDirArray(uuid) {
  return [uuid.slice(0, 2), uuid.slice(2, 4), uuid.slice(4, 6), uuid.slice(6, 8), uuid];
}

class Mound {
  static MOUND_DATA_DIR = '/tmp/mound';

  static setup(MOUND_DATA_DIR) {
    if (typeof MOUND_DATA_DIR != 'string') throw new TypeError('MOUND_DATA_DIR must be a string.');
    if (!MOUND_DATA_DIR) throw new RangeError('MOUND_DATA_DIR may not be an empty string.');

    Mound.MOUND_DATA_DIR = MOUND_DATA_DIR;
  }

  constructor(program, version) {
    if (typeof program != 'string') throw new TypeError('program must be a string.');
    if (!program) throw RangeError('program may not be an empty string.');
    if (typeof version != 'string') throw new TypeError('version must be a string.');
    if (!version) throw RangeError('version may not be an empty string.');

    this.did = randomUUID();
    this.program = program;
    this.version = version;
    this.status = -1; // -1 for unfinished, 0 for success, nonzero for fail.
    this.blobs = [];
    this.sources = [];

    this._ready = this._open(); // A promise...
  }

  async ready() {
    if (await this._ready) {
      this._ready = true;
    }
  }

  async _open() {
    const docDir = path.join(Mound.MOUND_DATA_DIR, ...uuidToDocDirArray(this.did));
    await fs.mkdir(docDir, { recursive: true });
    await this._writeDoc();
  }

  async close(status) {
    if (typeof status != 'number') throw new TypeError('status must be a number.'); // TODO more rigor for integers. range checks
    this.status = status;
    await this._writeDoc();
  }

  async link(sourceDID) {
    if (typeof sourceDID != 'string') throw new TypeError('sourceDID must be a string.');
    if (!isUUID(sourceDID)) throw new RangeError('sourceDID must be a UUID.');

    if (this.sources.indexOf(sourceDID) < 0) {
      this.sources.push(sourceDID);
    }

    await this._writeDoc();
  }

  async _writeDoc() {
    const docPath = path.join(Mound.MOUND_DATA_DIR, ...uuidToDocDirArray(this.did), 'doc');
    await fs.writeFile(docPath, JSON.stringify({
      did: this.did,
      program: this.program,
      version: this.version,
      status: this.status,
      blobs: this.blobs,
      sources: this.sources,
    }) + '\n');
  }

  async blob(optionalName) {
    if (typeof optionalName != 'undefined' && typeof optionalName != 'string') throw new TypeError('optionalName must be unspecified or a string.');
    if (typeof optionalName == 'string' && !optionalName) throw new RangeError('optionalName may not be an empty string.');

    const blobNumber = this.blobs.length;
    this.blobs.push(optionalName || blobNumber);
    await this._writeDoc();
    //return new Blob(this, blobNumber);
    return blobNumber;
  }

  async write(blobNumber, data) {
    if (typeof blobNumber != 'number') throw new TypeError('blobNumber must be a number.'); // TODO more rigor for integers. range checks

    const blobPath = path.join(Mound.MOUND_DATA_DIR, ...uuidToDocDirArray(this.did), '' + blobNumber);
    await fs.appendFile(blobPath, data);
  }

  async println(blobNumber, data) {
    if (typeof data == 'string') {
      await this.write(blobNumber, data + '\n');
    } else if (data instanceof Uint8Array || data instanceof Buffer) {
      await this.write(blobNumber, Buffer.concat([data, Buffer.from('\n')])); // TODO test this
    } else {
      throw new TypeError('data must be a string, Uint8Array, or Buffer');
    }
  }
}

class Blob {
  constructor(mound, blobNumber) {
    this.mound = mound;
    this.blobNumber = blobNumber;
  }

  async write(data) {
    return await this.mound.write(this.blobNumber, data);
  }

  async println(data) {
    return await this.mound.println(this.blobNumber, data);
  }
}

module.exports = {
  Mound,
  Blob,
};
