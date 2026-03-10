import { useEffect, useState } from 'react'
import { Table, Button, Space, Modal, Form, InputNumber, DatePicker, Input, message, Popconfirm, Select, Typography, Row, Col } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { lotsApi, accountsApi, type Lot, type Account } from '../api/client'
import dayjs from 'dayjs'

const { Title } = Typography

export default function Lots() {
  const [lots, setLots] = useState<Lot[]>([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Lot | null>(null)
  const [accounts, setAccounts] = useState<Account[]>([])
  const [form] = Form.useForm()

  const load = () => {
    setLoading(true)
    lotsApi.list().then(r => setLots(r.data || [])).finally(() => setLoading(false))
  }

  useEffect(() => {
    load()
    accountsApi.list({ per_page: 500 }).then(r => setAccounts(r.data.data || []))
  }, [])

  const openCreate = () => { setEditing(null); form.resetFields(); form.setFieldsValue({ t_date: dayjs() }); setModalOpen(true) }
  const openEdit = async (lot: Lot) => {
    setEditing(lot)
    form.setFieldsValue({ ...lot, t_date: dayjs(lot.t_date) })
    setModalOpen(true)
  }

  const handleSubmit = async () => {
    try {
      const v = await form.validateFields()
      v.t_date = v.t_date.format('YYYY-MM-DD')
      if (editing) { await lotsApi.update(editing.t_no, v); message.success('Lot updated') }
      else { await lotsApi.create(v); message.success('Lot created') }
      setModalOpen(false); load()
    } catch (err: any) { message.error(err.response?.data?.error ?? 'Error') }
  }

  const columns = [
    { title: 'T No', dataIndex: 't_no', key: 't_no' },
    { title: 'Date', dataIndex: 't_date', key: 't_date', render: (d: string) => new Date(d).toLocaleDateString('en-IN') },
    { title: 'Lot No', dataIndex: 'lot_no', key: 'lot_no' },
    { title: 'Party', dataIndex: 'party', key: 'party' },
    { title: 'Pieces', dataIndex: 'nos', key: 'nos' },
    { title: 'Gross Wt', dataIndex: 'gross_wt', key: 'gross_wt', render: (v: number) => `${(v ?? 0).toFixed(3)}g` },
    { title: 'Adj Wt', dataIndex: 'adj_wt', key: 'adj_wt', render: (v: number) => `${(v ?? 0).toFixed(3)}g` },
    { title: 'Amount', dataIndex: 'amt', key: 'amt', render: (v: number) => `₹${(v ?? 0).toFixed(2)}` },
    { title: 'Rate', dataIndex: 'g_rt', key: 'g_rt', render: (v: number) => v ? `₹${v.toFixed(2)}/g` : '—' },
    { title: 'Desc', dataIndex: 'desc', key: 'desc' },
    {
      title: 'Actions', key: 'actions',
      render: (_: any, r: Lot) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => openEdit(r)} />
          <Popconfirm title="Delete?" onConfirm={async () => { await lotsApi.delete(r.t_no); load() }}>
            <Button size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
        <Col><Title level={4} style={{ margin: 0 }}>Lots ({lots.length})</Title></Col>
        <Col><Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>New Lot</Button></Col>
      </Row>
      <Table dataSource={lots} columns={columns} rowKey="t_no" loading={loading} size="small" pagination={{ pageSize: 50 }} />

      <Modal title={editing ? `Edit Lot: ${editing.t_no}` : 'New Lot'} open={modalOpen}
        onOk={handleSubmit} onCancel={() => setModalOpen(false)} width={800}>
        <Form form={form} layout="vertical" size="small">
          <Row gutter={12}>
            <Col span={6}><Form.Item name="t_no" label="T No" rules={[{ required: true }]}><Input disabled={!!editing} /></Form.Item></Col>
            <Col span={6}><Form.Item name="lot_no" label="Lot No" rules={[{ required: true }]}><Input disabled={!!editing} /></Form.Item></Col>
            <Col span={6}><Form.Item name="t_date" label="Date" rules={[{ required: true }]}><DatePicker style={{ width: '100%' }} format="DD/MM/YYYY" /></Form.Item></Col>
            <Col span={6}>
              <Form.Item name="ac_code" label="Karigar">
                <Select showSearch allowClear placeholder="Select karigar" optionFilterProp="label"
                  options={accounts.map(a => ({ value: a.ac_code, label: `${a.ac_code} - ${a.desc}` }))} />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={12}>
            <Col span={12}><Form.Item name="party" label="Party Name"><Input /></Form.Item></Col>
            <Col span={12}><Form.Item name="desc" label="Description"><Input /></Form.Item></Col>
          </Row>
          <Row gutter={12}>
            {[['nos', 'Pieces', 0], ['gross_wt', 'Gross Wt', 3], ['kundan_wt', 'Kundan Wt', 3],
              ['nakash_wt', 'Nakash Wt', 3], ['stone_wt', 'Stone Wt', 3], ['pearl_wt', 'Pearl Wt', 3]].map(([n, l, p]) => (
              <Col span={4} key={n as string}>
                <Form.Item name={n} label={l} initialValue={0}>
                  <InputNumber style={{ width: '100%' }} precision={p as number} />
                </Form.Item>
              </Col>
            ))}
          </Row>
          <Row gutter={12}>
            {[['nos1', 'Return Pcs', 0], ['gross_wt1', 'Return Gross', 3], ['adj_wt', 'Adj Wt', 3],
              ['amt', 'Amount', 2], ['g_rt', 'Gold Rate', 2]].map(([n, l, p]) => (
              <Col span={4} key={n as string}>
                <Form.Item name={n} label={l} initialValue={0}>
                  <InputNumber style={{ width: '100%' }} precision={p as number} />
                </Form.Item>
              </Col>
            ))}
          </Row>
        </Form>
      </Modal>
    </div>
  )
}
