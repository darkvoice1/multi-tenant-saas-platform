const { test, expect } = require('@playwright/test')

test('bootstrap and login to dashboard', async ({ page }) => {
  await page.goto('/bootstrap')

  await page.getByRole('button', { name: '一键初始化' }).click()

  const successHeader = page.getByRole('heading', { name: '初始化成功' })
  const errorText = page.locator('.error')

  const result = await Promise.race([
    successHeader.waitFor({ state: 'visible', timeout: 15000 }).then(() => 'success'),
    errorText.waitFor({ state: 'visible', timeout: 15000 }).then(() => 'error')
  ])

  let tenantId = ''
  let email = ''
  let password = 'Admin123!'

  if (result === 'success') {
    const tenantRow = page.getByText('租户 ID：')
    await expect(tenantRow).toBeVisible()
    tenantId = await tenantRow.locator('strong').innerText()
    expect(tenantId).toMatch(/[0-9a-f-]{36}/)

    const emailRow = page.getByText('管理员邮箱：')
    if (await emailRow.isVisible()) {
      email = await emailRow.locator('strong').innerText()
    }

    await page.getByRole('button', { name: '去登录' }).click()
  } else {
    const err = (await errorText.innerText()).trim()
    const fallbackTenant = process.env.E2E_TENANT_ID
    const fallbackEmail = process.env.E2E_EMAIL || 'admin@example.com'
    const fallbackPassword = process.env.E2E_PASSWORD || 'Admin123!'

    if (!/already initialized/i.test(err)) {
      throw new Error(`bootstrap failed: ${err}`)
    }
    if (!fallbackTenant) {
      throw new Error('E2E_TENANT_ID is required when bootstrap already initialized')
    }

    tenantId = fallbackTenant
    email = fallbackEmail
    password = fallbackPassword

    const loginUrl = `/login?tenant=${encodeURIComponent(tenantId)}&email=${encodeURIComponent(email)}`
    await page.goto(loginUrl)
  }

  await page.getByLabel('密码').fill(password)
  await page.getByRole('button', { name: '登录' }).click()

  await expect(page.getByRole('heading', { name: '业务看板' })).toBeVisible()
})
