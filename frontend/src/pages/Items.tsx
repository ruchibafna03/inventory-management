import { useEffect, useState } from 'react'
import {
  Table, Button, Input, Space, Modal, Form, InputNumber,
  DatePicker, message, Popconfirm, Tag, Select, Typography, Row, Col,
} from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, SearchOutlined } from '@ant-design/icons'
import { itemsApi, type Item } from '../api/client'
import dayjs from 'dayjs'

const { Title } = Typography

const PURITY_OPTIONS = ['91', '75', '58', '92', '99']

export default function Items() {
  const [items, setItems] = useState<Item[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(false)
  const [search, setSearch] = useState('')
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Item | null>(null)
  const [form] = Form.useForm()

  const load = (p = page, s = search) => {
    setLoading(true)
    itemsApi.list({ page: p, per_page: 50, search: s })
      .then(r => { setItems(r.data.data); setTotal(r.data.total) })
      .finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  const openCreate = () => {
    setEditing(null)
    form.resetFields()
    setModalOpen(true)
  }

  const openEdit = (item: Item) => {
    setEditing(item)
    form.setFieldsValue({
      ...item,
      recpt_date: item.recpt_date ? dayjs(item.recpt_date) : null,
      issue_date: item.issue_date ? dayjs(item.issue_date) : null,
    })
    setModalOpen(true)
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (values.recpt_date) values.recpt_date = values.recpt_date.format('YYYY-MM-DD')
      if (values.issue_date) values.issue_date = values.issue_date.format('YYYY-MM-DD')

      if (editing) {
        await itemsApi.update(editing.itcd, values)
        message.success('Item updated')
      } else {
        await itemsApi.create(values)
        message.success('Item created')
      }
      setModalOpen(false)
      load()
    } catch (err: any) {
      message.error(err.response?.data?.error ?? 'Error saving item')
    }
  }

  const handleDelete = async (itcd: string) => {
    try {
      await itemsApi.delete(itcd)
      message.success('Item deleted')
      load()
    } catch (err: any) {
      message.error(err.response?.data?.error ?? 'Error deleting item')
    }
  }

  const columns = [
    { title: 'Code', dataIndex: 'itcd', key: 'itcd', width: 90, fixed: 'left' as const },
    { title: 'Description', dataIndex: 'desc', key: 'desc', ellipsis: true },
    { title: 'Purity', dataIndex: 'purity', key: 'purity', width: 70,
      render: (v: string) => v ? <Tag color="gold">{v}</Tag> : '—' },
    { title: 'Gross Wt', dataIndex: 'gross_wt', key: 'gross_wt', width: 100,
      render: (v: number) => `${(v ?? 0).toFixed(3)}g` },
    { title: 'Net Wt', dataIndex: 'net_wt', key: 'net_wt', width: 100,
      render: (v: number) => `${(v ?? 0).toFixed(3)}g` },
    { title: 'Gold Wt', dataIndex: 'gold_wt', key: 'gold_wt', width: 100,
      render: (v: number) => `${(v ?? 0).toFixed(3)}g` },
    { title: 'Stone Wt', dataIndex: 'stone_wt', key: 'stone_wt', width: 100,
      render: (v: number) => `${(v ?? 0).toFixed(3)}g` },
    { title: 'Making', dataIndex: 'mk_chrg', key: 'mk_chrg', width: 100,
      render: (v: number) => `₹${(v ?? 0).toFixed(2)}` },
    { title: 'Lot', dataIndex: 'lot_no', key: 'lot_no', width: 80 },
    { title: 'Cat', dataIndex: 'cat', key: 'cat', width: 60 },
    {
      title: 'Actions', key: 'actions', width: 100, fixed: 'right' as const,
      render: (_: any, record: Item) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => openEdit(record)} />
          <Popconfirm title="Delete this item?" onConfirm={() => handleDelete(record.itcd)}>
            <Button size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
        <Col><Title level={4} style={{ margin: 0 }}>Item Master ({total})</Title></Col>
        <Col>
          <Space>
            <Input.Search
              placeholder="Search code or description"
              allowClear
              onSearch={s => { setSearch(s); setPage(1); load(1, s) }}
              style={{ width: 280 }}
              prefix={<SearchOutlined />}
            />
            <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>New Item</Button>
          </Space>
        </Col>
      </Row>

      <Table
        dataSource={items}
        columns={columns}
        rowKey="itcd"
        loading={loading}
        scroll={{ x: 1000 }}
        size="small"
        pagination={{
          current: page,
          pageSize: 50,
          total,
          onChange: p => { setPage(p); load(p) },
          showSizeChanger: false,
          showTotal: t => `${t} items`,
        }}
      />

      <Modal
        title={editing ? `Edit Item: ${editing.itcd}` : 'New Item'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={() => setModalOpen(false)}
        width={900}
        okText={editing ? 'Update' : 'Create'}
      >
        <Form form={form} layout="vertical" size="small">
          <Row gutter={12}>
            <Col span={6}>
              <Form.Item name="itcd" label="Item Code" rules={[{ required: true }]}>
                <Input disabled={!!editing} placeholder="e.g. G001234" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="desc" label="Description" rules={[{ required: true }]}>
                <Input placeholder="Item description" />
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item name="purity" label="Purity">
                <Select placeholder="Select purity" allowClear>
                  {PURITY_OPTIONS.map(p => (
                    <Select.Option key={p} value={p}>{p}%</Select.Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
          </Row>

          <Title level={5} style={{ marginBottom: 8 }}>Weights (grams)</Title>
          <Row gutter={12}>
            {[
              ['gross_wt', 'Gross Wt'], ['net_wt', 'Net Wt'], ['gold_wt', 'Gold Wt'],
              ['ghat_wt', 'Ghat Wt'], ['kundan_wt', 'Kundan Wt'], ['extra_wt', 'Extra Wt'],
              ['nakash_wt', 'Nakash Wt'], ['stone_wt', 'Stone Wt'], ['prl_wt', 'Pearl Wt'],
            ].map(([name, label]) => (
              <Col span={6} key={name}>
                <Form.Item name={name} label={label} initialValue={0}>
                  <InputNumber style={{ width: '100%' }} min={0} precision={3} step={0.001} />
                </Form.Item>
              </Col>
            ))}
          </Row>

          <Title level={5} style={{ marginBottom: 8 }}>Stones</Title>
          <Row gutter={12}>
            {[
              ['rby_ct', 'Ruby CT'], ['rby_gm', 'Ruby GM'],
              ['eme_ct', 'Emerald CT'], ['eme_gm', 'Emerald GM'],
              ['plk_ct', 'Polki CT'], ['plk_gm', 'Polki GM'],
            ].map(([name, label]) => (
              <Col span={4} key={name}>
                <Form.Item name={name} label={label} initialValue={0}>
                  <InputNumber style={{ width: '100%' }} min={0} precision={3} step={0.001} />
                </Form.Item>
              </Col>
            ))}
          </Row>

          <Row gutter={12}>
            <Col span={6}>
              <Form.Item name="mk_chrg" label="Making Charge (₹)" initialValue={0}>
                <InputNumber style={{ width: '100%' }} min={0} precision={2} />
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item name="melt" label="Melt %" initialValue={0}>
                <InputNumber style={{ width: '100%' }} min={0} max={100} precision={2} />
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item name="cat" label="Category">
                <Input placeholder="e.g. RNG" />
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item name="lot_no" label="Lot No">
                <Input />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={12}>
            <Col span={12}>
              <Form.Item name="narr1" label="Narration 1">
                <Input />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="narr2" label="Narration 2">
                <Input />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Modal>
    </div>
  )
}
