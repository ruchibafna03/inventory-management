import { useEffect, useRef, useState } from 'react'
import {
  Tabs, Table, Button, Input, Space, Form, Typography, Row, Col,
  Card, Statistic, message, Popconfirm, Alert, Tag, Select, Modal,
  Descriptions, Timeline, Badge, Divider,
} from 'antd'
import {
  DatabaseOutlined, WarningOutlined, SwapOutlined, LockOutlined,
  UnlockOutlined, SearchOutlined, HistoryOutlined, BarcodeOutlined,
  KeyOutlined, CheckCircleOutlined, PrinterOutlined, ReloadOutlined,
} from '@ant-design/icons'
import {
  utilitiesApi, salesApi, issuesApi, receiptsApi,
  type TableCount, type OrphanDetail, type StockItem,
  type HistoryEvent, type CodeChange, type BlockedItem, type UserRecord,
  type Sale, type Issue, type Receipt,
} from '../api/client'

const { Title, Text } = Typography

// ─── 1. Organising Files ─────────────────────────────────────────────────────

function OrganisingFiles() {
  const [data, setData] = useState<TableCount[]>([])
  const [orphanIssues, setOrphanIssues] = useState<OrphanDetail[]>([])
  const [orphanReceipts, setOrphanReceipts] = useState<OrphanDetail[]>([])
  const [loading, setLoading] = useState(false)
  const [checkDone, setCheckDone] = useState(false)

  const loadSummary = () => {
    setLoading(true)
    utilitiesApi.summary()
      .then(r => setData(r.data || []))
      .finally(() => setLoading(false))
  }

  const runIntegrityCheck = () => {
    setLoading(true)
    Promise.all([utilitiesApi.orphanIssues(), utilitiesApi.orphanReceipts()])
      .then(([i, r]) => {
        setOrphanIssues(i.data || [])
        setOrphanReceipts(r.data || [])
        setCheckDone(true)
      })
      .finally(() => setLoading(false))
  }

  useEffect(() => { loadSummary() }, [])

  const orphanCols = [
    { title: 'ID', dataIndex: 'id', width: 70 },
    { title: 'Voucher', dataIndex: 't_no' },
    { title: 'Item Code', dataIndex: 'itcd', render: (v: string) => <Tag color="red">{v}</Tag> },
    { title: 'Description', dataIndex: 'itdesc' },
  ]

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Button icon={<ReloadOutlined />} onClick={loadSummary} loading={loading}>
          Refresh Counts
        </Button>
        <Button icon={<CheckCircleOutlined />} onClick={runIntegrityCheck} loading={loading} type="primary">
          Run Integrity Check
        </Button>
      </Space>

      <Row gutter={[12, 12]} style={{ marginBottom: 24 }}>
        {data.map(item => (
          <Col xs={12} sm={8} md={6} lg={4} key={item.table}>
            <Card size="small">
              <Statistic title={item.table} value={item.count} />
            </Card>
          </Col>
        ))}
      </Row>

      {checkDone && orphanIssues.length === 0 && orphanReceipts.length === 0 && (
        <Alert message="All records intact — no orphan references found." type="success" showIcon />
      )}
      {orphanIssues.length > 0 && (
        <>
          <Title level={5} style={{ color: '#ff4d4f' }}>
            <WarningOutlined /> Issue Lines with Missing Items ({orphanIssues.length})
          </Title>
          <Table rowKey="id" dataSource={orphanIssues} columns={orphanCols} size="small" pagination={false} style={{ marginBottom: 16 }} />
        </>
      )}
      {orphanReceipts.length > 0 && (
        <>
          <Title level={5} style={{ color: '#ff4d4f' }}>
            <WarningOutlined /> Receipt Lines with Missing Items ({orphanReceipts.length})
          </Title>
          <Table rowKey="id" dataSource={orphanReceipts} columns={orphanCols} size="small" pagination={false} />
        </>
      )}
    </div>
  )
}

// ─── 2. Password Changes ─────────────────────────────────────────────────────

