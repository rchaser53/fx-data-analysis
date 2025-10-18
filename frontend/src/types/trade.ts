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
