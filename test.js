'use strict';

const { Mound } = require('./nodejs/mound.js');

Mound.setup('/home/keith/mound_data');

async function main() {
  const m = new Mound('mound', '0.0.1-test');
  await m.ready();
  const b0 = await m.blob();
  await m.println(b0, 'Hello, nodejs!');
  const b1 = await m.blob();
  await m.println(b1, 'Goodbye, nodejs!');
  await m.close(0);
}

main();