function PasswordChanges() {
  const [form] = Form.useForm()
  const [users, setUsers] = useState<UserRecord[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    utilitiesApi.listUsers().then(r => setUsers(r.data || []))
  }, [])

  const handleSubmit = async () => {
    const values = await form.validateFields()
    if (values.new_password !== values.confirm_password) {
      message.error('Passwords do not match')
      return
    }
    setLoading(true)
    utilitiesApi.changePassword(values.username, values.new_password)
      .then(() => { message.success('Password changed successfully'); form.resetFields() })
      .catch(e => message.error(e.response?.data?.error ?? 'Failed to change password'))
      .finally(() => setLoading(false))
  }

  return (
    <Row justify="center">
      <Col xs={24} sm={18} md={12} lg={10}>
        <Card title={<><KeyOutlined /> Change Password</>}>
          <Form form={form} layout="vertical">
            <Form.Item name="username" label="Username" rules={[{ required: true }]}>
              <Select showSearch placeholder="Select user">
                {users.map(u => (
                  <Select.Option key={u.username} value={u.username}>
                    {u.username} {u.full_name ? `(${u.full_name})` : ''} — {u.role}
                  </Select.Option>
                ))}
              </Select>
            </Form.Item>
            <Form.Item name="new_password" label="New Password" rules={[{ required: true, min: 4 }]}>
              <Input.Password placeholder="Min 4 characters" />
            </Form.Item>
            <Form.Item name="confirm_password" label="Confirm Password" rules={[{ required: true }]}>
              <Input.Password placeholder="Repeat new password" />
            </Form.Item>
            <Button type="primary" block loading={loading} onClick={handleSubmit}>
              Change Password
            </Button>
          </Form>
        </Card>
      </Col>
    </Row>
  )
}

// ─── 3. Barcode Printing ─────────────────────────────────────────────────────

function BarcodePrinting() {
  const [search, setSearch] = useState('')
  const [items, setItems] = useState<StockItem[]>([])
  const [selected, setSelected] = useState<string[]>([])
  const [loading, setLoading] = useState(false)
  const printRef = useRef<HTMLDivElement>(null)

  const loadItems = () => {
    setLoading(true)
    utilitiesApi.stockPosition()
      .then(r => {
        const all = r.data || []
        setItems(search ? all.filter(i =>
          i.itcd.toLowerCase().includes(search.toLowerCase()) ||
          i.desc.toLowerCase().includes(search.toLowerCase())
        ) : all)
      })
      .finally(() => setLoading(false))
  }

  useEffect(() => { loadItems() }, [])

  const handlePrint = () => {
    const printItems = items.filter(i => selected.includes(i.itcd))
    if (!printItems.length) { message.warning('Select at least one item'); return }

    const html = `
      <html><head><style>
        body { font-family: monospace; margin: 0; }
        .label { border: 1px solid #000; display: inline-block; padding: 6px 10px; margin: 4px;
                 width: 160px; vertical-align: top; page-break-inside: avoid; }
        .code { font-size: 18px; font-weight: bold; letter-spacing: 3px; }
        .bars { font-size: 28px; letter-spacing: -1px; line-height: 1; }
        .desc { font-size: 9px; margin-top: 2px; }
        .wt   { font-size: 9px; }
        @media print { body { margin: 0; } }
      </style></head><body>
        ${printItems.map(i => `
          <div class="label">
            <div class="bars">|||||||||||||||</div>
            <div class="code">${i.itcd}</div>
            <div class="desc">${i.desc}</div>
            <div class="wt">Gr: ${i.gross_wt?.toFixed(3)}g${i.purity ? ` | ${i.purity}K` : ''}</div>
          </div>
        `).join('')}
      </body></html>`

    const win = window.open('', '_blank')
    if (win) { win.document.write(html); win.document.close(); win.print() }
  }

  const cols = [
    { title: 'Code', dataIndex: 'itcd', width: 90 },
    { title: 'Description', dataIndex: 'desc' },
    { title: 'Purity', dataIndex: 'purity', width: 60 },
    { title: 'Gross Wt', dataIndex: 'gross_wt', width: 90, render: (v: number) => v?.toFixed(3) },
    { title: 'Lot No', dataIndex: 'lot_no', width: 80 },
  ]

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Input.Search
          placeholder="Search item code or description"
          value={search}
          onChange={e => setSearch(e.target.value)}
          onSearch={loadItems}
          style={{ width: 280 }}
        />
        <Button icon={<ReloadOutlined />} onClick={loadItems} loading={loading}>Refresh</Button>
        <Button
          icon={<PrinterOutlined />}
          type="primary"
          onClick={handlePrint}
          disabled={!selected.length}
        >
          Print Labels ({selected.length})
        </Button>
        {selected.length > 0 && (
          <Button size="small" onClick={() => setSelected([])}>Clear Selection</Button>
        )}
      </Space>
      <div ref={printRef}>
        <Table
          rowKey="itcd"
          dataSource={items}
          columns={cols}
          size="small"
          loading={loading}
          pagination={{ pageSize: 50 }}
          rowSelection={{
            selectedRowKeys: selected,
            onChange: keys => setSelected(keys as string[]),
          }}
        />
      </div>
    </div>
  )
}

