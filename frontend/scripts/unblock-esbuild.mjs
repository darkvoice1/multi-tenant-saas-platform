import { spawnSync } from 'node:child_process'
import fs from 'node:fs'
import path from 'node:path'

function exists(p) {
  try {
    fs.accessSync(p)
    return true
  } catch {
    return false
  }
}

function tryRemoveZoneIdentifier(filePath) {
  // Best-effort: remove the alternate data stream that marks a file as downloaded from the internet.
  // This is what commonly causes `spawn EPERM` on Windows for esbuild.exe.
  try {
    fs.rmSync(`${filePath}:Zone.Identifier`, { force: true })
  } catch {
    // ignore
  }
}

function tryUnblockWithPowerShell(filePath) {
  const cmd = `Unblock-File -Path '${filePath.replace(/'/g, "''")}' -ErrorAction SilentlyContinue`
  spawnSync('powershell.exe', ['-NoProfile', '-ExecutionPolicy', 'Bypass', '-Command', cmd], {
    stdio: 'ignore',
    windowsHide: true
  })
}

if (process.platform === 'win32') {
  const root = process.cwd()
  const candidates = [
    path.join(root, 'node_modules', '@esbuild', 'win32-x64', 'esbuild.exe'),
    path.join(root, 'node_modules', 'esbuild', 'esbuild.exe')
  ]

  for (const p of candidates) {
    if (!exists(p)) continue
    tryRemoveZoneIdentifier(p)
    tryUnblockWithPowerShell(p)
  }
}

