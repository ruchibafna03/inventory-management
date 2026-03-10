import { useEffect, useState } from 'react'
import {
  Table, Button, Space, Modal, Form, InputNumber, DatePicker,
  Input, message, Popconfirm, Typography, Row, Col,
} from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { salesApi, type Sale } from '../api/client'
import dayjs from 'dayjs'

const { Title } = Typography

export default function Sales() {
  const [sales, setSales] = useState<Sale[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Sale | null>(null)
  const [form] = Form.useForm()

  const load = (p = page) => {
    setLoading(true)
    salesApi.list({ page: p, per_page: 50 })
      .then(r => { setSales(r.data.data); setTotal(r.data.total) })
      .finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  const openCreate = () => {
    setEditing(null)
    form.resetFields()
    form.setFieldsValue({ vouch_date: dayjs() })
    setModalOpen(true)
  }

  const openEdit = async (sale: Sale) => {
    setEditing(sale)
    form.setFieldsValue({ ...sale, vouch_date: dayjs(sale.vouch_date) })
    setModalOpen(true)
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      values.vouch_date = values.vouch_date.format('YYYY-MM-DD')
      if (editing) {
        await salesApi.update(editing.vouch_no, values)
        message.success('Sale updated')
      } else {
        await salesApi.create(values)
        message.success('Sale created')
      }
      setModalOpen(false)
      load()
    } catch (err: any) {
      message.error(err.response?.data?.error ?? 'Error saving sale')
    }
  }

  const columns = [
    { title: 'Voucher', dataIndex: 'vouch_no', key: 'vouch_no' },
    { title: 'Date', dataIndex: 'vouch_date', key: 'vouch_date',
      render: (d: string) => new Date(d).toLocaleDateString('en-IN') },
    { title: 'Customer', dataIndex: 'name', key: 'name' },
    { title: 'Phone', dataIndex: 'phone_no', key: 'phone_no' },
    { title: 'Pieces', dataIndex: 'nos', key: 'nos' },
    { title: 'Gross Wt', dataIndex: 'gross_wt', key: 'gross_wt',
      render: (v: number) => `${(v ?? 0).toFixed(3)}g` },
    { title: 'Net Wt', dataIndex: 'net_wt', key: 'net_wt',
      render: (v: number) => `${(v ?? 0).toFixed(3)}g` },
    { title: 'Discount', dataIndex: 'discount', key: 'discount',
      render: (v: number) => `${(v ?? 0).toFixed(2)}%` },
    { title: 'Pure Final', dataIndex: 'pure_final', key: 'pure_final',
      render: (v: number) => `${(v ?? 0).toFixed(3)}g` },
    {
      title: 'Actions', key: 'actions',
      render: (_: any, r: Sale) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => openEdit(r)} />
          <Popconfirm title="Delete?" onConfirm={async () => { await salesApi.delete(r.vouch_no); load() }}>
            <Button size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
        <Col><Title level={4} style={{ margin: 0 }}>Sales ({total})</Title></Col>
        <Col><Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>New Sale</Button></Col>
      </Row>

      <Table dataSource={sales} columns={columns} rowKey="vouch_no" loading={loading} size="small"
        pagination={{ current: page, pageSize: 50, total, onChange: p => { setPage(p); load(p) }, showTotal: t => `${t} sales` }} />

      <Modal title={editing ? `Edit Sale: ${editing.vouch_no}` : 'New Sale'}
        open={modalOpen} onOk={handleSubmit} onCancel={() => setModalOpen(false)} width={800}>
        <Form form={form} layout="vertical" size="small">
          <Row gutter={12}>
            <Col span={6}><Form.Item name="vouch_no" label="Voucher No" rules={[{ required: true }]}><Input disabled={!!editing} /></Form.Item></Col>
            <Col span={6}><Form.Item name="vouch_date" label="Date" rules={[{ required: true }]}><DatePicker style={{ width: '100%' }} format="DD/MM/YYYY" /></Form.Item></Col>
            <Col span={6}><Form.Item name="ac_code" label="Account Code"><Input /></Form.Item></Col>
            <Col span={6}><Form.Item name="sman" label="Salesman"><Input /></Form.Item></Col>
          </Row>
          <Row gutter={12}>
            <Col span={8}><Form.Item name="name" label="Customer Name"><Input /></Form.Item></Col>
            <Col span={8}><Form.Item name="add1" label="Address 1"><Input /></Form.Item></Col>
            <Col span={8}><Form.Item name="add2" label="Address 2"><Input /></Form.Item></Col>
          </Row>
          <Row gutter={12}>
            <Col span={6}><Form.Item name="phone_no" label="Phone"><Input /></Form.Item></Col>
            <Col span={6}><Form.Item name="phone_no1" label="Alt Phone"><Input /></Form.Item></Col>
            <Col span={6}><Form.Item name="nos" label="Pieces" initialValue={0}><InputNumber style={{ width: '100%' }} /></Form.Item></Col>
            <Col span={6}><Form.Item name="discount" label="Discount %" initialValue={0}><InputNumber style={{ width: '100%' }} precision={2} /></Form.Item></Col>
          </Row>
          <Row gutter={12}>
            {[['gross_wt', 'Gross Wt'], ['net_wt', 'Net Wt'], ['stone_wt', 'Stone Wt'], ['pure_wt', 'Pure Wt'], ['pure_final', 'Pure Final']].map(([n, l]) => (
              <Col span={4} key={n}>
                <Form.Item name={n} label={l} initialValue={0}>
                  <InputNumber style={{ width: '100%' }} precision={3} />
                </Form.Item>
              </Col>
            ))}
            <Col span={4}><Form.Item name="stn_amt" label="Stone Amt" initialValue={0}><InputNumber style={{ width: '100%' }} precision={2} /></Form.Item></Col>
          </Row>
          <Row gutter={12}>
            <Col span={12}><Form.Item name="narr" label="Narration"><Input /></Form.Item></Col>
            <Col span={12}><Form.Item name="narr1" label="Narration 2"><Input /></Form.Item></Col>
          </Row>
        </Form>
      </Modal>
    </div>
  )
}