// ─── 4. Stock Checking ───────────────────────────────────────────────────────

const TAG_LABELS: Record<string, { color: string; label: string }> = {
  S: { color: 'green',  label: 'Sold' },
  I: { color: 'orange', label: 'Issued' },
  B: { color: 'red',    label: 'Blocked' },
  R: { color: 'blue',   label: 'Returned' },
}

function StockChecking() {
  const [items, setItems] = useState<StockItem[]>([])
  const [filtered, setFiltered] = useState<StockItem[]>([])
  const [loading, setLoading] = useState(false)
  const [search, setSearch] = useState('')
  const [tagFilter, setTagFilter] = useState<string>('')

  const load = () => {
    setLoading(true)
    utilitiesApi.stockPosition()
      .then(r => { setItems(r.data || []); setFiltered(r.data || []) })
      .finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  useEffect(() => {
    let out = items
    if (search) out = out.filter(i =>
      i.itcd.toLowerCase().includes(search.toLowerCase()) ||
      i.desc.toLowerCase().includes(search.toLowerCase())
    )
    if (tagFilter === '__in_stock__') out = out.filter(i => !i.tag)
    else if (tagFilter) out = out.filter(i => i.tag === tagFilter)
    setFiltered(out)
  }, [search, tagFilter, items])

  // Summary counts
  const counts = {
    total:    items.length,
    inStock:  items.filter(i => !i.tag).length,
    sold:     items.filter(i => i.tag === 'S').length,
    issued:   items.filter(i => i.tag === 'I').length,
  }

  const cols = [
    { title: 'Code', dataIndex: 'itcd', width: 90 },
    { title: 'Description', dataIndex: 'desc' },
    { title: 'Purity', dataIndex: 'purity', width: 60 },
    { title: 'Gross Wt', dataIndex: 'gross_wt', width: 90, render: (v: number) => v?.toFixed(3) },
    { title: 'Gold Wt', dataIndex: 'gold_wt', width: 80, render: (v: number) => v?.toFixed(3) },
    { title: 'Lot No', dataIndex: 'lot_no', width: 80 },
    { title: 'Karigar', dataIndex: 'karigar_name', width: 120 },
    {
      title: 'Status', dataIndex: 'tag', width: 90,
      render: (v: string) => {
        const t = TAG_LABELS[v]
        return t
          ? <Tag color={t.color}>{t.label}</Tag>
          : <Tag color="default">In Stock</Tag>
      },
    },
  ]

  return (
    <div>
      <Row gutter={12} style={{ marginBottom: 16 }}>
        {[
          { label: 'Total Items', value: counts.total, color: '#1890ff' },
          { label: 'In Stock', value: counts.inStock, color: '#52c41a' },
          { label: 'Issued Out', value: counts.issued, color: '#fa8c16' },
          { label: 'Sold', value: counts.sold, color: '#389e0d' },
        ].map(s => (
          <Col key={s.label} xs={12} sm={6}>
            <Card size="small">
              <Statistic title={s.label} value={s.value} valueStyle={{ color: s.color }} />
            </Card>
          </Col>
        ))}
      </Row>
      <Space style={{ marginBottom: 12 }}>
        <Input.Search
          placeholder="Item code or description"
          value={search}
          onChange={e => setSearch(e.target.value)}
          style={{ width: 240 }}
        />
        <Select
          placeholder="Filter by status"
          allowClear
          style={{ width: 150 }}
          value={tagFilter || undefined}
          onChange={v => setTagFilter(v ?? '')}
        >
          <Select.Option value="__in_stock__">In Stock</Select.Option>
          <Select.Option value="I">Issued Out</Select.Option>
          <Select.Option value="S">Sold</Select.Option>
          <Select.Option value="B">Blocked</Select.Option>
          <Select.Option value="R">Returned</Select.Option>
        </Select>
        <Button icon={<ReloadOutlined />} onClick={load} loading={loading}>Refresh</Button>
      </Space>
      <Table
        rowKey="itcd"
        dataSource={filtered}
        columns={cols}
        size="small"
        loading={loading}
        pagination={{ pageSize: 50 }}
      />
    </div>
  )
}

// ─── 5 & 6. Sale Slip Reprint / Sale Return Reprint ──────────────────────────

function SaleReprint({ title, isReturn }: { title: string; isReturn: boolean }) {
  const [vouchNo, setVouchNo] = useState('')
  const [sale, setSale] = useState<Sale | null>(null)
  const [loading, setLoading] = useState(false)

  const fetch = () => {
    if (!vouchNo.trim()) return
    setLoading(true)
    salesApi.get(vouchNo.trim())
      .then(r => setSale(r.data))
      .catch(() => { message.error('Voucher not found'); setSale(null) })
      .finally(() => setLoading(false))
  }

  const handlePrint = () => {
    if (!sale) return
    const html = `
      <html><head><style>
        body { font-family: Arial, sans-serif; margin: 20px; font-size: 12px; }
        h2 { text-align: center; margin: 0; }
        .hdr { text-align: center; margin-bottom: 12px; }
        table { width: 100%; border-collapse: collapse; margin-top: 10px; }
        td, th { border: 1px solid #ccc; padding: 4px 8px; }
        th { background: #f0f0f0; }
        .row { display: flex; justify-content: space-between; margin: 2px 0; }
        @media print { body { margin: 10px; } }
      </style></head><body>
        <div class="hdr">
          <h2>VAL Inventory</h2>
          <div>${isReturn ? 'SALES RETURN SLIP' : 'SALES SLIP'}</div>
        </div>
        <div class="row"><span>Voucher No: <b>${sale.vouch_no}</b></span><span>Date: ${sale.vouch_date?.substring(0,10)}</span></div>
        <div class="row"><span>Customer: ${sale.name ?? sale.ac_code ?? '-'}</span><span>Phone: ${sale.phone_no ?? '-'}</span></div>
        <table>
          <tr><th>Particulars</th><th>Value</th></tr>
          <tr><td>Nos</td><td>${sale.nos}</td></tr>
          <tr><td>Gross Weight</td><td>${sale.gross_wt?.toFixed(3)} g</td></tr>
          <tr><td>Net Weight</td><td>${sale.net_wt?.toFixed(3)} g</td></tr>
          <tr><td>Net Pure</td><td>${sale.net_pure?.toFixed(3)}</td></tr>
          <tr><td>Discount</td><td>${sale.discount?.toFixed(3)}</td></tr>
          <tr><td><b>Pure Final</b></td><td><b>${sale.pure_final?.toFixed(3)}</b></td></tr>
        </table>
        ${sale.narr ? `<p>Narr: ${sale.narr}</p>` : ''}
      </body></html>`
    const win = window.open('', '_blank')
    if (win) { win.document.write(html); win.document.close(); win.print() }
  }

  return (
    <div>
      <Title level={5}>{title}</Title>
      <Space style={{ marginBottom: 16 }}>
        <Input
          placeholder="Enter voucher number"
          value={vouchNo}
          onChange={e => setVouchNo(e.target.value)}
          onPressEnter={fetch}
          style={{ width: 200 }}
        />
        <Button icon={<SearchOutlined />} onClick={fetch} loading={loading}>Fetch</Button>
        {sale && <Button icon={<PrinterOutlined />} type="primary" onClick={handlePrint}>Print</Button>}
      </Space>
      {sale && (
        <Card size="small" style={{ maxWidth: 500 }}>
          <Descriptions column={2} size="small" bordered>
            <Descriptions.Item label="Voucher No" span={2}>{sale.vouch_no}</Descriptions.Item>
            <Descriptions.Item label="Date">{sale.vouch_date?.substring(0, 10)}</Descriptions.Item>
            <Descriptions.Item label="Nos">{sale.nos}</Descriptions.Item>
            <Descriptions.Item label="Customer" span={2}>{sale.name ?? sale.ac_code ?? '-'}</Descriptions.Item>
            <Descriptions.Item label="Gross Wt">{sale.gross_wt?.toFixed(3)}</Descriptions.Item>
            <Descriptions.Item label="Net Wt">{sale.net_wt?.toFixed(3)}</Descriptions.Item>
            <Descriptions.Item label="Net Pure">{sale.net_pure?.toFixed(3)}</Descriptions.Item>
            <Descriptions.Item label="Discount">{sale.discount?.toFixed(3)}</Descriptions.Item>
            <Descriptions.Item label="Pure Final" span={2}>
              <Text strong>{sale.pure_final?.toFixed(3)}</Text>
            </Descriptions.Item>
          </Descriptions>
        </Card>
      )}
    </div>
  )
}

// ─── 7. Karigar Issue Reprint ────────────────────────────────────────────────

function IssueReprint() {
  const [tNo, setTNo] = useState('')
  const [issue, setIssue] = useState<Issue | null>(null)
  const [loading, setLoading] = useState(false)

  const fetch = () => {
    if (!tNo.trim()) return
    setLoading(true)
    issuesApi.get(tNo.trim())
      .then(r => setIssue(r.data))
      .catch(() => { message.error('Issue voucher not found'); setIssue(null) })
      .finally(() => setLoading(false))
  }

  const handlePrint = () => {
    if (!issue) return
    const html = `
      <html><head><style>
        body { font-family: Arial, sans-serif; margin: 20px; font-size: 12px; }
        h2 { text-align: center; }
        table { width: 100%; border-collapse: collapse; margin-top: 10px; }
        td, th { border: 1px solid #ccc; padding: 4px 8px; }
        th { background: #f0f0f0; }
        .row { display: flex; justify-content: space-between; margin: 2px 0; }
        @media print { body { margin: 10px; } }
      </style></head><body>
        <h2>VAL Inventory — KARIGAR ISSUE</h2>
        <div class="row"><span>Issue No: <b>${issue.t_no}</b></span><span>Date: ${issue.t_date?.substring(0,10)}</span></div>
        <div class="row"><span>Karigar: ${issue.party_name ?? issue.ac_code ?? '-'}</span><span>Due: ${issue.due_date?.substring(0,10) ?? '-'}</span></div>
        <table>
          <tr><th>Particulars</th><th>Value</th></tr>
          <tr><td>Nos</td><td>${issue.nos}</td></tr>
          <tr><td>Gross Weight</td><td>${issue.gross_wt?.toFixed(3)} g</td></tr>
          <tr><td>Net Weight</td><td>${issue.net_wt?.toFixed(3)} g</td></tr>
          <tr><td>Net Pure</td><td>${issue.net_pure?.toFixed(3)}</td></tr>
          <tr><td><b>Pure Final</b></td><td><b>${issue.pure_final?.toFixed(3)}</b></td></tr>
        </table>
        ${(issue.details ?? []).length > 0 ? `
          <h4>Items</h4>
          <table>
            <tr><th>#</th><th>Code</th><th>Description</th><th>Gross Wt</th><th>Net Wt</th></tr>
            ${(issue.details ?? []).map((d, i) => `
              <tr><td>${i+1}</td><td>${d.itcd??''}</td><td>${d.itdesc??''}</td>
              <td>${d.gross_wt?.toFixed(3)}</td><td>${d.net_wt?.toFixed(3)}</td></tr>
            `).join('')}
          </table>` : ''}
      </body></html>`
    const win = window.open('', '_blank')
    if (win) { win.document.write(html); win.document.close(); win.print() }
  }

  return (
    <div>
      <Title level={5}>Karigar Issue Reprint</Title>
      <Space style={{ marginBottom: 16 }}>
        <Input
          placeholder="Enter issue number"
          value={tNo}
          onChange={e => setTNo(e.target.value)}
          onPressEnter={fetch}
          style={{ width: 200 }}
        />
        <Button icon={<SearchOutlined />} onClick={fetch} loading={loading}>Fetch</Button>
        {issue && <Button icon={<PrinterOutlined />} type="primary" onClick={handlePrint}>Print</Button>}
      </Space>
      {issue && (
        <Card size="small" style={{ maxWidth: 500 }}>
          <Descriptions column={2} size="small" bordered>
            <Descriptions.Item label="Issue No" span={2}>{issue.t_no}</Descriptions.Item>
            <Descriptions.Item label="Date">{issue.t_date?.substring(0, 10)}</Descriptions.Item>
            <Descriptions.Item label="Due Date">{issue.due_date?.substring(0, 10) ?? '-'}</Descriptions.Item>
            <Descriptions.Item label="Karigar" span={2}>{issue.party_name ?? issue.ac_code ?? '-'}</Descriptions.Item>
            <Descriptions.Item label="Gross Wt">{issue.gross_wt?.toFixed(3)}</Descriptions.Item>
            <Descriptions.Item label="Net Wt">{issue.net_wt?.toFixed(3)}</Descriptions.Item>
            <Descriptions.Item label="Net Pure">{issue.net_pure?.toFixed(3)}</Descriptions.Item>
            <Descriptions.Item label="Pure Final"><Text strong>{issue.pure_final?.toFixed(3)}</Text></Descriptions.Item>
          </Descriptions>
          {(issue.details ?? []).length > 0 && (
            <Table
              rowKey="id"
              dataSource={issue.details}
              size="small"
              style={{ marginTop: 8 }}
              pagination={false}
              columns={[
                { title: 'Code', dataIndex: 'itcd', width: 80 },
                { title: 'Desc', dataIndex: 'itdesc' },
                { title: 'Gross Wt', dataIndex: 'gross_wt', render: (v: number) => v?.toFixed(3) },
                { title: 'Net Wt', dataIndex: 'net_wt', render: (v: number) => v?.toFixed(3) },
              ]}
            />
          )}
        </Card>
      )}
    </div>
  )
}

// ─── 8. History of Item ──────────────────────────────────────────────────────

function HistoryOfItem() {
  const [itcd, setItcd] = useState('')
  const [events, setEvents] = useState<HistoryEvent[]>([])
  const [loading, setLoading] = useState(false)
  const [searched, setSearched] = useState(false)

  const fetch = () => {
    if (!itcd.trim()) return
    setLoading(true)
    utilitiesApi.itemHistory(itcd.trim().toUpperCase())
      .then(r => { setEvents(r.data || []); setSearched(true) })
      .catch(() => { message.error('Error fetching history'); setEvents([]) })
      .finally(() => setLoading(false))
  }

  const eventColor: Record<string, string> = {
    'Lot Received':   'green',
    'Issued':         'orange',
    'Received Back':  'blue',
  }

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Input
          placeholder="Enter item code"
          value={itcd}
          onChange={e => setItcd(e.target.value)}
          onPressEnter={fetch}
          style={{ width: 180, textTransform: 'uppercase' }}
        />
        <Button icon={<HistoryOutlined />} type="primary" onClick={fetch} loading={loading}>
          Show History
        </Button>
      </Space>

      {searched && events.length === 0 && (
        <Alert message={`No history found for item "${itcd}"`} type="info" showIcon />
      )}

      {events.length > 0 && (
        <>
          <Title level={5}>History for Item: <Tag color="blue">{itcd.toUpperCase()}</Tag></Title>
          <Timeline
            items={events.map(e => ({
              color: eventColor[e.event_type] ?? 'gray',
              children: (
                <div>
                  <Space>
                    <Badge color={eventColor[e.event_type] ?? 'gray'} text={<Text strong>{e.event_type}</Text>} />
                    <Text type="secondary">{e.event_date?.substring(0, 10)}</Text>
                    <Tag>{e.vouch_no}</Tag>
                  </Space>
                  <div style={{ marginTop: 4, paddingLeft: 16, fontSize: 12, color: '#666' }}>
                    {e.ac_code && <span>A/C: {e.ac_code}{e.party ? ` (${e.party})` : ''} | </span>}
                    Gross: {e.gross_wt?.toFixed(3)}g | Net: {e.net_wt?.toFixed(3)}g
                    {e.narr && <span> | {e.narr}</span>}
                  </div>
                </div>
              ),
            }))}
          />
        </>
      )}
    </div>
  )
}

