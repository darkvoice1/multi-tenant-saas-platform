import { spawnSync } from 'node:child_process'

process.env.ANALYZE = '1'

const result = spawnSync('vite', ['build'], { stdio: 'inherit', shell: true })
process.exit(result.status ?? 1)

