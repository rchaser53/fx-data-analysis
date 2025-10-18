import React, { useState } from 'react';
import { Trade } from '../types/trade';

interface TradeListProps {
  trades: Trade[];
  onEdit: (trade: Trade) => void;
  onDelete: (id: number) => void;
}

const TradeList: React.FC<TradeListProps> = ({ trades, onEdit, onDelete }) => {
  const [deletingId, setDeletingId] = useState<number | null>(null);

  const handleDelete = async (id: number) => {
    if (window.confirm('この取引を削除してもよろしいですか？')) {
      setDeletingId(id);
      try {
        await onDelete(id);
      } catch (error) {
        console.error('Failed to delete:', error);
      } finally {
        setDeletingId(null);
      }
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString('ja-JP', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (trades.length === 0) {
    return <div style={styles.emptyState}>取引データがありません</div>;
  }

  return (
    <div style={styles.container}>
      <table style={styles.table}>
        <thead>
          <tr>
            <th style={styles.th}>ID</th>
            <th style={styles.th}>取引日時</th>
            <th style={styles.th}>ロット数</th>
            <th style={styles.th}>購入レート</th>
            <th style={styles.th}>操作</th>
          </tr>
        </thead>
        <tbody>
          {trades.map((trade) => (
            <tr key={trade.id} style={styles.tr}>
              <td style={styles.td}>{trade.id}</td>
              <td style={styles.td}>{formatDate(trade.trade_time)}</td>
              <td style={styles.td}>{trade.lot_size.toFixed(2)}</td>
              <td style={styles.td}>{trade.purchase_rate.toFixed(2)}</td>
              <td style={styles.td}>
                <button
                  onClick={() => onEdit(trade)}
                  style={styles.editButton}
                >
                  編集
                </button>
                <button
                  onClick={() => handleDelete(trade.id)}
                  disabled={deletingId === trade.id}
                  style={styles.deleteButton}
                >
                  {deletingId === trade.id ? '削除中...' : '削除'}
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

const styles = {
  container: {
    overflowX: 'auto' as const,
  },
  table: {
    width: '100%',
    borderCollapse: 'collapse' as const,
    backgroundColor: 'white',
    boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
  },
  th: {
    padding: '12px',
    textAlign: 'left' as const,
    backgroundColor: '#f8f9fa',
    borderBottom: '2px solid #dee2e6',
    fontWeight: 'bold' as const,
  },
  tr: {
    borderBottom: '1px solid #dee2e6',
  },
  td: {
    padding: '12px',
  },
  editButton: {
    padding: '6px 12px',
    marginRight: '8px',
    backgroundColor: '#28a745',
    color: 'white',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
    fontSize: '13px',
  },
  deleteButton: {
    padding: '6px 12px',
    backgroundColor: '#dc3545',
    color: 'white',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
    fontSize: '13px',
  },
  emptyState: {
    padding: '40px',
    textAlign: 'center' as const,
    color: '#6c757d',
    fontSize: '16px',
  },
};

export default TradeList;