// ─── 9. Change Item Code ─────────────────────────────────────────────────────

function ChangeItemCode() {
  const [form] = Form.useForm()
  const [history, setHistory] = useState<CodeChange[]>([])
  const [loading, setLoading] = useState(false)
  const [histLoading, setHistLoading] = useState(false)

  const loadHistory = () => {
    setHistLoading(true)
    utilitiesApi.codeChangeHistory()
      .then(r => setHistory(r.data || []))
      .finally(() => setHistLoading(false))
  }

  useEffect(() => { loadHistory() }, [])

  const handleSubmit = async () => {
    const values = await form.validateFields()
    setLoading(true)
    utilitiesApi.changeItemCode(values.from.trim().toUpperCase(), values.to.trim().toUpperCase())
      .then(() => { message.success(`Code changed: ${values.from} → ${values.to}`); form.resetFields(); loadHistory() })
      .catch(e => message.error(e.response?.data?.error ?? 'Failed to change code'))
      .finally(() => setLoading(false))
  }

  const histCols = [
    { title: '#', dataIndex: 'id', width: 60 },
    { title: 'From', dataIndex: 'itcdf', render: (v: string) => <Tag color="orange">{v}</Tag> },
    { title: 'To', dataIndex: 'itcdt', render: (v: string) => <Tag color="green">{v}</Tag> },
    { title: 'Date', dataIndex: 'date', render: (v: string) => v?.substring(0, 10) },
  ]

  return (
    <Row gutter={24}>
      <Col xs={24} md={10}>
        <Card title={<><SwapOutlined /> Change Item Code</>} size="small">
          <Alert
            type="warning"
            message="Updates item code across all tables (issues, receipts, blocked items) in a single transaction."
            style={{ marginBottom: 16 }}
            showIcon
          />
          <Form form={form} layout="vertical">
            <Form.Item name="from" label="Current Code" rules={[{ required: true }]}>
              <Input placeholder="e.g. A001" style={{ textTransform: 'uppercase' }} />
            </Form.Item>
            <Form.Item name="to" label="New Code" rules={[{ required: true }]}>
              <Input placeholder="e.g. A002" style={{ textTransform: 'uppercase' }} />
            </Form.Item>
            <Popconfirm
              title="Change Item Code"
              description="This updates all references. Continue?"
              onConfirm={handleSubmit}
              okText="Yes, Change"
              okButtonProps={{ danger: true }}
            >
              <Button type="primary" loading={loading} block>Change Code</Button>
            </Popconfirm>
          </Form>
        </Card>
      </Col>
      <Col xs={24} md={14}>
        <Title level={5}>Change History</Title>
        <Table rowKey="id" dataSource={history} columns={histCols} size="small" loading={histLoading} pagination={{ pageSize: 15 }} />
      </Col>
    </Row>
  )
}

