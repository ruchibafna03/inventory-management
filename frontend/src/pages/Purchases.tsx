import { useEffect, useState } from 'react'
import { Table, Button, Space, Modal, Form, InputNumber, DatePicker, Input, message, Popconfirm, Typography, Row, Col } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { salesApi as _s, type Sale as _Sale } from '../api/client'
import axios from 'axios'
import dayjs from 'dayjs'

const { Title } = Typography

interface Purchase {
  vouch_no: string; vouch_date: string; ac_code?: string; name?: string
  rate: number; gross_wt: number; net_wt: number; bill_amt: number; narr?: string
}

const api = {
  list: (p: any) => axios.get('/api/v1/purchases', { params: p }),
  get: (v: string) => axios.get(`/api/v1/purchases/${v}`),
  create: (d: any) => axios.post('/api/v1/purchases', d),
  update: (v: string, d: any) => axios.put(`/api/v1/purchases/${v}`, d),
  delete: (v: string) => axios.delete(`/api/v1/purchases/${v}`),
}

export default function Purchases() {
  const [data, setData] = useState<Purchase[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Purchase | null>(null)
  const [form] = Form.useForm()

  const load = (p = page) => {
    setLoading(true)
    api.list({ page: p, per_page: 50 }).then(r => { setData(r.data.data || []); setTotal(r.data.total) }).finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  const openCreate = () => { setEditing(null); form.resetFields(); form.setFieldsValue({ vouch_date: dayjs() }); setModalOpen(true) }
  const openEdit = (p: Purchase) => { setEditing(p); form.setFieldsValue({ ...p, vouch_date: dayjs(p.vouch_date) }); setModalOpen(true) }

  const handleSubmit = async () => {
    try {
      const v = await form.validateFields()
      v.vouch_date = v.vouch_date.format('YYYY-MM-DD')
      if (editing) { await api.update(editing.vouch_no, v); message.success('Updated') }
      else { await api.create(v); message.success('Created') }
      setModalOpen(false); load()
    } catch (err: any) { message.error(err.response?.data?.error ?? 'Error') }
  }

  return (
    <div>
      <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
        <Col><Title level={4} style={{ margin: 0 }}>Purchases ({total})</Title></Col>
        <Col><Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>New Purchase</Button></Col>
      </Row>
      <Table dataSource={data} rowKey="vouch_no" loading={loading} size="small"
        pagination={{ current: page, pageSize: 50, total, onChange: p => { setPage(p); load(p) }, showTotal: t => `${t} purchases` }}
        columns={[
          { title: 'Voucher', dataIndex: 'vouch_no' },
          { title: 'Date', dataIndex: 'vouch_date', render: (d: string) => new Date(d).toLocaleDateString('en-IN') },
          { title: 'Supplier', dataIndex: 'name' },
          { title: 'Rate', dataIndex: 'rate', render: (v: number) => `₹${v?.toFixed(2)}` },
          { title: 'Gross Wt', dataIndex: 'gross_wt', render: (v: number) => `${v?.toFixed(3)}g` },
          { title: 'Net Wt', dataIndex: 'net_wt', render: (v: number) => `${v?.toFixed(3)}g` },
          { title: 'Bill Amt', dataIndex: 'bill_amt', render: (v: number) => `₹${v?.toFixed(2)}` },
          { title: 'Actions', key: 'actions', render: (_: any, r: Purchase) => (
            <Space>
              <Button size="small" icon={<EditOutlined />} onClick={() => openEdit(r)} />
              <Popconfirm title="Delete?" onConfirm={async () => { await api.delete(r.vouch_no); load() }}>
                <Button size="small" danger icon={<DeleteOutlined />} />
              </Popconfirm>
            </Space>
          )},
        ]} />

      <Modal title={editing ? `Edit: ${editing.vouch_no}` : 'New Purchase'} open={modalOpen}
        onOk={handleSubmit} onCancel={() => setModalOpen(false)} width={700}>
        <Form form={form} layout="vertical" size="small">
          <Row gutter={12}>
            <Col span={6}><Form.Item name="vouch_no" label="Voucher No" rules={[{ required: true }]}><Input disabled={!!editing} /></Form.Item></Col>
            <Col span={6}><Form.Item name="vouch_date" label="Date" rules={[{ required: true }]}><DatePicker style={{ width: '100%' }} format="DD/MM/YYYY" /></Form.Item></Col>
            <Col span={6}><Form.Item name="ac_code" label="Account Code"><Input /></Form.Item></Col>
            <Col span={6}><Form.Item name="name" label="Supplier Name"><Input /></Form.Item></Col>
          </Row>
          <Row gutter={12}>
            <Col span={8}><Form.Item name="add1" label="Address 1"><Input /></Form.Item></Col>
            <Col span={8}><Form.Item name="add2" label="Address 2"><Input /></Form.Item></Col>
            <Col span={8}><Form.Item name="rate" label="Rate (₹/g)" initialValue={0}><InputNumber style={{ width: '100%' }} precision={2} /></Form.Item></Col>
          </Row>
          <Row gutter={12}>
            <Col span={8}><Form.Item name="gross_wt" label="Gross Wt" initialValue={0}><InputNumber style={{ width: '100%' }} precision={3} /></Form.Item></Col>
            <Col span={8}><Form.Item name="net_wt" label="Net Wt" initialValue={0}><InputNumber style={{ width: '100%' }} precision={3} /></Form.Item></Col>
            <Col span={8}><Form.Item name="bill_amt" label="Bill Amount" initialValue={0}><InputNumber style={{ width: '100%' }} precision={2} /></Form.Item></Col>
          </Row>
          <Form.Item name="narr" label="Narration"><Input /></Form.Item>
          <Form.Item name="rmk" label="Remarks"><Input /></Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
