import axios from 'axios'

const client = axios.create({
  baseURL: '/api/v1',
  headers: { 'Content-Type': 'application/json' },
})

export interface PagedResult<T> {
  data: T[]
  total: number
  page: number
  per_page: number
}

// ─── Items ───────────────────────────────────────────────────────────────────

export interface Item {
  itcd: string
  desc: string
  purity?: string
  gross_wt: number
  net_wt: number
  gold_wt: number
  kundan_wt: number
  nakash_wt: number
  stone_wt: number
  prl_wt: number
  rby_ct: number
  rby_gm: number
  eme_ct: number
  eme_gm: number
  plk_ct: number
  plk_gm: number
  mk_chrg: number
  ac_code?: string
  lot_no?: string
  issue_no?: string
  cat?: string
  tag?: string
  narr1?: string
  narr2?: string
  created_at: string
  updated_at: string
}

export const itemsApi = {
  list: (params?: Record<string, string | number>) =>
    client.get<PagedResult<Item>>('/items', { params }),
  get: (itcd: string) => client.get<Item>(`/items/${itcd}`),
  create: (data: Partial<Item>) => client.post<Item>('/items', data),
  update: (itcd: string, data: Partial<Item>) => client.put<Item>(`/items/${itcd}`, data),
  delete: (itcd: string) => client.delete(`/items/${itcd}`),
  tags: () => client.get<string[]>('/items/tags'),
}

// ─── Accounts ────────────────────────────────────────────────────────────────

export interface AccountGroup {
  g_code: string
  g_desc: string
}

export interface Account {
  ac_code: string
  desc: string
  g_code?: string
  opb: number
  amt: number
  rate: number
  cat?: string
  rmk?: string
}

export interface Address {
  ac_code: string
  add1?: string
  add2?: string
  add3?: string
  pin?: string
  mobile?: string
  tel_r1?: string
  tel_o1?: string
  panno?: string
  tinno?: string
  lst?: string
  cst?: string
}

export const accountsApi = {
  listGroups: () => client.get<AccountGroup[]>('/groups'),
  createGroup: (data: Partial<AccountGroup>) => client.post<AccountGroup>('/groups', data),
  updateGroup: (gCode: string, data: Partial<AccountGroup>) =>
    client.put<AccountGroup>(`/groups/${gCode}`, data),
  deleteGroup: (gCode: string) => client.delete(`/groups/${gCode}`),

  list: (params?: Record<string, string | number>) =>
    client.get<PagedResult<Account>>('/accounts', { params }),
  get: (acCode: string) => client.get<Account>(`/accounts/${acCode}`),
  getAddress: (acCode: string) => client.get<Address>(`/accounts/${acCode}/address`),
  create: (data: Partial<Account>) => client.post<Account>('/accounts', data),
  update: (acCode: string, data: Partial<Account>) =>
    client.put<Account>(`/accounts/${acCode}`, data),
  upsertAddress: (acCode: string, data: Partial<Address>) =>
    client.put<Address>(`/accounts/${acCode}/address`, data),
  delete: (acCode: string) => client.delete(`/accounts/${acCode}`),
}

// ─── Rates ───────────────────────────────────────────────────────────────────

export interface Rate {
  id: number
  date: string
  rate: number
  rate1: number
  s_rate: number
}

export const ratesApi = {
  list: () => client.get<Rate[]>('/rates'),
  latest: () => client.get<Rate>('/rates/latest'),
  create: (data: Partial<Rate>) => client.post<Rate>('/rates', data),
  delete: (id: number) => client.delete(`/rates/${id}`),
}

// ─── Issues ──────────────────────────────────────────────────────────────────

export interface IssueDetail {
  id?: number
  itcd?: string
  itdesc?: string
  nos: number
  gross_wt: number
  net_wt: number
  stone_wt: number
  kundan_wt: number
  nakash_wt: number
  narr?: string
}

export interface Issue {
  t_no: string
  t_date: string
  due_date?: string
  ac_code?: string
  party_name?: string
  nos: number
  gross_wt: number
  net_wt: number
  net_pure: number
  pure_final: number
  bill_amt: number
  rate: number
  narr?: string
  tag?: string
  details?: IssueDetail[]
}

export const issuesApi = {
  list: (params?: Record<string, string | number>) =>
    client.get<PagedResult<Issue>>('/issues', { params }),
  get: (tNo: string) => client.get<Issue>(`/issues/${tNo}`),
  create: (data: Partial<Issue> & { details?: IssueDetail[] }) =>
    client.post<Issue>('/issues', data),
  update: (tNo: string, data: Partial<Issue> & { details?: IssueDetail[] }) =>
    client.put<Issue>(`/issues/${tNo}`, data),
  delete: (tNo: string) => client.delete(`/issues/${tNo}`),
}

// ─── Receipts ─────────────────────────────────────────────────────────────────

export interface ReceiptDetail {
  id?: number
  itcd?: string
  itdesc?: string
  lot_no?: string
  gross_wt: number
  net_wt: number
  gold_wt: number
  stone_wt: number
  kundan_wt: number
  nakash_wt: number
}

export interface Receipt {
  t_no: string
  t_date: string
  ac_code?: string
  party_name?: string
  nos: number
  gross_wt: number
  net_wt: number
  net_pure: number
  pure_final: number
  bill_amt: number
  rate: number
  narr?: string
  details?: ReceiptDetail[]
}

export const receiptsApi = {
  list: (params?: Record<string, string | number>) =>
    client.get<PagedResult<Receipt>>('/receipts', { params }),
  get: (tNo: string) => client.get<Receipt>(`/receipts/${tNo}`),
  create: (data: Partial<Receipt> & { details?: ReceiptDetail[] }) =>
    client.post<Receipt>('/receipts', data),
  update: (tNo: string, data: Partial<Receipt> & { details?: ReceiptDetail[] }) =>
    client.put<Receipt>(`/receipts/${tNo}`, data),
  delete: (tNo: string) => client.delete(`/receipts/${tNo}`),
}

// ─── Sales ───────────────────────────────────────────────────────────────────

export interface Sale {
  vouch_no: string
  vouch_date: string
  ac_code?: string
  name?: string
  nos: number
  gross_wt: number
  net_wt: number
  net_pure: number
  pure_final: number
  discount: number
  narr?: string
  phone_no?: string
}

export const salesApi = {
  list: (params?: Record<string, string | number>) =>
    client.get<PagedResult<Sale>>('/sales', { params }),
  get: (vouchNo: string) => client.get<Sale>(`/sales/${vouchNo}`),
  create: (data: Partial<Sale>) => client.post<Sale>('/sales', data),
  update: (vouchNo: string, data: Partial<Sale>) =>
    client.put<Sale>(`/sales/${vouchNo}`, data),
  delete: (vouchNo: string) => client.delete(`/sales/${vouchNo}`),
}

// ─── Lots ────────────────────────────────────────────────────────────────────

export interface Lot {
  t_no: string
  t_date: string
  ac_code?: string
  party?: string
  lot_no: string
  nos: number
  gross_wt: number
  adj_wt: number
  amt: number
  g_rt: number
  desc?: string
}

export const lotsApi = {
  list: (params?: Record<string, string>) =>
    client.get<Lot[]>('/lots', { params }),
  get: (tNo: string) => client.get<Lot>(`/lots/${tNo}`),
  create: (data: Partial<Lot>) => client.post<Lot>('/lots', data),
  update: (tNo: string, data: Partial<Lot>) => client.put<Lot>(`/lots/${tNo}`, data),
  delete: (tNo: string) => client.delete(`/lots/${tNo}`),
}

export default client
