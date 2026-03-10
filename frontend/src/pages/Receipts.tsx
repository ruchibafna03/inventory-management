import { useEffect, useState } from 'react'
import {
  Table, Button, Space, Modal, Form, InputNumber, DatePicker,
  Input, message, Popconfirm, Select, Typography, Row, Col, Divider,
} from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, PlusCircleOutlined, DeleteFilled } from '@ant-design/icons'
import { receiptsApi, accountsApi, type Receipt, type ReceiptDetail, type Account } from '../api/client'
import dayjs from 'dayjs'

const { Title, Text } = Typography

export default function Receipts() {
  const [receipts, setReceipts] = useState<Receipt[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Receipt | null>(null)
  const [details, setDetails] = useState<ReceiptDetail[]>([])
  const [accounts, setAccounts] = useState<Account[]>([])
  const [form] = Form.useForm()

  const load = (p = page) => {
    setLoading(true)
    receiptsApi.list({ page: p, per_page: 50 })
      .then(r => { setReceipts(r.data.data); setTotal(r.data.total) })
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    load()
    accountsApi.list({ per_page: 500 }).then(r => setAccounts(r.data.data))
  }, [])

  const openCreate = () => {
    setEditing(null)
    form.resetFields()
    form.setFieldsValue({ t_date: dayjs() })
    setDetails([{ gross_wt: 0, net_wt: 0, gold_wt: 0, stone_wt: 0, kundan_wt: 0, nakash_wt: 0 }])
    setModalOpen(true)
  }

  const openEdit = async (recpt: Receipt) => {
    setEditing(recpt)
    const full = await receiptsApi.get(recpt.t_no)
    form.setFieldsValue({
      ...full.data,
      t_date: dayjs(full.data.t_date),
    })
    setDetails(full.data.details ?? [])
    setModalOpen(true)
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      const payload = { ...values, t_date: values.t_date.format('YYYY-MM-DD'), details }
      if (editing) {
        await receiptsApi.update(editing.t_no, payload)
        message.success('Receipt updated')
      } else {
        await receiptsApi.create(payload)
        message.success('Receipt created')
      }
      setModalOpen(false)
      load()
    } catch (err: any) {
      message.error(err.response?.data?.error ?? 'Error saving receipt')
    }
  }

  const updateDetail = (i: number, field: keyof ReceiptDetail, value: any) =>
    setDetails(d => d.map((row, idx) => idx === i ? { ...row, [field]: value } : row))

  const totalNetWt = details.reduce((s, d) => s + (d.net_wt ?? 0), 0)
  const totalGoldWt = details.reduce((s, d) => s + (d.gold_wt ?? 0), 0)

  const columns = [
    { title: 'Voucher', dataIndex: 't_no', key: 't_no' },
    { title: 'Party', dataIndex: 'party_name', key: 'party_name' },
    { title: 'Date', dataIndex: 't_date', key: 't_date',
      render: (d: string) => new Date(d).toLocaleDateString('en-IN') },
    { title: 'Pieces', dataIndex: 'nos', key: 'nos' },
    { title: 'Gross Wt', dataIndex: 'gross_wt', key: 'gross_wt',
      render: (v: number) => `${(v ?? 0).toFixed(3)}g` },
    { title: 'Net Wt', dataIndex: 'net_wt', key: 'net_wt',
      render: (v: number) => `${(v ?? 0).toFixed(3)}g` },
    { title: 'Pure Wt', dataIndex: 'pure_wt', key: 'pure_wt',
      render: (v: number) => `${(v ?? 0).toFixed(3)}g` },
    { title: 'Bill Amt', dataIndex: 'bill_amt', key: 'bill_amt',
      render: (v: number) => `₹${(v ?? 0).toFixed(2)}` },
    {
      title: 'Actions', key: 'actions',
      render: (_: any, record: Receipt) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => openEdit(record)} />
          <Popconfirm title="Delete?" onConfirm={async () => { await receiptsApi.delete(record.t_no); load() }}>
            <Button size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
        <Col><Title level={4} style={{ margin: 0 }}>Receipt Vouchers ({total})</Title></Col>
        <Col><Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>New Receipt</Button></Col>
      </Row>

      <Table dataSource={receipts} columns={columns} rowKey="t_no" loading={loading} size="small"
        pagination={{ current: page, pageSize: 50, total, onChange: p => { setPage(p); load(p) }, showTotal: t => `${t} receipts` }} />

      <Modal title={editing ? `Edit Receipt: ${editing.t_no}` : 'New Receipt Voucher'}
        open={modalOpen} onOk={handleSubmit} onCancel={() => setModalOpen(false)}
        width={1000} okText={editing ? 'Update' : 'Create'}>
        <Form form={form} layout="vertical" size="small">
          <Row gutter={12}>
            <Col span={6}><Form.Item name="t_no" label="Voucher No" rules={[{ required: true }]}><Input disabled={!!editing} /></Form.Item></Col>
            <Col span={6}>
              <Form.Item name="ac_code" label="Party" rules={[{ required: true }]}>
                <Select showSearch placeholder="Select party" optionFilterProp="label"
                  options={accounts.map(a => ({ value: a.ac_code, label: `${a.ac_code} - ${a.desc}` }))} />
              </Form.Item>
            </Col>
            <Col span={6}><Form.Item name="t_date" label="Date" rules={[{ required: true }]}><DatePicker style={{ width: '100%' }} format="DD/MM/YYYY" /></Form.Item></Col>
            <Col span={6}><Form.Item name="rate" label="Gold Rate" initialValue={0}><InputNumber style={{ width: '100%' }} precision={2} /></Form.Item></Col>
          </Row>
          <Row gutter={12}>
            <Col span={6}><Form.Item name="net_per" label="Net %" initialValue={0}><InputNumber style={{ width: '100%' }} precision={3} /></Form.Item></Col>
            <Col span={6}><Form.Item name="discount" label="Discount %" initialValue={0}><InputNumber style={{ width: '100%' }} precision={2} /></Form.Item></Col>
            <Col span={6}><Form.Item name="bill_amt" label="Bill Amount" initialValue={0}><InputNumber style={{ width: '100%' }} precision={2} /></Form.Item></Col>
            <Col span={6}><Form.Item name="narr" label="Narration"><Input /></Form.Item></Col>
          </Row>

          <Divider orientation="left">Item Details</Divider>
          <div style={{ overflowX: 'auto' }}>
            <table style={{ width: '100%', fontSize: 12, borderCollapse: 'collapse' }}>
              <thead>
                <tr style={{ background: '#fafafa' }}>
                  {['Item Code', 'Description', 'Lot No', 'Gross Wt', 'Net Wt', 'Gold Wt', 'Stone Wt', 'Kundan Wt', ''].map(h => (
                    <th key={h} style={{ padding: '4px 8px', border: '1px solid #eee' }}>{h}</th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {details.map((d, i) => (
                  <tr key={i}>
                    {(['itcd', 'itdesc', 'lot_no'] as const).map(f => (
                      <td key={f} style={{ border: '1px solid #eee', padding: 2 }}>
                        <Input size="small" value={d[f] ?? ''} onChange={e => updateDetail(i, f, e.target.value)} />
                      </td>
                    ))}
                    {(['gross_wt', 'net_wt', 'gold_wt', 'stone_wt', 'kundan_wt'] as const).map(f => (
                      <td key={f} style={{ border: '1px solid #eee', padding: 2 }}>
                        <InputNumber size="small" style={{ width: 85 }} value={d[f]} min={0} precision={3}
                          onChange={v => updateDetail(i, f, v ?? 0)} />
                      </td>
                    ))}
                    <td style={{ border: '1px solid #eee', padding: 2, textAlign: 'center' }}>
                      <Button size="small" danger icon={<DeleteFilled />}
                        onClick={() => setDetails(d => d.filter((_, idx) => idx !== i))} />
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div style={{ marginTop: 8, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Button size="small" icon={<PlusCircleOutlined />}
              onClick={() => setDetails(d => [...d, { gross_wt: 0, net_wt: 0, gold_wt: 0, stone_wt: 0, kundan_wt: 0, nakash_wt: 0 }])}>
              Add Row
            </Button>
            <Space size="large">
              <Text type="secondary">Total Net: <strong>{totalNetWt.toFixed(3)}g</strong></Text>
              <Text type="secondary">Total Gold: <strong>{totalGoldWt.toFixed(3)}g</strong></Text>
            </Space>
          </div>
        </Form>
      </Modal>
    </div>
  )
}