// ─── Main Page ───────────────────────────────────────────────────────────────

const TABS = [
  { key: 'organising',  label: <><DatabaseOutlined /> Organising Files</>,     children: <OrganisingFiles /> },
  { key: 'password',    label: <><KeyOutlined /> Password Changes</>,           children: <PasswordChanges /> },
  { key: 'barcode',     label: <><BarcodeOutlined /> Barcode Printing</>,       children: <BarcodePrinting /> },
  { key: 'stock',       label: <><CheckCircleOutlined /> Stock Checking</>,     children: <StockChecking /> },
  { key: 'sale-reprint',label: <><PrinterOutlined /> Sale Slip Reprint</>,      children: <SaleReprint title="Sale Slip Reprint" isReturn={false} /> },
  { key: 'return-reprint', label: <><PrinterOutlined /> Sale Return Reprint</>, children: <SaleReprint title="Sale Return Reprint" isReturn={true} /> },
  { key: 'issue-reprint',label: <><PrinterOutlined /> Karigar Issue Reprint</>, children: <IssueReprint /> },
  { key: 'history',     label: <><HistoryOutlined /> History of Item</>,        children: <HistoryOfItem /> },
  { key: 'codechange',  label: <><SwapOutlined /> Change Item Code</>,          children: <ChangeItemCode /> },
]

export default function Utilities() {
  return (
    <div>
      <Title level={4} style={{ marginBottom: 16 }}>Utilities</Title>
      <Tabs items={TABS} tabPosition="left" style={{ minHeight: 500 }} />
    </div>
  )
}
