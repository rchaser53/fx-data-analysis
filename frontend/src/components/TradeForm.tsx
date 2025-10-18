import React, { useState } from 'react';
import { CreateTradeRequest, Trade } from '../types/trade';

interface TradeFormProps {
  onSubmit: (data: CreateTradeRequest) => Promise<void>;
  initialData?: Trade;
  onCancel?: () => void;
}

const TradeForm: React.FC<TradeFormProps> = ({ onSubmit, initialData, onCancel }) => {
  const [tradeTime, setTradeTime] = useState(
    initialData?.trade_time 
      ? new Date(initialData.trade_time).toISOString().slice(0, 16)
      : new Date().toISOString().slice(0, 16)
  );
  const [lotSize, setLotSize] = useState(initialData?.lot_size?.toString() || '');
  const [purchaseRate, setPurchaseRate] = useState(initialData?.purchase_rate?.toString() || '');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      await onSubmit({
        trade_time: new Date(tradeTime).toISOString(),
        lot_size: parseFloat(lotSize),
        purchase_rate: parseFloat(purchaseRate),
      });

      // Reset form if it's a create form
      if (!initialData) {
        setTradeTime(new Date().toISOString().slice(0, 16));
        setLotSize('');
        setPurchaseRate('');
      }
    } catch (error) {
      console.error('Failed to submit:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} style={styles.form}>
      <div style={styles.formGroup}>
        <label style={styles.label}>
          取引日時:
          <input
            type="datetime-local"
            value={tradeTime}
            onChange={(e) => setTradeTime(e.target.value)}
            required
            style={styles.input}
          />
        </label>
      </div>

      <div style={styles.formGroup}>
        <label style={styles.label}>
          ロット数:
          <input
            type="number"
            step="0.01"
            value={lotSize}
            onChange={(e) => setLotSize(e.target.value)}
            required
            min="0.01"
            style={styles.input}
          />
        </label>
      </div>

      <div style={styles.formGroup}>
        <label style={styles.label}>
          購入レート:
          <input
            type="number"
            step="0.01"
            value={purchaseRate}
            onChange={(e) => setPurchaseRate(e.target.value)}
            required
            min="0.01"
            style={styles.input}
          />
        </label>
      </div>

      <div style={styles.buttonGroup}>
        <button type="submit" disabled={isSubmitting} style={styles.submitButton}>
          {isSubmitting ? '送信中...' : initialData ? '更新' : '作成'}
        </button>
        {onCancel && (
          <button type="button" onClick={onCancel} style={styles.cancelButton}>
            キャンセル
          </button>
        )}
      </div>
    </form>
  );
};

const styles = {
  form: {
    display: 'flex',
    flexDirection: 'column' as const,
    gap: '16px',
    padding: '20px',
    backgroundColor: '#f9f9f9',
    borderRadius: '8px',
    marginBottom: '20px',
  },
  formGroup: {
    display: 'flex',
    flexDirection: 'column' as const,
    gap: '8px',
  },
  label: {
    display: 'flex',
    flexDirection: 'column' as const,
    gap: '4px',
    fontSize: '14px',
    fontWeight: 'bold' as const,
  },
  input: {
    padding: '8px 12px',
    fontSize: '14px',
    border: '1px solid #ccc',
    borderRadius: '4px',
  },
  buttonGroup: {
    display: 'flex',
    gap: '8px',
  },
  submitButton: {
    padding: '10px 20px',
    backgroundColor: '#007bff',
    color: 'white',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
    fontSize: '14px',
    fontWeight: 'bold' as const,
  },
  cancelButton: {
    padding: '10px 20px',
    backgroundColor: '#6c757d',
    color: 'white',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
    fontSize: '14px',
  },
};

export default TradeForm;
