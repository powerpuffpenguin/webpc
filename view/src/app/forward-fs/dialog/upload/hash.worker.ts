/// <reference lib="webworker" />

import { buf } from "crc-32";

interface Request {
  file: File
  start: number
  end: number
}
interface Response {
  vals: Array<string>
  hash: string
}
function crc32ToHex(val: number): string {
  const hash = new Uint8Array(4)
  const dataView = new DataView(hash.buffer)
  dataView.setUint32(0, val, false)
  const strs = new Array<string>(4)
  for (let i = 0; i < 4; i++) {
    strs[i] = dataView.getUint8(i).toString(16).padStart(2, '0')
  }
  return strs.join('')
}
async function calculateHash(reqs: Array<Request>): Promise<Response> {
  const vals = new Array<string>(reqs.length)
  const hash = new Uint8Array(reqs.length * 4)
  const dataView = new DataView(hash.buffer)
  for (let i = 0; i < reqs.length; i++) {
    const req = reqs[i]
    const buffer = await req.file.slice(req.start, req.end).arrayBuffer()
    const val = buf(new Uint8Array(buffer), 0)
    vals[i] = crc32ToHex(val)
    dataView.setUint32(i * 4, val, false)
  }
  return {
    vals: vals,
    hash: crc32ToHex(buf(hash, 0)),
  }
}

addEventListener('message', ({ data }) => {
  calculateHash(data).then((resp) => {
    postMessage(resp)
  }).catch((e) => {
    console.warn(e)
    postMessage({
      error: e,
    })
  })
});
