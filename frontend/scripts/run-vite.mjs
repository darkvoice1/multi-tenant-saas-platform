import { spawn } from 'node:child_process'
import fs from 'node:fs'
import path from 'node:path'

const root = process.cwd()
const args = process.argv.slice(2)

function exists(p) {
  try {
    fs.accessSync(p)
    return true
  } catch {
    return false
  }
}

function ensureDir(p) {
  fs.mkdirSync(p, { recursive: true })
}

function tryRemoveZoneIdentifier(filePath) {
  try {
    fs.rmSync(`${filePath}:Zone.Identifier`, { force: true })
  } catch {
    // ignore
  }
}

function tryUnblockWithPowerShell(filePath) {
  const cmd = `Unblock-File -Path '${filePath.replace(/'/g, "''")}' -ErrorAction SilentlyContinue`
  spawn('powershell.exe', ['-NoProfile', '-ExecutionPolicy', 'Bypass', '-Command', cmd], {
    stdio: 'ignore',
    windowsHide: true
  })
}

function resolveEsbuildBinary() {
  const candidates = [
    path.join(root, 'node_modules', '@esbuild', 'win32-x64', 'esbuild.exe'),
    path.join(root, 'node_modules', 'esbuild', 'esbuild.exe')
  ]
  return candidates.find(exists) || ''
}

function prepareEsbuildBinary() {
  if (process.platform !== 'win32') return ''

  const src = resolveEsbuildBinary()
  if (!src) return ''

  const cacheDir = path.join(root, '.cache')
  ensureDir(cacheDir)
  const dst = path.join(cacheDir, 'esbuild.exe')

  try {
    fs.copyFileSync(src, dst)
  } catch {
    // ignore copy failure; we will still attempt to use src
  }

  tryRemoveZoneIdentifier(dst)
  tryUnblockWithPowerShell(dst)

  // Prefer the copied binary if it exists, fallback to original.
  if (exists(dst)) return dst
  tryRemoveZoneIdentifier(src)
  tryUnblockWithPowerShell(src)
  return src
}

const env = { ...process.env }
const esbuildBin = prepareEsbuildBinary()
if (esbuildBin) {
  env.ESBUILD_BINARY_PATH = esbuildBin
  env.ESBUILD_WORKER_THREADS = env.ESBUILD_WORKER_THREADS || '0'
}

const viteBin = path.join(root, 'node_modules', 'vite', 'bin', 'vite.js')
const child = spawn(process.execPath, [viteBin, ...args], {
  stdio: 'inherit',
  env
})

child.on('exit', (code) => {
  process.exit(code === null ? 0 : code)
})
