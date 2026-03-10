import { useEffect, useState } from 'react'
import { Row, Col, Card, Statistic, Table, Tag, Typography, Spin } from 'antd'
import { GoldOutlined, SwapOutlined, ShoppingCartOutlined, LineChartOutlined } from '@ant-design/icons'
import { itemsApi, issuesApi, salesApi, ratesApi, type Rate } from '../api/client'

const { Title } = Typography

export default function Dashboard() {
  const [loading, setLoading] = useState(true)
  const [stats, setStats] = useState({
    totalItems: 0,
    openIssues: 0,
    totalSales: 0,
  })
  const [latestRate, setLatestRate] = useState<Rate | null>(null)
  const [recentIssues, setRecentIssues] = useState<any[]>([])

  useEffect(() => {
    Promise.all([
      itemsApi.list({ per_page: 1 }),
      issuesApi.list({ per_page: 5, tag: 'O' }),
      salesApi.list({ per_page: 1 }),
      ratesApi.latest().catch(() => null),
    ]).then(([items, issues, sales, rate]) => {
      setStats({
        totalItems: items.data.total,
        openIssues: issues.data.total,
        totalSales: sales.data.total,
      })
      setLatestRate(rate?.data ?? null)
      setRecentIssues(issues.data.data)
    }).finally(() => setLoading(false))
  }, [])

  if (loading) return <Spin size="large" style={{ display: 'block', margin: '100px auto' }} />

  return (
    <div>
      <Title level={4} style={{ marginBottom: 24 }}>Overview</Title>

      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Total Items"
              value={stats.totalItems}
              prefix={<GoldOutlined />}
              valueStyle={{ color: '#d4a017' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Open Issues"
              value={stats.openIssues}
              prefix={<SwapOutlined />}
              valueStyle={{ color: '#1677ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Total Sales"
              value={stats.totalSales}
              prefix={<ShoppingCartOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Gold Rate (Today)"
              value={latestRate?.rate ?? '—'}
              prefix={<LineChartOutlined />}
              suffix={latestRate ? '₹/g' : ''}
              valueStyle={{ color: '#fa8c16' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
        <Col xs={24} lg={12}>
          <Card title="Recent Issues">
            <Table
              dataSource={recentIssues}
              rowKey="t_no"
              size="small"
              pagination={false}
              columns={[
                { title: 'Voucher', dataIndex: 't_no', key: 't_no' },
                { title: 'Party', dataIndex: 'party_name', key: 'party_name' },
                { title: 'Date', dataIndex: 't_date', key: 't_date',
                  render: (d: string) => d ? new Date(d).toLocaleDateString('en-IN') : '—' },
                { title: 'Net Wt', dataIndex: 'net_wt', key: 'net_wt',
                  render: (v: number) => `${v?.toFixed(3)} g` },
                { title: 'Status', dataIndex: 'tag', key: 'tag',
                  render: (t: string) => <Tag color={t === 'O' ? 'orange' : 'green'}>{t === 'O' ? 'Open' : 'Closed'}</Tag> },
              ]}
            />
          </Card>
        </Col>

        {latestRate && (
          <Col xs={24} lg={12}>
            <Card title="Current Gold Rates">
              <Row gutter={16}>
                <Col span={8}>
                  <Statistic title="22K Rate" value={latestRate.rate} suffix="₹/g" />
                </Col>
                <Col span={8}>
                  <Statistic title="18K Rate" value={latestRate.rate1} suffix="₹/g" />
                </Col>
                <Col span={8}>
                  <Statistic title="Silver Rate" value={latestRate.s_rate} suffix="₹/g" />
                </Col>
              </Row>
              <div style={{ marginTop: 16, color: '#888' }}>
                As of {new Date(latestRate.date).toLocaleDateString('en-IN')}
              </div>
            </Card>
          </Col>
        )}
      </Row>
    </div>
  )
}
