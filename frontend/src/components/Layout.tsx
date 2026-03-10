import { useState } from 'react'
import { Layout as AntLayout, Menu, Typography, theme } from 'antd'
import {
  DashboardOutlined,
  GoldOutlined,
  SwapOutlined,
  InboxOutlined,
  ShoppingCartOutlined,
  ShoppingOutlined,
  TeamOutlined,
  AppstoreOutlined,
  LineChartOutlined,
} from '@ant-design/icons'
import { useNavigate, useLocation } from 'react-router-dom'

const { Sider, Content, Header } = AntLayout
const { Title } = Typography

const menuItems = [
  { key: '/dashboard',  icon: <DashboardOutlined />,   label: 'Dashboard' },
  { key: '/items',      icon: <GoldOutlined />,         label: 'Items' },
  { key: '/issues',     icon: <SwapOutlined />,         label: 'Issues' },
  { key: '/receipts',   icon: <InboxOutlined />,        label: 'Receipts' },
  { key: '/sales',      icon: <ShoppingCartOutlined />, label: 'Sales' },
  { key: '/purchases',  icon: <ShoppingOutlined />,     label: 'Purchases' },
  { key: '/accounts',   icon: <TeamOutlined />,         label: 'Accounts' },
  { key: '/lots',       icon: <AppstoreOutlined />,     label: 'Lots' },
  { key: '/rates',      icon: <LineChartOutlined />,    label: 'Gold Rates' },
]

export default function Layout({ children }: { children: React.ReactNode }) {
  const [collapsed, setCollapsed] = useState(false)
  const navigate = useNavigate()
  const location = useLocation()
  const { token } = theme.useToken()

  const selectedKey = '/' + location.pathname.split('/')[1]

  return (
    <AntLayout style={{ minHeight: '100vh' }}>
      <Sider
        collapsible
        collapsed={collapsed}
        onCollapse={setCollapsed}
        theme="dark"
        width={220}
      >
        <div style={{
          padding: '16px',
          textAlign: 'center',
          borderBottom: '1px solid rgba(255,255,255,0.1)',
          marginBottom: 8,
        }}>
          {!collapsed && (
            <Title level={4} style={{ color: '#ffd700', margin: 0 }}>
              VAL Inventory
            </Title>
          )}
          {collapsed && <GoldOutlined style={{ color: '#ffd700', fontSize: 24 }} />}
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[selectedKey]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>

      <AntLayout>
        <Header style={{
          background: token.colorBgContainer,
          padding: '0 24px',
          display: 'flex',
          alignItems: 'center',
          borderBottom: `1px solid ${token.colorBorderSecondary}`,
        }}>
          <Title level={5} style={{ margin: 0, color: token.colorTextSecondary }}>
            {menuItems.find(m => m.key === selectedKey)?.label ?? 'VAL Inventory'}
          </Title>
        </Header>
        <Content style={{ margin: 24, minHeight: 280 }}>
          {children}
        </Content>
      </AntLayout>
    </AntLayout>
  )
}
