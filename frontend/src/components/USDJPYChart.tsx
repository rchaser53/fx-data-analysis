import React from 'react';
import { USDJPYRate } from '../types/trade';

interface USDJPYChartProps {
  points: USDJPYRate[];
}

const WIDTH = 980;
const HEIGHT = 320;
const PADDING = 40;

const USDJPYChart: React.FC<USDJPYChartProps> = ({ points }) => {
  if (points.length === 0) {
    return <div style={styles.message}>表示できるレートデータがありません</div>;
  }

  const closePrices = points.map((p) => p.close);
  const min = Math.min(...closePrices);
  const max = Math.max(...closePrices);
  const span = Math.max(max - min, 0.0001);
  const top = max + span * 0.1;
  const bottom = min - span * 0.1;
  const chartWidth = WIDTH - PADDING * 2;
  const chartHeight = HEIGHT - PADDING * 2;

  const xForIndex = (index: number) => {
    if (points.length === 1) {
      return WIDTH / 2;
    }
    return PADDING + (chartWidth * index) / (points.length - 1);
  };

  const yForValue = (value: number) => {
    return PADDING + ((top - value) / (top - bottom)) * chartHeight;
  };

  const polylinePoints = points
    .map((point, index) => `${xForIndex(index)},${yForValue(point.close)}`)
    .join(' ');

  const firstDate = points[0].date;
  const lastDate = points[points.length - 1].date;

  return (
    <div style={styles.wrapper}>
      <svg viewBox={`0 0 ${WIDTH} ${HEIGHT}`} width="100%" height="auto" role="img" aria-label="USDJPY close chart">
        <rect x={0} y={0} width={WIDTH} height={HEIGHT} fill="#ffffff" rx={10} />

        <line x1={PADDING} y1={PADDING} x2={PADDING} y2={HEIGHT - PADDING} stroke="#cfd6de" strokeWidth={1} />
        <line
          x1={PADDING}
          y1={HEIGHT - PADDING}
          x2={WIDTH - PADDING}
          y2={HEIGHT - PADDING}
          stroke="#cfd6de"
          strokeWidth={1}
        />

        <polyline fill="none" stroke="#0d6efd" strokeWidth={3} points={polylinePoints} />

        {points.map((point, index) => (
          <circle
            key={`${point.date}-${index}`}
            cx={xForIndex(index)}
            cy={yForValue(point.close)}
            r={3}
            fill="#0d6efd"
          />
        ))}

        <text x={PADDING} y={20} fontSize="13" fill="#2b2f33">{`High ${max.toFixed(3)}`}</text>
        <text x={PADDING} y={HEIGHT - 14} fontSize="13" fill="#2b2f33">{`Low ${min.toFixed(3)}`}</text>
        <text x={PADDING} y={HEIGHT - 4} fontSize="12" fill="#59636e">{firstDate}</text>
        <text x={WIDTH - PADDING} y={HEIGHT - 4} fontSize="12" fill="#59636e" textAnchor="end">{lastDate}</text>
      </svg>
    </div>
  );
};

const styles = {
  wrapper: {
    backgroundColor: '#ffffff',
    borderRadius: '10px',
    border: '1px solid #e2e8f0',
    padding: '12px',
    boxShadow: '0 2px 6px rgba(0, 0, 0, 0.06)',
  },
  message: {
    padding: '20px',
    textAlign: 'center' as const,
    color: '#6c757d',
    backgroundColor: '#ffffff',
    borderRadius: '8px',
    border: '1px solid #e2e8f0',
  },
};

export default USDJPYChart;
