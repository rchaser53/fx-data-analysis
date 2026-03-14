export interface Trade {
  id: number;
  trade_time: string;
  lot_size: number;
  purchase_rate: number;
  created_at: string;
  updated_at: string;
}

export interface CreateTradeRequest {
  trade_time: string;
  lot_size: number;
  purchase_rate: number;
}

export interface UpdateTradeRequest {
  trade_time?: string;
  lot_size?: number;
  purchase_rate?: number;
}

export interface USDJPYRate {
  date: string;
  pair: string;
  bid: number;
  ask: number;
  open: number;
  high: number;
  low: number;
  diff: number;
  close: number;
}

export interface USDJPYRatesResponse {
  pair: string;
  rates: USDJPYRate[];
}
