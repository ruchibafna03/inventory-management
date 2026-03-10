import { useEffect, useState } from 'react'
import {
  Table, Button, Space, Modal, Form, InputNumber, DatePicker,
  Input, message, Popconfirm, Tag, Select, Typography, Row, Col,
  Divider,
} from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, PlusCircleOutlined, DeleteFilled } from '@ant-design/icons'
import { issuesApi, accountsApi, type Issue, type IssueDetail, type Account } from '../api/client'
import dayjs from 'dayjs'

const { Title, Text } = Typography

export default function Issues() {
  const [issues, setIssues] = useState<Issue[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Issue | null>(null)
  const [details, setDetails] = useState<IssueDetail[]>([])
  const [accounts, setAccounts] = useState<Account[]>([])
  const [form] = Form.useForm()

  const load = (p = page) => {
    setLoading(true)
    issuesApi.list({ page: p, per_page: 50 })
      .then(r => { setIssues(r.data.data || []); setTotal(r.data.total) })
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    load()
    accountsApi.list({ per_page: 500 }).then(r => setAccounts(r.data.data || []))
  }, [])

  const openCreate = () => {
    setEditing(null)
    form.resetFields()
    form.setFieldsValue({ t_date: dayjs() })
    setDetails([{ nos: 1, gross_wt: 0, net_wt: 0, stone_wt: 0, kundan_wt: 0, nakash_wt: 0 }])
    setModalOpen(true)
  }

  const openEdit = async (issue: Issue) => {
    setEditing(issue)
    const full = await issuesApi.get(issue.t_no)
    form.setFieldsValue({
      ...full.data,
      t_date: dayjs(full.data.t_date),
      due_date: full.data.due_date ? dayjs(full.data.due_date) : null,
    })
    setDetails(full.data.details ?? [])
    setModalOpen(true)
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      const payload = {
        ...values,
        t_date: values.t_date.format('YYYY-MM-DD'),
        due_date: values.due_date?.format('YYYY-MM-DD'),
        details,
      }
      if (editing) {
        await issuesApi.update(editing.t_no, payload)
        message.success('Issue updated')
      } else {
        await issuesApi.create(payload)
        message.success('Issue created')
      }
      setModalOpen(false)
      load()
    } catch (err: any) {
      message.error(err.response?.data?.error ?? 'Error saving issue')
    }
  }

  const addDetail = () => setDetails(d => [
    ...d,
    { nos: 1, gross_wt: 0, net_wt: 0, stone_wt: 0, kundan_wt: 0, nakash_wt: 0 },
  ])

  const updateDetail = (i: number, field: keyof IssueDetail, value: any) => {
    setDetails(d => d.map((row, idx) => idx === i ? { ...row, [field]: value } : row))
  }

  const removeDetail = (i: number) => setDetails(d => d.filter((_, idx) => idx !== i))

  const totalNetWt = details.reduce((s, d) => s + (d.net_wt ?? 0), 0)
  const totalGrossWt = details.reduce((s, d) => s + (d.gross_wt ?? 0), 0)

  const columns = [
    { title: 'Voucher', dataIndex: 't_no', key: 't_no' },
    { title: 'Party', dataIndex: 'party_name', key: 'party_name' },
    { title: 'Date', dataIndex: 't_date', key: 't_date',
      render: (d: string) => new Date(d).toLocaleDateString('en-IN') },
    { title: 'Due Date', dataIndex: 'due_date', key: 'due_date',
      render: (d: string) => d ? new Date(d).toLocaleDateString('en-IN') : '—' },
    { title: 'Pieces', dataIndex: 'nos', key: 'nos' },
    { title: 'Gross Wt', dataIndex: 'gross_wt', key: 'gross_wt',
      render: (v: number) => `${(v ?? 0).toFixed(3)}g` },
    { title: 'Net Wt', dataIndex: 'net_wt', key: 'net_wt',
      render: (v: number) => `${(v ?? 0).toFixed(3)}g` },
    { title: 'Bill Amt', dataIndex: 'bill_amt', key: 'bill_amt',
      render: (v: number) => `₹${(v ?? 0).toFixed(2)}` },
    { title: 'Status', dataIndex: 'tag', key: 'tag',
      render: (t: string) => <Tag color={t === 'O' ? 'orange' : 'green'}>{t === 'O' ? 'Open' : 'Closed'}</Tag> },
    {
      title: 'Actions', key: 'actions', width: 100,
      render: (_: any, record: Issue) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => openEdit(record)} />
          <Popconfirm title="Delete this issue?" onConfirm={async () => {
            await issuesApi.delete(record.t_no); load()
          }}>
            <Button size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
        <Col><Title level={4} style={{ margin: 0 }}>Issue Vouchers ({total})</Title></Col>
        <Col>
          <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>New Issue</Button>
        </Col>
      </Row>

      <Table dataSource={issues} columns={columns} rowKey="t_no"
        loading={loading} size="small"
        pagination={{ current: page, pageSize: 50, total, onChange: p => { setPage(p); load(p) }, showTotal: t => `${t} issues` }} />

      <Modal
        title={editing ? `Edit Issue: ${editing.t_no}` : 'New Issue Voucher'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={() => setModalOpen(false)}
        width={1000}
        okText={editing ? 'Update' : 'Create'}
      >
        <Form form={form} layout="vertical" size="small">
          <Row gutter={12}>
            <Col span={6}>
              <Form.Item name="t_no" label="Voucher No" rules={[{ required: true }]}>
                <Input disabled={!!editing} />
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item name="ac_code" label="Party" rules={[{ required: true }]}>
                <Select showSearch placeholder="Select party" optionFilterProp="label"
                  options={accounts.map(a => ({ value: a.ac_code, label: `${a.ac_code} - ${a.desc}` }))} />
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item name="t_date" label="Issue Date" rules={[{ required: true }]}>
                <DatePicker style={{ width: '100%' }} format="DD/MM/YYYY" />
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item name="due_date" label="Due Date">
                <DatePicker style={{ width: '100%' }} format="DD/MM/YYYY" />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={12}>
            <Col span={6}>
              <Form.Item name="rate" label="Gold Rate (₹/g)" initialValue={0}>
                <InputNumber style={{ width: '100%' }} min={0} precision={2} />
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item name="net_per" label="Net %" initialValue={0}>
                <InputNumber style={{ width: '100%' }} min={0} max={100} precision={3} />
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item name="discount" label="Discount %" initialValue={0}>
                <InputNumber style={{ width: '100%' }} min={0} max={100} precision={2} />
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item name="bill_amt" label="Bill Amount (₹)" initialValue={0}>
                <InputNumber style={{ width: '100%' }} min={0} precision={2} />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={12}>
            <Col span={12}>
              <Form.Item name="narr" label="Narration">
                <Input />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="narr1" label="Narration 2">
                <Input />
              </Form.Item>
            </Col>
          </Row>

          <Divider orientation="left">Item Details</Divider>

          <div style={{ overflowX: 'auto' }}>
            <table style={{ width: '100%', fontSize: 12, borderCollapse: 'collapse' }}>
              <thead>
                <tr style={{ background: '#fafafa' }}>
                  {['Item Code', 'Description', 'Pcs', 'Gross Wt', 'Net Wt', 'Stone Wt', 'Kundan Wt', 'Nakash Wt', 'Narr', ''].map(h => (
                    <th key={h} style={{ padding: '4px 8px', border: '1px solid #eee', textAlign: 'left' }}>{h}</th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {details.map((d, i) => (
                  <tr key={i}>
                    {(['itcd', 'itdesc'] as const).map(f => (
                      <td key={f} style={{ border: '1px solid #eee', padding: 2 }}>
                        <Input size="small" value={d[f] ?? ''} onChange={e => updateDetail(i, f, e.target.value)} />
                      </td>
                    ))}
                    {(['nos', 'gross_wt', 'net_wt', 'stone_wt', 'kundan_wt', 'nakash_wt'] as const).map(f => (
                      <td key={f} style={{ border: '1px solid #eee', padding: 2 }}>
                        <InputNumber size="small" style={{ width: 85 }} value={d[f]} min={0} precision={3}
                          onChange={v => updateDetail(i, f, v ?? 0)} />
                      </td>
                    ))}
                    <td style={{ border: '1px solid #eee', padding: 2 }}>
                      <Input size="small" value={d.narr ?? ''} onChange={e => updateDetail(i, 'narr', e.target.value)} />
                    </td>
                    <td style={{ border: '1px solid #eee', padding: 2, textAlign: 'center' }}>
                      <Button size="small" danger icon={<DeleteFilled />} onClick={() => removeDetail(i)} />
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <div style={{ marginTop: 8, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Button size="small" icon={<PlusCircleOutlined />} onClick={addDetail}>Add Row</Button>
            <Space size="large">
              <Text type="secondary">Total Gross: <strong>{totalGrossWt.toFixed(3)}g</strong></Text>
              <Text type="secondary">Total Net: <strong>{totalNetWt.toFixed(3)}g</strong></Text>
            </Space>
          </div>
        </Form>
      </Modal>
    </div>
  )
}
