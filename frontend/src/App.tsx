import React, { useState, useEffect } from 'react';
import TradeList from './components/TradeList';
import TradeForm from './components/TradeForm';
import USDJPYChart from './components/USDJPYChart';
import { tradesApi } from './api/trades';
import { Trade, CreateTradeRequest, USDJPYRate } from './types/trade';

const App: React.FC = () => {
  const [trades, setTrades] = useState<Trade[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [editingTrade, setEditingTrade] = useState<Trade | null>(null);
  const [rates, setRates] = useState<USDJPYRate[]>([]);
  const [ratesLoading, setRatesLoading] = useState(true);
  const [ratesError, setRatesError] = useState<string | null>(null);

  const fetchTrades = async () => {
    try {
      setLoading(true);
      const data = await tradesApi.getAll();
      setTrades(data);
      setError(null);
    } catch (err) {
      setError('取引データの取得に失敗しました');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const fetchUSDJPYRates = async () => {
    try {
      setRatesLoading(true);
      const data = await tradesApi.getUSDJPYRates();
      setRates(data.rates || []);
      setRatesError(null);
    } catch (err) {
      setRatesError('USD/JPYレートデータの取得に失敗しました');
      console.error(err);
    } finally {
      setRatesLoading(false);
    }
  };

  useEffect(() => {
    fetchTrades();
    fetchUSDJPYRates();
  }, []);

  const handleCreate = async (data: CreateTradeRequest) => {
    try {
      await tradesApi.create(data);
      await fetchTrades();
    } catch (err) {
      console.error('Failed to create trade:', err);
      throw err;
    }
  };

  const handleUpdate = async (data: CreateTradeRequest) => {
    if (!editingTrade) return;

    try {
      await tradesApi.update(editingTrade.id, data);
      setEditingTrade(null);
      await fetchTrades();
    } catch (err) {
      console.error('Failed to update trade:', err);
      throw err;
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await tradesApi.delete(id);
      await fetchTrades();
    } catch (err) {
      console.error('Failed to delete trade:', err);
      throw err;
    }
  };

  const handleEdit = (trade: Trade) => {
    setEditingTrade(trade);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handleCancelEdit = () => {
    setEditingTrade(null);
  };

  return (
    <div style={styles.container}>
      <header style={styles.header}>
        <h1 style={styles.title}>FX取引データ管理</h1>
      </header>

      <main style={styles.main}>
        <section style={styles.section}>
          <h2 style={styles.sectionTitle}>USD/JPY 為替変動</h2>
          {ratesLoading && <div style={styles.message}>レートデータを読み込み中...</div>}
          {ratesError && <div style={styles.error}>{ratesError}</div>}
          {!ratesLoading && !ratesError && <USDJPYChart points={rates} />}
        </section>

        <section style={styles.section}>
          <h2 style={styles.sectionTitle}>
            {editingTrade ? '取引データ編集' : '新規取引データ作成'}
          </h2>
          <TradeForm
            onSubmit={editingTrade ? handleUpdate : handleCreate}
            initialData={editingTrade || undefined}
            onCancel={editingTrade ? handleCancelEdit : undefined}
          />
        </section>

        <section style={styles.section}>
          <h2 style={styles.sectionTitle}>取引データ一覧</h2>
          {loading && <div style={styles.message}>読み込み中...</div>}
          {error && <div style={styles.error}>{error}</div>}
          {!loading && !error && (
            <TradeList
              trades={trades}
              onEdit={handleEdit}
              onDelete={handleDelete}
            />
          )}
        </section>
      </main>
    </div>
  );
};

const styles = {
  container: {
    minHeight: '100vh',
    backgroundColor: '#f5f5f5',
  },
  header: {
    backgroundColor: '#007bff',
    color: 'white',
    padding: '20px',
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
  },
  title: {
    margin: 0,
    fontSize: '24px',
    textAlign: 'center' as const,
  },
  main: {
    maxWidth: '1200px',
    margin: '0 auto',
    padding: '20px',
  },
  section: {
    marginBottom: '40px',
  },
  sectionTitle: {
    fontSize: '20px',
    marginBottom: '16px',
    color: '#333',
  },
  message: {
    padding: '20px',
    textAlign: 'center' as const,
    color: '#6c757d',
  },
  error: {
    padding: '12px',
    backgroundColor: '#f8d7da',
    color: '#721c24',
    borderRadius: '4px',
    marginBottom: '16px',
  },
};

export default App;
