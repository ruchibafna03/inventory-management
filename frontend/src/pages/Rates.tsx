import { useEffect, useState } from 'react'
import { Table, Button, Space, Modal, Form, InputNumber, DatePicker, message, Popconfirm, Typography, Row, Col, Card, Statistic } from 'antd'
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons'
import { ratesApi, type Rate } from '../api/client'
import dayjs from 'dayjs'

const { Title } = Typography

export default function Rates() {
  const [rates, setRates] = useState<Rate[]>([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [form] = Form.useForm()

  const load = () => {
    setLoading(true)
    ratesApi.list().then(r => setRates(r.data)).finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  const handleSubmit = async () => {
    try {
      const v = await form.validateFields()
      v.date = v.date.format('YYYY-MM-DD')
      await ratesApi.create(v)
      message.success('Rate saved')
      setModalOpen(false)
      load()
    } catch (err: any) {
      message.error(err.response?.data?.error ?? 'Error saving rate')
    }
  }

  const latest = rates[0]

  const columns = [
    { title: 'Date', dataIndex: 'date', key: 'date',
      render: (d: string) => new Date(d).toLocaleDateString('en-IN') },
    { title: '22K Rate (₹/g)', dataIndex: 'rate', key: 'rate',
      render: (v: number) => <strong>₹{(v ?? 0).toFixed(2)}</strong> },
    { title: '18K Rate (₹/g)', dataIndex: 'rate1', key: 'rate1',
      render: (v: number) => `₹${(v ?? 0).toFixed(2)}` },
    { title: 'Silver (₹/g)', dataIndex: 's_rate', key: 's_rate',
      render: (v: number) => `₹${(v ?? 0).toFixed(2)}` },
    {
      title: 'Actions', key: 'actions',
      render: (_: any, r: Rate) => (
        <Popconfirm title="Delete this rate?" onConfirm={async () => { await ratesApi.delete(r.id); load() }}>
          <Button size="small" danger icon={<DeleteOutlined />} />
        </Popconfirm>
      ),
    },
  ]

  return (
    <div>
      <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
        <Col><Title level={4} style={{ margin: 0 }}>Gold Rates</Title></Col>
        <Col><Button type="primary" icon={<PlusOutlined />} onClick={() => { form.resetFields(); form.setFieldsValue({ date: dayjs() }); setModalOpen(true) }}>Update Rate</Button></Col>
      </Row>

      {latest && (
        <Row gutter={16} style={{ marginBottom: 16 }}>
          <Col xs={24} sm={8}>
            <Card size="small">
              <Statistic title="Current 22K Rate" value={latest.rate} suffix="₹/g" valueStyle={{ color: '#d4a017', fontSize: 24 }} />
            </Card>
          </Col>
          <Col xs={24} sm={8}>
            <Card size="small">
              <Statistic title="Current 18K Rate" value={latest.rate1} suffix="₹/g" valueStyle={{ color: '#aaa', fontSize: 24 }} />
            </Card>
          </Col>
          <Col xs={24} sm={8}>
            <Card size="small">
              <Statistic title="Silver Rate" value={latest.s_rate} suffix="₹/g" valueStyle={{ color: '#888', fontSize: 24 }} />
            </Card>
          </Col>
        </Row>
      )}

      <Table dataSource={rates} columns={columns} rowKey="id" loading={loading} size="small"
        pagination={{ pageSize: 20, showTotal: t => `${t} rate entries` }} />

      <Modal title="Update Gold Rate" open={modalOpen} onOk={handleSubmit} onCancel={() => setModalOpen(false)}>
        <Form form={form} layout="vertical">
          <Form.Item name="date" label="Date" rules={[{ required: true }]}>
            <DatePicker style={{ width: '100%' }} format="DD/MM/YYYY" />
          </Form.Item>
          <Row gutter={12}>
            <Col span={8}>
              <Form.Item name="rate" label="22K Rate (₹/g)" rules={[{ required: true }]} initialValue={0}>
                <InputNumber style={{ width: '100%' }} precision={2} min={0} />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="rate1" label="18K Rate (₹/g)" initialValue={0}>
                <InputNumber style={{ width: '100%' }} precision={2} min={0} />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="s_rate" label="Silver Rate (₹/g)" initialValue={0}>
                <InputNumber style={{ width: '100%' }} precision={2} min={0} />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Modal>
    </div>
  )
}
