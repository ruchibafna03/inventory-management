import { Routes, Route, Navigate } from 'react-router-dom'
import Layout from './components/Layout'
import Dashboard from './pages/Dashboard'
import Items from './pages/Items'
import Issues from './pages/Issues'
import Receipts from './pages/Receipts'
import Sales from './pages/Sales'
import Purchases from './pages/Purchases'
import Accounts from './pages/Accounts'
import Lots from './pages/Lots'
import Rates from './pages/Rates'

export default function App() {
  return (
    <Layout>
      <Routes>
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/items/*" element={<Items />} />
        <Route path="/issues/*" element={<Issues />} />
        <Route path="/receipts/*" element={<Receipts />} />
        <Route path="/sales/*" element={<Sales />} />
        <Route path="/purchases/*" element={<Purchases />} />
        <Route path="/accounts/*" element={<Accounts />} />
        <Route path="/lots/*" element={<Lots />} />
        <Route path="/rates/*" element={<Rates />} />
      </Routes>
    </Layout>
  )
}
