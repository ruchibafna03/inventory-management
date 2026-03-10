import { useEffect, useState } from 'react'
import {
  Table, Button, Space, Modal, Form, Input, InputNumber, DatePicker,
  message, Popconfirm, Select, Typography, Row, Col, Tabs,
} from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { accountsApi, type Account, type AccountGroup, type Address } from '../api/client'
import dayjs from 'dayjs'

const { Title } = Typography

export default function Accounts() {
  const [accounts, setAccounts] = useState<Account[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(false)
  const [search, setSearch] = useState('')
  const [groups, setGroups] = useState<AccountGroup[]>([])
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Account | null>(null)
  const [address, setAddress] = useState<Address | null>(null)
  const [form] = Form.useForm()
  const [addrForm] = Form.useForm()

  const loadGroups = () => accountsApi.listGroups().then(r => setGroups(r.data))

  const load = (p = page, s = search) => {
    setLoading(true)
    accountsApi.list({ page: p, per_page: 50, search: s })
      .then(r => { setAccounts(r.data.data); setTotal(r.data.total) })
      .finally(() => setLoading(false))
  }

  useEffect(() => { load(); loadGroups() }, [])

  const openCreate = () => {
    setEditing(null); setAddress(null)
    form.resetFields(); addrForm.resetFields()
    setModalOpen(true)
  }

  const openEdit = async (acc: Account) => {
    setEditing(acc)
    form.setFieldsValue(acc)
    try {
      const addr = await accountsApi.getAddress(acc.ac_code)
      setAddress(addr.data)
      addrForm.setFieldsValue(addr.data)
    } catch { addrForm.resetFields() }
    setModalOpen(true)
  }

  const handleSubmit = async () => {
    try {
      const accVals = await form.validateFields()
      const addrVals = await addrForm.validateFields()

      if (editing) {
        await accountsApi.update(editing.ac_code, accVals)
        await accountsApi.upsertAddress(editing.ac_code, addrVals)
        message.success('Account updated')
      } else {
        const created = await accountsApi.create(accVals)
        await accountsApi.upsertAddress(created.data.ac_code, addrVals)
        message.success('Account created')
      }
      setModalOpen(false); load()
    } catch (err: any) {
      message.error(err.response?.data?.error ?? 'Error saving account')
    }
  }

  const columns = [
    { title: 'Code', dataIndex: 'ac_code', key: 'ac_code', width: 80 },
    { title: 'Name', dataIndex: 'desc', key: 'desc' },
    { title: 'Group', dataIndex: 'g_code', key: 'g_code', width: 80 },
    { title: 'Opening Bal', dataIndex: 'opb', key: 'opb',
      render: (v: number) => `₹${(v ?? 0).toFixed(2)}` },
    { title: 'Current Bal', dataIndex: 'amt', key: 'amt',
      render: (v: number) => <span style={{ color: (v ?? 0) < 0 ? 'red' : 'green' }}>₹{(v ?? 0).toFixed(2)}</span> },
    { title: 'Rate', dataIndex: 'rate', key: 'rate',
      render: (v: number) => v ? `${v?.toFixed(3)}g` : '—' },
    { title: 'Cat', dataIndex: 'cat', key: 'cat', width: 60 },
    {
      title: 'Actions', key: 'actions', width: 100,
      render: (_: any, r: Account) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => openEdit(r)} />
          <Popconfirm title="Delete?" onConfirm={async () => { await accountsApi.delete(r.ac_code); load() }}>
            <Button size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
        <Col><Title level={4} style={{ margin: 0 }}>Account Master ({total})</Title></Col>
        <Col>
          <Space>
            <Input.Search placeholder="Search" allowClear onSearch={s => { setSearch(s); setPage(1); load(1, s) }} style={{ width: 240 }} />
            <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>New Account</Button>
          </Space>
        </Col>
      </Row>

      <Table dataSource={accounts} columns={columns} rowKey="ac_code" loading={loading} size="small"
        pagination={{ current: page, pageSize: 50, total, onChange: p => { setPage(p); load(p) }, showTotal: t => `${t} accounts` }} />

      <Modal title={editing ? `Edit: ${editing.ac_code} - ${editing.desc}` : 'New Account'}
        open={modalOpen} onOk={handleSubmit} onCancel={() => setModalOpen(false)} width={700}>
        <Tabs items={[
          {
            key: 'account', label: 'Account',
            children: (
              <Form form={form} layout="vertical" size="small">
                <Row gutter={12}>
                  <Col span={6}><Form.Item name="ac_code" label="Account Code" rules={[{ required: true }]}><Input disabled={!!editing} /></Form.Item></Col>
                  <Col span={10}><Form.Item name="desc" label="Name" rules={[{ required: true }]}><Input /></Form.Item></Col>
                  <Col span={8}>
                    <Form.Item name="g_code" label="Group">
                      <Select allowClear placeholder="Select group"
                        options={groups.map(g => ({ value: g.g_code, label: `${g.g_code} - ${g.g_desc}` }))} />
                    </Form.Item>
                  </Col>
                </Row>
                <Row gutter={12}>
                  <Col span={8}><Form.Item name="opb" label="Opening Balance" initialValue={0}><InputNumber style={{ width: '100%' }} precision={2} /></Form.Item></Col>
                  <Col span={8}><Form.Item name="rate" label="Gold Rate/Wt" initialValue={0}><InputNumber style={{ width: '100%' }} precision={3} /></Form.Item></Col>
                  <Col span={4}><Form.Item name="cat" label="Category"><Input maxLength={1} /></Form.Item></Col>
                  <Col span={4}><Form.Item name="rmk" label="Remark"><Input maxLength={3} /></Form.Item></Col>
                </Row>
              </Form>
            ),
          },
          {
            key: 'address', label: 'Address & Contact',
            children: (
              <Form form={addrForm} layout="vertical" size="small">
                <Row gutter={12}>
                  <Col span={8}><Form.Item name="add1" label="Address 1"><Input /></Form.Item></Col>
                  <Col span={8}><Form.Item name="add2" label="Address 2"><Input /></Form.Item></Col>
                  <Col span={8}><Form.Item name="add3" label="Address 3"><Input /></Form.Item></Col>
                </Row>
                <Row gutter={12}>
                  <Col span={6}><Form.Item name="pin" label="PIN"><Input /></Form.Item></Col>
                  <Col span={6}><Form.Item name="mobile" label="Mobile"><Input /></Form.Item></Col>
                  <Col span={6}><Form.Item name="tel_r1" label="Tel (Res)"><Input /></Form.Item></Col>
                  <Col span={6}><Form.Item name="tel_o1" label="Tel (Office)"><Input /></Form.Item></Col>
                </Row>
                <Row gutter={12}>
                  <Col span={8}><Form.Item name="panno" label="PAN No"><Input /></Form.Item></Col>
                  <Col span={8}><Form.Item name="tinno" label="TIN No"><Input /></Form.Item></Col>
                </Row>
                <Row gutter={12}>
                  <Col span={12}><Form.Item name="lst" label="LST No"><Input /></Form.Item></Col>
                  <Col span={12}><Form.Item name="cst" label="CST No"><Input /></Form.Item></Col>
                </Row>
              </Form>
            ),
          },
        ]} />
      </Modal>
    </div>
  )
}
