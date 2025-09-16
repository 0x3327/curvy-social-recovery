#  Wasm module

## Compilation
Compile using: `GOOS=js GOARCH=wasm go build -o main.wasm .`

## Using in browser

`wasm_exec.js` file can be found in:
- Go 1.23 and earlier: `$(go env GOROOT)/misc/wasm/wasm_exec.js`
- Otherwise: `$(go env GOROOT)/lib/wasm/wasm_exec.js`

For details, read the [Go Wiki](https://go.dev/wiki/WebAssembly).

### Functions provided

- `goSplit(t: number, n: number, sk: string, vk: string): Promise<Array<Share>>`
- `goRecover(t: number, shares: Array<Share>): Promise<Key>`
- `goGenRandomKeys(): Promise<Key>`

where

```ts
interface Share {
    x: string
    vkEval: string
    skEval: string
}

interface Key {
    k: string,
    v: string
}
```